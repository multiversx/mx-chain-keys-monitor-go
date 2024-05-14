package executors

import (
	"context"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"
	"github.com/multiversx/mx-chain-keys-monitor-go/mock"
	"github.com/stretchr/testify/assert"
)

func createStatusHandlerTriggerMockArgs() ArgsStatusHandlerTrigger {
	return ArgsStatusHandlerTrigger{
		TimeFunc:      time.Now,
		Executor:      &mock.ExecutorStub{},
		TriggerDay:    0,
		TriggerHour:   0,
		TriggerMinute: 0,
	}
}

func TestNewStatusHandlerTrigger(t *testing.T) {
	t.Parallel()

	t.Run("nil time function handler should error", func(t *testing.T) {
		t.Parallel()

		args := createStatusHandlerTriggerMockArgs()
		args.TimeFunc = nil
		trigger, err := NewStatusHandlerTrigger(args)
		assert.Nil(t, trigger)
		assert.Equal(t, errNilTimeFunc, err)
	})
	t.Run("nil executor should error", func(t *testing.T) {
		t.Parallel()

		args := createStatusHandlerTriggerMockArgs()
		args.Executor = nil
		trigger, err := NewStatusHandlerTrigger(args)
		assert.Nil(t, trigger)
		assert.Equal(t, errNilExecutor, err)
	})
	t.Run("invalid weekday should error", func(t *testing.T) {
		t.Parallel()

		args := createStatusHandlerTriggerMockArgs()
		args.TriggerDay = -200
		trigger, err := NewStatusHandlerTrigger(args)
		assert.Nil(t, trigger)
		assert.ErrorIs(t, err, errInvalidWeekDay)

		args = createStatusHandlerTriggerMockArgs()
		args.TriggerDay = -2
		trigger, err = NewStatusHandlerTrigger(args)
		assert.Nil(t, trigger)
		assert.ErrorIs(t, err, errInvalidWeekDay)

		args = createStatusHandlerTriggerMockArgs()
		args.TriggerDay = time.Saturday + 1
		trigger, err = NewStatusHandlerTrigger(args)
		assert.Nil(t, trigger)
		assert.ErrorIs(t, err, errInvalidWeekDay)

		args = createStatusHandlerTriggerMockArgs()
		args.TriggerDay = time.Saturday + 100
		trigger, err = NewStatusHandlerTrigger(args)
		assert.Nil(t, trigger)
		assert.ErrorIs(t, err, errInvalidWeekDay)
	})
	t.Run("invalid hour should error", func(t *testing.T) {
		t.Parallel()

		args := createStatusHandlerTriggerMockArgs()
		args.TriggerHour = -200
		trigger, err := NewStatusHandlerTrigger(args)
		assert.Nil(t, trigger)
		assert.ErrorIs(t, err, errInvalidHour)

		args = createStatusHandlerTriggerMockArgs()
		args.TriggerHour = -1
		trigger, err = NewStatusHandlerTrigger(args)
		assert.Nil(t, trigger)
		assert.ErrorIs(t, err, errInvalidHour)

		args = createStatusHandlerTriggerMockArgs()
		args.TriggerHour = 24
		trigger, err = NewStatusHandlerTrigger(args)
		assert.Nil(t, trigger)
		assert.ErrorIs(t, err, errInvalidHour)

		args = createStatusHandlerTriggerMockArgs()
		args.TriggerHour = 500
		trigger, err = NewStatusHandlerTrigger(args)
		assert.Nil(t, trigger)
		assert.ErrorIs(t, err, errInvalidHour)
	})
	t.Run("invalid minute should error", func(t *testing.T) {
		t.Parallel()

		args := createStatusHandlerTriggerMockArgs()
		args.TriggerMinute = -200
		trigger, err := NewStatusHandlerTrigger(args)
		assert.Nil(t, trigger)
		assert.ErrorIs(t, err, errInvalidMinute)

		args = createStatusHandlerTriggerMockArgs()
		args.TriggerMinute = -1
		trigger, err = NewStatusHandlerTrigger(args)
		assert.Nil(t, trigger)
		assert.ErrorIs(t, err, errInvalidMinute)

		args = createStatusHandlerTriggerMockArgs()
		args.TriggerMinute = 60
		trigger, err = NewStatusHandlerTrigger(args)
		assert.Nil(t, trigger)
		assert.ErrorIs(t, err, errInvalidMinute)

		args = createStatusHandlerTriggerMockArgs()
		args.TriggerMinute = 500
		trigger, err = NewStatusHandlerTrigger(args)
		assert.Nil(t, trigger)
		assert.ErrorIs(t, err, errInvalidMinute)
	})
	t.Run("invalid minute should error", func(t *testing.T) {
		t.Parallel()

		args := createStatusHandlerTriggerMockArgs()
		args.TriggerMinute = -200
		trigger, err := NewStatusHandlerTrigger(args)
		assert.Nil(t, trigger)
		assert.ErrorIs(t, err, errInvalidMinute)

		args = createStatusHandlerTriggerMockArgs()
		args.TriggerMinute = -1
		trigger, err = NewStatusHandlerTrigger(args)
		assert.Nil(t, trigger)
		assert.ErrorIs(t, err, errInvalidMinute)

		args = createStatusHandlerTriggerMockArgs()
		args.TriggerMinute = 60
		trigger, err = NewStatusHandlerTrigger(args)
		assert.Nil(t, trigger)
		assert.ErrorIs(t, err, errInvalidMinute)

		args = createStatusHandlerTriggerMockArgs()
		args.TriggerMinute = 500
		trigger, err = NewStatusHandlerTrigger(args)
		assert.Nil(t, trigger)
		assert.ErrorIs(t, err, errInvalidMinute)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createStatusHandlerTriggerMockArgs()
		args.TriggerDay = core.EveryWeekDay
		trigger, err := NewStatusHandlerTrigger(args)
		assert.NotNil(t, trigger)
		assert.Nil(t, err)

		args = createStatusHandlerTriggerMockArgs()
		args.TriggerDay = time.Saturday
		args.TriggerHour = 23
		args.TriggerMinute = 59
		trigger, err = NewStatusHandlerTrigger(args)
		assert.NotNil(t, trigger)
		assert.Nil(t, err)

		args = createStatusHandlerTriggerMockArgs()
		args.TriggerDay = time.Monday
		args.TriggerHour = 16
		args.TriggerMinute = 23
		trigger, err = NewStatusHandlerTrigger(args)
		assert.NotNil(t, trigger)
		assert.Nil(t, err)
	})
}

