package factory

import (
	"fmt"
	"testing"

	"github.com/multiversx/mx-chain-keys-monitor-go/config"
	"github.com/stretchr/testify/assert"
)

func TestNewBLSKeysFilter(t *testing.T) {
	t.Parallel()

	t.Run("disabled should work", func(t *testing.T) {
		t.Parallel()

		instance, err := NewBLSKeysFilter(config.AlarmSnoozeConfig{
			Enabled: false,
		})
		assert.Nil(t, err)
		assert.Equal(t, "*disabled.disabledBLSKeysFilter", fmt.Sprintf("%T", instance))
	})
	t.Run("enabled should work", func(t *testing.T) {
		t.Parallel()

		instance, err := NewBLSKeysFilter(config.AlarmSnoozeConfig{
			Enabled: true,
		})
		assert.Nil(t, err)
		assert.Equal(t, "*executors.blsKeysTimeCache", fmt.Sprintf("%T", instance))
	})
}
