package networking

import (
	"../message"
	"../message/types"
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

func Connect(ip string, port uint16) error {
	log.Printf("Attempting to connect to %s:%d", ip, port)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}
	log.Printf("Successfully connected to %s:%d", ip, port)

	for {
		err = messageLoop(conn)
		if err != nil {
			return err
		}
	}
}

func messageLoop(conn net.Conn) error {	// Note: Only returns an error if there is something wrong with the communication itself
										// Will not return an error due to malformed data
	reader := bufio.NewReader(conn)
	buffer, err := reader.ReadBytes('\n')	// Read message until we encounter a new line
	if err != nil {
		return err
	}
	buffer = buffer[:len(buffer)-1]		// Remove the last \n from the buffer
	packet := handlePacket(buffer)
	_, err = conn.Write(append(packet, '\n'))		// Write the response we get from message.HandleMessage()
	if err != nil {		// Return an error if we can't write to the connection
		return err
	}
	return nil
}

/*
 * Handles a packet sent to the stub
 * Gets the raw bytes sent and returns the bytes to respond with
 */
func handlePacket(buffer []byte) []byte {
	if len(buffer) == 0 {
		return []byte{byte(types.ERR_MESSAGE_TOO_SMALL)}
	}
	messageID := buffer[0]

	var responseID byte
	var response interface{}

	if len(buffer) == 1 {	// If we have just a message ID
		responseID, response = message.HandleMessage(messageID, nil)
	} else {	// Else we must decode the buffer
		buff := bytes.NewBuffer(buffer[1:]) // Don't use the first element of the buffer
		decoder := gob.NewDecoder(buff)
		var msg interface{} // The structured message being sent
		err := decoder.Decode(&msg)
		if err != nil {
			log.Printf("error decoding: %s", err)
			return []byte{byte(types.ERR_GOB)}
		}

		fmt.Printf("%s", msg)	// TODO: remove
		responseID, response = message.HandleMessage(messageID, msg)
	}

	if response == nil {	// If there is no actual response data
		return []byte{responseID}	// Return just the response ID
	}
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(response)
	if err != nil {		// If there is an error encoding, return it
		log.Printf("error encoding: %s", err)
		return []byte{byte(types.ERR_GOB)}
	}

	return append([]byte{responseID}, buff.Bytes()...)	// Return our response
}