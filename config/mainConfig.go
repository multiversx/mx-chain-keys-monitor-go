package config

// MainConfig defines the main configuration file
type MainConfig struct {
	General           GeneralConfigs
	OutputNotifiers   OutputNotifiersConfig
	BLSKeysMonitoring []BLSKeysMonitorConfig
}

// GeneralConfigs defines the general configurations for the app
type GeneralConfigs struct {
	ApplicationName string
	SystemSelfCheck SystemSelfCheckConfig
	Logs            LogsConfig
}

// SystemSelfCheckConfig defines the configuration for the self check system
type SystemSelfCheckConfig struct {
	Enabled              bool
	DayOfWeek            string
	Hour                 int
	Minute               int
	PollingIntervalInSec int
}

// LogsConfig defines the logs configuration
type LogsConfig struct {
	LogFileLifeSpanInSec int
	LogFileLifeSpanInMB  int
}

// OutputNotifiersConfig specifies the implemented types of output notifiers
type OutputNotifiersConfig struct {
	NumRetries            uint32
	SecondsBetweenRetries int
	Pushover              PushoverNotifierConfig
	Smtp                  SmtpNotifierConfig
}

// PushoverNotifierConfig specifies the options for the Pushover service
type PushoverNotifierConfig struct {
	Enabled bool
	URL     string
}

// SmtpNotifierConfig specifies the options for the SMTP email service
type SmtpNotifierConfig struct {
	Enabled  bool
	To       string
	SmtpPort int
	SmtpHost string
}

// BLSKeysMonitorConfig defines the configuration for the BLS keys monitor
type BLSKeysMonitorConfig struct {
	AlarmDeltaRatingDrop     float64
	Name                     string
	ApiURL                   string
	ExplorerURL              string
	PollingIntervalInSeconds int
	ListFile                 string
}
