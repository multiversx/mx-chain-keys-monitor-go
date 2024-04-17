package checkers

import "errors"

var (
	errInvalidAlarmDeltaRatingDrop = errors.New("invalid value for AlarmDeltaRatingDrop")
	errNilMapProvided              = errors.New("nil map provided")
)
