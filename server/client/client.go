package client

import (
	"./types"
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var ConnectedClients = make(map[Client]bool)		// Client set

type Client struct {
	Conn			net.Conn
	mux				sync.Mutex
	IP				string
	BasicSystemInfo	types.BasicSystemInfo
	SystemInfo		types.SystemInfo
}

func (c Client) OnConnect() {
	ms, err := c.Ping()		// Ping the client to see if it is working
	if err != nil {
		log.Printf("Error pinging client with IP %s: %s", c.IP, err)
		return
	}
	ConnectedClients[c] = true		// Append the new client to the client list
	log.Printf("New client connected with IP %s and ping %d", c.IP, ms)
}

// Closes the socket and removes the client from the client list
func (c Client) Remove() {
	_ = c.Conn.Close()
	delete(ConnectedClients, c)
}

// Send data and wait for a response
func (c Client) Send(b []byte) ([]byte, error) {
	log.Printf("Sending %v to client %s", b, c.IP)
	b = append(b, '\n')	// Add a \n to the data so the client knows the message is finished

	c.mux.Lock()			// Lock the mutex so only one thread can access the client at once
	defer c.mux.Unlock()	// Unlock it when the function is done
	_, err := c.Conn.Write(b)	// Try to write data
	if err != nil {		// If there is an error, print an error and delete the client
		log.Printf("Error writing to client %s: %s Removing the client", c.IP, err)
		c.Remove()
		return nil, err
	}

	/// Get response
	err = c.Conn.SetReadDeadline(time.Now().Add(5 * time.Second)) // Make the read time out after 5 seconds
	reader := bufio.NewReader(c.Conn)		// Create a new reader
	resp, err := reader.ReadBytes('\n')
	if err != nil {
		// For some reason there was an error reading
		log.Printf("Error reading from client %s: %s Removing the client", c.IP, err)
		c.Remove()	// Remove the client
		return nil, err
	}
	// If we successfully read from the client, return the response and no error
	// The last item of the response is cut off because it is a \n
	log.Printf("Received %v from client %s", resp[:len(resp)-1], c.IP)
	return resp[:len(resp)-1], nil
}

// Send a message to the client with a message ID and data
// Returns the response ID, the response data, and any error
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
	// If we don't get any errors, return the first byte of the response as the responseID and the rest as the responseData
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
 * If the ping is successful, the function will return the time in ms it took to ping the client and nil
 * If it is not, it will return an error and the client will be removed
 * from the client list in the Client.Send() function
 */
func (c Client) Ping() (int, error) {
	start := time.Now()		// Time how long the ping takes
	responseID, _, err := c.SendMessage(types.REQ_PING, nil)
	timeMs := int(time.Now().Sub(start).Nanoseconds() / 1000000)
	if err != nil {		// If there is an error pinging the client, remove the client and return the error
		c.Remove()
		return 0, err
	}
	if responseID != types.RES_PING {		// Return an error if we don't get a RES_PING response
		log.Printf("For some reason, we did not get a ping repsonse when pinging client %s and instead got a response with ID %d", c.IP, responseID)
		c.Remove()
		return 0, errors.New("invalid response type")
	}
	return timeMs, nil	// If we get here we successfully pinged the client
}

/*
 * Returns a struct containing basic system info of the client
 */
func (c Client) GetBasicSystemInfo() (types.BasicSystemInfo, error) {
	if c.BasicSystemInfo != (types.BasicSystemInfo{}) {		// check if we have basic system info cached already
		if _, err := c.Ping(); err != nil {					// Ping the client in case it is offline
			return types.BasicSystemInfo{}, err				// If it is offline, return the error
		}
		return c.BasicSystemInfo, nil				// We have the system info cached, so just return that
	}
	responseID, responseData, err := c.SendMessage(types.REQ_BASIC_SYSTEM_INFO, nil)	// Send the request
	if err != nil {		// Return the error if there is one
		return types.BasicSystemInfo{}, err
	}
	if responseID == types.ERR_GOB_ENCODING {		// Return an error if there is an error with gob encoding on the client
		log.Printf("Error encoding buffer on client with IP %s", c.IP)
		return types.BasicSystemInfo{}, errors.New("error encoding struct on client")
	}
	if responseID != types.RES_BASIC_SYSTEM_INFO {
		log.Printf("For some reason, we did not get a repsonse when getting basic system info from client %s and instead got a response with ID %d", c.IP, responseID)
		return types.BasicSystemInfo{}, errors.New("invalid response type")
	}
	gobBuff := bytes.NewBuffer(responseData)
	tmpStruct := new(types.BasicSystemInfo)
	gobObj := gob.NewDecoder(gobBuff)
	err = gobObj.Decode(tmpStruct)
	if err != nil {
		log.Printf("error decoding gob buffer")
		return types.BasicSystemInfo{}, errors.New(fmt.Sprintf("error deserilzing data: %s", err))
	}
	c.BasicSystemInfo = *tmpStruct		// Cache the response so we don't have to request it again
	return *tmpStruct, nil
}