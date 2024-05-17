package config

// CredentialsConfig defines the credentials configuration file
type CredentialsConfig struct {
	Pushover PushoverCredentialsConfig
	Smtp     EmailPasswordConfig
	Telegram TelegramCredentialsConfig
}

// PushoverCredentialsConfig defines the Pushover service credentials
type PushoverCredentialsConfig struct {
	Token   string
	UserKey string
}

// EmailPasswordConfig defines an email-password credential
type EmailPasswordConfig struct {
	Email    string
	Password string
}

// TelegramCredentialsConfig defines the Telegram service credentials
type TelegramCredentialsConfig struct {
	Token  string
	ChatID string
}
