package client

import (
	"./types"
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
	"time"
)

var ConnectedClients = make(map[Client]bool)		// Client set

// TODO: Add a consecutive ID
type Client struct {
	Conn		net.Conn
	mux			sync.Mutex
	IP			string
}

// This function runs when the client connects for the first time
func (c Client) OnConnect() {
	ConnectedClients[c] = true		// Append the new client to the client list
	err := c.Ping()		// Ping the client to see if it is working
	if err != nil {
		log.Printf("Error pinging client %s", c)
		c.Delete()	// If it doesn't ping, remove the client
	}
	log.Printf("New client connected with IP %s", c.IP)
	info, _ := c.GetBasicSystemInfo()
	log.Println(info)
}

// Deletes the client from the client list and closes the connection
func (c Client) Delete() {
	log.Printf("Client %s disconnected", c)
	delete(ConnectedClients, c)
	_ = c.Conn.Close()	// Close the connection
}

func (c Client) String() string {
	return c.IP
}

// Send data and wait for a response
func (c Client) Send(b []byte) ([]byte, error) {
	log.Debugf("Sending %v to client %s", b, c)
	b = append(b, '\n')	// Add a \n to the data so the client knows the message is finished

	c.mux.Lock()			// Lock the mutex so only one thread can access the client at once
	defer c.mux.Unlock()	// Unlock it when the function is done
	_, err := c.Conn.Write(b)	// Try to write data
	if err != nil {		// If there is an error, print an error and delete the client
		log.Debugf("Error writing to client %s: %s Removing the client", c, err)
		c.Delete()
		return nil, err
	}

	/// Get response
	err = c.Conn.SetReadDeadline(time.Now().Add(5 * time.Second)) // Make the read time out after 5 seconds
	reader := bufio.NewReader(c.Conn)		// Create a new reader
	resp, err := reader.ReadBytes('\n')
	if err != nil {
		// For some reason there was an error reading
		log.Debugf("Error reading from client %s: %s Removing the client", c, err)
		c.Delete()	// Remove the client
		return nil, err
	}
	// If we successfully read from the client, return the response and no error
	// The last item of the response is cut off because it is a \n
	log.Debugf("Received %v from client %s", resp[:len(resp)-1], c)
	return resp[:len(resp)-1], nil
}

// Send a message to the client with a message ID and data
// Returns the response ID, the response data, and any error
// Note that although the messageID is an int this function
// will return an error if it doesn't cast to a byte
func (c Client) SendMessage(messageID int, data interface{}) (responseID byte, response interface{}, err error) {
	// Check if messageID is too large
	if messageID > 255 {
		return 0, nil, errors.New("message ID is too large")
	}

	var payloadBuffer bytes.Buffer
	if data != nil {	// Add the payload to the buffer if there is a payload
		encoder := gob.NewEncoder(&payloadBuffer)
		err = encoder.Encode(data)
		if err != nil { // If there is an error encoding, return it
			log.Printf("error encoding: %s", err)
			return 0, nil, err
		}
	}

	resp, err := c.Send(append([]byte{byte(messageID)}, payloadBuffer.Bytes()...))		// Create a new byte array with the messageID as the
																		// first byte and the data as the rest
	if err != nil {		// Return the error if we get one
		return 0, nil, err
	}
	// If we don't get any errors, return the first byte of the response as the responseID and the rest as the responseData
	// If the response is only one byte, return just the response ID
	if len(resp) == 0 {	// If the length of the response is 0, return nothing
		return 0, nil, nil
	}
	if len(resp) == 1 {	// If the length of the response is 1, return just the message ID
		return resp[0], nil, nil
	}
	// Otherwise, deserialize the data and return that
	buff := bytes.NewBuffer(resp[1:]) // Don't use the first element of the buffer
	decoder := gob.NewDecoder(buff)
	var msg interface{} // The structured message being received
	err = decoder.Decode(&msg)
	if err != nil {
		log.Printf("error decoding: %s", err)
		return 0, nil, err
	}
	return resp[0], msg, nil
}

/*
 * the Client.Ping() function pings the client.
 * If the ping is successful, the function will return nil
 * If it is not, it will return an error and the client will be reomved
 * from the client list in the Client.Send() function
 */
func (c Client) Ping() error {
	responseID, _, err := c.SendMessage(types.REQ_PING, nil)
	if err != nil {		// If there is an error pinging the client, return it
		return err
	}
	if responseID != types.RES_PING {		// Return an error if we don't get a RES_PING response
		log.Debugf("For some reason, we did not get a ping response when pinging client %s and instead got a response with ID %d", c.IP, responseID)
		return errors.New("invalid response type")
	}
	return nil	// If we get here we successfully pinged the client
}

/*
 * Returns a struct containing basic system info of the client
 */
func (c Client) GetBasicSystemInfo() (types.SystemInfo, error) {
	responseID, response, err := c.SendMessage(types.REQ_BASIC_SYSTEM_INFO, nil)
	if err != nil {		// Return the error if there is one
		return types.SystemInfo{}, err
	}
	if responseID != types.RES_BASIC_SYSTEM_INFO {
		log.Debugf("For some reason, we did not get a response when getting basic system info from client %s and instead got a response with ID %d", c.IP, responseID)
		return types.SystemInfo{}, errors.New("invalid response type")
	}

	if info, ok := response.(types.SystemInfo); ok {
		return info, nil
	}
	return types.SystemInfo{}, errors.New("gob type error")
}