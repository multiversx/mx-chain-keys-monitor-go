package monitor

import (
	"fmt"
	"time"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-sdk-go/core/polling"
)

const minInterval = time.Second
const intervalWhenError = time.Second * 10

var log = logger.GetOrCreate("monitor")

type blsKeysMonitor struct {
	pollingHandler pollingHandler
	name           string
}

// NewBLSKeysMonitor creates a new BLS keys monitor instance
func NewBLSKeysMonitor(
	executor polling.Executor,
	interval time.Duration,
	name string,
) (*blsKeysMonitor, error) {
	if interval < minInterval {
		return nil, fmt.Errorf("%w, got %v, minimum %v", errInvalidInterval, interval, minInterval)
	}

	args := polling.ArgsPollingHandler{
		Log:              log,
		Name:             name,
		PollingInterval:  interval,
		PollingWhenError: intervalWhenError,
		Executor:         executor,
	}

	pollingHandlerInstance, err := polling.NewPollingHandler(args)
	if err != nil {
		return nil, err
	}

	err = pollingHandlerInstance.StartProcessingLoop()
	if err != nil {
		return nil, err
	}
	log.Debug("started new BLS keys monitor", "monitor name", name)

	return &blsKeysMonitor{
		pollingHandler: pollingHandlerInstance,
		name:           name,
	}, nil
}

// Close closes the polling handler
func (monitor *blsKeysMonitor) Close() error {
	log.Debug("closed the BLS keys monitor", "monitor name", monitor.name)
	return monitor.pollingHandler.Close()
}

// IsInterfaceNil returns true if there is no value under the interface
func (monitor *blsKeysMonitor) IsInterfaceNil() bool {
	return monitor == nil
}
