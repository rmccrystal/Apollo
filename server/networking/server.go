package networking

import (
	"../client"
	"fmt"
	"log"
	"net"
)

func Listen(port int) error {
	log.Printf("Attempting to listen on port %d", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))	// Start listening
	if err != nil {
		return err
	}
	log.Printf("Listening on port %d", port)

	defer lis.Close()	// Close the listener when we're done
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Printf("Error accepting connection from %s", conn.RemoteAddr())
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	log.Printf("New connection: %s", conn.RemoteAddr())
	newClient := client.Client{Conn:conn, IP:conn.RemoteAddr().String()}
	client.ConnectedClients[newClient] = true		// Append the new client to the client list
	newClient.OnConnect()		// Run the OnConnect function
}