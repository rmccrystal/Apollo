package client

import (
	"./types"
)

/*
 * the Client.Ping() function pings the client.
 * If the ping is successful, the function will return nil
 * If it is not, it will return an error and the client will be reomved
 * from the client list in the Client.Send() function
 */
func (c Client) Ping() error {
	return c.SendMessage(types.REQ_PING, nil, nil, types.RES_PING)
}

/*
 * Returns a struct containing basic system info of the client
 */
func (c Client) GetBasicSystemInfo() (types.BasicSystemInfo, error) {
	var response types.BasicSystemInfo
	err := c.SendMessage(types.REQ_BASIC_SYSTEM_INFO, nil, &response, types.RES_BASIC_SYSTEM_INFO)
	if err == nil {	// if there is no error, cache the basic system info
		c.BasicSystemInfo = response
	}
	return response, err
}

/*
 * Returns a struct containing all system info of the client
 */
func (c Client) GetSystemInfo() (types.SystemInfo, error) {
	var response types.SystemInfo
	err := c.SendMessage(types.REQ_SYSTEM_INFO, nil, &response, types.RES_SYSTEM_INFO)
	if err == nil {	// if there is no error, cache the system info
		c.SystemInfo = response
	}
	return response, err
}