package types

import (
	"encoding/gob"
	"time"
)

// Registers all of the data types for gob
func Register() {
	gob.Register(BasicSystemInfo{})
	gob.Register(SystemInfo{})
}

// Basic system info
type BasicSystemInfo struct {
	OS				string
	InstallDate		time.Time
	Username		string
	Administrator	bool
	Language		int
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
}