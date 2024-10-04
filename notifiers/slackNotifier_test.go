package notifiers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSlackSecret = "test-token"

func createHttpTestServerThatRespondsOKForSlack(
	t *testing.T,
	expectedMessage string,
	expectedTitle string,
	numCalls *uint32,
) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		body := req.Body
		defer func() {
			errClose := body.Close()
			assert.Nil(t, errClose)
		}()
		buff := make([]byte, 524288)
		numRead, _ := body.Read(buff)
		buff = buff[:numRead]

		request := &slackRequest{}
		err := json.Unmarshal(buff, request)
		assert.Nil(t, err)

		messageString := fmt.Sprintf("%s\n\n%s", expectedTitle, expectedMessage)
		assert.Equal(t, messageString, request.Text)

		assert.Contains(t, req.URL.Path, testSlackSecret)

		rw.WriteHeader(http.StatusOK)
		atomic.AddUint32(numCalls, 1)
	}))
}

func TestNewSlackNotifier(t *testing.T) {
	t.Parallel()

	notifier := NewSlackNotifier("", "")
	require.NotNil(t, notifier)
}

func TestSlackNotifier_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var instance *slackNotifier
	assert.True(t, instance.IsInterfaceNil())

	instance = &slackNotifier{}
	assert.False(t, instance.IsInterfaceNil())
}

func TestSlackNotifier_Name(t *testing.T) {
	t.Parallel()

	notifier := NewSlackNotifier("url", "")
	assert.Equal(t, "*notifiers.slackNotifier", notifier.Name())
}

func TestSlackNotifier_OutputMessages(t *testing.T) {
	t.Parallel()

	t.Run("sending empty slice of messages should not call the service", func(t *testing.T) {
		t.Parallel()

		numCalls := uint32(0)
		expectedTitle := ""
		expectedMessage := ""
		testServer := createHttpTestServerThatRespondsOKForSlack(t, expectedMessage, expectedTitle, &numCalls)
		defer testServer.Close()

		notifier := NewSlackNotifier(testServer.URL, testSlackSecret)
		err := notifier.OutputMessages()
		assert.Nil(t, err)

		time.Sleep(time.Second)
		assert.Equal(t, uint32(0), atomic.LoadUint32(&numCalls))
	})
	t.Run("post method fails should error", func(t *testing.T) {
		t.Parallel()

		notifier := NewSlackNotifier("not-a-server-URL", "")
		err := notifier.OutputMessages(testInfoMessage)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "not-a-server-URL")
	})
	t.Run("server errors should error", func(t *testing.T) {
		t.Parallel()

		testHttpServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		}))

		notifier := NewSlackNotifier(testHttpServer.URL, "")
		err := notifier.OutputMessages(testInfoMessage)
		assert.ErrorIs(t, err, errReturnCodeIsNotOk)
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

		numCalls := uint32(0)
		expectedTitle := "ⓘ Info for executor"
		expectedMessage := `✅ info1 *<https://examples.com/info3|info3>*: problem1

✅ info10 

✅  *info20*

`

		testServer := createHttpTestServerThatRespondsOKForSlack(t, expectedMessage, expectedTitle, &numCalls)
		defer testServer.Close()

		notifier := NewSlackNotifier(testServer.URL, testSlackSecret)
		err := notifier.OutputMessages(msg1, msg2, msg3)
		assert.Nil(t, err)

		time.Sleep(time.Second)
		assert.Equal(t, uint32(1), atomic.LoadUint32(&numCalls))
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

		numCalls := uint32(0)
		expectedTitle := "⚠️ Warnings occurred on executor"
		expectedMessage := `✅ info1 *<https://examples.com/info3|info3>*: problem1

✅ info10 

⚠️  *info20*

`

		testServer := createHttpTestServerThatRespondsOKForSlack(t, expectedMessage, expectedTitle, &numCalls)
		defer testServer.Close()

		notifier := NewSlackNotifier(testServer.URL, testSlackSecret)
		err := notifier.OutputMessages(msg1, msg2, msg3)
		assert.Nil(t, err)

		time.Sleep(time.Second)
		assert.Equal(t, uint32(1), atomic.LoadUint32(&numCalls))
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

		numCalls := uint32(0)
		expectedTitle := "🚨 Problems occurred on executor"
		expectedMessage := `🚨 info1 *<https://examples.com/info3|info3>*: problem1

✅ info10 

⚠️  *info20*

`

		testServer := createHttpTestServerThatRespondsOKForSlack(t, expectedMessage, expectedTitle, &numCalls)
		defer testServer.Close()

		notifier := NewSlackNotifier(testServer.URL, testSlackSecret)
		err := notifier.OutputMessages(msg1, msg2, msg3)
		assert.Nil(t, err)

		time.Sleep(time.Second)
		assert.Equal(t, uint32(1), atomic.LoadUint32(&numCalls))
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

		numCalls := uint32(0)
		expectedTitle := "executor"
		expectedMessage := ` info1 *<https://examples.com/info3|info3>*: problem1

 info10 

  *info20*

`

		testServer := createHttpTestServerThatRespondsOKForSlack(t, expectedMessage, expectedTitle, &numCalls)
		defer testServer.Close()

		notifier := NewSlackNotifier(testServer.URL, testSlackSecret)
		err := notifier.OutputMessages(msg1, msg2, msg3)
		assert.Nil(t, err)

		time.Sleep(time.Second)
		assert.Equal(t, uint32(1), atomic.LoadUint32(&numCalls))
	})
}

func TestSlackNotifier_FunctionalTest(t *testing.T) {
	slackAppSecret := os.Getenv("SLACK_APP_SECRET")
	if len(slackAppSecret) == 0 {
		t.Skip("this is a functional test, will need real credentials. Please define your environment variable SLACK_APP_SECRET so this test can work")
	}

	_ = logger.SetLogLevel("*:DEBUG")

	notifier := NewSlackNotifier(
		"https://hooks.slack.com/services",
		slackAppSecret,
	)

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
		err := notifier.OutputMessages(message1, message2)
		assert.Nil(t, err)
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
		err := notifier.OutputMessages(message1, message2, message3)
		assert.Nil(t, err)
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
		err := notifier.OutputMessages(message1, message2)
		assert.Nil(t, err)
	})
}
