package notifiers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"multiversx/mvx-keys-monitor/core"

	logger "github.com/multiversx/mx-chain-logger-go"
	httpSDK "github.com/multiversx/mx-sdk-go/core/http"
)

const maxSendTimeout = time.Second * 30

var log = logger.GetOrCreate("notifiers")

type pushoverRequest struct {
	Token   string `json:"token"`
	User    string `json:"user"`
	Title   string `json:"title"`
	Message string `json:"message"`
	HTML    int    `json:"html"`
}

type pushoverResponse struct {
	Status  int    `json:"status"`
	Request string `json:"request"`
}

type pushoverNotifier struct {
	token             string
	userKey           string
	httpClientWrapper HTTPClientWrapper
}

// NewPushoverNotifier will create a new Pushover notifier
func NewPushoverNotifier(url string, token string, userKey string) *pushoverNotifier {
	return &pushoverNotifier{
		httpClientWrapper: httpSDK.NewHttpClientWrapper(nil, url),
		token:             token,
		userKey:           userKey,
	}
}

// OutputMessages will push the provided messages as error
func (notifier *pushoverNotifier) OutputMessages(messages ...core.OutputMessage) {
	log.Debug("notifier.OutputMessage pushing error notification", "num messages", len(messages))
	if len(messages) == 0 {
		return
	}

	msgString := ""
	maxMessageOutputType := core.MessageOutputType(0)
	for _, msg := range messages {
		if msg.Type > maxMessageOutputType {
			maxMessageOutputType = msg.Type
		}

		msgString += createMessageString(msg)
	}

	title := createTitle(maxMessageOutputType, messages[0].ExecutorName)

	err := notifier.pushNotification(msgString, title)
	if err != nil {
		log.Error("notifier.OutputMessage pushing notification", "error", err)
	}
}

func createMessageString(msg core.OutputMessage) string {
	identifier := processIdentifier(msg)
	iconString := getIconString(msg)

	if len(msg.ProblemEncountered) == 0 {
		return fmt.Sprintf("%s %s %s\n\n",
			iconString, msg.IdentifierType, identifier)
	}

	return fmt.Sprintf("%s %s %s: %s\n\n",
		iconString, msg.IdentifierType, identifier, msg.ProblemEncountered)
}

func getIconString(message core.OutputMessage) string {
	switch message.Type {
	case core.InfoMessageOutputType:
		return "✅"
	case core.ErrorMessageOutputType:
		return "🚨"
	case core.WarningMessageOutputType:
		return "⚠️"
	default:
		return ""
	}
}

func createTitle(maxMessageOutputType core.MessageOutputType, executor string) string {
	switch maxMessageOutputType {
	case core.InfoMessageOutputType:
		return "ⓘ Info for " + executor
	case core.ErrorMessageOutputType:
		return "🚨 Problems occurred on " + executor
	case core.WarningMessageOutputType:
		return "⚠️ Warnings occurred on " + executor
	default:
		return executor
	}
}

func processIdentifier(msg core.OutputMessage) string {
	if len(msg.IdentifierURL) == 0 {
		return fmt.Sprintf("<b>%s</b>", msg.ShortIdentifier)
	}

	return fmt.Sprintf(`<b><a href="%s">%s</a></b>`, msg.IdentifierURL, msg.ShortIdentifier)
}

func (notifier *pushoverNotifier) pushNotification(msgString string, title string) error {
	req := pushoverRequest{
		Token:   notifier.token,
		User:    notifier.userKey,
		Message: msgString,
		Title:   title,
		HTML:    1,
	}

	data, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), maxSendTimeout)
	defer cancel()

	responseBytes, statusCode, err := notifier.httpClientWrapper.PostHTTP(ctx, "", data)
	if err != nil {
		return err
	}
	if !core.IsStatusCodeIs2xx(statusCode) {
		return fmt.Errorf("%w, but %d", errReturnCodeIsNotOk, statusCode)
	}

	resp := &pushoverResponse{}
	err = json.Unmarshal(responseBytes, &resp)
	if err != nil {
		return err
	}

	log.Debug("notifier.pushNotification: sent notification",
		"status", resp.Status, "request ID", resp.Request)

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (notifier *pushoverNotifier) IsInterfaceNil() bool {
	return notifier == nil
}
