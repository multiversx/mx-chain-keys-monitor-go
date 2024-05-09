package mock

import (
	"github.com/multiversx/mx-chain-keys-monitor-go/core"
)

// OutputNotifierStub -
type OutputNotifierStub struct {
	OutputMessagesHandler func(messages ...core.OutputMessage)
}

// OutputMessages -
func (stub *OutputNotifierStub) OutputMessages(messages ...core.OutputMessage) {
	if stub.OutputMessagesHandler != nil {
		stub.OutputMessagesHandler(messages...)
	}
}

// IsInterfaceNil -
func (stub *OutputNotifierStub) IsInterfaceNil() bool {
	return stub == nil
}
