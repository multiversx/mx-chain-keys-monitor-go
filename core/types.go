package core

import "fmt"

// MessageOutputType defines the output type for a message
type MessageOutputType int

// defined constants for the MessageOutputType
const (
	InfoMessageOutputType    MessageOutputType = 3
	WarningMessageOutputType MessageOutputType = 4
	ErrorMessageOutputType   MessageOutputType = 5
)

// String converts the type to its string representation
func (messageOutputType MessageOutputType) String() string {
	switch messageOutputType {
	case InfoMessageOutputType:
		return "info"
	case WarningMessageOutputType:
		return "warn"
	case ErrorMessageOutputType:
		return "error"
	default:
		return fmt.Sprintf("unknown MessageOutputType: %d", messageOutputType)
	}
}
