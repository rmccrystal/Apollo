package client

import "net"

var ConnectedClients []Client

type Client struct {
	Conn		net.Conn

}

func (c Client) Ping() error {
	c.Conn.Write()
}