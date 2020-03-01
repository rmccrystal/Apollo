package modules

import (
	"fmt"
	"os/exec"
	"runtime"
)

/*
 * Runs a console command
 * if `background` is true the command will be ran in the background and it will be immidately returned
 */
func RunCommand(command string, background bool) (success bool, response string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("sh", "-c", command)
	} else {
		return false, fmt.Sprintf("invalid os: %s", runtime.GOOS)
	}
	if !background { // If we're not running the background, run the command and get its output
		out, err := cmd.CombinedOutput()
		if err != nil {
			return false, err.Error()
		}
		return true, fmt.Sprintf("%s", out)
	} else {	// Run the command in the background
		err := cmd.Start()
		if err != nil {
			return false, err.Error()
		}
		return true, ""
	}
}
