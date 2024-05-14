package mock

import (
	"github.com/multiversx/mx-chain-keys-monitor-go/core"
)

// OutputNotifiersHandlerStub -
type OutputNotifiersHandlerStub struct {
	NotifyWithRetryHandler func(caller string, messages ...core.OutputMessage) error
}

// NotifyWithRetry -
func (stub *OutputNotifiersHandlerStub) NotifyWithRetry(caller string, messages ...core.OutputMessage) error {
	if stub.NotifyWithRetryHandler != nil {
		return stub.NotifyWithRetryHandler(caller, messages...)
	}

	return nil
}

// IsInterfaceNil -
func (stub *OutputNotifiersHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}
