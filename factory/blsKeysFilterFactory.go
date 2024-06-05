package factory

import (
	"time"

	"github.com/multiversx/mx-chain-keys-monitor-go/config"
	"github.com/multiversx/mx-chain-keys-monitor-go/executors"
	"github.com/multiversx/mx-chain-keys-monitor-go/executors/disabled"
)

// NewBLSKeysFilter creates a new instance of type BLSKeysFilter
func NewBLSKeysFilter(cfg config.AlarmSnoozeConfig) (executors.BLSKeysFilter, error) {
	if !cfg.Enabled {
		return disabled.NewDisabledBLSKeysFilter(), nil
	}

	args := executors.ArgsBlsKeysTimeCache{
		MaxSnoozeEvents:     cfg.NumNotificationsForEachFaultyKey,
		CacheExpiration:     cfg.SnoozeTimeInSec,
		GetCurrentTimestamp: time.Now().Unix,
	}

	return executors.NewBLSKeysTimeCache(args)
}
