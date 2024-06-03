package disabled

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDisabledStatusHandler(t *testing.T) {
	t.Parallel()

	handler := NewDisabledStatusHandler()
	assert.NotNil(t, handler)
}

func TestDisabledStatusHandler_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var instance *disabledStatusHandler
	assert.True(t, instance.IsInterfaceNil())

	instance = &disabledStatusHandler{}
	assert.False(t, instance.IsInterfaceNil())
}

func TestDisabledStatusHandler_MethodsShouldNotPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r != nil {
			assert.Fail(t, fmt.Sprintf("should have not panicked %v", r))
		}
	}()

	handler := NewDisabledStatusHandler()
	handler.NotifyAppStart()
	handler.CollectKeysProblems(nil)
	handler.ErrorEncountered(nil)
	handler.SendCloseMessage()
	assert.Nil(t, handler.Execute(context.Background()))
}
