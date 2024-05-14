package disabled

import (
	"context"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"
)

type disabledStatusHandler struct {
}

// NewDisabledStatusHandler will create a new instance of type disabledStatusHandler
func NewDisabledStatusHandler() *disabledStatusHandler {
	return &disabledStatusHandler{}
}

// NotifyAppStart does nothing
func (disabled *disabledStatusHandler) NotifyAppStart() {
}

// ErrorEncountered does nothing
func (disabled *disabledStatusHandler) ErrorEncountered(_ error) {
}

// CollectKeysProblems does nothing
func (disabled *disabledStatusHandler) CollectKeysProblems(_ []core.OutputMessage) {
}

// Execute returns nil
func (disabled *disabledStatusHandler) Execute(_ context.Context) error {
	return nil
}

// SendCloseMessage does nothing
func (disabled *disabledStatusHandler) SendCloseMessage() {
}

// IsInterfaceNil returns true if there is no value under the interface
func (disabled *disabledStatusHandler) IsInterfaceNil() bool {
	return disabled == nil
}
