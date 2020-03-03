package cli

import (
	"fmt"
	"strings"
)

// The command type
type command struct {
	Name     string                // The Name of the command
	Aliases  []string              // List of aliases the command can be referred to by
	MinArgs  int                   // The number of args required for this command. Note that 0 args means just the command
	Help     string                // The help string for the command which will be printed if there is an error or the help command is used
	Usage	 string				   // The usage for the command which will be printed out if the command is invalid or the help command is used
	Function func([]string) string // The function which will be ran. This function takes the args and returns the output
}

var commandList = []command{
	{
		Name:    "test",
		Aliases: []string{"p"},
		MinArgs: 0,
		Help:    "Prints out test and any args",
		Usage:	 "test [args...]",
		Function: func(args []string) string {
			return fmt.Sprintf("test %v", args)
		},
	},
}

func getHelpString(args []string) string {
	if len(args) >= 1 {		// If we have more than one args, find the command and print its help
		for _, command := range commandList {
			if command.Name == strings.ToLower(args[0]) {
				return fmt.Sprintf("%s %v: %s\nUsage: %s", command.Name, command.Aliases, command.Help, command.Usage)
			}
		}
		return fmt.Sprintf("Unknown command: %s", args[0])
	}

	helpString := "Available commands:\n"
	for _, command := range commandList {
		helpString += "----------------------------------------\n"
		helpString += fmt.Sprintf("%s %v: %s\nUsage: %s\n", command.Name, command.Aliases, command.Help, command.Usage)
	}
	return helpString
}

// This function is here so we can add the help command which refrences itself
func InitCommands() {
	commandList = append(commandList, command{
		Name:     "help",
		Aliases:  []string{"h", "?"},
		MinArgs:  0,
		Help:	  "Prints out help for all commands or a specified command",
		Usage:    "help [command]",
		Function: getHelpString,
	},)
}