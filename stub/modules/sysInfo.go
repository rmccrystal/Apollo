package modules

import (
	"../message/types"
	"./os"
	"time"
)

var startTime = time.Now()

func GetBasicSystemInfo() types.BasicSystemInfo {
	return types.BasicSystemInfo{
		OS:            os.OSName(),
		InstallDate:   startTime,
		Username:      os.Username(),
		Administrator: os.Administrator(),
		Language:      os.Language(),
	}
}

func GetSystemInfo() types.SystemInfo {
	return types.SystemInfo{
		Username:              os.Username(),
		InstallDate:           startTime,
		OS:                    os.OSName(),
		OSVersion:			   os.OSVersion(),
		Administrator:         os.Administrator(),
		ClientVersion:         1,
		DeviceName:            os.Hostname(),
		Language:              os.Language(),
		MBRam:                 os.Ram(),
		CoreCount:             os.Cores(),
		LogicalProcessorCount: os.LogicalProcessors(),
		Architecture:          os.Arcitecture(),
		CPU:                   os.Cpu(),
		GPU:                   os.Gpu(),
		Device:				   os.Device(),
	}
}