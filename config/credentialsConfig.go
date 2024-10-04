package config

// CredentialsConfig defines the credentials configuration file
type CredentialsConfig struct {
	Pushover PushoverCredentialsConfig
	Smtp     EmailPasswordConfig
	Telegram TelegramCredentialsConfig
	Slack    SlackCredentialsConfig
}

// TokenUserKeyConfig defines a struct that contains one token and one user key
type TokenUserKeyConfig struct {
	Token   string
	UserKey string
}

// TokenChatIDConfig defines a struct that contains one token and one chat ID
type TokenChatIDConfig struct {
	Token  string
	ChatID string
}

// PushoverCredentialsConfig defines the Pushover service credentials
type PushoverCredentialsConfig struct {
	TokenUserKeyConfig
	Additional []TokenUserKeyConfig
}

// EmailPasswordConfig defines an email-password credential
type EmailPasswordConfig struct {
	Email    string
	Password string
}

// TelegramCredentialsConfig defines the Telegram service credentials
type TelegramCredentialsConfig struct {
	TokenChatIDConfig
	Additional []TokenChatIDConfig
}

// SlackSecretConfig defines a slack secret credential
type SlackSecretConfig struct {
	Secret string
}

// SlackCredentialsConfig defines the Slack service credentials
type SlackCredentialsConfig struct {
	SlackSecretConfig
	Additional []SlackSecretConfig
}
