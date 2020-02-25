package client

import (
	"./types"
)

/*
 * the Client.Ping() function pings the client.
 * If the ping is successful, the function will return nil
 * If it is not, it will return an error and the client will be removed
 */
func (c Client) Ping() error {
	err := c.SendMessage(types.REQ_PING, nil, nil, types.RES_PING)
	if err != nil {
		c.Delete()
	}
	return err
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

/*
 * Runs a console command
 * if `background` is true the command will be ran in the background and it will be immidately returned
 */
func (c Client) RunCommand(command string, background bool) (success bool, response string, err error) {
	var res types.RunCommandReponse
	request := types.RunCommandRequest{
		Command:   command,
		Backround: background,
	}
	err = c.SendMessage(types.REQ_RUN_COMMAND, request, &res, types.RES_RUN_COMMAND)
	return res.Success, res.Response, err
}