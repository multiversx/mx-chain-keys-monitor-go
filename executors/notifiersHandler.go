package executors

import (
	"fmt"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-keys-monitor-go/core"
)

const minTimeBetweenRetries = time.Millisecond * 10

type notifiersHandler struct {
	notifiers          []OutputNotifier
	numRetries         uint32
	timeBetweenRetries time.Duration
}

// ArgsNotifiersHandler defines the DTO struct for the NewNotifiersHandler constructor function
type ArgsNotifiersHandler struct {
	Notifiers          []OutputNotifier
	NumRetries         uint32
	TimeBetweenRetries time.Duration
}

// NewNotifiersHandler creates a new instance of type notifiersHandler
func NewNotifiersHandler(args ArgsNotifiersHandler) (*notifiersHandler, error) {
	for index, notifier := range args.Notifiers {
		if check.IfNil(notifier) {
			return nil, fmt.Errorf("%w at index %d", errNilOutputNotifier, index)
		}
	}
	if args.TimeBetweenRetries < minTimeBetweenRetries {
		return nil, fmt.Errorf("%w provided: %v, minimum: %v", errInvalidTimeBetweenRetries, args.TimeBetweenRetries, minTimeBetweenRetries)
	}

	return &notifiersHandler{
		notifiers:          args.Notifiers,
		numRetries:         args.NumRetries,
		timeBetweenRetries: args.TimeBetweenRetries,
	}, nil
}

// NotifyWithRetry will try to send the messages to the inner list of notifiers. If one or more notifiers errored, it will retry
// for the maximum number of retries provided. There is a delay between the retries, to maximize the chances of recovery
func (handler *notifiersHandler) NotifyWithRetry(caller string, messages ...core.OutputMessage) error {
	log.Debug("notifiersHandler.NotifyWithRetry",
		"executor", caller, "num messages", len(messages), "num retries", handler.numRetries,
		"num notifiers", len(handler.notifiers))

	notifiersThatErrored := make(map[string]OutputNotifier)
	handler.handleFirstSend(notifiersThatErrored, messages)
	return handler.handleRetries(notifiersThatErrored, messages)
}

func (handler *notifiersHandler) handleFirstSend(notifiersThatErrored map[string]OutputNotifier, messages []core.OutputMessage) {
	if len(messages) == 0 {
		return
	}

	for _, notifier := range handler.notifiers {
		err := notifier.OutputMessages(messages...)
		if err != nil {
			uniqueNotifierNameStr := uniqueNotifierName(notifier)
			log.Error("error sending notification", "notifier", uniqueNotifierNameStr, "error", err)
			notifiersThatErrored[uniqueNotifierNameStr] = notifier
		}
	}
}

func (handler *notifiersHandler) handleRetries(notifiersThatErrored map[string]OutputNotifier, messages []core.OutputMessage) error {
	if handler.numRetries == 0 && len(notifiersThatErrored) == 0 {
		log.Debug("notifiersHandler.NotifyWithRetry no errors found")
		return nil
	}

	for i := uint32(0); i < handler.numRetries; i++ {
		time.Sleep(handler.timeBetweenRetries)
		handler.notifyInRetryMode(i, notifiersThatErrored, messages)

		if len(notifiersThatErrored) == 0 {
			return nil
		}
	}

	return fmt.Errorf("%w: num notifiers with problems: %d", errNotificationsSendingProblems, len(notifiersThatErrored))
}

func (handler *notifiersHandler) notifyInRetryMode(retryNumber uint32, notifiersThatErrored map[string]OutputNotifier, messages []core.OutputMessage) {
	for identifier, notifier := range notifiersThatErrored {
		err := notifier.OutputMessages(messages...)
		if err == nil {
			delete(notifiersThatErrored, identifier)
			continue
		}

		log.Error("error sending notification",
			"notifier", identifier, "retry no.", retryNumber, "error", err)
	}
}

// this function will prevent the abnormal operation of this component when we have, by mistake, 2 or more notifiers that return the same name
func uniqueNotifierName(notifier OutputNotifier) string {
	return fmt.Sprintf("%s: %p", notifier.Name(), notifier)
}

// IsInterfaceNil returns true if there is no value under the interface
func (handler *notifiersHandler) IsInterfaceNil() bool {
	return handler == nil
}
