package aogo

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/types"
	"github.com/stretchr/testify/assert"
)

func NewMUMock(URL string) MU {
	return MU{
		client: http.DefaultClient,
		url:    URL,
	}
}

func TestSendMessage0(t *testing.T) {
	process := "yugMfaR-u_11GkAuZhqeChPuzoxVYuJW8RnNCIby-D8"
	data := ""
	tags := []types.Tag{{Name: "Action", Value: "Stakers"}}

	s, err := signer.FromPath("./keys/wallet.json")
	assert.NoError(t, err)

	ts := httptest.NewServer(nil)
	defer ts.Close()

	mu := NewMUMock(ts.URL)
	res, err := mu.SendMessage(process, data, tags, "", s)
	assert.NoError(t, err)
	assert.True(t, res != "")
}

func TestSendMessage1(t *testing.T) {
	muServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "mockMessageID"}`))
	}))
	defer muServer.Close()

	ao := &AO{mu: newMU(muServer.URL)}

	signer := signer.FromPath("./keys/arweave.json") // Mock signer or use a real one for the test
	id, err := ao.SendMessage("process", "data", nil, "", signer)
	assert.NoError(t, err)
	assert.Equal(t, "mockMessageID", id)
}

func TestSpawnProcess0(t *testing.T) {
	data := ""
	tags := []types.Tag{{Name: "Action", Value: "Stakers"}}

	s, err := signer.FromPath("./keys/wallet.json")
	assert.NoError(t, err)

	ts := httptest.NewServer(nil)
	defer ts.Close()

	mu := NewMUMock(ts.URL)
	res, err := mu.SpawnProcess("", data, tags, s)

	assert.NoError(t, err)
	assert.True(t, res != "")
}

func TestSpawnProcess1(t *testing.T) {
	muServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"id": "mockProcessID"}`))
	}))
	defer muServer.Close()

	ao := &AO{mu: newMU(muServer.URL)}

	signer := signer.FromPath("./keys/arweave.json") // Mock signer or use a real one for the test
	id, err := ao.SpawnProcess("module", "data", nil, signer)
	assert.NoError(t, err)
	assert.Equal(t, "mockProcessID", id)
}
