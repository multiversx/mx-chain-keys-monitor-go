package mock

import (
	"github.com/multiversx/mx-chain-keys-monitor-go/core"
)

// OutputNotifierStub -
type OutputNotifierStub struct {
	NameHandler           func() string
	OutputMessagesHandler func(messages ...core.OutputMessage) error
}

// OutputMessages -
func (stub *OutputNotifierStub) OutputMessages(messages ...core.OutputMessage) error {
	if stub.OutputMessagesHandler != nil {
		return stub.OutputMessagesHandler(messages...)
	}

	return nil
}

// Name -
func (stub *OutputNotifierStub) Name() string {
	if stub.NameHandler != nil {
		return stub.NameHandler()
	}

	return ""
}

// IsInterfaceNil -
func (stub *OutputNotifierStub) IsInterfaceNil() bool {
	return stub == nil
}
