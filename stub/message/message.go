package message

import (
	"../modules"
	"./types"
	"bytes"
	"encoding/gob"
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

	/// Basic system info
	if messageID == types.REQ_BASIC_SYSTEM_INFO {
		responseBuffer := new(bytes.Buffer)		// Create a buffer for our encoding
		responseBuffer.Write([]byte{byte(types.RES_BASIC_SYSTEM_INFO)})		// Add the repsonse ID to the beginning of the response
		gobObj := gob.NewEncoder(responseBuffer)		// Create the gob encoder
		err := gobObj.Encode(modules.GetSystemInfo())
		if err != nil {		// Return a gob encoding error if we can't encode
			return []byte{byte(types.ERR_GOB_ENCODING)}
		}

		return responseBuffer.Bytes()	// Return the buffer
	}

	/// Messages with payload
	// payload := buffer[1:]

	// The message wasn't handled for some reason
	return []byte{byte(types.ERR_MESSAGE_NOT_HANDLED)}
}