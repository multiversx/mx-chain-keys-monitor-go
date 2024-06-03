package executors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBLSKeysTimeCache(t *testing.T) {
	t.Parallel()

	t.Run("nil timestamp handler should error", func(t *testing.T) {
		t.Parallel()

		args := ArgsBlsKeysTimeCache{}
		timeCache, err := NewBLSKeysTimeCache(args)
		assert.Nil(t, timeCache)
		assert.Equal(t, errNilCurrentTimestampHandler, err)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := ArgsBlsKeysTimeCache{
			GetCurrentTimestamp: func() int64 {
				return 0
			},
		}
		timeCache, err := NewBLSKeysTimeCache(args)
		assert.NotNil(t, timeCache)
		assert.Nil(t, err)
	})
}

func TestBlsKeysTimeCache_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var instance *blsKeysTimeCache
	assert.True(t, instance.IsInterfaceNil())

	instance = &blsKeysTimeCache{}
	assert.False(t, instance.IsInterfaceNil())
}

func TestBlsKeysTimeCache_ShouldNotify(t *testing.T) {
	t.Parallel()

	timestampCounter := int64(0)
	args := ArgsBlsKeysTimeCache{
		GetCurrentTimestamp: func() int64 {
			timestampCounter++
			return timestampCounter
		},
		CacheExpiration: 100,
		MaxSnoozeEvents: 3,
	}
	timeCache, _ := NewBLSKeysTimeCache(args)

	// these tests can be run in parallel provided that each test uses a new key
	t.Run("first fault should return true", func(t *testing.T) {
		t.Parallel()

		key := "key1"
		assert.True(t, timeCache.ShouldNotify(key))
	})
	t.Run("first & second fault should return true", func(t *testing.T) {
		t.Parallel()

		key := "key2"
		assert.True(t, timeCache.ShouldNotify(key))
		assert.True(t, timeCache.ShouldNotify(key))
	})
	t.Run("first, second and third fault should return true", func(t *testing.T) {
		t.Parallel()

		key := "key3"
		assert.True(t, timeCache.ShouldNotify(key))
		assert.True(t, timeCache.ShouldNotify(key))
		assert.True(t, timeCache.ShouldNotify(key))
	})
	t.Run("first, second and third fault should return true, forth one should return false", func(t *testing.T) {
		t.Parallel()

		key := "key4"
		assert.True(t, timeCache.ShouldNotify(key))
		assert.True(t, timeCache.ShouldNotify(key))
		assert.True(t, timeCache.ShouldNotify(key))
		assert.False(t, timeCache.ShouldNotify(key))
	})
	t.Run("first, second and third fault should return true, forth and fifth should return false, "+
		"after expiration time should return true", func(t *testing.T) {
		t.Parallel()

		// we need a new instance for this test
		localTimestampCounter := int64(0)
		localArgs := ArgsBlsKeysTimeCache{
			GetCurrentTimestamp: func() int64 {
				localTimestampCounter++
				return localTimestampCounter
			},
			CacheExpiration: 100,
			MaxSnoozeEvents: 3,
		}
		localTimeCache, _ := NewBLSKeysTimeCache(localArgs)

		key := "key"
		assert.True(t, localTimeCache.ShouldNotify(key))
		assert.True(t, localTimeCache.ShouldNotify(key))
		assert.True(t, localTimeCache.ShouldNotify(key))
		assert.False(t, localTimeCache.ShouldNotify(key))
		assert.False(t, localTimeCache.ShouldNotify(key))
		// time passes, the info gets expired
		localTimestampCounter = int64(args.CacheExpiration) // the automatic increment in the handler stands for the initial current timestamp of 1
		assert.True(t, localTimeCache.ShouldNotify(key))
	})
}
