package disabled

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDisabledCloser(t *testing.T) {
	t.Parallel()

	instance := NewDisabledCloser()
	assert.NotNil(t, instance)
}

func TestDisabledCloser_Close(t *testing.T) {
	t.Parallel()

	instance := NewDisabledCloser()
	assert.Nil(t, instance.Close())
}
