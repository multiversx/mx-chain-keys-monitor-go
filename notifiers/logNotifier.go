package notifiers

import (
	"fmt"

	"multiversx/mx-chain-keys-monitor-go/core"

	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
)

type logNotifier struct {
	log logger.Logger
}

// NewLogNotifier creates a new log notifier instance
func NewLogNotifier(log logger.Logger) (*logNotifier, error) {
	if check.IfNil(log) {
		return nil, errNilLogger
	}

	return &logNotifier{
		log: log,
	}, nil
}

// OutputMessages will send the output messages to the current logger
func (notifier *logNotifier) OutputMessages(messages ...core.OutputMessage) {
	for _, msg := range messages {
		notifier.output(msg)
	}
}

func (notifier *logNotifier) output(message core.OutputMessage) {
	msg := ""
	if len(message.IdentifierType) > 0 {
		msg += message.IdentifierType + " "
	}
	if len(message.Identifier) > 0 {
		msg += message.Identifier + " "
	}
	if len(message.ProblemEncountered) > 0 {
		msg += "-> " + message.ProblemEncountered + " "
	}
	msg += fmt.Sprintf("called by %s", message.ExecutorName)

	logLevel := logger.LogInfo
	switch message.Type {
	case core.ErrorMessageOutputType:
		logLevel = logger.LogError
	case core.WarningMessageOutputType:
		logLevel = logger.LogWarning
	}

	notifier.log.Log(logLevel, msg)
}

// IsInterfaceNil returns true if there is no value under the interface
func (notifier *logNotifier) IsInterfaceNil() bool {
	return notifier == nil
}
