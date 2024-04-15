package aogo

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/stretchr/testify/assert"
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
	tags := []types.Tag{{Name: "Action", Value: "Stakers"}}

	s, err := goar.NewSignerFromPath("./keys/wallet.json")
	assert.NoError(t, err)

	itemSigner, err := goar.NewItemSigner(s)
	assert.NoError(t, err)

	ts := httptest.NewServer(nil)
	defer ts.Close()

	mu := NewMUMock(ts.URL)
	res, err := mu.SendMessage(process, data, tags, "", itemSigner)
	assert.NoError(t, err)
	assert.True(t, res != "")
}

func TestSpawnProcess(t *testing.T) {
	data := ""
	tags := []types.Tag{{Name: "Action", Value: "Stakers"}}

	s, err := goar.NewSignerFromPath("./keys/wallet.json")
	assert.NoError(t, err)

	itemSigner, err := goar.NewItemSigner(s)
	assert.NoError(t, err)

	ts := httptest.NewServer(nil)
	defer ts.Close()

	mu := NewMUMock(ts.URL)
	res, err := mu.SpawnProcess(data, tags, itemSigner)

	assert.NoError(t, err)
	assert.True(t, res != "")
}
