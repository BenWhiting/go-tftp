package constants

// ErrorCodes provided by the RFC TFTP PROTOCOL page 10
var ErrorCodes = map[uint16]string{
	0: "Not defined, see error message (if any).",
	1: "File not found.",
	2: "Access violation.",
	3: "Disk full or allocation exceeded.",
	4: "Illegal TFTP operation.",
	5: "Unknown transfer ID.",
	6: "File already exists.",
	7: "No such user.",
}

// ServerStartMsg for the TFTP server
const ServerStartMsg string = "Listening on %s \n"

// SendErrMsg template
const SendErrMsg string = "Sending err %q to %v\n"

// FailedSendErrMsg template
const FailedSendErrMsg string = "Failed to send err to %v %s\n"

// SimpleErrMsg template
const SimpleErrMsg string = "Error: %s\n"

// ParsingPackErrMsg template
const ParsingPackErrMsg string = "Error parsing a packet: %s"

// FileRequestMsg prefix
const FileRequestMsg string = "File request: "

// ReceivedMsg template
const ReceivedMsg string = "[%d] Read received - %+v.\n"

// ReceivedErrMsg template
const ReceivedErrMsg string = "Received an error from %v: %d : %s \n"

// RRQFileNotFoundErrMsg template
const RRQFileNotFoundErrMsg = "[RRQ] File %s was not found. "

// WRQFileAlreadyExistsErrMsg template
const WRQFileAlreadyExistsErrMsg = "[WRQ] File %s already exists."

// ModeNotSupportedMsg template
const ModeNotSupportedMsg = "[%d] Mode %s is not supported."

// UnrecOpCodeMsg template
const UnrecOpCodeMsg string = "Unrecognized OpCode - %d."

// UnknownLogFlagMsg for log formatting
const UnknownLogFlagMsg string = "Unknown log flag, must be set to 1 - 7."

// FailedSendPacketMsg template
const FailedSendPacketMsg string = "Error: failed sending a packet. %s"

// UnsupportedTransderMode template
const UnsupportedTransderMode string = "Unsupported transder mode %q"

// ServerMessageWrapper that logs where the request comes from
const ServerMessageWrapper string = "%s From client - %s"

// UnknownTransferMsg log
const UnknownTransferMsg string = "Unknown transfer"

// FileTransferCompleteMsg template
const FileTransferCompleteMsg string = "Transfer is complete %s\n"

// FileReceivingCompleteMsg template
const FileReceivingCompleteMsg string = "Finished receiving %s from %v\n"

// TransferTimeoutMsg template
const TransferTimeoutMsg string = "Transfer %s timed out\n."

// RetryLastPktMsg template
const RetryLastPktMsg string = "Retrying last packet for %s"
