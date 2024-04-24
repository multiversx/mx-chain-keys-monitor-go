package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsHttpStatusCodeSuccess(t *testing.T) {
	t.Parallel()

	assert.False(t, IsHttpStatusCodeSuccess(0))
	assert.False(t, IsHttpStatusCodeSuccess(100))
	assert.False(t, IsHttpStatusCodeSuccess(199))
	assert.True(t, IsHttpStatusCodeSuccess(200))
	assert.True(t, IsHttpStatusCodeSuccess(201))
	assert.True(t, IsHttpStatusCodeSuccess(250))
	assert.True(t, IsHttpStatusCodeSuccess(298))
	assert.True(t, IsHttpStatusCodeSuccess(299))
	assert.False(t, IsHttpStatusCodeSuccess(300))
	assert.False(t, IsHttpStatusCodeSuccess(301))
	assert.False(t, IsHttpStatusCodeSuccess(404))
	assert.False(t, IsHttpStatusCodeSuccess(500))
}
