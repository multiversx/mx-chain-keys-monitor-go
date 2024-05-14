package executors

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"
	"github.com/multiversx/mx-chain-keys-monitor-go/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewNotifiersHandler(t *testing.T) {
	t.Parallel()

	testArgs := ArgsNotifiersHandler{
		Notifiers:          nil,
		NumRetries:         0,
		TimeBetweenRetries: minTimeBetweenRetries,
	}

	t.Run("nil notifier should error", func(t *testing.T) {
		t.Parallel()

		localArgs := testArgs
		localArgs.Notifiers = []OutputNotifier{&mock.OutputNotifierStub{}, nil}
		handler, err := NewNotifiersHandler(localArgs)
		assert.Nil(t, handler)
		assert.ErrorIs(t, err, errNilOutputNotifier)
		assert.Contains(t, err.Error(), "at index 1")
	})
	t.Run("invalid time between retries should error", func(t *testing.T) {
		t.Parallel()

		localArgs := testArgs
		localArgs.TimeBetweenRetries = minTimeBetweenRetries - time.Nanosecond
		handler, err := NewNotifiersHandler(localArgs)
		assert.Nil(t, handler)
		assert.ErrorIs(t, err, errInvalidTimeBetweenRetries)
		assert.Contains(t, err.Error(), "9.999999ms, minimum: 10ms")
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		localArgs := testArgs
		handler, err := NewNotifiersHandler(localArgs)
		assert.NotNil(t, handler)
		assert.Nil(t, err)
	})
}

func TestNotifiersHandler_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var instance *notifiersHandler
	assert.True(t, instance.IsInterfaceNil())

	instance = &notifiersHandler{}
	assert.False(t, instance.IsInterfaceNil())
}

