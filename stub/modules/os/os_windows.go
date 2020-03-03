//+build windows

package os

import (
	"fmt"
	"github.com/StackExchange/wmi"
	"github.com/denisbrodbeck/machineid"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/tools/go/ssa/interp/testdata/src/runtime"
	"log"
	"math"
	"os"
	user2 "os/user"
)

/*
 * Returns the OS name from the ProductName registry key
 * Example: Windows 10 Pro
 */
func OSName() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		log.Printf("Error opening reg key to get OS name: %s", err)
		return ""
	}
	defer k.Close()

	pn , _, err := k.GetStringValue("ProductName")
	if err != nil {
		log.Printf("Error reading ProductName: %s", err)
		return ""
	}
	return pn
}

/*
 * Returns the windows releaseID.currentBuild
 * Example: 1809.17763
 */
func OSVersion() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		log.Printf("Error opening reg key to get OS name: %s", err)
		return ""
	}
	defer k.Close()

	rid , _, err := k.GetStringValue("ReleaseId")
	if err != nil {
		log.Printf("Error getting ReleaseId: %s", err)
		return ""
	}

	cb, _, err := k.GetStringValue("CurrentBuild")
	if err != nil {
		log.Printf("Error getting CurrentBuild: %s", err)
		return rid		// If there is an error with getting the CurrentBuild, just return the release ID
	}

	return fmt.Sprintf("%s.%s", rid, cb)		// Return releaseID.currentBuild
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
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")		// Open the physical drive
	return err == nil			// If there is no error, return true
}

/*
 * Returns the default language ID
 * To decode these IDs, use this link:
 * https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-lcid/a9eac961-e77d-41a6-90a5-ce1a8b0cdb9c?redirectedfrom=MSDN
 */
func Language() int {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Nls\Language`, registry.QUERY_VALUE)
	if err != nil {
		log.Printf("Error opening reg key to get language: %s", err)
		return 0
	}
	defer k.Close()

	lang, _, err := k.GetIntegerValue("Default")
	if err != nil {
		log.Printf("Error reading Default for language: %s", err)
		return 0
	}
	if lang > math.MaxInt32 {		// If for some reason the value is above the max size of an int32, return 0 so we don't panic while casting
		return 0
	}
	return int(lang)
}

/*
 * Gets the computer's hostname
 */
func Hostname() string {
	host, err := os.Hostname()
	if err != nil {
		return ""
	}
	return host
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
	return cpuInfo.String()
}

// This struct is needed to get data from the wmi
type gpuInfo struct { Name string }
/*
 * Returns info about the GPU
 */
func Gpu() string {
	var gpuinfo []gpuInfo
	err := wmi.Query("Select * from Win32_VideoController", &gpuinfo)
	if err != nil {
		return ""
	}
	if len(gpuinfo) == 0 {	// Return nothing if we get nothing returned
		return ""
	}
	return gpuinfo[0].Name
}

type computerInfo struct { Vendor string; Name string }
/*
 * Returns info about the current device
 */
func Device() string {
	var computerinfo []computerInfo
	err := wmi.Query("Select * from Win32_ComputerSystemProduct", &computerinfo)
	if err != nil {
		return ""
	}
	if len(computerinfo) == 0 {	// Return nothing if we get nothing returned
		return ""
	}
	ci := computerinfo[0]		// Shorthand for our respons
	return fmt.Sprintf("%s %s", ci.Vendor, ci.Name)

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