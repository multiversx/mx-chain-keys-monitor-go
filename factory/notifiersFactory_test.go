package factory

import (
	"fmt"
	"testing"

	"github.com/multiversx/mx-chain-keys-monitor-go/config"
	"github.com/stretchr/testify/assert"
)

func TestCreateOutputNotifiers(t *testing.T) {
	t.Parallel()

	t.Run("nothing enabled should create a log notifier", func(t *testing.T) {
		notifiers, err := CreateOutputNotifiers(config.AllConfigs{})

		assert.Nil(t, err)
		assert.Equal(t, 1, len(notifiers))
		assert.Equal(t, "*notifiers.logNotifier", fmt.Sprintf("%T", notifiers[0]))
	})
	t.Run("should create a pushover and a log notifier", func(t *testing.T) {
		notifiers, err := CreateOutputNotifiers(config.AllConfigs{
			Config: config.MainConfig{
				OutputNotifiers: config.OutputNotifiersConfig{
					Pushover: config.PushoverNotifierConfig{
						Enabled: true,
					},
				},
			},
		})

		assert.Nil(t, err)
		assert.Equal(t, 2, len(notifiers))
		assert.Equal(t, "*notifiers.logNotifier", fmt.Sprintf("%T", notifiers[0]))
		assert.Equal(t, "*notifiers.pushoverNotifier", fmt.Sprintf("%T", notifiers[1]))
	})
	t.Run("should create a smtp and a log notifier", func(t *testing.T) {
		notifiers, err := CreateOutputNotifiers(config.AllConfigs{
			Config: config.MainConfig{
				OutputNotifiers: config.OutputNotifiersConfig{
					Smtp: config.SmtpNotifierConfig{
						Enabled: true,
					},
				},
			},
		})

		assert.Nil(t, err)
		assert.Equal(t, 2, len(notifiers))
		assert.Equal(t, "*notifiers.logNotifier", fmt.Sprintf("%T", notifiers[0]))
		assert.Equal(t, "*notifiers.smtpNotifier", fmt.Sprintf("%T", notifiers[1]))
	})
	t.Run("should create a telegram and a log notifier", func(t *testing.T) {
		notifiers, err := CreateOutputNotifiers(config.AllConfigs{
			Config: config.MainConfig{
				OutputNotifiers: config.OutputNotifiersConfig{
					Telegram: config.TelegramNotifierConfig{
						Enabled: true,
					},
				},
			},
		})

		assert.Nil(t, err)
		assert.Equal(t, 2, len(notifiers))
		assert.Equal(t, "*notifiers.logNotifier", fmt.Sprintf("%T", notifiers[0]))
		assert.Equal(t, "*notifiers.telegramNotifier", fmt.Sprintf("%T", notifiers[1]))
	})
}
