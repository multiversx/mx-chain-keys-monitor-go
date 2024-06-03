package executors

import (
	"context"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"
)

// RatingsChecker defines the operation of a component able to check the ratings
type RatingsChecker interface {
	Check(statistics map[string]*core.ValidatorStatistics, extraBLSKeys []string) ([]core.CheckResponse, error)
	IsInterfaceNil() bool
}

// ValidatorStatisticsQuerier defines the operations of a component able to query the validators statistics
type ValidatorStatisticsQuerier interface {
	Query(ctx context.Context) (map[string]*core.ValidatorStatistics, error)
	IsInterfaceNil() bool
}

// BLSKeysFetcher is able to get all staked BLS keys of an identity
type BLSKeysFetcher interface {
	GetAllBLSKeys(ctx context.Context, sender string) ([]string, error)
	IsInterfaceNil() bool
}

// StatusHandler defines the operations of a component able to keep the status of the app
type StatusHandler interface {
	NotifyAppStart()
	ErrorEncountered(err error)
	CollectKeysProblems(messages []core.OutputMessage)
	Execute(ctx context.Context) error
	SendCloseMessage()
	IsInterfaceNil() bool
}

// BLSKeysFilter is able to decide if a provided BLS key needs to be filtered out or not
type BLSKeysFilter interface {
	ShouldNotify(blsKey string) bool
	IsInterfaceNil() bool
}

// OutputNotifier defines the operations supported by an output notifier instance
type OutputNotifier interface {
	OutputMessages(messages ...core.OutputMessage) error
	Name() string
	IsInterfaceNil() bool
}

// OutputNotifiersHandler defines the behavior of a component that is able to notify all notifiers
type OutputNotifiersHandler interface {
	NotifyWithRetry(caller string, messages ...core.OutputMessage) error
	IsInterfaceNil() bool
}

// Executor defines the behavior of a component able to execute a certain task
type Executor interface {
	Execute(ctx context.Context) error
	IsInterfaceNil() bool
}
