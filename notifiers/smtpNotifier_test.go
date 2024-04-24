package notifiers

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"testing"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/stretchr/testify/assert"
)

func TestNewSmtpNotifier(t *testing.T) {
	t.Parallel()

	notifier := NewSmtpNotifier(ArgsSmtpNotifier{})
	assert.NotNil(t, notifier)
}

func TestSmtpNotifier_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var instance *smtpNotifier
	assert.True(t, instance.IsInterfaceNil())

	instance = &smtpNotifier{}
	assert.False(t, instance.IsInterfaceNil())
}

func TestSmtpNotifier_OutputMessages(t *testing.T) {
	testArgs := ArgsSmtpNotifier{
		To:       "to@email.com",
		SmtpPort: 37,
		SmtpHost: "host.email.com",
		From:     "from@email.com",
		Password: "pass",
	}
	expectedErr := errors.New("expected error")

	t.Run("sending empty slice of messages should not call the service", func(t *testing.T) {
		t.Parallel()

		notifier := NewSmtpNotifier(testArgs)
		notifier.sendMail = func(host string, auth smtp.Auth, from string, to []string, msgBytes []byte) error {
			assert.Fail(t, "should have not called sendMail function")

			return nil
		}
		notifier.OutputMessages()
	})
	t.Run("send mail function fails, should error", func(t *testing.T) {
		t.Parallel()

		notifier := NewSmtpNotifier(testArgs)
		notifier.sendMail = func(host string, auth smtp.Auth, from string, to []string, msgBytes []byte) error {
			return expectedErr
		}
		err := notifier.pushNotification("test", "title")
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, expectedErr)
	})
	t.Run("sending info messages should work", func(t *testing.T) {
		t.Parallel()

		msg1 := core.OutputMessage{
			Type:               core.InfoMessageOutputType,
			IdentifierType:     "info1",
			ExecutorName:       "executor",
			Identifier:         "info2",
			ShortIdentifier:    "info3",
			IdentifierURL:      "https://examples.com/info3",
			ProblemEncountered: "problem1",
		}
		msg2 := core.OutputMessage{
			Type:           core.InfoMessageOutputType,
			IdentifierType: "info10",
			ExecutorName:   "executor",
		}
		msg3 := core.OutputMessage{
			Type:            core.InfoMessageOutputType,
			ShortIdentifier: "info20",
			ExecutorName:    "executor",
		}

		expectedBody := `Subject: ⓘ Info for executor 
MIME-version: 1.0;
Content-Type: text/html; charset="UTF-8";




<!DOCTYPE html>
<html lang="en">
<body>
   ✅ info1 <b><a href="https://examples.com/info3">info3</a></b>: problem1

<br>✅ info10 <b></b>

<br>✅  <b>info20</b>

<br>
</body>
</html>
`
		var sentMsgBytes []byte
		notifier := NewSmtpNotifier(testArgs)
		notifier.sendMail = func(host string, auth smtp.Auth, from string, to []string, msgBytes []byte) error {
			assert.Equal(t, fmt.Sprintf("%s:%d", testArgs.SmtpHost, testArgs.SmtpPort), host)
			assert.Equal(t, testArgs.From, from)
			assert.Equal(t, []string{testArgs.To}, to)
			sentMsgBytes = msgBytes

			return nil
		}
		notifier.OutputMessages(msg1, msg2, msg3)
		assert.Equal(t, expectedBody, string(sentMsgBytes))
	})
	t.Run("sending info messages and warn messages should work", func(t *testing.T) {
		t.Parallel()

		msg1 := core.OutputMessage{
			Type:               core.InfoMessageOutputType,
			IdentifierType:     "info1",
			ExecutorName:       "executor",
			Identifier:         "info2",
			ShortIdentifier:    "info3",
			IdentifierURL:      "https://examples.com/info3",
			ProblemEncountered: "problem1",
		}
		msg2 := core.OutputMessage{
			Type:           core.InfoMessageOutputType,
			IdentifierType: "info10",
			ExecutorName:   "executor",
		}
		msg3 := core.OutputMessage{
			Type:            core.WarningMessageOutputType,
			ShortIdentifier: "info20",
			ExecutorName:    "executor",
		}

		expectedBody := `Subject: ⚠️ Warnings occurred on executor 
MIME-version: 1.0;
Content-Type: text/html; charset="UTF-8";




<!DOCTYPE html>
<html lang="en">
<body>
   ✅ info1 <b><a href="https://examples.com/info3">info3</a></b>: problem1

<br>✅ info10 <b></b>

<br>⚠️  <b>info20</b>

<br>
</body>
</html>
`

		var sentMsgBytes []byte
		notifier := NewSmtpNotifier(testArgs)
		notifier.sendMail = func(host string, auth smtp.Auth, from string, to []string, msgBytes []byte) error {
			assert.Equal(t, fmt.Sprintf("%s:%d", testArgs.SmtpHost, testArgs.SmtpPort), host)
			assert.Equal(t, testArgs.From, from)
			assert.Equal(t, []string{testArgs.To}, to)
			sentMsgBytes = msgBytes

			return nil
		}
		notifier.OutputMessages(msg1, msg2, msg3)
		assert.Equal(t, expectedBody, string(sentMsgBytes))
	})
	t.Run("sending info, warn and error messages should work", func(t *testing.T) {
		t.Parallel()

		msg1 := core.OutputMessage{
			Type:               core.ErrorMessageOutputType,
			IdentifierType:     "info1",
			ExecutorName:       "executor",
			Identifier:         "info2",
			ShortIdentifier:    "info3",
			IdentifierURL:      "https://examples.com/info3",
			ProblemEncountered: "problem1",
		}
		msg2 := core.OutputMessage{
			Type:           core.InfoMessageOutputType,
			IdentifierType: "info10",
			ExecutorName:   "executor",
		}
		msg3 := core.OutputMessage{
			Type:            core.WarningMessageOutputType,
			ShortIdentifier: "info20",
			ExecutorName:    "executor",
		}

		expectedBody := `Subject: 🚨 Problems occurred on executor 
MIME-version: 1.0;
Content-Type: text/html; charset="UTF-8";




<!DOCTYPE html>
<html lang="en">
<body>
   🚨 info1 <b><a href="https://examples.com/info3">info3</a></b>: problem1

<br>✅ info10 <b></b>

<br>⚠️  <b>info20</b>

<br>
</body>
</html>
`

		var sentMsgBytes []byte
		notifier := NewSmtpNotifier(testArgs)
		notifier.sendMail = func(host string, auth smtp.Auth, from string, to []string, msgBytes []byte) error {
			assert.Equal(t, fmt.Sprintf("%s:%d", testArgs.SmtpHost, testArgs.SmtpPort), host)
			assert.Equal(t, testArgs.From, from)
			assert.Equal(t, []string{testArgs.To}, to)
			sentMsgBytes = msgBytes

			return nil
		}
		notifier.OutputMessages(msg1, msg2, msg3)
		assert.Equal(t, expectedBody, string(sentMsgBytes))
	})
	t.Run("sending unknown type of messages should work", func(t *testing.T) {
		t.Parallel()

		msg1 := core.OutputMessage{
			Type:               0,
			IdentifierType:     "info1",
			ExecutorName:       "executor",
			Identifier:         "info2",
			ShortIdentifier:    "info3",
			IdentifierURL:      "https://examples.com/info3",
			ProblemEncountered: "problem1",
		}
		msg2 := core.OutputMessage{
			Type:           0,
			IdentifierType: "info10",
			ExecutorName:   "executor",
		}
		msg3 := core.OutputMessage{
			Type:            0,
			ShortIdentifier: "info20",
			ExecutorName:    "executor",
		}

		expectedBody := `Subject: executor 
MIME-version: 1.0;
Content-Type: text/html; charset="UTF-8";




<!DOCTYPE html>
<html lang="en">
<body>
    info1 <b><a href="https://examples.com/info3">info3</a></b>: problem1

<br> info10 <b></b>

<br>  <b>info20</b>

<br>
</body>
</html>
`

		var sentMsgBytes []byte
		notifier := NewSmtpNotifier(testArgs)
		notifier.sendMail = func(host string, auth smtp.Auth, from string, to []string, msgBytes []byte) error {
			assert.Equal(t, fmt.Sprintf("%s:%d", testArgs.SmtpHost, testArgs.SmtpPort), host)
			assert.Equal(t, testArgs.From, from)
			assert.Equal(t, []string{testArgs.To}, to)
			sentMsgBytes = msgBytes

			return nil
		}
		notifier.OutputMessages(msg1, msg2, msg3)
		assert.Equal(t, expectedBody, string(sentMsgBytes))
	})
}

