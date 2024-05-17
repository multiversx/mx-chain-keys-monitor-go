package factory

import (
	"fmt"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-keys-monitor-go/config"
	"github.com/multiversx/mx-chain-keys-monitor-go/core"
	"github.com/multiversx/mx-chain-keys-monitor-go/mock"
	"github.com/stretchr/testify/assert"
)

func TestCreateStatusHandler(t *testing.T) {
	t.Parallel()

	t.Run("disabled components", func(t *testing.T) {
		cfg := config.GeneralConfigs{
			SystemSelfCheck: config.SystemSelfCheckConfig{
				Enabled: false,
			},
		}

		handler, closer, err := CreateStatusHandler(cfg, nil)
		assert.Nil(t, err)
		assert.Equal(t, "*disabled.disabledStatusHandler", fmt.Sprintf("%T", handler))
		assert.Equal(t, "*disabled.disabledCloser", fmt.Sprintf("%T", closer))
	})
	t.Run("nil notifiers handler should error", func(t *testing.T) {
		cfg := config.GeneralConfigs{
			ApplicationName: "test",
			SystemSelfCheck: config.SystemSelfCheckConfig{
				Enabled:              true,
				DayOfWeek:            "every day",
				Hour:                 0,
				Minute:               0,
				PollingIntervalInSec: 30,
			},
		}

		handler, closer, err := CreateStatusHandler(cfg, nil)
		assert.NotNil(t, err)
		assert.Equal(t, "nil output notifiers handler", err.Error())
		assert.Nil(t, handler)
		assert.Nil(t, closer)
	})
	t.Run("invalid day of week should error", func(t *testing.T) {
		cfg := config.GeneralConfigs{
			ApplicationName: "test",
			SystemSelfCheck: config.SystemSelfCheckConfig{
				Enabled:              true,
				DayOfWeek:            "not-a-day-of-week",
				Hour:                 0,
				Minute:               0,
				PollingIntervalInSec: 30,
			},
		}

		handler, closer, err := CreateStatusHandler(cfg, &mock.OutputNotifiersHandlerStub{})
		assert.NotNil(t, err)
		assert.Equal(t, "unknown day of week not-a-day-of-week", err.Error())
		assert.Nil(t, handler)
		assert.Nil(t, closer)
	})
	t.Run("invalid hour should error", func(t *testing.T) {
		cfg := config.GeneralConfigs{
			ApplicationName: "test",
			SystemSelfCheck: config.SystemSelfCheckConfig{
				Enabled:              true,
				DayOfWeek:            "every day",
				Hour:                 -1,
				Minute:               0,
				PollingIntervalInSec: 30,
			},
		}

		handler, closer, err := CreateStatusHandler(cfg, &mock.OutputNotifiersHandlerStub{})
		assert.NotNil(t, err)
		assert.Equal(t, "invalid hour, provided -1, interval allowed [0, 23]", err.Error())
		assert.Nil(t, handler)
		assert.Nil(t, closer)
	})
	t.Run("should work", func(t *testing.T) {
		cfg := config.GeneralConfigs{
			ApplicationName: "test",
			SystemSelfCheck: config.SystemSelfCheckConfig{
				Enabled:              true,
				DayOfWeek:            "every day",
				Hour:                 0,
				Minute:               0,
				PollingIntervalInSec: 30,
			},
		}

		handler, closer, err := CreateStatusHandler(cfg, &mock.OutputNotifiersHandlerStub{})
		assert.Nil(t, err)
		assert.Equal(t, "*executors.statusHandler", fmt.Sprintf("%T", handler))
		assert.Equal(t, "*polling.pollingHandler", fmt.Sprintf("%T", closer))
	})
}

func TestParseWeekday(t *testing.T) {
	t.Parallel()

	dow, err := parseWeekday("eVeRy DaY")
	assert.Equal(t, core.EveryWeekDay, dow)
	assert.Nil(t, err)

	dow, err = parseWeekday("MoNdAy")
	assert.Equal(t, time.Monday, dow)
	assert.Nil(t, err)

	dow, err = parseWeekday("TuEsDaY")
	assert.Equal(t, time.Tuesday, dow)
	assert.Nil(t, err)

	dow, err = parseWeekday("WeDnEsDaY")
	assert.Equal(t, time.Wednesday, dow)
	assert.Nil(t, err)

	dow, err = parseWeekday("ThUrSdAy")
	assert.Equal(t, time.Thursday, dow)
	assert.Nil(t, err)

	dow, err = parseWeekday("FrIdAy")
	assert.Equal(t, time.Friday, dow)
	assert.Nil(t, err)

	dow, err = parseWeekday("SaTuRdAy")
	assert.Equal(t, time.Saturday, dow)
	assert.Nil(t, err)

	dow, err = parseWeekday("SuNdAy")
	assert.Equal(t, time.Sunday, dow)
	assert.Nil(t, err)
}
