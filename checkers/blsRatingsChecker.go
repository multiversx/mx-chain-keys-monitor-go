package checkers

import (
	"fmt"

	"multiversx/mvx-keys-monitor/core"

	logger "github.com/multiversx/mx-chain-logger-go"
)

const minAlarmDeltaRatingDrop = float64(0.0)
const maxAlarmDeltaRatingDrop = float64(100.0)
const jailThreshold = float32(10)
const ratingDropMessageFormat = "Rating drop detected: temp rating: %0.2f, rating: %0.2f"
const imminentJailMessageFormat = "Imminent jail: temp rating: %0.2f, rating: %0.2f"

var log = logger.GetOrCreate("checkers")

type blsRatingsChecker struct {
	name                 string
	hexBlsKeys           []string
	alarmDeltaRatingDrop float32
}

// NewBLSRatingsChecker creates a new instance of type blsRatingsChecker
func NewBLSRatingsChecker(hexBlsKeys []string, name string, alarmDeltaRatingDrop float64) (*blsRatingsChecker, error) {
	log.Debug("NewBLSRatingsChecker", "checker name", name, "num initial keys", len(hexBlsKeys))

	if alarmDeltaRatingDrop < minAlarmDeltaRatingDrop || alarmDeltaRatingDrop > maxAlarmDeltaRatingDrop {
		return nil, fmt.Errorf("%w, allowed interval [%0.2f, %0.2f]", errInvalidAlarmDeltaRatingDrop, minAlarmDeltaRatingDrop, maxAlarmDeltaRatingDrop)
	}

	return &blsRatingsChecker{
		name:                 name,
		hexBlsKeys:           hexBlsKeys,
		alarmDeltaRatingDrop: float32(alarmDeltaRatingDrop),
	}, nil
}

// Check will check all containing BLS keys if there is a rating drop for that key. It returns the list of BLS keys with rating drop or if an error occurred
func (checker *blsRatingsChecker) Check(statistics map[string]*core.ValidatorStatistics, extraBLSKeys []string) ([]core.CheckResponse, error) {
	if statistics == nil {
		return nil, errNilMapProvided
	}

	allKeys := checker.hexBlsKeys
	if len(extraBLSKeys) > 0 {
		allKeys = append(allKeys, extraBLSKeys...)
	}

	log.Debug("blsRatingsChecker.Check", "checker name", checker.name, "num keys", len(allKeys))

	response := make([]core.CheckResponse, 0)
	for _, blsKey := range allKeys {
		if len(blsKey) == 0 {
			continue
		}

		log.Trace("checking", "checker", checker.name, "bls key", blsKey)
		stats, found := statistics[blsKey]
		if !found {
			continue
		}
		if stats.TempRating < jailThreshold {
			log.Debug("found imminent jail node", "checker", checker.name, "bls key", blsKey)
			response = append(response, core.CheckResponse{
				HexBLSKey: blsKey,
				Status:    fmt.Sprintf(imminentJailMessageFormat, stats.TempRating, stats.Rating),
			})
			continue
		}
		if stats.TempRating >= stats.Rating {
			continue
		}
		if stats.TempRating+checker.alarmDeltaRatingDrop > stats.Rating {
			continue
		}

		log.Debug("found node with rating drop", "checker", checker.name, "bls key", blsKey)
		response = append(response, core.CheckResponse{
			HexBLSKey: blsKey,
			Status:    fmt.Sprintf(ratingDropMessageFormat, stats.TempRating, stats.Rating),
		})
	}

	return response, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (checker *blsRatingsChecker) IsInterfaceNil() bool {
	return checker == nil
}
