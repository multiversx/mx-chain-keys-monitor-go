package executors

import (
	"context"
	"fmt"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-keys-monitor-go/core"
)

type statusHandlerTrigger struct {
	timeFunc      func() time.Time
	executor      Executor
	triggerDay    time.Weekday
	triggerHour   int
	triggerMinute int
	shouldTrigger bool
}

// ArgsStatusHandlerTrigger represents the DTO used in the NewStatusHandlerTrigger constructor function
type ArgsStatusHandlerTrigger struct {
	TimeFunc      func() time.Time
	Executor      Executor
	TriggerDay    time.Weekday
	TriggerHour   int
	TriggerMinute int
}

// NewStatusHandlerTrigger will create a new instance of type statusHandlerTrigger
func NewStatusHandlerTrigger(args ArgsStatusHandlerTrigger) (*statusHandlerTrigger, error) {
	if args.TimeFunc == nil {
		return nil, errNilTimeFunc
	}
	if check.IfNil(args.Executor) {
		return nil, errNilExecutor
	}
	if args.TriggerDay < core.EveryWeekDay || args.TriggerDay > time.Saturday {
		return nil, fmt.Errorf("%w for provided trigger day %d", errInvalidWeekDay, args.TriggerDay)
	}
	if args.TriggerHour < 0 || args.TriggerHour > 23 {
		return nil, fmt.Errorf("%w, provided %d, interval allowed [0, 23]", errInvalidHour, args.TriggerHour)
	}
	if args.TriggerMinute < 0 || args.TriggerMinute > 59 {
		return nil, fmt.Errorf("%w, provided %d, interval allowed [0, 23]", errInvalidMinute, args.TriggerMinute)
	}

	return &statusHandlerTrigger{
		timeFunc:      args.TimeFunc,
		executor:      args.Executor,
		triggerDay:    args.TriggerDay,
		triggerHour:   args.TriggerHour,
		triggerMinute: args.TriggerMinute,
		shouldTrigger: true,
	}, nil
}

// Execute will evaluate if the trigger function should be called or not
func (trigger *statusHandlerTrigger) Execute(ctx context.Context) error {
	currentTime := trigger.timeFunc()

	// setting the trigger.shouldTrigger to true whenever there is a mismatch on the trigger day, hour or minute will help prepare
	// the triggering when the trigger conditions are met. Setting this member to false after the Execute call will ensure that
	// the execution is done only once.

	if !trigger.isCorrectTriggerDay(currentTime) {
		trigger.shouldTrigger = true
		return nil
	}
	if currentTime.Hour() != trigger.triggerHour {
		trigger.shouldTrigger = true
		return nil
	}
	if currentTime.Minute() != trigger.triggerMinute {
		trigger.shouldTrigger = true
		return nil
	}

	if trigger.shouldTrigger {
		trigger.shouldTrigger = false
		return trigger.executor.Execute(ctx)
	}

	return nil
}

func (trigger *statusHandlerTrigger) isCorrectTriggerDay(currentTime time.Time) bool {
	if trigger.triggerDay == core.EveryWeekDay {
		return true
	}

	return currentTime.Weekday() == trigger.triggerDay
}

// IsInterfaceNil returns true if there is no value under the interface
func (trigger *statusHandlerTrigger) IsInterfaceNil() bool {
	return trigger == nil
}
