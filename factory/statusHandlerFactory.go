package factory

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/multiversx/mx-chain-keys-monitor-go/config"
	"github.com/multiversx/mx-chain-keys-monitor-go/core"
	"github.com/multiversx/mx-chain-keys-monitor-go/executors"
	"github.com/multiversx/mx-chain-keys-monitor-go/executors/disabled"
	"github.com/multiversx/mx-sdk-go/core/polling"
)

// CreateStatusHandler will create an instance of type executors.StatusHandler and an instance of io.Closer (the polling component)
func CreateStatusHandler(
	generalConfig config.GeneralConfigs,
	notifiersHandler OutputNotifiersHandler,
) (executors.StatusHandler, io.Closer, error) {
	if !generalConfig.SystemSelfCheck.Enabled {
		return disabled.NewDisabledStatusHandler(), disabled.NewDisabledCloser(), nil
	}

	statusHandler, err := executors.NewStatusHandler(generalConfig.ApplicationName, notifiersHandler)
	if err != nil {
		return nil, nil, err
	}

	dayOfWeek, err := convertStringToDayOfWeek(generalConfig.SystemSelfCheck.DayOfWeek)
	if err != nil {
		return nil, nil, err
	}

	argsTrigger := executors.ArgsStatusHandlerTrigger{
		TimeFunc:      time.Now,
		Executor:      statusHandler,
		TriggerDay:    dayOfWeek,
		TriggerHour:   generalConfig.SystemSelfCheck.Hour,
		TriggerMinute: generalConfig.SystemSelfCheck.Minute,
	}

	trigger, err := executors.NewStatusHandlerTrigger(argsTrigger)
	if err != nil {
		return nil, nil, err
	}

	argsPolling := polling.ArgsPollingHandler{
		Log:              log,
		Name:             generalConfig.ApplicationName,
		PollingInterval:  time.Second * time.Duration(generalConfig.SystemSelfCheck.PollingIntervalInSec),
		PollingWhenError: time.Second * time.Duration(generalConfig.SystemSelfCheck.PollingIntervalInSec),
		Executor:         trigger,
	}

	pollingHandler, err := polling.NewPollingHandler(argsPolling)
	if err != nil {
		return nil, nil, err
	}

	return statusHandler, pollingHandler, pollingHandler.StartProcessingLoop()
}

func convertStringToDayOfWeek(dayOfWeek string) (time.Weekday, error) {
	dayOfWeek = strings.ToLower(dayOfWeek)
	switch dayOfWeek {
	case "every day":
		return core.EveryWeekDay, nil
	case "monday":
		return time.Monday, nil
	case "tuesday":
		return time.Tuesday, nil
	case "wednesday":
		return time.Wednesday, nil
	case "thursday":
		return time.Thursday, nil
	case "friday":
		return time.Friday, nil
	case "saturday":
		return time.Saturday, nil
	case "sunday":
		return time.Sunday, nil
	}

	return -2, fmt.Errorf("unknown day of week %s", dayOfWeek)
}
