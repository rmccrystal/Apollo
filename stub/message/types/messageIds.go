package types

/// Message IDS
const (
	// Requests
	REQ_PING = iota + 1
	REQ_SYS_INFO
	REQ_DOWNLOAD_EXECUTE

	// Responses
	RES_PING
	RES_SYS_INFO
	RES_DOWNLOAD_EXECUTE

	// Errors
	ERR_MESSAGE_TOO_SMALL
	ERR_MESSAGE_NOT_HANDLED
)

