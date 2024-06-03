package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-keys-monitor-go/config"
	"github.com/multiversx/mx-chain-keys-monitor-go/executors"
	"github.com/multiversx/mx-chain-keys-monitor-go/factory"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-logger-go/file"
	"github.com/pelletier/go-toml"
	"github.com/urfave/cli"
)

const (
	defaultLogsPath = "logs"
	logFilePrefix   = "monitor"
)

var (
	nodeHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}
VERSION:
   {{.Version}}
   {{end}}
`
)

// appVersion should be populated at build time using ldflags
// Usage examples:
// linux/mac:
//
//	go build -v -ldflags="-X main.appVersion=$(git describe --tags --long --dirty)"
//
// windows:
//
//	for /f %i in ('git describe --tags --long --dirty') do set VERS=%i
//	go build -v -ldflags="-X main.appVersion=%VERS%"
var appVersion = "undefined"
var closers = make([]io.Closer, 0)
var statusHandler executors.StatusHandler

func main() {
	_ = logger.SetDisplayByteSlice(logger.ToHexShort)
	log := logger.GetOrCreate("main")

	app := cli.NewApp()
	cli.AppHelpTemplate = nodeHelpTemplate
	app.Name = "MultiversX keys monitor tool"
	machineID := core.GetAnonymizedMachineID(app.Name)

	baseVersion := fmt.Sprintf("%s/%s/%s-%s", appVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	app.Version = fmt.Sprintf("%s/%s", baseVersion, machineID)
	app.Usage = "This is the entry point for starting a new MultiversX keys monitor"
	app.Flags = []cli.Flag{
		configurationFile,
		credentialsFile,
		logLevel,
		logSaveFile,
	}
	app.Authors = []cli.Author{
		{
			Name:  "The MultiversX Team",
			Email: "contact@multiversx.com",
		},
	}

	app.Action = func(c *cli.Context) error {
		return startMonitorTool(c, baseVersion, log)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func startMonitorTool(c *cli.Context, baseVersion string, log logger.Logger) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	allConfigs, err := readConfiguration(c)
	if err != nil {
		return err
	}

	fileLogging, errLogger := attachFileLogger(log, workingDir, c.GlobalBool(logSaveFile.Name), c.GlobalString(logLevel.Name))
	if errLogger != nil {
		return errLogger
	}

	if !check.IfNil(fileLogging) {
		timeLogLifeSpan := time.Second * time.Duration(allConfigs.Config.General.Logs.LogFileLifeSpanInSec)
		sizeLogLifeSpanInMB := uint64(allConfigs.Config.General.Logs.LogFileLifeSpanInMB)

		errLifespan := fileLogging.ChangeFileLifeSpan(timeLogLifeSpan, sizeLogLifeSpanInMB)
		if errLifespan != nil {
			return errLifespan
		}
	}

	log.Info("starting application...", "version", baseVersion)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	err = startMonitoring(allConfigs)
	if err != nil {
		return err
	}

	<-sigs
	log.Info("terminating at user's signal...")
	statusHandler.SendCloseMessage()

	time.Sleep(time.Second * 5)

	closeAll(log)

	if !check.IfNil(fileLogging) {
		err = fileLogging.Close()
		log.LogIfError(err)
	}

	return err
}

func readConfiguration(c *cli.Context) (config.AllConfigs, error) {
	allConfigs := config.AllConfigs{}

	err := readTomlConfiguration(c.GlobalString(configurationFile.Name), &allConfigs.Config)
	if err != nil {
		return allConfigs, err
	}

	err = readTomlConfiguration(c.GlobalString(credentialsFile.Name), &allConfigs.Credentials)
	return allConfigs, err
}

func readTomlConfiguration(filename string, object interface{}) error {
	tomlBytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return toml.Unmarshal(tomlBytes, object)
}

func attachFileLogger(log logger.Logger, workingDir string, isSaveLogFile bool, logLevel string) (factory.FileLoggingHandler, error) {
	var fileLogging factory.FileLoggingHandler
	var err error
	if isSaveLogFile {
		args := file.ArgsFileLogging{
			WorkingDir:      workingDir,
			DefaultLogsPath: defaultLogsPath,
			LogFilePrefix:   logFilePrefix,
		}
		fileLogging, err = file.NewFileLogging(args)
		if err != nil {
			return nil, fmt.Errorf("%w creating a log file", err)
		}
	}

	err = logger.SetDisplayByteSlice(logger.ToHex)
	log.LogIfError(err)
	logLevelFlagValue := logLevel
	err = logger.SetLogLevel(logLevelFlagValue)
	if err != nil {
		return nil, err
	}

	return fileLogging, nil
}

func startMonitoring(allConfigs config.AllConfigs) error {
	notifiersList, err := factory.CreateOutputNotifiers(allConfigs)
	if err != nil {
		return err
	}

	argsNotifiersHandler := executors.ArgsNotifiersHandler{
		Notifiers:          notifiersList,
		NumRetries:         allConfigs.Config.OutputNotifiers.NumRetries,
		TimeBetweenRetries: time.Duration(allConfigs.Config.OutputNotifiers.SecondsBetweenRetries) * time.Second,
	}

	notifiersHandler, err := executors.NewNotifiersHandler(argsNotifiersHandler)
	if err != nil {
		return err
	}

	var polling io.Closer
	statusHandler, polling, err = factory.CreateStatusHandler(allConfigs.Config.General, notifiersHandler)
	if err != nil {
		return err
	}

	closers = append(closers, polling)

	for _, blsKeysConfig := range allConfigs.Config.BLSKeysMonitoring {
		monitor, errCreate := factory.NewBLSKeysMonitor(
			blsKeysConfig,
			allConfigs.Config.General.AlarmSnooze,
			notifiersHandler,
			statusHandler,
		)
		if errCreate != nil {
			return errCreate
		}

		closers = append(closers, monitor)
	}

	return nil
}

func closeAll(log logger.Logger) {
	for _, closer := range closers {
		err := closer.Close()
		log.LogIfError(err)
	}
}
