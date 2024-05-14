package executors

import (
	"context"
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"
	"github.com/multiversx/mx-chain-keys-monitor-go/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewStatusHandler(t *testing.T) {
	t.Parallel()

	t.Run("nil notifiers handler should error", func(t *testing.T) {
		t.Parallel()

		handler, err := NewStatusHandler("app", nil)
		assert.Nil(t, handler)
		assert.Equal(t, errNilOutputNotifiersHandler, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		handler, err := NewStatusHandler("app", &mock.OutputNotifiersHandlerStub{})
		assert.NotNil(t, handler)
		assert.Nil(t, err)
	})
}

func TestStatusHandler_NotifyAppStart(t *testing.T) {
	t.Parallel()

	var returnedErr error
	sentMessages := make([]core.OutputMessage, 0)
	outputNotifiersHandler := &mock.OutputNotifiersHandlerStub{
		NotifyWithRetryHandler: func(caller string, messages ...core.OutputMessage) error {
			sentMessages = append(sentMessages, messages...)

			return returnedErr
		},
	}

	handler, _ := NewStatusHandler("app", outputNotifiersHandler)
	assert.Equal(t, 0, len(sentMessages)) // should not notify at startup

	t.Run("notifiers handler does not error", func(t *testing.T) {
		returnedErr = nil

		handler.NotifyAppStart()

		assert.Equal(t, 1, len(sentMessages))
		assert.Equal(t, "app", sentMessages[0].ExecutorName)
		assert.Equal(t, core.InfoMessageOutputType, sentMessages[0].Type)
		assert.Contains(t, sentMessages[0].IdentifierType, "Application started on")
		assert.Zero(t, handler.NumErrorsEncountered())
	})
	t.Run("notifiers handler errors", func(t *testing.T) {
		returnedErr = errors.New("expected error")

		handler.NotifyAppStart()

		assert.Equal(t, 2, len(sentMessages))
		assert.Equal(t, "app", sentMessages[0].ExecutorName)
		assert.Equal(t, core.InfoMessageOutputType, sentMessages[0].Type)
		assert.Contains(t, sentMessages[0].IdentifierType, "Application started on")
		assert.Equal(t, uint32(1), handler.NumErrorsEncountered())
	})
}

func TestStatusHandler_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var instance *statusHandler
	assert.True(t, instance.IsInterfaceNil())

	instance = &statusHandler{}
	assert.False(t, instance.IsInterfaceNil())
}

func TestStatusHandler_Execute(t *testing.T) {
	t.Parallel()

	sentMessages := make([]core.OutputMessage, 0)
	outputNotifiersHandler := &mock.OutputNotifiersHandlerStub{
		NotifyWithRetryHandler: func(caller string, messages ...core.OutputMessage) error {
			sentMessages = append(sentMessages, messages...)

			return nil
		},
	}

	notifier, _ := NewStatusHandler("app", outputNotifiersHandler)
	sentMessages = make([]core.OutputMessage, 0) // reset the constructor sent messages

	t.Run("empty state should return info messages", func(t *testing.T) {
		err := notifier.Execute(context.Background())
		assert.Nil(t, err)

		expectedMessageErr := core.OutputMessage{
			Type:           core.InfoMessageOutputType,
			IdentifierType: "No application errors occurred",
			ExecutorName:   "app",
		}
		expectedMessageKeys := core.OutputMessage{
			Type:           core.InfoMessageOutputType,
			IdentifierType: "All monitored keys are performing as expected",
			ExecutorName:   "app",
		}

		assert.Equal(t, []core.OutputMessage{expectedMessageErr, expectedMessageKeys}, sentMessages)

		sentMessages = make([]core.OutputMessage, 0)
	})
	t.Run("2 errors should return warn messages", func(t *testing.T) {
		notifier.ErrorEncountered(nil)
		notifier.ErrorEncountered(errors.New("error 1"))
		notifier.ErrorEncountered(errors.New("error 2"))
		err := notifier.Execute(context.Background())
		assert.Nil(t, err)

		expectedMessageErr := core.OutputMessage{
			Type:            core.WarningMessageOutputType,
			ShortIdentifier: "2 application error(s) occurred, please check the app logs",
			ExecutorName:    "app",
		}
		expectedMessageKeys := core.OutputMessage{
			Type:           core.InfoMessageOutputType,
			IdentifierType: "All monitored keys are performing as expected",
			ExecutorName:   "app",
		}

		assert.Equal(t, []core.OutputMessage{expectedMessageErr, expectedMessageKeys}, sentMessages)

		sentMessages = make([]core.OutputMessage, 0)
	})
	t.Run("2 problematic keys found should return warn messages", func(t *testing.T) {
		key1 := core.OutputMessage{
			Identifier: "bls1",
		}
		key2 := core.OutputMessage{
			Identifier: "bls2",
		}

		notifier.CollectKeysProblems([]core.OutputMessage{key2})
		notifier.CollectKeysProblems([]core.OutputMessage{key1, key2})
		notifier.CollectKeysProblems([]core.OutputMessage{key1})
		notifier.CollectKeysProblems(nil)

		err := notifier.Execute(context.Background())
		assert.Nil(t, err)

		expectedMessageErr := core.OutputMessage{
			Type:           core.InfoMessageOutputType,
			IdentifierType: "No application errors occurred",
			ExecutorName:   "app",
		}
		expectedMessageKeys := core.OutputMessage{
			Type:            core.WarningMessageOutputType,
			ShortIdentifier: "2 monitored key(s) encountered problems",
			ExecutorName:    "app",
		}

		assert.Equal(t, []core.OutputMessage{expectedMessageErr, expectedMessageKeys}, sentMessages)

		sentMessages = make([]core.OutputMessage, 0)
	})
	t.Run("empty state should return info messages", func(t *testing.T) {
		err := notifier.Execute(context.Background())
		assert.Nil(t, err)

		expectedMessageErr := core.OutputMessage{
			Type:           core.InfoMessageOutputType,
			IdentifierType: "No application errors occurred",
			ExecutorName:   "app",
		}
		expectedMessageKeys := core.OutputMessage{
			Type:           core.InfoMessageOutputType,
			IdentifierType: "All monitored keys are performing as expected",
			ExecutorName:   "app",
		}

		assert.Equal(t, []core.OutputMessage{expectedMessageErr, expectedMessageKeys}, sentMessages)

		sentMessages = make([]core.OutputMessage, 0)
	})
}

func TestStatusHandler_SendCloseMessage(t *testing.T) {
	t.Parallel()

	sentMessages := make([]core.OutputMessage, 0)
	outputNotifiersHandler := &mock.OutputNotifiersHandlerStub{
		NotifyWithRetryHandler: func(caller string, messages ...core.OutputMessage) error {
			sentMessages = append(sentMessages, messages...)

			return nil
		},
	}

	notifier, _ := NewStatusHandler("app", outputNotifiersHandler)
	sentMessages = make([]core.OutputMessage, 0) // reset the constructor sent messages

	notifier.SendCloseMessage()

	expectedMessage := core.OutputMessage{
		Type:            core.WarningMessageOutputType,
		ShortIdentifier: "Application closing",
		ExecutorName:    "app",
	}

	assert.Equal(t, []core.OutputMessage{expectedMessage}, sentMessages)
}
