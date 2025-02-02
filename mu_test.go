package aogo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/tag"
	"github.com/stretchr/testify/assert"
)

func NewMUMock(URL string) MU {
	return MU{
		client: http.DefaultClient,
		url:    URL,
	}
}

func TestSendMessage(t *testing.T) {
	t.Run("0", func(t *testing.T) {
		process := "yugMfaR-u_11GkAuZhqeChPuzoxVYuJW8RnNCIby-D8"
		data := ""
		tags := &[]tag.Tag{{Name: "Action", Value: "Stakers"}}

		s, err := signer.FromPath("./keys/wallet.json")
		assert.NoError(t, err)
		muServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"id": "mockMessageID"}`))
			assert.NoError(t, err)
		}))
		defer muServer.Close()

		mu := NewMUMock(muServer.URL)
		res, err := mu.SendMessage(process, data, tags, "", s)
		assert.NoError(t, err)
		assert.True(t, res != "")
	})

	t.Run("1", func(t *testing.T) {
		muServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"id": "mockMessageID"}`))
			assert.NoError(t, err)
		}))
		defer muServer.Close()

		ao := &AO{mu: newMU(muServer.URL)}

		s, err := signer.FromPath("./keys/wallet.json") // Mock signer or use a real one for the test
		assert.NoError(t, err)

		id, err := ao.SendMessage("process", "data", nil, "", s)
		assert.NoError(t, err)
		assert.Equal(t, "mockMessageID", id)
	})
}

func TestSpawnProcess(t *testing.T) {
	t.Run("0", func(t *testing.T) {
		muServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			w.WriteHeader(http.StatusAccepted)
			_, err := w.Write([]byte(`{"id": "mockMessageID"}`))
			assert.NoError(t, err)
		}))
		defer muServer.Close()

		mu := NewMUMock(muServer.URL)

		tags := []tag.Tag{{Name: "Action", Value: "Stakers"}}

		s, err := signer.FromPath("./keys/wallet.json")
		assert.NoError(t, err)

		res, err := mu.SpawnProcess("", nil, tags, s)

		assert.NoError(t, err)
		assert.True(t, res != "")
	})

	t.Run("1", func(t *testing.T) {
		muServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			w.WriteHeader(http.StatusAccepted)
			_, err := w.Write([]byte(`{"id": "mockProcessID"}`))
			assert.NoError(t, err)
		}))
		defer muServer.Close()

		ao := &AO{mu: newMU(muServer.URL)}

		signer, err := signer.FromPath("./keys/wallet.json") // Mock signer or use a real one for the test
		assert.NoError(t, err)

		id, err := ao.SpawnProcess("module", []byte("data"), nil, signer)
		assert.NoError(t, err)
		assert.Equal(t, "mockProcessID", id)
	})
}
