package executors

import (
	"context"
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-keys-monitor-go/core"
	"github.com/multiversx/mx-chain-keys-monitor-go/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewBLSKeysExecutor(t *testing.T) {
	t.Parallel()

	testArgs := ArgsBLSKeysExecutor{
		OutputNotifiersHandler:     &mock.OutputNotifiersHandlerStub{},
		RatingsChecker:             &mock.RatingsCheckerStub{},
		ValidatorStatisticsQuerier: &mock.ValidatorStatisticsQuerierStub{},
		StatusHandler:              &mock.StatusHandlerStub{},
		BlsKeysFetcher:             &mock.BLSKEysFetcherStub{},
		BLSKeysFilter:              &mock.BLSKeysFilterStub{},
	}

	t.Run("nil checker should error", func(t *testing.T) {
		t.Parallel()

		localArgs := testArgs
		localArgs.RatingsChecker = nil
		executor, err := NewBLSKeysExecutor(localArgs)
		assert.Nil(t, executor)
		assert.Equal(t, errNilRatingsChecker, err)
	})
	t.Run("nil output notifiers handler should error", func(t *testing.T) {
		t.Parallel()

		localArgs := testArgs
		localArgs.OutputNotifiersHandler = nil
		executor, err := NewBLSKeysExecutor(localArgs)
		assert.Nil(t, executor)
		assert.Equal(t, errNilOutputNotifiersHandler, err)
	})
	t.Run("nil querier should error", func(t *testing.T) {
		t.Parallel()

		localArgs := testArgs
		localArgs.ValidatorStatisticsQuerier = nil
		executor, err := NewBLSKeysExecutor(localArgs)
		assert.Nil(t, executor)
		assert.Equal(t, errNilValidatorStatisticsQuerier, err)
	})
	t.Run("nil status handler should error", func(t *testing.T) {
		t.Parallel()

		localArgs := testArgs
		localArgs.StatusHandler = nil
		executor, err := NewBLSKeysExecutor(localArgs)
		assert.Nil(t, executor)
		assert.Equal(t, errNilStatusHandler, err)
	})
	t.Run("nil BLS keys fetcher should error", func(t *testing.T) {
		t.Parallel()

		localArgs := testArgs
		localArgs.BlsKeysFetcher = nil
		executor, err := NewBLSKeysExecutor(localArgs)
		assert.Nil(t, executor)
		assert.Equal(t, errNilBLSKeysFetcher, err)
	})
	t.Run("nil BLS keys filter should error", func(t *testing.T) {
		t.Parallel()

		localArgs := testArgs
		localArgs.BLSKeysFilter = nil
		executor, err := NewBLSKeysExecutor(localArgs)
		assert.Nil(t, executor)
		assert.Equal(t, errNilBLSKeysFilter, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		localArgs := testArgs
		executor, err := NewBLSKeysExecutor(localArgs)
		assert.NotNil(t, executor)
		assert.Nil(t, err)
	})
}

func TestBlsKeysExecutor_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var instance *blsKeysExecutor
	assert.True(t, instance.IsInterfaceNil())

	instance = &blsKeysExecutor{}
	assert.False(t, instance.IsInterfaceNil())
}

