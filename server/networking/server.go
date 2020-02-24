package networking

import (
	"../client"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

func Listen(port int) error {
	log.Debugf("Attempting to listen on port %d", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))	// Start listening
	if err != nil {
		return err
	}
	log.Printf("Listening on port %d", port)

	defer lis.Close()	// Close the listener when we're done
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Errorf("Error accepting connection from %s", conn.RemoteAddr())
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	log.Debugf("New connection: %s", conn.RemoteAddr())
	newClient := client.Client{Conn:conn, IP:conn.RemoteAddr().String()}
	newClient.OnConnect()		// Run the OnConnect function
}