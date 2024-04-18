package factory

import "time"

// FileLoggingHandler will handle log file rotation
type FileLoggingHandler interface {
	ChangeFileLifeSpan(newDuration time.Duration, newSizeInMB uint64) error
	Close() error
	IsInterfaceNil() bool
}
