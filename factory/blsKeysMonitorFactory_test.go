package factory

import (
	"testing"

	"github.com/multiversx/mx-chain-keys-monitor-go/config"
	"github.com/multiversx/mx-chain-keys-monitor-go/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewBLSKeysMonitor(t *testing.T) {
	t.Parallel()

	t.Run("invalid BLS keys should error", func(t *testing.T) {
		t.Parallel()

		cfg := config.BLSKeysMonitorConfig{
			AlarmDeltaRatingDrop:     1,
			ApiURL:                   "url",
			PollingIntervalInSeconds: 1,
			ListFile:                 "./testdata/invalidKeys.list",
			Name:                     "test",
		}

		monitor, err := NewBLSKeysMonitor(
			cfg,
			config.AlarmSnoozeConfig{},
			&mock.OutputNotifiersHandlerStub{},
			&mock.StatusHandlerStub{},
		)
		assert.NotNil(t, err)
		assert.Nil(t, monitor)
	})
	t.Run("invalid interval should error", func(t *testing.T) {
		t.Parallel()

		cfg := config.BLSKeysMonitorConfig{
			AlarmDeltaRatingDrop:     1,
			ApiURL:                   "url",
			PollingIntervalInSeconds: 0,
			ListFile:                 "./testdata/keys.list",
			Name:                     "test",
		}

		monitor, err := NewBLSKeysMonitor(
			cfg,
			config.AlarmSnoozeConfig{},
			&mock.OutputNotifiersHandlerStub{},
			&mock.StatusHandlerStub{},
		)
		assert.NotNil(t, err)
		assert.Nil(t, monitor)
	})
	t.Run("nil error handler should error", func(t *testing.T) {
		t.Parallel()

		cfg := config.BLSKeysMonitorConfig{
			AlarmDeltaRatingDrop:     1,
			ApiURL:                   "url",
			PollingIntervalInSeconds: 1,
			ListFile:                 "./testdata/keys.list",
			Name:                     "test",
		}

		monitor, err := NewBLSKeysMonitor(
			cfg,
			config.AlarmSnoozeConfig{},
			&mock.OutputNotifiersHandlerStub{},
			nil,
		)
		assert.NotNil(t, err)
		assert.Nil(t, monitor)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		cfg := config.BLSKeysMonitorConfig{
			AlarmDeltaRatingDrop:     1,
			ApiURL:                   "url",
			PollingIntervalInSeconds: 1,
			ListFile:                 "./testdata/keys.list",
			Name:                     "test",
		}

		monitor, err := NewBLSKeysMonitor(
			cfg,
			config.AlarmSnoozeConfig{},
			&mock.OutputNotifiersHandlerStub{},
			&mock.StatusHandlerStub{},
		)
		assert.Nil(t, err)
		assert.NotNil(t, monitor)

		err = monitor.Close()
		assert.Nil(t, err)
	})
}
