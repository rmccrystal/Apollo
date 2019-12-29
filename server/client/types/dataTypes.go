package types

import "time"

// Basic system info
type BasicSystemInfo struct {
	OS				string
	InstallDate		time.Time
	Username		string
	Administrator	bool
	Language		int
}

type Display struct {		// Used for Displays in SystemInfo
	Width	int		// Width in pixels
	Height	int
	RefreshRate	int
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