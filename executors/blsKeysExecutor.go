package executors

import (
	"context"
	"fmt"
	"net/url"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-keys-monitor-go/core"
	logger "github.com/multiversx/mx-chain-logger-go"
)

const identifierType = "BLS key"
const explorerURLNodesPathName = "nodes"
const numPrefixCharactersForKey = 6
const numSuffixCharactersForKey = 6

var log = logger.GetOrCreate("executors")

type blsKeysExecutor struct {
	notifiers                  []OutputNotifier
	ratingsChecker             RatingsChecker
	blsKeysFetcher             BLSKeysFetcher
	validatorStatisticsQuerier ValidatorStatisticsQuerier
	statusHandler              StatusHandler
	name                       string
	explorerURL                string
}

// ArgsBLSKeysExecutor defines the DTO struct for the NewBLSKeysExecutor constructor function
type ArgsBLSKeysExecutor struct {
	Notifiers                  []OutputNotifier
	RatingsChecker             RatingsChecker
	ValidatorStatisticsQuerier ValidatorStatisticsQuerier
	BlsKeysFetcher             BLSKeysFetcher
	StatusHandler              StatusHandler
	Name                       string
	ExplorerURL                string
}

// NewBLSKeysExecutor creates a new instance of type blsKeysExecutor
func NewBLSKeysExecutor(args ArgsBLSKeysExecutor) (*blsKeysExecutor, error) {
	if check.IfNil(args.RatingsChecker) {
		return nil, errNilRatingsChecker
	}
	for index, notifier := range args.Notifiers {
		if check.IfNil(notifier) {
			return nil, fmt.Errorf("%w at index %d", errNilOutputNotifier, index)
		}
	}
	if check.IfNil(args.ValidatorStatisticsQuerier) {
		return nil, errNilValidatorStatisticsQuerier
	}
	if check.IfNil(args.StatusHandler) {
		return nil, errNilStatusHandler
	}
	if check.IfNil(args.BlsKeysFetcher) {
		return nil, errNilBLSKeysFetcher
	}

	return &blsKeysExecutor{
		notifiers:                  args.Notifiers,
		ratingsChecker:             args.RatingsChecker,
		validatorStatisticsQuerier: args.ValidatorStatisticsQuerier,
		statusHandler:              args.StatusHandler,
		name:                       args.Name,
		explorerURL:                args.ExplorerURL,
		blsKeysFetcher:             args.BlsKeysFetcher,
	}, nil
}

// Execute executes one checking cycle
func (executor *blsKeysExecutor) Execute(ctx context.Context) error {
	log.Debug("executing query-check-notify cycle", "executor", executor.name)
	statistics, err := executor.validatorStatisticsQuerier.Query(ctx)
	if err != nil {
		log.Debug("blsKeysExecutor.Execute", "executor", executor.name, "error calling query", err.Error())
		executor.statusHandler.ErrorEncountered(err)

		return err
	}

	extraBLSKeys, err := executor.blsKeysFetcher.GetAllBLSKeys(ctx, executor.name)
	if err != nil {
		log.Debug("blsKeysExecutor.Execute", "executor", executor.name, "error calling GetAllBLSKeys", err.Error())
		executor.statusHandler.ErrorEncountered(err)

		return err
	}

	problematicKeys, err := executor.ratingsChecker.Check(statistics, extraBLSKeys)
	if err != nil {
		log.Debug("blsKeysExecutor.Execute", "executor", executor.name, "error calling check", err.Error())
		executor.statusHandler.ErrorEncountered(err)

		return err
	}

	if len(problematicKeys) == 0 {
		log.Debug("all keys are performing normally", "executor", executor.name)

		return nil
	}

	messages := executor.createMessages(problematicKeys)
	executor.notify(messages)

	return nil
}

func (executor *blsKeysExecutor) createMessages(problematicKeys []core.CheckResponse) []core.OutputMessage {
	result := make([]core.OutputMessage, 0, len(problematicKeys))
	for _, key := range problematicKeys {
		message := core.OutputMessage{
			Type:               core.ErrorMessageOutputType,
			IdentifierType:     identifierType,
			Identifier:         key.HexBLSKey,
			ShortIdentifier:    shortIdentifier(key.HexBLSKey),
			IdentifierURL:      executor.createIdentifierURL(key.HexBLSKey),
			ExecutorName:       executor.name,
			ProblemEncountered: key.Status,
		}

		result = append(result, message)
	}

	return result
}

func (executor *blsKeysExecutor) notify(messages []core.OutputMessage) {
	log.Debug("blsKeysExecutor.notify", "executor", executor.name, "num messages", len(messages))

	executor.statusHandler.ProblematicBLSKeysFound(messages)
	for _, notifier := range executor.notifiers {
		notifier.OutputMessages(messages...)
	}
}

func shortIdentifier(identifier string) string {
	ellipsisString := "..."
	minNumCharactersToTrim := numPrefixCharactersForKey + len(ellipsisString) + numSuffixCharactersForKey
	if len(identifier) <= minNumCharactersToTrim {
		return identifier
	}

	return identifier[:numPrefixCharactersForKey] + ellipsisString + identifier[len(identifier)-numSuffixCharactersForKey:]
}

func (executor *blsKeysExecutor) createIdentifierURL(hexKey string) string {
	if len(executor.explorerURL) == 0 {
		return ""
	}

	result, err := url.JoinPath(executor.explorerURL, explorerURLNodesPathName, hexKey)
	if err != nil {
		log.Debug("blsKeysExecutor.createIdentifierURL",
			"executor.explorerURL", executor.explorerURL, "hexKey", hexKey, "error", err)
	}

	return result
}

// IsInterfaceNil returns true if there is no value under the interface
func (executor *blsKeysExecutor) IsInterfaceNil() bool {
	return executor == nil
}
