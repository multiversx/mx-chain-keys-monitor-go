package notifiers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"
	httpSDK "github.com/multiversx/mx-sdk-go/core/http"
)

const (
	slackBoldFormat       = "*%s*"
	slackBoldedLinkFormat = "*<%s|%s>*"
)

type slackRequest struct {
	Text string `json:"text"`
}

type slackNotifier struct {
	httpClientWrapper HTTPClientWrapper
	secret            string
}

// NewSlackNotifier will create a new Slack notifier
func NewSlackNotifier(url string, secret string) *slackNotifier {
	return &slackNotifier{
		httpClientWrapper: httpSDK.NewHttpClientWrapper(nil, url),
		secret:            secret,
	}
}

// OutputMessages will send the provided messages to Slack
func (notifier *slackNotifier) OutputMessages(messages ...core.OutputMessage) error {
	log.Debug("slackNotifier.OutputMessages sending messages", "num messages", len(messages))
	if len(messages) == 0 {
		return nil
	}

	msgString := ""
	maxMessageOutputType := core.MessageOutputType(0)
	for _, msg := range messages {
		if msg.Type > maxMessageOutputType {
			maxMessageOutputType = msg.Type
		}

		msgString += createMessageString(msg, slackBoldFormat, slackBoldedLinkFormat)
	}

	title := createTitle(maxMessageOutputType, messages[0].ExecutorName)

	err := notifier.pushNotification(msgString, title)
	if err != nil {
		return fmt.Errorf("%w in slackNotifier.OutputMessages", err)
	}

	return nil
}

func (notifier *slackNotifier) pushNotification(msgString string, title string) error {
	ctx, cancel := context.WithTimeout(context.Background(), maxSendTimeout)
	defer cancel()

	request := &slackRequest{
		Text: fmt.Sprintf("%s\n\n%s", title, msgString),
	}
	requestBuff, err := json.Marshal(request)
	if err != nil {
		return err
	}

	_, statusCode, err := notifier.httpClientWrapper.PostHTTP(ctx, notifier.secret, requestBuff)
	if err != nil {
		return err
	}
	if !core.IsHttpStatusCodeSuccess(statusCode) {
		return fmt.Errorf("%w, but %d", errReturnCodeIsNotOk, statusCode)
	}

	log.Debug("slackNotifier.pushNotification: sent notification",
		"status", statusCode)

	return nil
}

// Name returns the name of the notifier
func (notifier *slackNotifier) Name() string {
	return fmt.Sprintf("%T", notifier)
}

// IsInterfaceNil returns true if there is no value under the interface
func (notifier *slackNotifier) IsInterfaceNil() bool {
	return notifier == nil
}
