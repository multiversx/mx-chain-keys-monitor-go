package notifiers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"
	logger "github.com/multiversx/mx-chain-logger-go"
	httpSDK "github.com/multiversx/mx-sdk-go/core/http"
)

const (
	maxSendTimeout       = time.Second * 30
	httpBoldFormat       = "<b>%s</b>"
	httpBoldedLinkFormat = `<b><a href="%s">%s</a></b>`
)

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

// OutputMessages will send the provided messages to the Pushover service
func (notifier *pushoverNotifier) OutputMessages(messages ...core.OutputMessage) error {
	log.Debug("pushoverNotifier.OutputMessages sending messages", "num messages", len(messages))
	if len(messages) == 0 {
		return nil
	}

	msgString := ""
	highestMessageOutputType := core.MessageOutputType(0)
	for _, msg := range messages {
		if msg.Type > highestMessageOutputType {
			highestMessageOutputType = msg.Type
		}

		msgString += createMessageString(msg, httpBoldFormat, httpBoldedLinkFormat)
	}

	title := createTitle(highestMessageOutputType, messages[0].ExecutorName)

	err := notifier.pushNotification(msgString, title)
	if err != nil {
		return fmt.Errorf("%w in pushoverNotifier.OutputMessages", err)
	}

	return nil
}

func createMessageString(
	msg core.OutputMessage,
	boldFormat string,
	linkFormat string,
) string {
	identifier := processIdentifier(msg, boldFormat, linkFormat)
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

func processIdentifier(
	msg core.OutputMessage,
	boldFormat string,
	linkFormat string,
) string {
	if len(msg.IdentifierURL) == 0 {
		if len(msg.ShortIdentifier) == 0 {
			return ""
		}

		return fmt.Sprintf(boldFormat, msg.ShortIdentifier)
	}

	return fmt.Sprintf(linkFormat, msg.IdentifierURL, msg.ShortIdentifier)
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
	if !core.IsHttpStatusCodeSuccess(statusCode) {
		return fmt.Errorf("%w, but %d", errReturnCodeIsNotOk, statusCode)
	}

	resp := &pushoverResponse{}
	err = json.Unmarshal(responseBytes, &resp)
	if err != nil {
		return err
	}

	log.Debug("pushoverNotifier.pushNotification: sent notification",
		"status", resp.Status, "request ID", resp.Request)

	return nil
}

// Name returns the name of the notifier
func (notifier *pushoverNotifier) Name() string {
	return fmt.Sprintf("%T", notifier)
}

// IsInterfaceNil returns true if there is no value under the interface
func (notifier *pushoverNotifier) IsInterfaceNil() bool {
	return notifier == nil
}
