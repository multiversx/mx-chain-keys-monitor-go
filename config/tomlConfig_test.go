package config

import (
	"testing"

	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	testString := `
[General]
    ApplicationName = "Keys monitoring app"
    # the application can send messages about the internal status at regular intervals
    [General.SystemSelfCheck]
        Enabled = true
        DayOfWeek = "every day" # can also be "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday" and "Sunday"
        Hour = 12 # valid interval 0-23
        Minute = 01 # valid interval 0-59
        PollingIntervalInSec = 30
    [General.Logs]
        LogFileLifeSpanInMB = 1024 # 1GB
        LogFileLifeSpanInSec = 86400 # 1 day


[OutputNotifiers]
    # Uses Pushover service that can notify Desktop, Android or iOS devices. Requires a valid subscription.
    # If you enable this notifier, remember to specify the credentials in credentials.toml file
    [OutputNotifiers.Pushover]
        Enabled = true
        URL = "https://api.pushover.net/1/messages.json"
    
    # SMTP (email) based notification
    # If you enable this notifier, remember to specify the credentials in credentials.toml file
    # If you are using gmail server, please make sure you activate the IMAP server and use App passwords instead of the account's password
    [OutputNotifiers.Smtp]
        Enabled = true
        To = "to@email.com"
	    SmtpPort = 587
	    SmtpHost = "smtp.gmail.com"


[[BLSKeysMonitoring]]
    AlarmDeltaRatingDrop = 1.0 # maximum Rating-TempRating value that will trigger an alarm, for the public testnet might use a higher value (2 or 3)
    Name = "test 1"
    ApiURL = "test URL 1"
    ExplorerURL = "explorer URL 1"
    PollingIntervalInSeconds = 300  # 5 minutes
    ListFile = "./config/network1.list"

[[BLSKeysMonitoring]]
    AlarmDeltaRatingDrop = 2.0 # maximum Rating-TempRating value that will trigger an alarm, for the public testnet might use a higher value (2 or 3)
    Name = "test 2"
    ApiURL = "test URL 2"
    ExplorerURL = "explorer URL 2"
    PollingIntervalInSeconds = 301   # 5 minutes and 1 second
    ListFile = "./config/network2.list"
`

	expectedCfg := MainConfig{
		General: GeneralConfigs{
			ApplicationName: "Keys monitoring app",
			SystemSelfCheck: SystemSelfCheckConfig{
				Enabled:              true,
				DayOfWeek:            "every day",
				Hour:                 12,
				Minute:               1,
				PollingIntervalInSec: 30,
			},
			Logs: LogsConfig{
				LogFileLifeSpanInSec: 86400,
				LogFileLifeSpanInMB:  1024,
			},
		},
		OutputNotifiers: OutputNotifiersConfig{
			Pushover: PushoverNotifierConfig{
				Enabled: true,
				URL:     "https://api.pushover.net/1/messages.json",
			},
			Smtp: SmtpNotifierConfig{
				Enabled:  true,
				To:       "to@email.com",
				SmtpPort: 587,
				SmtpHost: "smtp.gmail.com",
			},
		},
		BLSKeysMonitoring: []BLSKeysMonitorConfig{
			{
				AlarmDeltaRatingDrop:     1.0,
				Name:                     "test 1",
				ApiURL:                   "test URL 1",
				ExplorerURL:              "explorer URL 1",
				PollingIntervalInSeconds: 300,
				ListFile:                 "./config/network1.list",
			},
			{
				AlarmDeltaRatingDrop:     2.0,
				Name:                     "test 2",
				ApiURL:                   "test URL 2",
				ExplorerURL:              "explorer URL 2",
				PollingIntervalInSeconds: 301,
				ListFile:                 "./config/network2.list",
			},
		},
	}

	cfg := MainConfig{}

	err := toml.Unmarshal([]byte(testString), &cfg)
	assert.Nil(t, err)
	assert.Equal(t, expectedCfg, cfg)
}

func TestLoadCredentialsConfig(t *testing.T) {
	t.Parallel()

	testString := `
[Pushover]
    Token="token1"
    UserKey="userKey1"

[Smtp]
    Email="from@email.com"
    Password="password"
`

	expectedCfg := CredentialsConfig{
		Pushover: PushoverCredentialsConfig{
			Token:   "token1",
			UserKey: "userKey1",
		},
		Smtp: EmailPasswordConfig{
			Email:    "from@email.com",
			Password: "password",
		},
	}

	cfg := CredentialsConfig{}

	err := toml.Unmarshal([]byte(testString), &cfg)
	assert.Nil(t, err)
	assert.Equal(t, expectedCfg, cfg)
}
