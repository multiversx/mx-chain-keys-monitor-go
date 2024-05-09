package mock

import (
	"context"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"
)

// StatusHandlerStub -
type StatusHandlerStub struct {
	NotifyAppStartHandler      func()
	ErrorEncounteredHandler    func(err error)
	CollectKeysProblemsHandler func(messages []core.OutputMessage)
	ExecuteHandler             func(ctx context.Context) error
	SendCloseMessageHandler    func()
}

// NotifyAppStart -
func (stub *StatusHandlerStub) NotifyAppStart() {
	if stub.NotifyAppStartHandler != nil {
		stub.NotifyAppStartHandler()
	}
}

// ErrorEncountered -
func (stub *StatusHandlerStub) ErrorEncountered(err error) {
	if stub.ErrorEncounteredHandler != nil {
		stub.ErrorEncounteredHandler(err)
	}
}

// CollectKeysProblems -
func (stub *StatusHandlerStub) CollectKeysProblems(messages []core.OutputMessage) {
	if stub.CollectKeysProblemsHandler != nil {
		stub.CollectKeysProblemsHandler(messages)
	}
}

// SendCloseMessage -
func (stub *StatusHandlerStub) SendCloseMessage() {
	if stub.SendCloseMessageHandler != nil {
		stub.SendCloseMessageHandler()
	}
}

// Execute -
func (stub *StatusHandlerStub) Execute(ctx context.Context) error {
	if stub.ExecuteHandler != nil {
		return stub.ExecuteHandler(ctx)
	}

	return nil
}

// IsInterfaceNil -
func (stub *StatusHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}