func TestBlsKeysExecutor_Execute(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("expected error")
	t.Run("query errors, should error", func(t *testing.T) {
		t.Parallel()

		encounteredError := false
		args := ArgsBLSKeysExecutor{
			OutputNotifiersHandler: &mock.OutputNotifiersHandlerStub{
				NotifyWithRetryHandler: func(caller string, messages ...core.OutputMessage) error {
					assert.Fail(t, "should have not called the output notifiers handler")
					return nil
				},
			},
			RatingsChecker: &mock.RatingsCheckerStub{
				CheckHandler: func(statistics map[string]*core.ValidatorStatistics, extraBLSKeys []string) ([]core.CheckResponse, error) {
					return make([]core.CheckResponse, 0), nil
				},
			},
			ValidatorStatisticsQuerier: &mock.ValidatorStatisticsQuerierStub{
				QueryHandler: func(ctx context.Context) (map[string]*core.ValidatorStatistics, error) {
					return nil, expectedErr
				},
			},
			StatusHandler: &mock.StatusHandlerStub{
				ErrorEncounteredHandler: func(err error) {
					assert.Equal(t, expectedErr, err)
					encounteredError = true
				},
				CollectKeysProblemsHandler: func(messages []core.OutputMessage) {
					assert.Fail(t, "should have not called the status handler")
				},
			},
			BlsKeysFetcher: &mock.BLSKEysFetcherStub{},
			BLSKeysFilter:  &mock.BLSKeysFilterStub{},
		}
		executor, _ := NewBLSKeysExecutor(args)
		err := executor.Execute(context.Background())
		assert.Equal(t, expectedErr, err)
		assert.True(t, encounteredError)
	})
	t.Run("BLS keys fetcher errors, should error", func(t *testing.T) {
		t.Parallel()

		encounteredError := false
		args := ArgsBLSKeysExecutor{
			OutputNotifiersHandler: &mock.OutputNotifiersHandlerStub{
				NotifyWithRetryHandler: func(caller string, messages ...core.OutputMessage) error {
					assert.Fail(t, "should have not called the output notifiers handler")
					return nil
				},
			},
			RatingsChecker: &mock.RatingsCheckerStub{
				CheckHandler: func(statistics map[string]*core.ValidatorStatistics, extraBLSKeys []string) ([]core.CheckResponse, error) {
					return make([]core.CheckResponse, 0), nil
				},
			},
			ValidatorStatisticsQuerier: &mock.ValidatorStatisticsQuerierStub{
				QueryHandler: func(ctx context.Context) (map[string]*core.ValidatorStatistics, error) {
					return make(map[string]*core.ValidatorStatistics), nil
				},
			},
			StatusHandler: &mock.StatusHandlerStub{
				ErrorEncounteredHandler: func(err error) {
					assert.Equal(t, expectedErr, err)
					encounteredError = true
				},
				CollectKeysProblemsHandler: func(messages []core.OutputMessage) {
					assert.Fail(t, "should have not called the status handler")
				},
			},
			BlsKeysFetcher: &mock.BLSKEysFetcherStub{
				GetAllBLSKeysHandler: func(ctx context.Context, sender string) ([]string, error) {
					return nil, expectedErr
				},
			},
			BLSKeysFilter: &mock.BLSKeysFilterStub{},
		}
		executor, _ := NewBLSKeysExecutor(args)
		err := executor.Execute(context.Background())
		assert.Equal(t, expectedErr, err)
		assert.True(t, encounteredError)
	})
	t.Run("ratings check errors, should error", func(t *testing.T) {
		t.Parallel()

		encounteredError := false
		args := ArgsBLSKeysExecutor{
			OutputNotifiersHandler: &mock.OutputNotifiersHandlerStub{
				NotifyWithRetryHandler: func(caller string, messages ...core.OutputMessage) error {
					assert.Fail(t, "should have not called the output notifiers handler")
					return nil
				},
			},
			RatingsChecker: &mock.RatingsCheckerStub{
				CheckHandler: func(statistics map[string]*core.ValidatorStatistics, extraBLSKeys []string) ([]core.CheckResponse, error) {
					return nil, expectedErr
				},
			},
			ValidatorStatisticsQuerier: &mock.ValidatorStatisticsQuerierStub{
				QueryHandler: func(ctx context.Context) (map[string]*core.ValidatorStatistics, error) {
					return nil, nil
				},
			},
			StatusHandler: &mock.StatusHandlerStub{
				ErrorEncounteredHandler: func(err error) {
					assert.Equal(t, expectedErr, err)
					encounteredError = true
				},
				CollectKeysProblemsHandler: func(messages []core.OutputMessage) {
					assert.Fail(t, "should have not called the status handler")
				},
			},
			BlsKeysFetcher: &mock.BLSKEysFetcherStub{},
			BLSKeysFilter:  &mock.BLSKeysFilterStub{},
		}
		executor, _ := NewBLSKeysExecutor(args)
		err := executor.Execute(context.Background())
		assert.Equal(t, expectedErr, err)
		assert.True(t, encounteredError)
	})
	t.Run("should work for 0 problematic keys", func(t *testing.T) {
		t.Parallel()

		args := ArgsBLSKeysExecutor{
			OutputNotifiersHandler: &mock.OutputNotifiersHandlerStub{
				NotifyWithRetryHandler: func(caller string, messages ...core.OutputMessage) error {
					assert.Fail(t, "should have not called the output notifiers handler")
					return nil
				},
			},
			RatingsChecker: &mock.RatingsCheckerStub{
				CheckHandler: func(statistics map[string]*core.ValidatorStatistics, extraBLSKeys []string) ([]core.CheckResponse, error) {
					return make([]core.CheckResponse, 0), nil
				},
			},
			ValidatorStatisticsQuerier: &mock.ValidatorStatisticsQuerierStub{
				QueryHandler: func(ctx context.Context) (map[string]*core.ValidatorStatistics, error) {
					return nil, nil
				},
			},
			StatusHandler: &mock.StatusHandlerStub{
				ErrorEncounteredHandler: func(err error) {
					assert.Fail(t, "should have not called error encountered")
				},
				CollectKeysProblemsHandler: func(messages []core.OutputMessage) {
					assert.Fail(t, "should have not called the status handler")
				},
			},
			BlsKeysFetcher: &mock.BLSKEysFetcherStub{},
			BLSKeysFilter:  &mock.BLSKeysFilterStub{},
		}
		executor, _ := NewBLSKeysExecutor(args)
		err := executor.Execute(context.Background())
		assert.Nil(t, err)
	})
	t.Run("should work for 2 problematic keys", func(t *testing.T) {
		t.Parallel()

		bls1 := "bls1"
		bls2 := "bls2"

		var outputNotifierMessages []core.OutputMessage
		var statusHandlerMessages []core.OutputMessage
		getAllBLSKeysCalled := false
		args := ArgsBLSKeysExecutor{
			OutputNotifiersHandler: &mock.OutputNotifiersHandlerStub{
				NotifyWithRetryHandler: func(caller string, messages ...core.OutputMessage) error {
					outputNotifierMessages = messages
					assert.Equal(t, blsExecutorName, caller)

					return nil
				},
			},
			RatingsChecker: &mock.RatingsCheckerStub{
				CheckHandler: func(statistics map[string]*core.ValidatorStatistics, extraBLSKeys []string) ([]core.CheckResponse, error) {
					assert.Equal(t, []string{"extra key"}, extraBLSKeys)

					return []core.CheckResponse{
						{
							HexBLSKey: bls1,
							Status:    "status1",
						},
						{
							HexBLSKey: bls2,
							Status:    "status2",
						},
					}, nil
				},
			},
			ValidatorStatisticsQuerier: &mock.ValidatorStatisticsQuerierStub{
				QueryHandler: func(ctx context.Context) (map[string]*core.ValidatorStatistics, error) {
					return nil, nil
				},
			},
			StatusHandler: &mock.StatusHandlerStub{
				ErrorEncounteredHandler: func(err error) {
					assert.Fail(t, "should have not called error encountered")
				},
				CollectKeysProblemsHandler: func(messages []core.OutputMessage) {
					statusHandlerMessages = messages
				},
			},
			BlsKeysFetcher: &mock.BLSKEysFetcherStub{
				GetAllBLSKeysHandler: func(ctx context.Context, sender string) ([]string, error) {
					getAllBLSKeysCalled = true

					return []string{"extra key"}, nil
				},
			},
			BLSKeysFilter: &mock.BLSKeysFilterStub{},
			Name:          "executor test name",
			ExplorerURL:   "https://explorer.com",
		}
		executor, _ := NewBLSKeysExecutor(args)

		err := executor.Execute(context.Background())
		assert.Nil(t, err)

		expectedSlice := []core.OutputMessage{
			{
				Type:               core.ErrorMessageOutputType,
				IdentifierType:     "BLS key",
				Identifier:         "bls1",
				ShortIdentifier:    "bls1",
				IdentifierURL:      "https://explorer.com/nodes/bls1",
				ExecutorName:       "executor test name",
				ProblemEncountered: "status1",
			},
			{
				Type:               core.ErrorMessageOutputType,
				IdentifierType:     "BLS key",
				Identifier:         "bls2",
				ShortIdentifier:    "bls2",
				IdentifierURL:      "https://explorer.com/nodes/bls2",
				ExecutorName:       "executor test name",
				ProblemEncountered: "status2",
			},
		}

		assert.True(t, getAllBLSKeysCalled)
		assert.Equal(t, expectedSlice, outputNotifierMessages)
		assert.Equal(t, expectedSlice, statusHandlerMessages)
	})
	t.Run("should work for 2 problematic keys that are filtered out", func(t *testing.T) {
		t.Parallel()

		bls1 := "bls1"
		bls2 := "bls2"

		var outputNotifierMessages []core.OutputMessage
		var statusHandlerMessages []core.OutputMessage
		getAllBLSKeysCalled := false
		args := ArgsBLSKeysExecutor{
			OutputNotifiersHandler: &mock.OutputNotifiersHandlerStub{
				NotifyWithRetryHandler: func(caller string, messages ...core.OutputMessage) error {
					outputNotifierMessages = messages
					assert.Equal(t, blsExecutorName, caller)

					return nil
				},
			},
			RatingsChecker: &mock.RatingsCheckerStub{
				CheckHandler: func(statistics map[string]*core.ValidatorStatistics, extraBLSKeys []string) ([]core.CheckResponse, error) {
					assert.Equal(t, []string{"extra key"}, extraBLSKeys)

					return []core.CheckResponse{
						{
							HexBLSKey: bls1,
							Status:    "status1",
						},
						{
							HexBLSKey: bls2,
							Status:    "status2",
						},
					}, nil
				},
			},
			ValidatorStatisticsQuerier: &mock.ValidatorStatisticsQuerierStub{
				QueryHandler: func(ctx context.Context) (map[string]*core.ValidatorStatistics, error) {
					return nil, nil
				},
			},
			StatusHandler: &mock.StatusHandlerStub{
				ErrorEncounteredHandler: func(err error) {
					assert.Fail(t, "should have not called error encountered")
				},
				CollectKeysProblemsHandler: func(messages []core.OutputMessage) {
					statusHandlerMessages = messages
				},
			},
			BlsKeysFetcher: &mock.BLSKEysFetcherStub{
				GetAllBLSKeysHandler: func(ctx context.Context, sender string) ([]string, error) {
					getAllBLSKeysCalled = true

					return []string{"extra key"}, nil
				},
			},
			BLSKeysFilter: &mock.BLSKeysFilterStub{
				ShouldNotifyCalled: func(blsKey string) bool {
					return false
				},
			},
			Name:        "executor test name",
			ExplorerURL: "https://explorer.com",
		}
		executor, _ := NewBLSKeysExecutor(args)

		err := executor.Execute(context.Background())
		assert.Nil(t, err)

		assert.True(t, getAllBLSKeysCalled)
		assert.Empty(t, outputNotifierMessages)
		assert.Empty(t, statusHandlerMessages)
	})
	t.Run("should work for 2 problematic keys while the notifiers handler errors", func(t *testing.T) {
		t.Parallel()

		bls1 := "bls1"
		bls2 := "bls2"

		var outputNotifierMessages []core.OutputMessage
		var statusHandlerMessages []core.OutputMessage
		getAllBLSKeysCalled := false
		errorEncountered := false
		args := ArgsBLSKeysExecutor{
			OutputNotifiersHandler: &mock.OutputNotifiersHandlerStub{
				NotifyWithRetryHandler: func(caller string, messages ...core.OutputMessage) error {
					outputNotifierMessages = messages
					assert.Equal(t, blsExecutorName, caller)

					return expectedErr
				},
			},
			RatingsChecker: &mock.RatingsCheckerStub{
				CheckHandler: func(statistics map[string]*core.ValidatorStatistics, extraBLSKeys []string) ([]core.CheckResponse, error) {
					assert.Equal(t, []string{"extra key"}, extraBLSKeys)

					return []core.CheckResponse{
						{
							HexBLSKey: bls1,
							Status:    "status1",
						},
						{
							HexBLSKey: bls2,
							Status:    "status2",
						},
					}, nil
				},
			},
			ValidatorStatisticsQuerier: &mock.ValidatorStatisticsQuerierStub{
				QueryHandler: func(ctx context.Context) (map[string]*core.ValidatorStatistics, error) {
					return nil, nil
				},
			},
			StatusHandler: &mock.StatusHandlerStub{
				ErrorEncounteredHandler: func(err error) {
					assert.Equal(t, expectedErr, err)
					errorEncountered = true
				},
				CollectKeysProblemsHandler: func(messages []core.OutputMessage) {
					statusHandlerMessages = messages
				},
			},
			BlsKeysFetcher: &mock.BLSKEysFetcherStub{
				GetAllBLSKeysHandler: func(ctx context.Context, sender string) ([]string, error) {
					getAllBLSKeysCalled = true

					return []string{"extra key"}, nil
				},
			},
			BLSKeysFilter: &mock.BLSKeysFilterStub{},
			Name:          "executor test name",
			ExplorerURL:   "https://explorer.com",
		}
		executor, _ := NewBLSKeysExecutor(args)

		err := executor.Execute(context.Background())
		assert.Equal(t, expectedErr, err)

		expectedSlice := []core.OutputMessage{
			{
				Type:               core.ErrorMessageOutputType,
				IdentifierType:     "BLS key",
				Identifier:         "bls1",
				ShortIdentifier:    "bls1",
				IdentifierURL:      "https://explorer.com/nodes/bls1",
				ExecutorName:       "executor test name",
				ProblemEncountered: "status1",
			},
			{
				Type:               core.ErrorMessageOutputType,
				IdentifierType:     "BLS key",
				Identifier:         "bls2",
				ShortIdentifier:    "bls2",
				IdentifierURL:      "https://explorer.com/nodes/bls2",
				ExecutorName:       "executor test name",
				ProblemEncountered: "status2",
			},
		}

		assert.True(t, getAllBLSKeysCalled)
		assert.Equal(t, expectedSlice, outputNotifierMessages)
		assert.Equal(t, expectedSlice, statusHandlerMessages)
		assert.True(t, errorEncountered)
	})
}

