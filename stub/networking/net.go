package networking

import (
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

func messageLoop(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	bytes, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}

	
}

func handleMessage(messageId byte, buffer []byte) {

}