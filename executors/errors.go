package executors

import "errors"

var (
	errNilRatingsChecker             = errors.New("nil ratings checker instance")
	errNilOutputNotifier             = errors.New("nil output notifier")
	errNilOutputNotifiersHandler     = errors.New("nil output notifiers handler")
	errNilValidatorStatisticsQuerier = errors.New("nil validator statistics querier")
	errNilStatusHandler              = errors.New("nil status handler")
	errNilBLSKeysFetcher             = errors.New("nil BLS Keys fetcher")
	errNilTimeFunc                   = errors.New("nil pointer for the current time function")
	errNilExecutor                   = errors.New("nil executor")
	errInvalidWeekDay                = errors.New("invalid week day")
	errInvalidHour                   = errors.New("invalid hour")
	errInvalidMinute                 = errors.New("invalid minute")
	errInvalidTimeBetweenRetries     = errors.New("invalid time between retries")
	errNotificationsSendingProblems  = errors.New("notification sending problems")
	errNilCurrentTimestampHandler    = errors.New("nil current timestamp")
	errNilBLSKeysFilter              = errors.New("nil BLS keys filter")
)
