package mock

import (
	"context"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"
)

// StatusHandlerStub -
type StatusHandlerStub struct {
	ErrorEncounteredHandler        func(err error)
	ProblematicBLSKeysFoundHandler func(messages []core.OutputMessage)
	ExecuteHandler                 func(ctx context.Context) error
	SendCloseMessageHandler        func()
}

// ErrorEncountered -
func (stub *StatusHandlerStub) ErrorEncountered(err error) {
	if stub.ErrorEncounteredHandler != nil {
		stub.ErrorEncounteredHandler(err)
	}
}

// ProblematicBLSKeysFound -
func (stub *StatusHandlerStub) ProblematicBLSKeysFound(messages []core.OutputMessage) {
	if stub.ProblematicBLSKeysFoundHandler != nil {
		stub.ProblematicBLSKeysFoundHandler(messages)
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
