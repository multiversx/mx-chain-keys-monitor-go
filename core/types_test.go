package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageOutputType_String(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "unknown MessageOutputType: 0", MessageOutputType(0).String())
	assert.Equal(t, "info", InfoMessageOutputType.String())
	assert.Equal(t, "warn", WarningMessageOutputType.String())
	assert.Equal(t, "error", ErrorMessageOutputType.String())
	assert.Equal(t, "unknown MessageOutputType: 100", MessageOutputType(100).String())
}
