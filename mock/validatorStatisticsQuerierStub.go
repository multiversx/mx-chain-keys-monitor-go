package mock

import (
	"context"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"
)

// ValidatorStatisticsQuerierStub -
type ValidatorStatisticsQuerierStub struct {
	QueryHandler func(ctx context.Context) (map[string]*core.ValidatorStatistics, error)
}

// Query -
func (stub *ValidatorStatisticsQuerierStub) Query(ctx context.Context) (map[string]*core.ValidatorStatistics, error) {
	if stub.QueryHandler != nil {
		return stub.QueryHandler(ctx)
	}

	return nil, nil
}

// IsInterfaceNil -
func (stub *ValidatorStatisticsQuerierStub) IsInterfaceNil() bool {
	return stub == nil
}
