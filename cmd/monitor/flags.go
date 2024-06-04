package main

import (
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/urfave/cli"
)

var (
	// configurationFile defines a flag for the path to the main toml configuration file
	configurationFile = cli.StringFlag{
		Name:  "config",
		Usage: "The main configuration file",
		Value: "./config/config.toml",
	}
	// credentialsFile defines a flag for the path to the credentials toml file
	credentialsFile = cli.StringFlag{
		Name:  "credentials",
		Usage: "The credentials configuration file",
		Value: "./config/credentials.toml",
	}
	// logLevel defines the logger level
	logLevel = cli.StringFlag{
		Name: "log-level",
		Usage: "This flag specifies the logger `level(s)`. It can contain multiple comma-separated value. For example" +
			", if set to *:INFO the logs for all packages will have the INFO level. However, if set to *:INFO,api:DEBUG" +
			" the logs for all packages will have the INFO level, excepting the api package which will receive a DEBUG" +
			" log level.",
		Value: "*:" + logger.LogInfo.String(),
	}
	// logFile is used when the log output needs to be logged in a file
	logSaveFile = cli.BoolFlag{
		Name:  "log-save",
		Usage: "Boolean option for enabling log saving. If set, it will automatically save all the logs into a file.",
	}
	// testNotifiers is used when trying to put the application in testing mode to try the configured notifiers
	testNotifiers = cli.BoolFlag{
		Name:  "test-notifiers",
		Usage: "Boolean option testing out the notifiers. The application will send a message to all configured notifiers and will close.",
	}
)
