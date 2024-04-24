package notifiers

import (
	"testing"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"
	"github.com/multiversx/mx-chain-keys-monitor-go/mock"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/stretchr/testify/assert"
)

func TestNewLogNotifier(t *testing.T) {
	t.Parallel()

	t.Run("nil logger should error", func(t *testing.T) {
		t.Parallel()

		notifier, err := NewLogNotifier(nil)
		assert.Nil(t, notifier)
		assert.Equal(t, errNilLogger, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		notifier, err := NewLogNotifier(&mock.LoggerStub{})
		assert.NotNil(t, notifier)
		assert.Nil(t, err)
	})
}

func TestLogNotifier_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var notifier *logNotifier
	assert.True(t, notifier.IsInterfaceNil())

	notifier = &logNotifier{}
	assert.False(t, notifier.IsInterfaceNil())
}

func TestLogNotifier_OutputErrorMessages(t *testing.T) {
	t.Parallel()

	messages := make(map[string]logger.LogLevel)
	logInstance := &mock.LoggerStub{
		LogHandler: func(logLevel logger.LogLevel, message string, args ...interface{}) {
			messages[message] = logLevel
		},
	}

	notifier, _ := NewLogNotifier(logInstance)
	message1 := core.OutputMessage{
		Type:               core.ErrorMessageOutputType,
		IdentifierType:     "BLS key",
		Identifier:         "0295e29aef11c30391a70c35781",
		ExecutorName:       "test executor",
		ProblemEncountered: "Rating drop detected: temp rating: 90.70, rating: 100.00",
	}

	message2 := core.OutputMessage{
		Type:               core.WarningMessageOutputType,
		IdentifierType:     "BLS key",
		Identifier:         "0295e29aef11c30391a70c35782",
		ExecutorName:       "test executor",
		ProblemEncountered: "Rating drop detected: temp rating: 90.70, rating: 100.00",
	}

	message3 := core.OutputMessage{
		Type:               core.InfoMessageOutputType,
		IdentifierType:     "BLS key",
		Identifier:         "0295e29aef11c30391a70c35783",
		ExecutorName:       "test executor",
		ProblemEncountered: "Rating drop detected: temp rating: 90.70, rating: 100.00",
	}

	notifier.OutputMessages(message1, message2, message3)

	expectedMap := map[string]logger.LogLevel{
		"BLS key 0295e29aef11c30391a70c35781 -> Rating drop detected: temp rating: 90.70, rating: 100.00 called by test executor": logger.LogError,
		"BLS key 0295e29aef11c30391a70c35782 -> Rating drop detected: temp rating: 90.70, rating: 100.00 called by test executor": logger.LogWarning,
		"BLS key 0295e29aef11c30391a70c35783 -> Rating drop detected: temp rating: 90.70, rating: 100.00 called by test executor": logger.LogInfo,
	}

	assert.Equal(t, expectedMap, messages)
}
