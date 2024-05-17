package factory

import (
	"time"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"
)

// FileLoggingHandler will handle log file rotation
type FileLoggingHandler interface {
	ChangeFileLifeSpan(newDuration time.Duration, newSizeInMB uint64) error
	Close() error
	IsInterfaceNil() bool
}

// Monitor defines the monitor interface used in polling mechanism to check some conditions
type Monitor interface {
	Close() error
	IsInterfaceNil() bool
}

// OutputNotifiersHandler defines the behavior of a component that is able to notify all notifiers
type OutputNotifiersHandler interface {
	NotifyWithRetry(caller string, messages ...core.OutputMessage) error
	IsInterfaceNil() bool
}
