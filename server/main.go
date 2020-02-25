package main

import (
	"./client/types"
	"./networking"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.InfoLevel)
	types.Register()
	port := 4422
	err := networking.Listen(port)
	if err != nil {
		log.Fatalf("Error listening on port %d: %s", port, err)
	}
}