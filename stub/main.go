package main

import (
	"apollo/stub/networking"
	"apollo/types"
	"log"
	"time"
)

func main() {
	types.Register() // Register all of the types
	for {
		err := networking.Connect("localhost", 4422)
		if err == nil {
			// No errors, we can exit the process
			return
		}
		// Print out the error and reconnect
		log.Printf("Lost connection to the server: %s Trying again in 5 seconds", err)
		time.Sleep(5 * time.Second)
	}
}
