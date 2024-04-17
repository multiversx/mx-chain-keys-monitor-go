package config

// CredentialsConfig defines the credentials configuration file
type CredentialsConfig struct {
	Pushover PushoverCredentialsConfig
	Smtp     EmailPasswordConfig
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
