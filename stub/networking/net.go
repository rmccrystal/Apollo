package networking

import (
	"../message"
	"bufio"
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
	_, err = conn.Write(append(message.HandleMessage(buffer), '\n'))		// Write the response we get from message.HandleMessage()
	if err != nil {		// Return an error if we can't write to the connection
		return err
	}
	return nil
}