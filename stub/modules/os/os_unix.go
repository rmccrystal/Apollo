//+build !windows

package os

import (
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"math"
	"os"
	user2 "os/user"
	"runtime"
)

var cachedInfoStat host.InfoStat
/*
 * Returns the host.InfoStat object from gopsutil
 * caches if it can
 */
func getHostInfo() host.InfoStat {
	if (cachedInfoStat != host.InfoStat{}) {
		return cachedInfoStat
	}
	info, err := host.Info()
	if err != nil {
		return host.InfoStat{}
	}
	cachedInfoStat = *info
	return *info
}

/*
 * Returns the OS name
 * Example: Windows 10 Pro
 */
func OSName() string {
	return getHostInfo().Platform
}

/*
 * returns the platform version
 */
func OSVersion() string {
	return getHostInfo().PlatformVersion
}

/*
 * Returns the current logged in username
 * If the username and the name are the same, it will return the username
 * If the name of the logged in user and the username are different
 * it will return name (username)
 */
func Username() string {
	user, err := user2.Current()
	if err != nil {
		return ""
	}
	if user.Username == user.Name {		// If the username and the name are the same, return just the username
		return user.Username
	}
	return fmt.Sprintf("%s (%s)", user.Name, user.Username)		// Else return the username and the name
}

/*
 * Checks if the current process is running as admin
 * TODO: This could be made better by not attempting to access a physical drive
 */
func Administrator() bool {
	return os.Getuid() == 0
}

// doesn't work on linux
func Language() int {
	return 0
}

/*
 * Gets the computer's hostname
 */
func Hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}

/*
 * Returns the MB of ram installed on the computer
 */
func Ram() int {
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	return int(vmem.Total/(uint64(math.Pow(2, 20))))
}

/*
 * Returns the number of physical cores
 */
func Cores() int {
	count, err := cpu.Counts(false)
	if err != nil {
		return 0		// Return 0 if there is an error
	}
	return count
}

/*
 * Returns the number of physical cores
 */
func LogicalProcessors() int {
	count, err := cpu.Counts(true)
	if err != nil {
		return 0		// Return 0 if there is an error
	}
	return count
}

/*
 * Returns the CPU architecture
 */
func Arcitecture() string {
	return runtime.GOARCH
}

/*
 * Returns a JSON containing CPU information
 */
func Cpu() string {
	cpuInfoArr, err := cpu.Info()
	if err != nil {
		return ""
	}
	if len(cpuInfoArr) == 0 {	// If there are no cores for some reason, return nothing
		return ""
	}
	cpuInfo := cpuInfoArr[0]	// cpu.Info returns an array for each core; get only the first one
	return cpuInfo.ModelName
}

/*
 * Returns info about the GPU
 */
func Gpu() string {
	return ""
}

/*
 * Returns info about the current device
 * right now only returns if it is a vm
 */
func Device() string {
	return getHostInfo().VirtualizationRole
}

/*
 * Returns a unique ID for the machine which can be used to
 * uniquely identify the computer
 */
func MachineID() string {
	id, err := machineid.ID()
	if err != nil {
		return fmt.Sprintf("error getting machine id: %s", err)
	}
	return id
}
