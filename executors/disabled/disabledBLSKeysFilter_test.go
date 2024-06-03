package disabled

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDisabledBLSKeysFilter(t *testing.T) {
	t.Parallel()

	handler := NewDisabledBLSKeysFilter()
	assert.NotNil(t, handler)
}

func TestDisabledBLSKeysFilter_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var instance *disabledBLSKeysFilter
	assert.True(t, instance.IsInterfaceNil())

	instance = &disabledBLSKeysFilter{}
	assert.False(t, instance.IsInterfaceNil())
}

func TestDisabledBLSKeysFilter_ShouldNotify(t *testing.T) {
	t.Parallel()

	handler := NewDisabledBLSKeysFilter()
	assert.True(t, handler.ShouldNotify(""))
	assert.True(t, handler.ShouldNotify("random string"))
}
