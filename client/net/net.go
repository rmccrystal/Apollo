package net

import (
	"../proto"
	"../rpc"
	"fmt"
	"github.com/hashicorp/yamux"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

var lastHost string
var lastPort int

func ConnectToServer(host string, port int) {
	lastHost = host
	lastPort = port
	c, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Second * 5)
	if err != nil {
		log.Println(err)
		time.Sleep(5 * time.Second)
		ConnectToServer(lastHost, lastPort)
	}
	srvConn, err := yamux.Server(c, yamux.DefaultConfig())
	if err != nil {
		log.Println(err)
		time.Sleep(5 * time.Second)
		ConnectToServer(lastHost, lastPort)
	}
	
	// Create a grpc server
	grpcServer := grpc.NewServer()
	
	// Register the apollo object
	proto.RegisterApolloServer(grpcServer, &rpc.Server{})
	
	// start the gRPC server
	go func() {
		err := grpcServer.Serve(srvConn)
		if err != nil {
			fmt.Println(err)
			ConnectToServer(lastHost, lastPort)
		}
	}()
}