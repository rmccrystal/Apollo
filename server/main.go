package main

import "./net"

func main() {
	err := net.Listen("127.0.0.1", 7878)
	if err != nil {
		panic(err)
	}
	for{}
}