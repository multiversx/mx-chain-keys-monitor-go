package executors

import "sync/atomic"

// NumErrorsEncountered -
func (handler *statusHandler) NumErrorsEncountered() uint32 {
	return atomic.LoadUint32(&handler.numErrors)
}
