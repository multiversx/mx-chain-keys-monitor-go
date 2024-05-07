package mock

import (
	"github.com/multiversx/mx-chain-keys-monitor-go/core"
)

// RatingsCheckerStub -
type RatingsCheckerStub struct {
	CheckHandler func(statistics map[string]*core.ValidatorStatistics, extraBLSKeys []string) ([]core.CheckResponse, error)
}

// Check -
func (stub *RatingsCheckerStub) Check(statistics map[string]*core.ValidatorStatistics, extraBLSKeys []string) ([]core.CheckResponse, error) {
	if stub.CheckHandler != nil {
		return stub.CheckHandler(statistics, extraBLSKeys)
	}

	return nil, nil
}

// IsInterfaceNil -
func (stub *RatingsCheckerStub) IsInterfaceNil() bool {
	return stub == nil
}
