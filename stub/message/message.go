package message

import (
	"../modules"
	"./types"
)

/*
 * HandleMessage() function handles commands sent from the server
 * First byte of the buffer is the message ID, the rest of it is the payload
 * This function returns a byte array containing the response that should be sent back to the server
 */
func HandleMessage(messageID byte, message interface{}) (responseID byte, response interface{}) {

	/// Messages with no payload
	/// Ping
	if messageID == types.REQ_PING {
		return types.RES_PING, nil
	}

	/// Basic system info
	if messageID == types.REQ_BASIC_SYSTEM_INFO {
		return types.RES_BASIC_SYSTEM_INFO, modules.GetSystemInfo()
	}

	// The message wasn't handled for some reason
	return types.ERR_MESSAGE_NOT_HANDLED, nil
}