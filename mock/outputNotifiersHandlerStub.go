package mock

import (
	"github.com/multiversx/mx-chain-keys-monitor-go/core"
)

// OutputNotifiersHandlerStub -
type OutputNotifiersHandlerStub struct {
	NotifyWithRetryHandler func(caller string, messages ...core.OutputMessage)
}

// NotifyWithRetry -
func (stub *OutputNotifiersHandlerStub) NotifyWithRetry(caller string, messages ...core.OutputMessage) {
	if stub.NotifyWithRetryHandler != nil {
		stub.NotifyWithRetryHandler(caller, messages...)
	}
}

// IsInterfaceNil -
func (stub *OutputNotifiersHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}
