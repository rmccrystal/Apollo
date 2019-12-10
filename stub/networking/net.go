package networking

import (
	"../message"
	"bufio"
	"fmt"
	"net"
)

var ServerConn net.Conn

func Connect(ip string, port uint16) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}

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
	go message.HandleMessage(buffer)	// Handle the message TODO: Add xor encryption
	return nil
}