package cli

import (
	"../client"
	"fmt"
	"github.com/logrusorgru/aurora"
	"strings"
	"time"
)

// The command type
type command struct {
	Name     string                // The Name of the command
	Aliases  []string              // List of aliases the command can be referred to by
	MinArgs  int                   // The number of args required for this command. Note that 0 args means just the command
	Help     string                // The help string for the command which will be printed if there is an error or the help command is used
	Usage	 string				   // The usage for the command which will be printed out if the command is invalid or the help command is used
	Function func(Cli, []string) string // The function which will be ran. This function takes cli refrence and the args and returns the output
}

var commandList = []command{
	{
		Name:     "clear",
		Aliases:  []string{"cls"},
		MinArgs:  0,
		Help:     "Clears the screen",
		Function: func(c Cli, args []string) string {
			c.Clear()
			return ""
		},
	},
	{
		Name:     "exit",
		Aliases:  []string{"quit"},
		MinArgs:  0,
		Help:     "Exists the cli",
		Function: func(c Cli, args []string) string {
			c.remove()
			return ""
		},
	},
	{
		Name:     "clients",
		Aliases:  []string{"list", "bots", "c"},
		MinArgs:  0,
		Help:     "Lists the connected clients",
		Function: func(c Cli, args []string) string {
			if len(client.Clients) == 0 {
				return "No clients"
			}
			str := ""
			for client, connected := range client.Clients {
				if !connected {
					str += c.au.Gray(11, fmt.Sprintf("%s (not connected)", client.String())).String()
				} else {
					str += client.String()
				}
				str += "\n"
			}
			return str
		},
	},


	{
		Name:     "ping",
		Aliases:  nil,
		MinArgs:  1,
		Help:     "Pings a client and returns the response time",
		Usage:    "ping (clientID)",
		Function: func(c Cli, args []string) string {
			clients, err := getClientsFromCapture(args[0])
			if err != nil {
				return c.au.Red(err).String()
			}
			for _, client := range clients {
				start := time.Now()
				err := client.Ping()
				if err != nil {
					return c.au.Red(fmt.Sprintf("%s not connected", client.String())).String()
				}
				ms := time.Since(start).Nanoseconds()/1e6
				return c.au.Green(fmt.Sprintf("Client %s responded in %dms", client.String(), ms)).String()
			}
			return ""
		},
	},

}

/*
 * Gets the description of the command used for the help command
 * Takes an aurora object as an argument for color support
 */
func (c command) description(au aurora.Aurora) string {
	usageString := ""
	if len(c.Usage) > 0 {
		usageString = au.Gray(11, fmt.Sprintf("\n\t└─Usage: %s", c.Usage)).String()
	}
	text := au.Sprintf("%s: %s%s", au.Bold(c.Name), au.Gray(15, c.Help), usageString)
	return text
}

func getHelpString(c Cli, args []string) string {
	if len(args) >= 1 {		// If we have more than one args, find the command and print its help
		for _, command := range commandList {
			if command.Name == strings.ToLower(args[0]) {
				aliasesString := ""				// Only print aliases if we get help for specific command
				if len(command.Aliases) > 0 {
					aliasesString = aurora.Gray(11, fmt.Sprintf("\n\t└─Aliases: %v", command.Aliases)).String()
				}
				return command.description(c.au) + aliasesString
			}
		}
		return fmt.Sprintf("Unknown command: %s", args[0])
	}

	helpString := "Available commands:\n"
	for _, command := range commandList {
		helpString += command.description(c.au) + "\n"
	}
	return helpString
}

// This function is here so we can add the help command which refrences itself
func InitCommands() {
	commandList = append([]command{{
		Name:     "help",
		Aliases:  []string{"h", "?"},
		MinArgs:  0,
		Help:	  "Prints out help for all commands or a specified command",
		Usage:    "help [command]",
		Function: getHelpString,
	}}, commandList...)
}