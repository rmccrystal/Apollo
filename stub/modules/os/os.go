package os

import (
	"fmt"
	"os/exec"
	"strings"
)

/*
 * Runs a console command
 * if `background` is true the command will be ran in the background and it will be immidately returned
 */
func RunCommand(command string, background bool) (success bool, response string) {
	commandArgs := strings.Split(command, " ")
	if len(commandArgs) == 0 {
		return false, "error: no command specified"
	}
	var cmd *exec.Cmd
	if len(commandArgs) == 1 {
		cmd = exec.Command(commandArgs[0])
	} else {
		cmd = exec.Command(commandArgs[0], commandArgs[1:]...)
	}
	if !background {	// If we're not running the background, run the command and get its output
		out, err := cmd.CombinedOutput()
		if err != nil {
			return false, err.Error()
		}
		return true, fmt.Sprintf("%s", out)
	} else {
		err := cmd.Start()
		if err != nil {
			return false, err.Error()
		}
		return true, ""
	}
}
