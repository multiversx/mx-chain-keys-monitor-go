package factory

import (
	"github.com/multiversx/mx-chain-keys-monitor-go/config"
	"github.com/multiversx/mx-chain-keys-monitor-go/executors"
	"github.com/multiversx/mx-chain-keys-monitor-go/notifiers"
	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("factory")

// CreateOutputNotifiers will create the output notifiers based on the provided configuration
func CreateOutputNotifiers(allConfig config.AllConfigs) ([]executors.OutputNotifier, error) {
	outputNotifiers := make([]executors.OutputNotifier, 0)

	logForNotifier := logger.GetOrCreate("notifiers")
	logNotifier, err := notifiers.NewLogNotifier(logForNotifier)
	if err != nil {
		return nil, err
	}
	log.Debug("created log notifier")

	// the default notifier is the one that logs in a file
	outputNotifiers = append(outputNotifiers, logNotifier)

	if allConfig.Config.OutputNotifiers.Pushover.Enabled {
		outputNotifiers = append(outputNotifiers, notifiers.NewPushoverNotifier(
			allConfig.Config.OutputNotifiers.Pushover.URL,
			allConfig.Credentials.Pushover.Token,
			allConfig.Credentials.Pushover.UserKey,
		))
		log.Debug("created pushover notifier")
	}
	if allConfig.Config.OutputNotifiers.Smtp.Enabled {
		args := notifiers.ArgsSmtpNotifier{
			To:       allConfig.Config.OutputNotifiers.Smtp.To,
			SmtpPort: allConfig.Config.OutputNotifiers.Smtp.SmtpPort,
			SmtpHost: allConfig.Config.OutputNotifiers.Smtp.SmtpHost,
			From:     allConfig.Credentials.Smtp.Email,
			Password: allConfig.Credentials.Smtp.Password,
		}

		outputNotifiers = append(outputNotifiers, notifiers.NewSmtpNotifier(args))
		log.Debug("created smtp notifier")
	}
	if allConfig.Config.OutputNotifiers.Telegram.Enabled {
		outputNotifiers = append(outputNotifiers, notifiers.NewTelegramNotifier(
			allConfig.Config.OutputNotifiers.Telegram.URL,
			allConfig.Credentials.Telegram.Token,
			allConfig.Credentials.Telegram.ChatID,
		))
		log.Debug("created telegram notifier")
	}

	return outputNotifiers, nil
}
