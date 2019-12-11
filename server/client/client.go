package client

import (
	"./types"
	"bufio"
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

var ConnectedClients = make(map[Client]bool)		// Client set

type Client struct {
	Conn		net.Conn
	mux			sync.Mutex
	IP			string
}

func (c Client) OnConnect() {
	err := c.Ping()		// Ping the client to see if it is working
	if err != nil {
		log.Printf("Error pinging client with IP %s: %s", c.IP, err)
		delete(ConnectedClients, c)		// If it doesn't ping, remove the client
	}
	log.Printf("New client connected with IP %s", c.IP)
}

// Send data and wait for a response
func (c Client) Send(b []byte) ([]byte, error) {
	b = append(b, '\n')	// Add a \n to the data so the client knows the message is finished

	c.mux.Lock()			// Lock the mutex so only one thread can access the client at once
	defer c.mux.Unlock()	// Unlock it when the function is done
	_, err := c.Conn.Write(b)	// Try to write data
	if err != nil {		// If there is an error, print an error and delete the client
		log.Printf("Error writing to client %s: %s Removing the client", c.IP, err)
		delete(ConnectedClients, c)
		return nil, err
	}

	/// Get response
	err = c.Conn.SetReadDeadline(time.Now().Add(5 * time.Second)) // Make the read time out after 5 seconds
	reader := bufio.NewReader(c.Conn)		// Create a new reader
	resp, err := reader.ReadBytes('\n')
	if err != nil {
		// For some reason there was an error reading
		log.Printf("Error reading from client %s: %s Removing the client", c.IP, err)
		delete(ConnectedClients, c)	// Remove the client
		return nil, err
	}
	// If we successfully read from the client, return the response and no error
	return resp, nil
}

// Send a message to the client with a message ID and data
// Returns the respnose ID, the response data, and any error
// Note that although the messageID is an int this function
// will return an error if it doesn't cast to a byte
func (c Client) SendMessage(messageID int, data []byte) (responseID byte, responseData []byte, err error) {
	// Check if messageID is too large
	if messageID > 255 {
		return 0, nil, errors.New("message ID is too large")
	}

	resp, err := c.Send(append([]byte{byte(messageID)}, data...))		// Create a new byte array with the messageID as the
																		// first byte and the data as the rest
	if err != nil {		// Return the error if we get one
		return 0, nil, err
	}
	// If we don't get any errors, return the first byte of the repsonse as the responseID and the rest as the responseData
	// If the response is only one byte, return just the response ID
	if len(resp) == 0 {	// If the length of the response is 0, return nothing
		return 0, nil, nil
	}
	if len(resp) == 1 {	// If the length of the response is 1, return just the message ID
		return resp[0], nil, nil
	}
	// Otherwise, return the respose ID and the data
	return resp[0], resp[1:], nil
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
		log.Printf("For some reason, we did not get a ping repsonse when pinging client %s and instead got a response with ID %d", c.IP, responseID)
		return errors.New("invalid response type")
	}
	return nil	// If we get here we successfully pinged the client
}