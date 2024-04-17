package notifiers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"multiversx/mvx-keys-monitor/core"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/stretchr/testify/assert"
)

const testToken = "test-token"
const testUserKey = "test-user-key"

func createHttpTestServerThatRespondsOK(
	t *testing.T,
	expectedMessage string,
	expectedTitle string,
	numCalls *uint32,
) *httptest.Server {
	response := &pushoverResponse{
		Status:  1,
		Request: "e43a9e0f-6836-42f1-8b06-e8bc56012637",
	}
	responseBytes, _ := json.Marshal(response)

	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		body := req.Body
		defer func() {
			errClose := body.Close()
			assert.Nil(t, errClose)
		}()
		buff := make([]byte, 524288)
		numRead, _ := body.Read(buff)
		buff = buff[:numRead]

		request := &pushoverRequest{}
		err := json.Unmarshal(buff, request)
		assert.Nil(t, err)

		assert.Equal(t, testToken, request.Token)
		assert.Equal(t, testUserKey, request.User)
		assert.Equal(t, 1, request.HTML)
		assert.Equal(t, expectedMessage, request.Message)
		assert.Equal(t, expectedTitle, request.Title)

		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write(responseBytes)
		atomic.AddUint32(numCalls, 1)
	}))
}

func TestNewPushoverNotifier(t *testing.T) {
	t.Parallel()

	notifier := NewPushoverNotifier("url", "", "")
	assert.NotNil(t, notifier)
}

func TestPushoverNotifier_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var instance *pushoverNotifier
	assert.True(t, instance.IsInterfaceNil())

	instance = &pushoverNotifier{}
	assert.False(t, instance.IsInterfaceNil())
}

func TestPushoverNotifier_OutputMessages(t *testing.T) {
	t.Parallel()

	t.Run("sending empty slice of messages should not call the service", func(t *testing.T) {
		t.Parallel()

		numCalls := uint32(0)
		expectedTitle := ""
		expectedMessage := ""
		testServer := createHttpTestServerThatRespondsOK(t, expectedMessage, expectedTitle, &numCalls)
		defer testServer.Close()

		notifier := NewPushoverNotifier(testServer.URL, testToken, testUserKey)
		notifier.OutputMessages()

		time.Sleep(time.Second)
		assert.Equal(t, uint32(0), atomic.LoadUint32(&numCalls))
	})
	t.Run("post method fails should error", func(t *testing.T) {
		t.Parallel()

		notifier := NewPushoverNotifier("not-a-server-URL", "", "")
		err := notifier.pushNotification("test", "title")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "not-a-server-URL")
	})
	t.Run("server errors should error", func(t *testing.T) {
		t.Parallel()

		testHttpServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		}))

		notifier := NewPushoverNotifier(testHttpServer.URL, "", "")
		err := notifier.pushNotification("test", "title")
		assert.ErrorIs(t, err, errReturnCodeIsNotOk)
	})
	t.Run("http post response is not OK, should error", func(t *testing.T) {
		t.Parallel()

		testHttpServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write([]byte("not-a-valid-json"))
		}))

		notifier := NewPushoverNotifier(testHttpServer.URL, "", "")
		err := notifier.pushNotification("test", "title")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid character")

		// make sure any accidental calls on API endpoint routes are caught by the test server
		time.Sleep(time.Second)
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
		expectedMessage := `✅ info1 <b><a href="https://examples.com/info3">info3</a></b>: problem1

✅ info10 <b></b>

✅  <b>info20</b>

`

		testServer := createHttpTestServerThatRespondsOK(t, expectedMessage, expectedTitle, &numCalls)
		defer testServer.Close()

		notifier := NewPushoverNotifier(testServer.URL, testToken, testUserKey)
		notifier.OutputMessages(msg1, msg2, msg3)

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
		expectedMessage := `✅ info1 <b><a href="https://examples.com/info3">info3</a></b>: problem1

✅ info10 <b></b>

⚠️  <b>info20</b>

`

		testServer := createHttpTestServerThatRespondsOK(t, expectedMessage, expectedTitle, &numCalls)
		defer testServer.Close()

		notifier := NewPushoverNotifier(testServer.URL, testToken, testUserKey)
		notifier.OutputMessages(msg1, msg2, msg3)

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
		expectedMessage := `🚨 info1 <b><a href="https://examples.com/info3">info3</a></b>: problem1

✅ info10 <b></b>

⚠️  <b>info20</b>

`

		testServer := createHttpTestServerThatRespondsOK(t, expectedMessage, expectedTitle, &numCalls)
		defer testServer.Close()

		notifier := NewPushoverNotifier(testServer.URL, testToken, testUserKey)
		notifier.OutputMessages(msg1, msg2, msg3)

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
		expectedMessage := ` info1 <b><a href="https://examples.com/info3">info3</a></b>: problem1

 info10 <b></b>

  <b>info20</b>

`

		testServer := createHttpTestServerThatRespondsOK(t, expectedMessage, expectedTitle, &numCalls)
		defer testServer.Close()

		notifier := NewPushoverNotifier(testServer.URL, testToken, testUserKey)
		notifier.OutputMessages(msg1, msg2, msg3)

		time.Sleep(time.Second)
		assert.Equal(t, uint32(1), atomic.LoadUint32(&numCalls))
	})
}

func TestPushoverNotifier_FunctionalTest(t *testing.T) {
	t.Skip("this is a functional test, will need real credentials")

	_ = logger.SetLogLevel("*:DEBUG")

	notifier := NewPushoverNotifier(
		"https://api.pushover.net/1/messages.json",
		"", // TODO: replace here your token value
		"", // TODO: replace here your userKey value
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
