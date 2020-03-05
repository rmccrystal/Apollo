package main

import (
	"./cli"
	"./client/types"
	"./networking"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	types.Register()
	cli.InitCommands()

	go func() {
		err := cli.ListenRaw(2300, "password420")
		if err != nil {
			log.Errorf("error listening for cli: %s", err)
		}
	}()
	port := 4422
	err := networking.Listen(port)
	if err != nil {
		log.Fatalf("Error listening on port %d: %s", port, err)
	}
}
