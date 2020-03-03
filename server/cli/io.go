package cli

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

/*
 * Writes a string to the Cli ReadWrite interface
 */
func (c Cli) writeString(text string) error {
	return c.write([]byte(text))
}

/*
 * Writes a byte array to the Cli ReadWrite interface
 */
func (c Cli) write(data []byte) error {
	_, err := c.Writer.Write(data)
	return err
}

/*
 * Reads from the Cli ReadWrite interface
 * This function will return the data returned
 * from the Cli and any errors returned
 *
 * TODO: Add a timeout to this function
 */
func (c Cli) read() (string, error) {
	str, err := bufio.NewReader(c.Writer).ReadString('\n')
	if err != nil {
		return "", err
	}
	str = strings.ReplaceAll(str, "\r", "") 	// Remove \r in the string
	str = strings.ReplaceAll(str, "\n", "") 	// Remove \r in the string
	return str, nil
}

/*
 * Closes the ReadWriteCloser and removes the cli from the list
 */
func (c Cli) remove() {
	_ = c.Writer.Close()
	CliList[c] = false
}

/*
 * Prints formatted text to the client.
 * If there is an error writing, the client
 * will be removed
 * Automatically prints a new line to the end of the string if it's not there
 */
func (c Cli) Printf(format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	if len(text) == 0 {
		return
	}
	if text[len(text)-1] != '\n' {	// Add a newline if it's not there
		text = text + "\n"
	}
	text = strings.ReplaceAll(text, "\n", "\r\n")		// Replaces newlines with \r\n
	err := c.writeString(text)
	if err != nil {
		log.Debugf("error writing to cli: %s", err)
		c.remove()
	}
}

/*
 * Prints the text to the client.
 * If there is an error writing, the client
 * will be removed
 */
func (c Cli) Print(data ...interface{}) {
	text := fmt.Sprint(data...)
	err := c.writeString(text)
	if err != nil {
		log.Debugf("error writing to cli: %s", err)
		c.remove()
	}
}

/*
 * Clears the cli by sending a bunch of newlines
 */
func (c Cli) Clear() {
	text := ""
	for i := 1; i<=50; i++ {
		text += "\n"
	}
	c.Print(text)
}