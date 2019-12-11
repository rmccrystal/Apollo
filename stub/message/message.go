package message

import (
	"./types"
)

/*
 * HandleMessage() function handles commands sent from the server
 * First byte of the buffer is the message ID, the rest of it is the payload
 * This function returns a byte array containing the response that should be sent back to the server
 */
func HandleMessage(buffer []byte) []byte {
	if len(buffer) == 0 {
		return []byte{byte(types.ERR_MESSAGE_TOO_SMALL)}
	}
	messageID := buffer[0]

	/// Messages with no payload
	/// Ping
	if messageID == types.REQ_PING {
		return []byte{byte(types.RES_PING)}
	}

	/// Messages with payload
	//payload := buffer[1:]

	// The message wasn't handled for some reason
	return []byte{byte(types.ERR_MESSAGE_NOT_HANDLED)}
}