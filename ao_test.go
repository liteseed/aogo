package aogo

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/stretchr/testify/assert"
)

func TestNewAO(t *testing.T) {
	ao, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, ao)
	assert.NotNil(t, ao.mu)
	assert.NotNil(t, ao.cu)
}

func TestSpawnProcess(t *testing.T) {
	muServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"id": "mockProcessID"}`))
	}))
	defer muServer.Close()

	ao := &AO{mu: newMU(muServer.URL)}

	signer := goar.NewItemSigner(...) // Mock signer or use a real one for the test
	id, err := ao.SpawnProcess("module", "data", nil, signer)
	assert.NoError(t, err)
	assert.Equal(t, "mockProcessID", id)
}

func TestSendMessage(t *testing.T) {
	muServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "mockMessageID"}`))
	}))
	defer muServer.Close()

	ao := &AO{mu: newMU(muServer.URL)}

	signer := goar.NewItemSigner(...) // Mock signer or use a real one for the test
	id, err := ao.SendMessage("process", "data", nil, "", signer)
	assert.NoError(t, err)
	assert.Equal(t, "mockMessageID", id)
}

func TestLoadResult(t *testing.T) {
	cuServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Messages": [], "Spawns": [], "Outputs": [], "Error": "", "GasUsed": 0}`))
	}))
	defer cuServer.Close()

	ao := &AO{cu: newCU(cuServer.URL)}

	resp, err := ao.LoadResult("process", "message")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 0, resp.GasUsed)
}

func TestDryRun(t *testing.T) {
	cuServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Messages": [], "Spawns": [], "Outputs": [], "Error": "", "GasUsed": 0}`))
	}))
	defer cuServer.Close()

	ao := &AO{cu: newCU(cuServer.URL)}

	message := Message{
		ID:     "testID",
		Target: "testTarget",
		Owner:  "testOwner",
		Data:   "testData",
		Tags:   []types.Tag{},
	}
	resp, err := ao.DryRun(message)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 0, resp.GasUsed)
}