func TestShortIdentifier(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "", shortIdentifier(""))
	assert.Equal(t, "1", shortIdentifier("1"))
	assert.Equal(t, "1234567890", shortIdentifier("1234567890"))
	assert.Equal(t, "1234567890A", shortIdentifier("1234567890A"))
	assert.Equal(t, "1234567890ABCD", shortIdentifier("1234567890ABCD"))
	assert.Equal(t, "1234567890ABCDE", shortIdentifier("1234567890ABCDE"))
	assert.Equal(t, "123456...ABCDEF", shortIdentifier("1234567890ABCDEF"))
	assert.Equal(t, "123456...BCDEFG", shortIdentifier("1234567890ABCDEFG"))
	assert.Equal(t, "0295e2...7fde80", shortIdentifier("0295e29aef11c30391a70c3578d3c3dea23da84b2465fe8bbb17cbf2d4e87ca4e416a32626f2c51e1f125054ed8720077df8daa475857a35129e8772a39112c252e67dd783acb83f6fffc70dd8a7830e599995ac4c7dd35f08664c479f7fde80"))
}

func TestCreateIdentifierURL(t *testing.T) {
	t.Parallel()

	testArgs := ArgsBLSKeysExecutor{
		OutputNotifiersHandler:     &mock.OutputNotifiersHandlerStub{},
		RatingsChecker:             &mock.RatingsCheckerStub{},
		ValidatorStatisticsQuerier: &mock.ValidatorStatisticsQuerierStub{},
		StatusHandler:              &mock.StatusHandlerStub{},
		BlsKeysFetcher:             &mock.BLSKEysFetcherStub{},
		BLSKeysFilter:              &mock.BLSKeysFilterStub{},
	}

	t.Run("no explorer URL string defined should return empty", func(t *testing.T) {
		t.Parallel()

		localArgs := testArgs
		localArgs.ExplorerURL = ""
		executor, _ := NewBLSKeysExecutor(localArgs)
		assert.Equal(t, "", executor.createIdentifierURL("bls1"))
	})
	t.Run("invalid character in explorer URL should return empty", func(t *testing.T) {
		localArgs := testArgs
		localArgs.ExplorerURL = string(append([]byte("https://example.com/"), 0x7f))
		executor, _ := NewBLSKeysExecutor(localArgs)
		assert.Equal(t, "", executor.createIdentifierURL("bls1"))
	})
	t.Run("with explorer URL string defined should return the correct path - finished in slash", func(t *testing.T) {
		t.Parallel()

		localArgs := testArgs
		localArgs.ExplorerURL = "https://example.com/"
		executor, _ := NewBLSKeysExecutor(localArgs)
		assert.Equal(t, "https://example.com/nodes/bls1", executor.createIdentifierURL("bls1"))
	})
	t.Run("with explorer URL string defined should return the correct path", func(t *testing.T) {
		t.Parallel()

		localArgs := testArgs
		localArgs.ExplorerURL = "https://example.com"
		executor, _ := NewBLSKeysExecutor(localArgs)
		assert.Equal(t, "https://example.com/nodes/bls1", executor.createIdentifierURL("bls1"))
	})
}
