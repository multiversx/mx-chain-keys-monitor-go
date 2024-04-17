package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsStatusCodeIs2xx(t *testing.T) {
	t.Parallel()

	assert.False(t, IsStatusCodeIs2xx(0))
	assert.False(t, IsStatusCodeIs2xx(100))
	assert.False(t, IsStatusCodeIs2xx(199))
	assert.True(t, IsStatusCodeIs2xx(200))
	assert.True(t, IsStatusCodeIs2xx(201))
	assert.True(t, IsStatusCodeIs2xx(250))
	assert.True(t, IsStatusCodeIs2xx(298))
	assert.True(t, IsStatusCodeIs2xx(299))
	assert.False(t, IsStatusCodeIs2xx(300))
	assert.False(t, IsStatusCodeIs2xx(301))
	assert.False(t, IsStatusCodeIs2xx(404))
	assert.False(t, IsStatusCodeIs2xx(500))
}