func TestStatusHandlerTrigger_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var instance *statusHandlerTrigger
	assert.True(t, instance.IsInterfaceNil())

	instance = &statusHandlerTrigger{}
	assert.False(t, instance.IsInterfaceNil())
}

func TestStatusHandlerTrigger_Execute(t *testing.T) {
	t.Parallel()

	t.Run("execute 1 week execution computation should trigger once", func(t *testing.T) {
		t.Parallel()

		startTime := time.Date(2024, 04, 12, 11, 27, 00, 0, time.UTC) // Friday
		endTime := startTime.Add(time.Hour*24*7 + time.Second)

		var currentTime *time.Time

		args := createStatusHandlerTriggerMockArgs()
		args.TriggerDay = time.Monday
		args.TriggerHour = 16
		args.TriggerMinute = 23
		numCalls := 0
		args.TimeFunc = func() time.Time {
			return *currentTime
		}
		args.Executor = &mock.ExecutorStub{
			ExecuteHandler: func(ctx context.Context) error {
				numCalls++
				return nil
			},
		}
		trigger, _ := NewStatusHandlerTrigger(args)

		for sweepTime := startTime; sweepTime.Unix() < endTime.Unix(); sweepTime = sweepTime.Add(time.Second * 30) {
			currentTime = &sweepTime

			err := trigger.Execute(context.Background())
			assert.Nil(t, err)
		}

		assert.Equal(t, 1, numCalls)
	})
	t.Run("execute 1 week execution computation should trigger 7 times", func(t *testing.T) {
		t.Parallel()

		startTime := time.Date(2024, 04, 12, 11, 27, 00, 0, time.UTC) // Friday
		endTime := startTime.Add(time.Hour*24*7 + time.Second)

		var currentTime *time.Time

		args := createStatusHandlerTriggerMockArgs()
		args.TriggerDay = core.EveryWeekDay
		args.TriggerHour = 16
		args.TriggerMinute = 23
		numCalls := 0
		args.TimeFunc = func() time.Time {
			return *currentTime
		}
		args.Executor = &mock.ExecutorStub{
			ExecuteHandler: func(ctx context.Context) error {
				numCalls++
				return nil
			},
		}
		trigger, _ := NewStatusHandlerTrigger(args)

		for sweepTime := startTime; sweepTime.Unix() < endTime.Unix(); sweepTime = sweepTime.Add(time.Second * 30) {
			currentTime = &sweepTime

			err := trigger.Execute(context.Background())
			assert.Nil(t, err)
		}

		assert.Equal(t, 7, numCalls)
	})
}
