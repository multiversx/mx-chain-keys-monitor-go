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
		pushoverNotifiers := createPushoverNotifiers(allConfig)
		outputNotifiers = append(outputNotifiers, pushoverNotifiers...)

		log.Debug("created pushover notifier(s)", "num pushover notifiers", len(pushoverNotifiers))
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
		telegramNotifiers := createTelegramNotifiers(allConfig)
		outputNotifiers = append(outputNotifiers, telegramNotifiers...)

		log.Debug("created telegram notifier(s)", "num telegram notifiers", len(telegramNotifiers))
	}

	return outputNotifiers, nil
}

func createPushoverNotifiers(allConfig config.AllConfigs) []executors.OutputNotifier {
	defaultNotifier := notifiers.NewPushoverNotifier(
		allConfig.Config.OutputNotifiers.Pushover.URL,
		allConfig.Credentials.Pushover.Token,
		allConfig.Credentials.Pushover.UserKey,
	)

	notifierInstances := []executors.OutputNotifier{defaultNotifier}
	for _, credentials := range allConfig.Credentials.Pushover.Additional {
		notifierInstance := notifiers.NewPushoverNotifier(
			allConfig.Config.OutputNotifiers.Pushover.URL,
			credentials.Token,
			credentials.UserKey,
		)

		notifierInstances = append(notifierInstances, notifierInstance)
	}

	return notifierInstances
}

func createTelegramNotifiers(allConfig config.AllConfigs) []executors.OutputNotifier {
	defaultNotifier := notifiers.NewTelegramNotifier(
		allConfig.Config.OutputNotifiers.Telegram.URL,
		allConfig.Credentials.Telegram.Token,
		allConfig.Credentials.Telegram.ChatID,
	)

	notifierInstances := []executors.OutputNotifier{defaultNotifier}
	for _, credentials := range allConfig.Credentials.Telegram.Additional {
		notifierInstance := notifiers.NewTelegramNotifier(
			allConfig.Config.OutputNotifiers.Telegram.URL,
			credentials.Token,
			credentials.ChatID,
		)

		notifierInstances = append(notifierInstances, notifierInstance)
	}

	return notifierInstances
}
