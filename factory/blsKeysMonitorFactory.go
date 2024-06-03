package factory

import (
	"time"

	"github.com/multiversx/mx-chain-keys-monitor-go/checkers"
	"github.com/multiversx/mx-chain-keys-monitor-go/config"
	"github.com/multiversx/mx-chain-keys-monitor-go/executors"
	"github.com/multiversx/mx-chain-keys-monitor-go/interactors"
	"github.com/multiversx/mx-chain-keys-monitor-go/monitor"
	"github.com/multiversx/mx-chain-keys-monitor-go/parsers"
	"github.com/multiversx/mx-sdk-go/core/http"
)

const timeBetweenBLSKeysFetch = time.Second

// NewBLSKeysMonitor will create a BLS keys monitor based on the configs & other internal components
func NewBLSKeysMonitor(
	cfg config.BLSKeysMonitorConfig,
	snoozeConfig config.AlarmSnoozeConfig,
	notifiersHandler OutputNotifiersHandler,
	statusHandler executors.StatusHandler,
) (Monitor, error) {
	parser := parsers.NewListParser()
	intentitiesHolder, err := parser.ParseFile(cfg.ListFile)
	if err != nil {
		return nil, err
	}

	ratingsChecker, err := checkers.NewBLSRatingsChecker(intentitiesHolder.BlsHexKeys, cfg.Name, cfg.AlarmDeltaRatingDrop)
	if err != nil {
		return nil, err
	}

	httpWrapper := http.NewHttpClientWrapper(nil, cfg.ApiURL)
	interactor, err := interactors.NewValidatorStatisticsInteractor(httpWrapper)
	if err != nil {
		return nil, err
	}

	fetcher, err := interactors.NewBLSKeysFetcher(httpWrapper, intentitiesHolder.Addresses, timeBetweenBLSKeysFetch)
	if err != nil {
		return nil, err
	}

	blsKeysFilter, err := NewBLSKeysFilter(snoozeConfig)
	if err != nil {
		return nil, err
	}

	argsExecutor := executors.ArgsBLSKeysExecutor{
		OutputNotifiersHandler:     notifiersHandler,
		RatingsChecker:             ratingsChecker,
		ValidatorStatisticsQuerier: interactor,
		BlsKeysFetcher:             fetcher,
		StatusHandler:              statusHandler,
		BLSKeysFilter:              blsKeysFilter,
		Name:                       cfg.Name,
		ExplorerURL:                cfg.ExplorerURL,
	}

	executor, err := executors.NewBLSKeysExecutor(argsExecutor)
	if err != nil {
		return nil, err
	}

	return monitor.NewBLSKeysMonitor(
		executor,
		time.Duration(cfg.PollingIntervalInSeconds)*time.Second,
		cfg.Name)
}
