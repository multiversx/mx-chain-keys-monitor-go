package notifiers

import (
	"context"
	"fmt"
	"net/url"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"
	httpSDK "github.com/multiversx/mx-sdk-go/core/http"
)

type telegramNotifier struct {
	token             string
	chatID            string
	httpClientWrapper HTTPClientWrapper
}

// NewTelegramNotifier will create a new Telegram notifier
func NewTelegramNotifier(url string, token string, chatID string) *telegramNotifier {
	return &telegramNotifier{
		httpClientWrapper: httpSDK.NewHttpClientWrapper(nil, url),
		token:             token,
		chatID:            chatID,
	}
}

// OutputMessages will push the provided messages as error
func (notifier *telegramNotifier) OutputMessages(messages ...core.OutputMessage) {
	log.Debug("telegramNotifier.OutputMessages sending messages", "num messages", len(messages))
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
		log.Error("telegramNotifier.OutputMessages: error sending notification", "error", err)
	}
}

func (notifier *telegramNotifier) pushNotification(msgString string, title string) error {
	ctx, cancel := context.WithTimeout(context.Background(), maxSendTimeout)
	defer cancel()

	urlVal := url.Values{
		"chat_id":    {notifier.chatID},
		"parse_mode": {"html"},
		"text":       {fmt.Sprintf("%s\n\n%s", title, msgString)},
	}

	encodedURL := fmt.Sprintf("bot%s/sendMessage?%s", notifier.token, urlVal.Encode())
	_, statusCode, err := notifier.httpClientWrapper.PostHTTP(ctx, encodedURL, make([]byte, 0))
	if err != nil {
		return err
	}
	if !core.IsHttpStatusCodeSuccess(statusCode) {
		return fmt.Errorf("%w, but %d", errReturnCodeIsNotOk, statusCode)
	}

	log.Debug("telegramNotifier.pushNotification: sent notification",
		"status", statusCode)

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (notifier *telegramNotifier) IsInterfaceNil() bool {
	return notifier == nil
}
