package types

import (
	"encoding/gob"
	"time"
)

// Registers all of the data types for gob
func Register() {
	gob.Register(BasicSystemInfo{})
	gob.Register(SystemInfo{})
	gob.Register(RunCommandRequest{})
	gob.Register(RunCommandReponse{})
}

type RunCommandRequest struct {
	Command	  string
	Backround bool	// If this is true the command will be ran in the background
}
type RunCommandReponse struct {
	Success	 bool
	Response string
}

// Basic system info
type BasicSystemInfo struct {
	OS				string
	InstallDate		time.Time
	Username		string
	Administrator	bool
	Language		int
	MachineID		string
}

type SystemInfo struct {
	Username		string
	InstallDate		time.Time
	OS				string
	OSVersion		string
	Administrator	bool
	ClientVersion	int
	DeviceName		string
	Language		int

	MBRam			int		// MB of ram installed
	CoreCount		int
	LogicalProcessorCount	int
	Architecture	string
	CPU				string

	GPU				string

	Device			string		// Info about device name and model

	MachineID		string
}