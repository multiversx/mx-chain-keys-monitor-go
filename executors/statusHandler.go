package executors

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-keys-monitor-go/core"
)

const statusHandlerName = "statusHandler"

type statusHandler struct {
	name             string
	mut              sync.Mutex
	numErrors        uint32
	problematicKeys  map[string]struct{}
	notifiersHandler OutputNotifiersHandler
}

// NewStatusHandler creates a new instance of type statusHandler
func NewStatusHandler(name string, notifiersHandler OutputNotifiersHandler) (*statusHandler, error) {
	if check.IfNil(notifiersHandler) {
		return nil, errNilOutputNotifiersHandler
	}

	handler := &statusHandler{
		notifiersHandler: notifiersHandler,
		name:             name,
		problematicKeys:  make(map[string]struct{}),
	}

	return handler, nil
}

// NotifyAppStart will send a message on all notifiers stating that the application was started
func (handler *statusHandler) NotifyAppStart() {
	msg := core.OutputMessage{
		Type:           core.InfoMessageOutputType,
		ExecutorName:   handler.name,
		IdentifierType: fmt.Sprintf("Application started on %s", time.Now().Format("01-02-2006 15:04:05")),
	}

	handler.notifiersHandler.NotifyWithRetry(statusHandlerName, msg)
}

// ErrorEncountered will log the error and record the number of errors encountered
func (handler *statusHandler) ErrorEncountered(err error) {
	if err == nil {
		return
	}

	log.LogIfError(err)

	atomic.AddUint32(&handler.numErrors, 1)
}

// CollectKeysProblems will record the problem on the BLS keys so the periodic report will include the number of identity failures
func (handler *statusHandler) CollectKeysProblems(messages []core.OutputMessage) {
	handler.mut.Lock()
	for _, msg := range messages {
		handler.problematicKeys[msg.Identifier] = struct{}{}
	}
	handler.mut.Unlock()
}

// Execute will push the notification messages and will clear the internal state of the status handler
func (handler *statusHandler) Execute(_ context.Context) error {
	handler.mut.Lock()
	numErrors := atomic.SwapUint32(&handler.numErrors, 0)

	numProblematicKeys := len(handler.problematicKeys)
	handler.problematicKeys = make(map[string]struct{})
	handler.mut.Unlock()

	messages := []core.OutputMessage{
		handler.createErrorsMessage(numErrors),
		handler.createProblematicKeysMessage(numProblematicKeys),
	}

	handler.notifiersHandler.NotifyWithRetry(statusHandlerName, messages...)

	return nil
}

func (handler *statusHandler) createErrorsMessage(numErrors uint32) core.OutputMessage {
	msg := core.OutputMessage{
		Type:           core.InfoMessageOutputType,
		ExecutorName:   handler.name,
		IdentifierType: "No application errors occurred",
	}

	if numErrors == 0 {
		return msg
	}

	msg.Type = core.WarningMessageOutputType
	msg.ShortIdentifier = fmt.Sprintf("%d application error(s) occurred, please check the app logs", numErrors)
	msg.IdentifierType = ""

	return msg
}

func (handler *statusHandler) createProblematicKeysMessage(numProblematicKeys int) core.OutputMessage {
	msg := core.OutputMessage{
		Type:           core.InfoMessageOutputType,
		ExecutorName:   handler.name,
		IdentifierType: "All monitored keys are performing as expected",
	}

	if numProblematicKeys == 0 {
		return msg
	}

	msg.Type = core.WarningMessageOutputType
	msg.ShortIdentifier = fmt.Sprintf("%d monitored key(s) encountered problems", numProblematicKeys)
	msg.IdentifierType = ""

	return msg
}

func (handler *statusHandler) createCloseMessage() core.OutputMessage {
	msg := core.OutputMessage{
		Type:            core.WarningMessageOutputType,
		ExecutorName:    handler.name,
		ShortIdentifier: "Application closing",
	}

	return msg
}

// SendCloseMessage will send to all notifiers a close message
func (handler *statusHandler) SendCloseMessage() {
	closeMessage := handler.createCloseMessage()

	handler.notifiersHandler.NotifyWithRetry(statusHandlerName, closeMessage)
}

// IsInterfaceNil returns true if there is no value under the interface
func (handler *statusHandler) IsInterfaceNil() bool {
	return handler == nil
}