func TestNotifiersHandler_NotifyWithRetry(t *testing.T) {
	t.Parallel()

	testMessages := []core.OutputMessage{
		{
			Type:               core.ErrorMessageOutputType,
			IdentifierType:     "BLS key",
			Identifier:         "bls1",
			ShortIdentifier:    "bls1",
			IdentifierURL:      "https://explorer.com/nodes/bls1",
			ExecutorName:       "executor test name",
			ProblemEncountered: "status1",
		},
		{
			Type:               core.ErrorMessageOutputType,
			IdentifierType:     "BLS key",
			Identifier:         "bls2",
			ShortIdentifier:    "bls2",
			IdentifierURL:      "https://explorer.com/nodes/bls2",
			ExecutorName:       "executor test name",
			ProblemEncountered: "status2",
		},
		{
			Type:            core.InfoMessageOutputType,
			IdentifierType:  "ok",
			Identifier:      "ok",
			ShortIdentifier: "ok",
			ExecutorName:    "executor test name",
		},
	}
	expectedErr := errors.New("expected error")

	t.Run("no notifiers should not panic", func(t *testing.T) {
		defer func() {
			r := recover()
			if r != nil {
				assert.Fail(t, fmt.Sprintf("should have not panicked %v", r))
			}
		}()

		testArgs := ArgsNotifiersHandler{
			Notifiers:          nil,
			NumRetries:         0,
			TimeBetweenRetries: minTimeBetweenRetries,
		}
		handler, _ := NewNotifiersHandler(testArgs)

		err := handler.NotifyWithRetry("test")
		assert.Nil(t, err)
		err = handler.NotifyWithRetry("test", make([]core.OutputMessage, 0)...)
		assert.Nil(t, err)
		err = handler.NotifyWithRetry("test", testMessages...)
		assert.Nil(t, err)
	})
	t.Run("should not notify if no messages are to be sent", func(t *testing.T) {
		resultMap := make(map[string][]core.OutputMessage)

		notifiers := []OutputNotifier{
			&mock.OutputNotifierStub{
				OutputMessagesHandler: func(messages ...core.OutputMessage) error {
					resultMap["notifier1"] = append(resultMap["notifier1"], messages...)

					return nil
				},
			},
			&mock.OutputNotifierStub{
				OutputMessagesHandler: func(messages ...core.OutputMessage) error {
					resultMap["notifier2"] = append(resultMap["notifier2"], messages...)

					return nil
				},
			},
		}

		testArgs := ArgsNotifiersHandler{
			Notifiers:          notifiers,
			NumRetries:         0,
			TimeBetweenRetries: minTimeBetweenRetries,
		}
		handler, _ := NewNotifiersHandler(testArgs)

		err := handler.NotifyWithRetry("test")
		assert.Nil(t, err)
		err = handler.NotifyWithRetry("test", make([]core.OutputMessage, 0)...)
		assert.Nil(t, err)

		assert.Equal(t, 0, len(resultMap))
	})
	t.Run("should work if no errors were found", func(t *testing.T) {
		resultMap := make(map[string][]core.OutputMessage)

		notifiers := []OutputNotifier{
			&mock.OutputNotifierStub{
				OutputMessagesHandler: func(messages ...core.OutputMessage) error {
					resultMap["notifier1"] = append(resultMap["notifier1"], messages...)

					return nil
				},
			},
			&mock.OutputNotifierStub{
				OutputMessagesHandler: func(messages ...core.OutputMessage) error {
					resultMap["notifier2"] = append(resultMap["notifier2"], messages...)

					return nil
				},
			},
		}

		testArgs := ArgsNotifiersHandler{
			Notifiers:          notifiers,
			NumRetries:         0,
			TimeBetweenRetries: minTimeBetweenRetries,
		}
		handler, _ := NewNotifiersHandler(testArgs)
		err := handler.NotifyWithRetry("test", testMessages...)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(resultMap))
		assert.Equal(t, testMessages, resultMap["notifier1"])
		assert.Equal(t, testMessages, resultMap["notifier2"])
	})
	t.Run("should work if errors were found but with 0 retries", func(t *testing.T) {
		resultMap := make(map[string][]core.OutputMessage)

		notifiers := []OutputNotifier{
			&mock.OutputNotifierStub{
				OutputMessagesHandler: func(messages ...core.OutputMessage) error {
					resultMap["notifier1"] = append(resultMap["notifier1"], messages...)

					return expectedErr
				},
			},
			&mock.OutputNotifierStub{
				OutputMessagesHandler: func(messages ...core.OutputMessage) error {
					resultMap["notifier2"] = append(resultMap["notifier2"], messages...)

					return nil
				},
			},
		}

		testArgs := ArgsNotifiersHandler{
			Notifiers:          notifiers,
			NumRetries:         0,
			TimeBetweenRetries: minTimeBetweenRetries,
		}
		handler, _ := NewNotifiersHandler(testArgs)
		err := handler.NotifyWithRetry("test", testMessages...)
		assert.ErrorIs(t, err, errNotificationsSendingProblems)
		assert.Contains(t, err.Error(), "num notifiers with problems: 1")
		assert.Equal(t, 2, len(resultMap))
		assert.Equal(t, testMessages, resultMap["notifier1"])
		assert.Equal(t, testMessages, resultMap["notifier2"])
	})
	t.Run("should work with retries if errors were found", func(t *testing.T) {
		resultMap := make(map[string][]core.OutputMessage)

		notifiers := []OutputNotifier{
			&mock.OutputNotifierStub{
				OutputMessagesHandler: func(messages ...core.OutputMessage) error {
					resultMap["notifier1"] = append(resultMap["notifier1"], messages...)

					return expectedErr
				},
			},
			&mock.OutputNotifierStub{
				OutputMessagesHandler: func(messages ...core.OutputMessage) error {
					resultMap["notifier2"] = append(resultMap["notifier2"], messages...)

					return nil
				},
			},
		}

		testArgs := ArgsNotifiersHandler{
			Notifiers:          notifiers,
			NumRetries:         1,
			TimeBetweenRetries: minTimeBetweenRetries,
		}
		handler, _ := NewNotifiersHandler(testArgs)
		err := handler.NotifyWithRetry("test", testMessages...)
		assert.ErrorIs(t, err, errNotificationsSendingProblems)
		assert.Contains(t, err.Error(), "num notifiers with problems: 1")
		assert.Equal(t, 2, len(resultMap))
		assert.Equal(t, append(testMessages, testMessages...), resultMap["notifier1"])
		assert.Equal(t, testMessages, resultMap["notifier2"])
	})
	t.Run("should work with retries if errors were found on all notifiers", func(t *testing.T) {
		resultMap := make(map[string][]core.OutputMessage)

		notifiers := []OutputNotifier{
			&mock.OutputNotifierStub{
				OutputMessagesHandler: func(messages ...core.OutputMessage) error {
					resultMap["notifier1"] = append(resultMap["notifier1"], messages...)

					return expectedErr
				},
			},
			&mock.OutputNotifierStub{
				OutputMessagesHandler: func(messages ...core.OutputMessage) error {
					resultMap["notifier2"] = append(resultMap["notifier2"], messages...)

					return expectedErr
				},
			},
		}

		testArgs := ArgsNotifiersHandler{
			Notifiers:          notifiers,
			NumRetries:         2,
			TimeBetweenRetries: minTimeBetweenRetries,
		}
		handler, _ := NewNotifiersHandler(testArgs)
		err := handler.NotifyWithRetry("test", testMessages...)
		assert.ErrorIs(t, err, errNotificationsSendingProblems)
		assert.Contains(t, err.Error(), "num notifiers with problems: 2")
		assert.Equal(t, 2, len(resultMap))
		assert.Equal(t, append(testMessages, append(testMessages, testMessages...)...), resultMap["notifier1"])
		assert.Equal(t, append(testMessages, append(testMessages, testMessages...)...), resultMap["notifier2"])
	})
	t.Run("should work with retries & notifier recovers after first try", func(t *testing.T) {
		resultMap := make(map[string][]core.OutputMessage)

		numErrors := 0
		notifiers := []OutputNotifier{
			&mock.OutputNotifierStub{
				OutputMessagesHandler: func(messages ...core.OutputMessage) error {
					resultMap["notifier1"] = append(resultMap["notifier1"], messages...)

					numErrors++
					if numErrors == 2 {
						return nil
					}

					return expectedErr
				},
			},
			&mock.OutputNotifierStub{
				OutputMessagesHandler: func(messages ...core.OutputMessage) error {
					resultMap["notifier2"] = append(resultMap["notifier2"], messages...)

					return nil
				},
			},
		}

		testArgs := ArgsNotifiersHandler{
			Notifiers:          notifiers,
			NumRetries:         2,
			TimeBetweenRetries: minTimeBetweenRetries,
		}
		handler, _ := NewNotifiersHandler(testArgs)
		err := handler.NotifyWithRetry("test", testMessages...)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(resultMap))
		assert.Equal(t, append(testMessages, testMessages...), resultMap["notifier1"])
		assert.Equal(t, testMessages, resultMap["notifier2"])
	})
	t.Run("should work with retries & notifier recovers after second try", func(t *testing.T) {
		resultMap := make(map[string][]core.OutputMessage)

		numErrors := 0
		notifiers := []OutputNotifier{
			&mock.OutputNotifierStub{
				OutputMessagesHandler: func(messages ...core.OutputMessage) error {
					resultMap["notifier1"] = append(resultMap["notifier1"], messages...)

					numErrors++
					if numErrors == 3 {
						return nil
					}

					return expectedErr
				},
			},
			&mock.OutputNotifierStub{
				OutputMessagesHandler: func(messages ...core.OutputMessage) error {
					resultMap["notifier2"] = append(resultMap["notifier2"], messages...)

					return nil
				},
			},
		}

		testArgs := ArgsNotifiersHandler{
			Notifiers:          notifiers,
			NumRetries:         2,
			TimeBetweenRetries: minTimeBetweenRetries,
		}
		handler, _ := NewNotifiersHandler(testArgs)
		err := handler.NotifyWithRetry("test", testMessages...)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(resultMap))
		assert.Equal(t, append(testMessages, append(testMessages, testMessages...)...), resultMap["notifier1"])
		assert.Equal(t, testMessages, resultMap["notifier2"])
	})
}
