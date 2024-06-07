package aogo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/types"
	"github.com/stretchr/testify/assert"
)

func TestNewAO(t *testing.T) {
	ao, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, ao)
	assert.NotNil(t, ao.mu)
	assert.NotNil(t, ao.cu)
}