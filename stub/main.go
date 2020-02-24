package main

import (
	"./message/types"
	"./networking"
	"log"
	"time"
)

func main() {
	//log.Println(os.OSVersion())
	//log.Println(os.OSName())
	//log.Println(os.Username())
	//log.Println(os.Administrator())
	//log.Println(os.Ram())
	//log.Println(os.Cores())
	//log.Println(os.LogicalProcessors())
	//log.Println(os.Arcitecture())
	//log.Println(os.Cpu())
	//log.Println(os.Gpu())
	//log.Println(os.Device())
	types.Register()	// Register all of the types
	for {
		err := 	networking.Connect("localhost", 4422)
		if err == nil {
			// No errors, we can exit the process
			return
		}
		// Print out the error and reconnect
		log.Printf("Lost connection to the server: %s Trying again in 5 seconds", err)
		time.Sleep(5 * time.Second)
	}
}