package net

import (
	"../proto"
	"context"
	"fmt"
	"github.com/hashicorp/yamux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"log"
	"net"
	"time"
)

type Listener struct {
	listener net.Listener
	Host string
	Port int
	Clients []proto.ApolloClient
	NewClinetCallback func(client proto.ApolloClient)
}

func CreateListener(host string, port int) (Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return Listener{}, err
	}
	lisObj := Listener{Host:host, Port:port, listener:lis}
	go lisObj.listenLoop()
	return lisObj, nil
}

func (l *Listener) listenLoop() {
	lis := l.listener
	defer lis.Close()
	log.Println("Started listener")
	for {
		incoming, err := lis.Accept()
		if err != nil {
			log.Printf("error accepting a connection, continuing loop: %s\n", err)
			continue
		}
		log.Printf("New client connected: %s\n", incoming.RemoteAddr())
		incomingConn, err := yamux.Client(incoming, yamux.DefaultConfig())
		if err != nil {
			log.Printf("Error creating a yamux client, continuing loop: %s\n", err)
		}

		var conn *grpc.ClientConn

		conn, err = grpc.Dial(":7777", grpc.WithInsecure(),
			grpc.WithDialer(func(target string, timeout time.Duration) (net.Conn, error) {
				return incomingConn.Open()
			}),
		)

		if err != nil {
			log.Printf("Error connecting to client: %s\n", err)
		}

		go l.handleConn(conn)
	}
}

func (l *Listener) handleConn(conn *grpc.ClientConn) {
	defer conn.Close()
	client := proto.NewApolloClient(conn)
	l.Clients = append(l.Clients, client)
	l.NewClinetCallback(client)
	conn.WaitForStateChange(context.Background(), connectivity.Shutdown)
	l.Clients =
}