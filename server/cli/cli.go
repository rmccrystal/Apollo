package cli

import (
	"../client"
	"bufio"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

const commandPrefix = "> "

var CliList = make(map[Cli]bool) // CliList set

type Cli struct {
	Writer io.ReadWriteCloser
	au     aurora.Aurora // This is to enable and disable the colors
}

/*
 * Listens on the specified port
 * If the password is correct, a new cli will be created
 * with the raw socket.
 *
 * TODO: This code is really messy. Take some time to clean it up
 */
func ListenRaw(port int, password string) error {
	log.Debugf("Attempting to listen for Clis on port %d", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	log.Printf("Listening for clis on port %d", port)

	defer lis.Close() // Close the listener when we're done
	for {
		conn, err := lis.Accept()
		log.Debugf("New cli connecting: %s", conn.RemoteAddr())
		if err != nil {
			log.Errorf("Error accepting connection from %s", conn.RemoteAddr())
			continue
		}
		go func(c net.Conn) {
			if _, err := c.Write([]byte("password: ")); err != nil {
				log.Debugf("error writing to cli connection: %s", err)
				c.Close()
				return
			}
			if err := c.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
				log.Debugf("error setting read deadline: %s", err)
				c.Close()
				return
			}
			passwordAttempt, err := bufio.NewReader(c).ReadString('\n')
			passwordAttempt = strings.ReplaceAll(passwordAttempt, "\r", "") // Remove \r and the last element in the string
			passwordAttempt = strings.ReplaceAll(passwordAttempt, "\n", "")
			if err != nil {
				log.Debugf("error reading from cli connection: %s", err)
				c.Close()
				return
			}
			if passwordAttempt != password {
				log.Debugf("Cli client %s used the wrong password: %s", c.RemoteAddr(), passwordAttempt)
				c.Write([]byte("Incorrect password."))
				_, _ = c.Read(nil) // Wait for input and close the conn
				c.Close()
				return
			}
			log.Printf("New Cli connected: %s", c.RemoteAddr())
			// If we get here we got a successful connection with the right password
			err = c.SetReadDeadline(time.Now().Add(2 * time.Hour)) // Two hour timeout
			if err != nil {
				log.Debugf("error setting read deadline: %s", err)
				c.Close()
				return
			}

			err = NewCli(c, true) // Create the cli
			if err != nil {
				log.Debugf("error creating a cli client for %s: %s", c.RemoteAddr(), err)
				c.Close()
				return
			}
		}(conn)
	}
}

/*
 * This function will start a Cli to manage all of the clients
 * If the enableColor arg is true, colors will be enabled
 */
func NewCli(writer io.ReadWriteCloser, enableColor bool) error {
	cli := Cli{Writer: writer, au: aurora.NewAurora(enableColor)}
	if enableColor { // Set the title to Apollo only if enableColor is true
		cli.SetTitle("Apollo")
	}
	CliList[cli] = true  // Add the client to the client list
	go cli.messageLoop() // Start the message loop
	return nil
}

func (c Cli) messageLoop() {
	c.onConnect()
	for {
		c.Print(commandPrefix)
		text, err := c.read()
		if err != nil {
			log.Debugf("error reading from cli: %s", err)
			c.remove()
			return
		}
		// Split the text into args separated by spaces
		args := strings.Fields(text)
		if len(args) == 0 { // Continue if we get a length of 0 for our args
			continue
		}
		cmd := args[0]  // The cmd is the first element of the args
		args = args[1:] // Remove the first element from the args
		c.handleCommand(cmd, args)
	}
}

/*
 * Handles the specified `cmd` with arguments `args`
 * Returns if the command was successfully handled or not
 */
func (c Cli) handleCommand(cmd string, args []string) bool {
	for _, command := range commandList {
		// If the command has the same Name as the inputted command or its aliases
		if command.Name == strings.ToLower(cmd) || stringInSlice(strings.ToLower(cmd), command.Aliases) {
			if len(args) < command.MinArgs { // if we have too little args
				if command.Usage != "" { // If we have a usage
					c.Printf("Usage: %s", command.Usage)
				} else {
					c.Printf("%s takes a minimum of %d args", command.Name, command.MinArgs)
				}
				return false // We have too little args, return false
			}
			c.Printf(command.Function(c, args)) // Run the command
			return true                         // The command was successful, return true
		}
	}

	c.Printf("%s: command not found. Use 'help' to list available commands", cmd)
	return false
}

func (c Cli) onConnect() {
	c.Printf(c.au.Black(`

      :::     :::::::::   ::::::::  :::        :::        :::::::: 
    :+: :+:   :+:    :+: :+:    :+: :+:        :+:       :+:    :+:
   +:+   +:+  +:+    +:+ +:+    +:+ +:+        +:+       +:+    +:+
  +#++:++#++: +#++:++#+  +#+    +:+ +#+        +#+       +#+    +:+
  +#+     +#+ +#+        +#+    +#+ +#+        +#+       +#+    +#+
  #+#     #+# #+#        #+#    #+# #+#        #+#       #+#    #+#
  ###     ### ###         ########  ########## ########## ######## 
`).BgWhite().String())
	c.Printf("") // print a blank new line
}

// Util functions
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

/*
 * Returns a list of clients from a capture
 * A capture is either a client number or "all"
 * TODO: Add a capture to select multiple client IDs
 */
func getClientsFromCapture(capture string) ([]*client.Client, error) {
	var clients []*client.Client

	if strings.ToLower(capture) == "all" {	// If we want to use all clients
		clients = client.GetOnlineClients()	// Get all online clients
		if len(clients) == 0 {		// If there are no online clients
			return nil, errors.New("No clients are online")
		}
	} else {		// Else get the client by the ID
		id, err := strconv.Atoi(capture)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error: %s is not a number", capture))
		}
		cl := client.GetClientById(id)
		if cl == nil {
			return nil, errors.New(fmt.Sprintf("Client with ID %d not found", id))
		}
		clients = append(clients, cl)
	}
	return clients, nil
}