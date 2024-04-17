package checkers

import (
	"fmt"
	"testing"

	"multiversx/mvx-keys-monitor/core"

	"github.com/stretchr/testify/assert"
)

func TestNewBLSRatingsChecker(t *testing.T) {
	t.Parallel()

	t.Run("invalid values for alarmDeltaRatingDrop", func(t *testing.T) {
		t.Parallel()

		hexBlsKeys := []string{"bls1", "bls2"}
		instance, err := NewBLSRatingsChecker(hexBlsKeys, "test", -0.00001)
		assert.Nil(t, instance)
		assert.ErrorIs(t, err, errInvalidAlarmDeltaRatingDrop)

		instance, err = NewBLSRatingsChecker(hexBlsKeys, "test", -100)
		assert.Nil(t, instance)
		assert.ErrorIs(t, err, errInvalidAlarmDeltaRatingDrop)

		instance, err = NewBLSRatingsChecker(hexBlsKeys, "test", 100.000001)
		assert.Nil(t, instance)
		assert.ErrorIs(t, err, errInvalidAlarmDeltaRatingDrop)

		instance, err = NewBLSRatingsChecker(hexBlsKeys, "test", 50000)
		assert.Nil(t, instance)
		assert.ErrorIs(t, err, errInvalidAlarmDeltaRatingDrop)
	})
	t.Run("should work with 0 keys", func(t *testing.T) {
		t.Parallel()

		instance, err := NewBLSRatingsChecker(nil, "test", 1.0)
		assert.NotNil(t, instance)
		assert.Nil(t, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		hexBlsKeys := []string{"bls1", "bls2"}
		instance, err := NewBLSRatingsChecker(hexBlsKeys, "test", 1.0)
		assert.NotNil(t, instance)
		assert.Nil(t, err)
	})
}

func TestBlsRatingsChecker_Check(t *testing.T) {
	t.Parallel()

	testMap := map[string]*core.ValidatorStatistics{
		"bls1": {
			TempRating: 0,
			Rating:     0,
		},
		"bls2": {
			TempRating: 9.99,
			Rating:     50,
		},
		"bls3": {
			TempRating: 10,
			Rating:     50,
		},
		"bls4": {
			TempRating: 100,
			Rating:     100,
		},
		"bls5": {
			TempRating: 100,
			Rating:     50,
		},
		"bls6": {
			TempRating: 90.651,
			Rating:     50,
		},
		"bls7": {
			TempRating: 49.651,
			Rating:     50,
		},
		"bls8": {
			TempRating: 48.999,
			Rating:     50,
		},
		"bls9": {
			TempRating: 49,
			Rating:     50,
		},
		"bls10": {
			TempRating: 50,
			Rating:     50,
		},
		"": { // this should not be triggered on any test
			TempRating: 0,
			Rating:     0,
		},
	}

	t.Run("nil map should error", func(t *testing.T) {
		t.Parallel()

		instance, _ := NewBLSRatingsChecker([]string{"bls1", "bls2"}, "test", 1.0)
		response, err := instance.Check(nil, nil)
		assert.Equal(t, errNilMapProvided, err)
		assert.Nil(t, response)
	})
	t.Run("bls key is not a validator should not signal", func(t *testing.T) {
		t.Parallel()

		instance, _ := NewBLSRatingsChecker([]string{"bls_not_found"}, "test", 1.0)
		response, err := instance.Check(testMap, nil)
		assert.Nil(t, err)
		assert.Empty(t, response)
	})
	t.Run("under the jail threshold should signal", func(t *testing.T) {
		t.Parallel()

		instance, _ := NewBLSRatingsChecker([]string{"bls1", "bls2"}, "test", 1.0)
		response, err := instance.Check(testMap, nil)
		assert.Nil(t, err)
		expectedCheckResponse := []core.CheckResponse{
			{
				HexBLSKey: "bls1",
				Status:    fmt.Sprintf(imminentJailMessageFormat, 0.0, 0.0),
			},
			{
				HexBLSKey: "bls2",
				Status:    fmt.Sprintf(imminentJailMessageFormat, 9.99, 50.0),
			},
		}
		assert.Equal(t, expectedCheckResponse, response)
	})
	t.Run("bls key same rating should not signal", func(t *testing.T) {
		t.Parallel()

		instance, _ := NewBLSRatingsChecker([]string{"bls4", "bls10"}, "test", 1.0)
		response, err := instance.Check(testMap, nil)
		assert.Nil(t, err)
		assert.Empty(t, response)
	})
	t.Run("bls key rating increase not signal", func(t *testing.T) {
		t.Parallel()

		instance, _ := NewBLSRatingsChecker([]string{"bls5", "bls6"}, "test", 1.0)
		response, err := instance.Check(testMap, nil)
		assert.Nil(t, err)
		assert.Empty(t, response)
	})
	t.Run("bls key rating drop above limit should not signal", func(t *testing.T) {
		t.Parallel()

		instance, _ := NewBLSRatingsChecker([]string{"bls7"}, "test", 1.0)
		response, err := instance.Check(testMap, nil)
		assert.Nil(t, err)
		assert.Empty(t, response)
	})
	t.Run("bls key rating drop beyond limit should signal", func(t *testing.T) {
		t.Parallel()

		instance, _ := NewBLSRatingsChecker([]string{"bls3", "bls8", "bls9"}, "test", 1.0)
		response, err := instance.Check(testMap, nil)
		assert.Nil(t, err)
		expectedCheckResponse := []core.CheckResponse{
			{
				HexBLSKey: "bls3",
				Status:    fmt.Sprintf(ratingDropMessageFormat, 10.0, 50.0),
			},
			{
				HexBLSKey: "bls8",
				Status:    fmt.Sprintf(ratingDropMessageFormat, 48.999, 50.0),
			},
			{
				HexBLSKey: "bls9",
				Status:    fmt.Sprintf(ratingDropMessageFormat, 49.0, 50.0),
			},
		}
		assert.Equal(t, expectedCheckResponse, response)
	})
	t.Run("bls key rating drop beyond limit should signal on extra keys", func(t *testing.T) {
		t.Parallel()

		instance, _ := NewBLSRatingsChecker(nil, "test", 1.0)
		response, err := instance.Check(testMap, []string{"bls3", "bls8", "bls9", ""})
		assert.Nil(t, err)
		expectedCheckResponse := []core.CheckResponse{
			{
				HexBLSKey: "bls3",
				Status:    fmt.Sprintf(ratingDropMessageFormat, 10.0, 50.0),
			},
			{
				HexBLSKey: "bls8",
				Status:    fmt.Sprintf(ratingDropMessageFormat, 48.999, 50.0),
			},
			{
				HexBLSKey: "bls9",
				Status:    fmt.Sprintf(ratingDropMessageFormat, 49.0, 50.0),
			},
		}
		assert.Equal(t, expectedCheckResponse, response)
	})
}

func TestBlsRatingsChecker_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var instance *blsRatingsChecker
	assert.True(t, instance.IsInterfaceNil())

	instance = &blsRatingsChecker{}
	assert.False(t, instance.IsInterfaceNil())
}
