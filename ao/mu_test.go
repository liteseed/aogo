package ao

import (
	"testing"

	"net/http"
	"net/http/httptest"

	Data "github.com/liteseed/argo/data"
	"github.com/liteseed/argo/signer"
	"gotest.tools/v3/assert"
)

func NewMUMock(URL string) MU {
	return MU{
		client: http.DefaultClient,
		url:    URL,
	}
}

func TestSendMessage(t *testing.T) {
	process := "yugMfaR-u_11GkAuZhqeChPuzoxVYuJW8RnNCIby-D8"
	data := ""
	tags := []Data.Tag{{Name: "Action", Value: "Stakers"}}
	s, err := signer.New("../data/wallet.json")
	assert.NilError(t, err)

	ts := httptest.NewServer(nil)
	defer ts.Close()

	mu := NewMUMock(ts.URL)
	res, err := mu.SendMessage(process, data, tags, "", s)
	assert.NilError(t, err)
	assert.Check(t, res != "", true)
}

func TestSpawnProcess(t *testing.T) {
	data := ""
	tags := []Data.Tag{{Name: "Action", Value: "Stakers"}}
	s, err := signer.New("../data/wallet.json")
	assert.NilError(t, err)

	ts := httptest.NewServer(nil)
	defer ts.Close()

	mu := NewMUMock(ts.URL)
	res, err := mu.SpawnProcess(data, tags, s)
	assert.NilError(t, err)
	assert.Check(t, res != "", true)
}
