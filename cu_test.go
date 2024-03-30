package aogo

import (
	"encoding/json"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func NewCUMock(URL string) CU {
	return CU{
		client: http.DefaultClient,
		url:    URL,
	}
}
func TestReadResult(t *testing.T) {
	process := "W7Ax6G1i3C4ksRRNP4Urxvq9bcSmwBK9J0S3QBt9J70"
	message := "ahcFiWM5RMcXDA-OrAdpjK10Afty6qxvELa83mMbxI0"
	messages := []map[string]any{{
		"Target": "F7fmxSBJx5RlIRrt825iIEAL110cKP2Bf8tYd0Q1STU",
		"Anchor": "00000000000000000000000000000043",
		"Data":   "{\"wjvyv-Z36LbY8y0UZ21dhzygU56GdqaDqFdT9rq-GPc\":{\"stakedAt\":1363029,\"amount\":200}}",
		"Tags": []any{
			map[string]any{
				"value": "ao",
				"name":  "Data-Protocol",
			},
			map[string]any{
				"value": "ao.TN.1",
				"name":  "Variant",
			},
			map[string]any{
				"value": "Message",
				"name":  "Type",
			},
			map[string]any{
				"value": "W7Ax6G1i3C4ksRRNP4Urxvq9bcSmwBK9J0S3QBt9J70",
				"name":  "From-Process",
			},
			map[string]any{
				"value": "2rEYpGAF-zuvgKh8-7fie7TLUdXCS1ZHa7GJ_lw3jpo",
				"name":  "From-Module",
			},
			map[string]any{
				"value": "43",
				"name":  "Ref_",
			}},
	}}

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(map[string]any{"Messages": messages, "Spawns": []any{}, "Outputs": []any{}, "GasUsed": 599159077})
		assert.NoError(t, err)
		_, err = w.Write(data)
		assert.NoError(t, err)
	})

	http.HandleFunc("/result/"+message, handlerFunc)

	ts := httptest.NewServer(nil)
	defer ts.Close()

	cu := NewCUMock(ts.URL)

	res, err := cu.ReadResult(process, message)
	assert.NoError(t, err)
	assert.Equal(t, messages[0]["Target"], res.Messages[0]["Target"].(string))
	assert.Equal(t, messages[0]["Anchor"], res.Messages[0]["Anchor"].(string))
	assert.Equal(t, messages[0]["Data"], res.Messages[0]["Data"].(string))
	assert.ElementsMatch(t, messages[0]["Tags"], res.Messages[0]["Tags"])
	assert.Equal(t, res.GasUsed, 599159077)
}
