package integrationTests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-keys-monitor-go/config"
	"github.com/multiversx/mx-chain-keys-monitor-go/core"
	"github.com/multiversx/mx-chain-keys-monitor-go/executors"
	"github.com/multiversx/mx-chain-keys-monitor-go/factory"
	"github.com/multiversx/mx-chain-keys-monitor-go/mock"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlsKeysNotifier(t *testing.T) {
	data, err := os.ReadFile("testdata/response.json")
	require.Nil(t, err)

	_ = logger.SetLogLevel("*:DEBUG")

	testHttpServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write(data)
	}))

	cfg := config.BLSKeysMonitorConfig{
		AlarmDeltaRatingDrop:     1,
		ApiURL:                   testHttpServer.URL,
		PollingIntervalInSeconds: 2,
		ListFile:                 "./testdata/keys.list",
		Name:                     "integration-test",
		ExplorerURL:              "https://examples.com",
	}
	errorHandler := &mock.StatusHandlerStub{
		ErrorEncounteredHandler: func(err error) {
			require.Fail(t, "should have not called ErrorEncountered")
		},
	}

	mut := sync.RWMutex{}
	result := make([]core.OutputMessage, 0)
	notifier := &mock.OutputNotifierStub{
		OutputMessagesHandler: func(messages ...core.OutputMessage) error {
			mut.Lock()
			result = append(result, messages...)
			mut.Unlock()

			return nil
		},
	}

	allConfigs := config.AllConfigs{
		Config: config.MainConfig{
			OutputNotifiers: config.OutputNotifiersConfig{
				NumRetries:            0,
				SecondsBetweenRetries: 1,
				Pushover: config.PushoverNotifierConfig{
					Enabled: false,
				},
				Smtp: config.SmtpNotifierConfig{
					Enabled: false,
				},
				Telegram: config.TelegramNotifierConfig{
					Enabled: false,
				},
			},
		},
	}

	// test also the output notifiers factory
	notifiers, err := factory.CreateOutputNotifiers(allConfigs)
	assert.Nil(t, err)

	notifiers = append(notifiers, notifier) // add the notifier used in this test

	argsNotifiersHandler := executors.ArgsNotifiersHandler{
		Notifiers:          notifiers,
		NumRetries:         allConfigs.Config.OutputNotifiers.NumRetries,
		TimeBetweenRetries: time.Second * time.Duration(allConfigs.Config.OutputNotifiers.SecondsBetweenRetries),
	}

	notifiersHandler, err := executors.NewNotifiersHandler(argsNotifiersHandler)
	assert.Nil(t, err)

	monitor, err := factory.NewBLSKeysMonitor(
		cfg,
		config.AlarmSnoozeConfig{},
		notifiersHandler,
		errorHandler,
	)
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)

	err = monitor.Close()
	assert.Nil(t, err)

	expectedMessages := []core.OutputMessage{
		{
			IdentifierType:     "BLS key",
			Identifier:         "0026a4b6d8f4b6a2e22141341efb5dddf4db130a7e04d539dfd8c70bf3139d016ed958ccfd0bdcaf6aa866d11a09e21058f9dcb96fab9b863fe832cbed7f1705970ab9ad8c4de9da69e59a890751740064bfd84b7eb9e714e0e03fa8d776d004",
			ShortIdentifier:    "0026a4...76d004",
			IdentifierURL:      "https://examples.com/nodes/0026a4b6d8f4b6a2e22141341efb5dddf4db130a7e04d539dfd8c70bf3139d016ed958ccfd0bdcaf6aa866d11a09e21058f9dcb96fab9b863fe832cbed7f1705970ab9ad8c4de9da69e59a890751740064bfd84b7eb9e714e0e03fa8d776d004",
			ExecutorName:       "integration-test",
			ProblemEncountered: "Rating drop detected: temp rating: 90.70, rating: 100.00",
			Type:               core.ErrorMessageOutputType,
		},
		{
			IdentifierType:     "BLS key",
			Identifier:         "0295e29aef11c30391a70c3578d3c3dea23da84b2465fe8bbb17cbf2d4e87ca4e416a32626f2c51e1f125054ed8720077df8daa475857a35129e8772a39112c252e67dd783acb83f6fffc70dd8a7830e599995ac4c7dd35f08664c479f7fde80",
			ShortIdentifier:    "0295e2...7fde80",
			IdentifierURL:      "https://examples.com/nodes/0295e29aef11c30391a70c3578d3c3dea23da84b2465fe8bbb17cbf2d4e87ca4e416a32626f2c51e1f125054ed8720077df8daa475857a35129e8772a39112c252e67dd783acb83f6fffc70dd8a7830e599995ac4c7dd35f08664c479f7fde80",
			ExecutorName:       "integration-test",
			ProblemEncountered: "Rating drop detected: temp rating: 48.70, rating: 50.00",
			Type:               core.ErrorMessageOutputType,
		},
		{
			IdentifierType:     "BLS key",
			Identifier:         "fdd9e63fe5317de782e3e5268e4f0645319cda34c34d85b235602e795ffdc1cce62a9936b6288d1fe288923ac675c51195150ad696a7fa7ddbf5dc643412f1ae13502518d9869279f59d106b4f0ced3d12a1bd19d38e7762c856c680335bd015",
			ShortIdentifier:    "fdd9e6...5bd015",
			IdentifierURL:      "https://examples.com/nodes/fdd9e63fe5317de782e3e5268e4f0645319cda34c34d85b235602e795ffdc1cce62a9936b6288d1fe288923ac675c51195150ad696a7fa7ddbf5dc643412f1ae13502518d9869279f59d106b4f0ced3d12a1bd19d38e7762c856c680335bd015",
			ExecutorName:       "integration-test",
			ProblemEncountered: "Imminent jail: temp rating: 8.70, rating: 50.00",
			Type:               core.ErrorMessageOutputType,
		},
		// second iteration
		{
			IdentifierType:     "BLS key",
			Identifier:         "0026a4b6d8f4b6a2e22141341efb5dddf4db130a7e04d539dfd8c70bf3139d016ed958ccfd0bdcaf6aa866d11a09e21058f9dcb96fab9b863fe832cbed7f1705970ab9ad8c4de9da69e59a890751740064bfd84b7eb9e714e0e03fa8d776d004",
			ShortIdentifier:    "0026a4...76d004",
			IdentifierURL:      "https://examples.com/nodes/0026a4b6d8f4b6a2e22141341efb5dddf4db130a7e04d539dfd8c70bf3139d016ed958ccfd0bdcaf6aa866d11a09e21058f9dcb96fab9b863fe832cbed7f1705970ab9ad8c4de9da69e59a890751740064bfd84b7eb9e714e0e03fa8d776d004",
			ExecutorName:       "integration-test",
			ProblemEncountered: "Rating drop detected: temp rating: 90.70, rating: 100.00",
			Type:               core.ErrorMessageOutputType,
		},
		{
			IdentifierType:     "BLS key",
			Identifier:         "0295e29aef11c30391a70c3578d3c3dea23da84b2465fe8bbb17cbf2d4e87ca4e416a32626f2c51e1f125054ed8720077df8daa475857a35129e8772a39112c252e67dd783acb83f6fffc70dd8a7830e599995ac4c7dd35f08664c479f7fde80",
			ShortIdentifier:    "0295e2...7fde80",
			IdentifierURL:      "https://examples.com/nodes/0295e29aef11c30391a70c3578d3c3dea23da84b2465fe8bbb17cbf2d4e87ca4e416a32626f2c51e1f125054ed8720077df8daa475857a35129e8772a39112c252e67dd783acb83f6fffc70dd8a7830e599995ac4c7dd35f08664c479f7fde80",
			ExecutorName:       "integration-test",
			ProblemEncountered: "Rating drop detected: temp rating: 48.70, rating: 50.00",
			Type:               core.ErrorMessageOutputType,
		},
		{
			IdentifierType:     "BLS key",
			Identifier:         "fdd9e63fe5317de782e3e5268e4f0645319cda34c34d85b235602e795ffdc1cce62a9936b6288d1fe288923ac675c51195150ad696a7fa7ddbf5dc643412f1ae13502518d9869279f59d106b4f0ced3d12a1bd19d38e7762c856c680335bd015",
			ShortIdentifier:    "fdd9e6...5bd015",
			IdentifierURL:      "https://examples.com/nodes/fdd9e63fe5317de782e3e5268e4f0645319cda34c34d85b235602e795ffdc1cce62a9936b6288d1fe288923ac675c51195150ad696a7fa7ddbf5dc643412f1ae13502518d9869279f59d106b4f0ced3d12a1bd19d38e7762c856c680335bd015",
			ExecutorName:       "integration-test",
			ProblemEncountered: "Imminent jail: temp rating: 8.70, rating: 50.00",
			Type:               core.ErrorMessageOutputType,
		},
	}

	mut.RLock()
	assert.Equal(t, expectedMessages, result)
	mut.RUnlock()
}
