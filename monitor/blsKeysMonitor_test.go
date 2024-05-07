package monitor

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-keys-monitor-go/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewBLSKeysMonitor(t *testing.T) {
	t.Parallel()

	t.Run("invalid interval should error", func(t *testing.T) {
		t.Parallel()

		monitor, err := NewBLSKeysMonitor(&mock.ExecutorStub{}, time.Second-time.Nanosecond, "test")
		assert.Nil(t, monitor)
		assert.ErrorIs(t, err, errInvalidInterval)
		assert.Contains(t, err.Error(), "got 999.999999ms, minimum 1s")
	})
	t.Run("nil executor should error", func(t *testing.T) {
		t.Parallel()

		monitor, err := NewBLSKeysMonitor(nil, time.Second, "test")
		assert.Nil(t, monitor)
		assert.NotNil(t, err)
		assert.Equal(t, "nil executor", err.Error())
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		monitor, err := NewBLSKeysMonitor(&mock.ExecutorStub{}, time.Second, "test")
		assert.NotNil(t, monitor)
		assert.Nil(t, err)

		_ = monitor.Close()
	})
}

func TestBlsKeysMonitor_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var instance *blsKeysMonitor
	assert.True(t, instance.IsInterfaceNil())

	instance = &blsKeysMonitor{}
	assert.False(t, instance.IsInterfaceNil())
}

func TestBlsKeysMonitor_Close(t *testing.T) {
	t.Parallel()

	numExecution := uint32(0)
	monitor, _ := NewBLSKeysMonitor(&mock.ExecutorStub{
		ExecuteHandler: func(ctx context.Context) error {
			atomic.AddUint32(&numExecution, 1)
			return nil
		},
	}, time.Second, "test")
	assert.NotNil(t, monitor)

	time.Sleep(time.Second + time.Millisecond*300)
	assert.Equal(t, uint32(2), atomic.LoadUint32(&numExecution))

	time.Sleep(time.Second)
	assert.Equal(t, uint32(3), atomic.LoadUint32(&numExecution))

	err := monitor.Close()
	assert.Nil(t, err)

	// no more increment calls
	time.Sleep(time.Second * 2)
	assert.Equal(t, uint32(3), atomic.LoadUint32(&numExecution))
}