func TestSmtpNotifier_FunctionalTest(t *testing.T) {
	// before running this test, please define your environment variables SMTP_TO, SMTP_FROM and SMTP_PASSWORD so this test can work

	t.Skip("this is a functional test, will need real credentials")

	_ = logger.SetLogLevel("*:DEBUG")

	args := ArgsSmtpNotifier{
		To:       os.Getenv("SMTP_TO"),
		SmtpPort: 587,
		SmtpHost: "smtp.gmail.com",
		From:     os.Getenv("SMTP_FROM"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}

	notifier := NewSmtpNotifier(args)

	t.Run("info messages", func(t *testing.T) {
		message1 := core.OutputMessage{
			Type:           core.InfoMessageOutputType,
			IdentifierType: "this is an info line",
			ExecutorName:   "Keys monitoring app",
		}
		message2 := core.OutputMessage{
			Type:            core.InfoMessageOutputType,
			ShortIdentifier: "this is a bold info line",
			ExecutorName:    "Keys monitoring app",
		}
		notifier.OutputMessages(message1, message2)

	})
	t.Run("info and warn messages", func(t *testing.T) {
		message1 := core.OutputMessage{
			Type:           core.InfoMessageOutputType,
			IdentifierType: "this is an info line",
			ExecutorName:   "Keys monitoring app",
		}
		message2 := core.OutputMessage{
			Type:            core.InfoMessageOutputType,
			ShortIdentifier: "this is a bold info line",
			ExecutorName:    "Keys monitoring app",
		}
		message3 := core.OutputMessage{
			Type:            core.WarningMessageOutputType,
			ShortIdentifier: "internal app errors occurred: 45",
			ExecutorName:    "Keys monitoring app",
		}
		notifier.OutputMessages(message1, message2, message3)
	})
	t.Run("error messages", func(t *testing.T) {
		message1 := core.OutputMessage{
			Type:               core.ErrorMessageOutputType,
			IdentifierType:     "BLS key",
			Identifier:         "0295e29aef11c30391a70c3578d3c3dea23da84b2465fe8bbb17cbf2d4e87ca4e416a32626f2c51e1f125054ed8720077df8daa475857a35129e8772a39112c252e67dd783acb83f6fffc70dd8a7830e599995ac4c7dd35f08664c479f7fde80",
			ShortIdentifier:    "0295e2...7fde80",
			IdentifierURL:      "https://testnet-explorer.multiversx.com/nodes/0295e29aef11c30391a70c3578d3c3dea23da84b2465fe8bbb17cbf2d4e87ca4e416a32626f2c51e1f125054ed8720077df8daa475857a35129e8772a39112c252e67dd783acb83f6fffc70dd8a7830e599995ac4c7dd35f08664c479f7fde80",
			ExecutorName:       "testnet - set 1",
			ProblemEncountered: "Rating drop detected: temp rating: 90.70, rating: 100.00",
		}
		message2 := core.OutputMessage{
			Type:               core.ErrorMessageOutputType,
			IdentifierType:     "BLS key",
			Identifier:         "fdd9e63fe5317de782e3e5268e4f0645319cda34c34d85b235602e795ffdc1cce62a9936b6288d1fe288923ac675c51195150ad696a7fa7ddbf5dc643412f1ae13502518d9869279f59d106b4f0ced3d12a1bd19d38e7762c856c680335bd015",
			ShortIdentifier:    "fdd9e6...5bd015",
			IdentifierURL:      "",
			ExecutorName:       "testnet - set 1",
			ProblemEncountered: "Rating drop detected: temp rating: 95.37, rating: 100.00",
		}
		notifier.OutputMessages(message1, message2)
	})
}
