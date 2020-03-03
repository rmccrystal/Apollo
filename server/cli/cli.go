package cli

import (
	"bufio"
	"fmt"
	"github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"strings"
	"time"
)

const commandPrefix = "> "

var CliList = make(map[Cli]bool)	// CliList set

type Cli struct {
	Writer io.ReadWriteCloser
	au	   aurora.Aurora	// This is to enable and disable the colors
}

/*
 * Listens on the specified port
 * If the password is correct, a new cli will be created
 * with the raw socket.
 *
 * TODO: This code is really messy. Take some time to clean it up
 */
func ListenRaw(port int, password string) error {
	log.Debugf("Attempting to listen fro Clis on port %d", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	log.Printf("Listening for clis on port %d", port)

	defer lis.Close()       // Close the listener when we're done
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
			passwordAttempt = strings.ReplaceAll(passwordAttempt, "\r", "")	// Remove \r and the last element in the string
			passwordAttempt = strings.ReplaceAll(passwordAttempt, "\n", "")
			if err != nil {
				log.Debugf("error reading from cli connection: %s", err)
				c.Close()
				return
			}
			if passwordAttempt != password {
				log.Debugf("Cli client %s used the wrong password: %s", c.RemoteAddr(), passwordAttempt)
				c.Write([]byte("Incorrect password."))
				_, _ = c.Read(nil)		// Wait for input and close the conn
				c.Close()
				return
			}
			log.Printf("New Cli connected: %s", c.RemoteAddr())
			// If we get here we got a successful connection with the right password
			err = c.SetReadDeadline(time.Now().Add(2 * time.Hour))		// Two hour timeout
			if err := c.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
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
	cli := Cli{Writer:writer,au:aurora.NewAurora(enableColor)}
	if enableColor { // Set the title to Apollo only if enableColor is true
		err := cli.writeString("\033]0;Apollo\007")
		if err != nil {
			return err
		}
	}
	CliList[cli] = true	// Add the client to the client list
	go cli.messageLoop()			// Start the message loop
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
		if len(args) == 0 {		// Continue if we get a length of 0 for our args
			continue
		}
		cmd := args[0]			// The cmd is the first element of the args
		args = args[1:]			// Remove the first element from the args
		for _, command := range commandList {
			// If the command has the same Name as the inputted command or its aliases
			if command.Name == strings.ToLower(cmd) || stringInSlice(strings.ToLower(cmd), command.Aliases) {
				if len(args) < command.MinArgs {	// if we have too little args
					c.Printf("%s takes a minimum of %d args.\n\n%s",
						command.Name, command.MinArgs, command.Help)
				}
				c.Printf(command.Function(c, args)) // Run the command
			}
		}
	}
}

func (c Cli) onConnect() {
	c.Printf(`
    :::     :::::::::   ::::::::  :::        :::        :::::::: 
  :+: :+:   :+:    :+: :+:    :+: :+:        :+:       :+:    :+:
 +:+   +:+  +:+    +:+ +:+    +:+ +:+        +:+       +:+    +:+
+#++:++#++: +#++:++#+  +#+    +:+ +#+        +#+       +#+    +:+
+#+     +#+ +#+        +#+    +#+ +#+        +#+       +#+    +#+
#+#     #+# #+#        #+#    #+# #+#        #+#       #+#    #+#
###     ### ###         ########  ########## ########## ######## `)
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