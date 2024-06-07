package aogo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAO(t *testing.T) {
	ao, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, ao)
	assert.NotNil(t, ao.mu)
	assert.NotNil(t, ao.cu)
}
