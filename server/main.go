package main

import (
	"./networking"
	"log"
)

func main() {
	port := 4422
	err := networking.Listen(port)
	if err != nil {
		log.Fatalf("Error listening on port %d: %s", port, err)
	}
}