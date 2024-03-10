package ao

import (
	"encoding/json"
	"io"
	"net/http"
)



type ICU interface {
	ReadResult(process string, message string) (*ReadResultResponse, error)
}

type CU struct {
	client *http.Client
	url    string
}

func NewCU(URL string) CU {
	return CU{
		client: http.DefaultClient,
		url:    URL,
	}
}


type ReadResultResponse struct {
	Messages []map[string]interface{} `json:"Messages"`
	Spawns   []any                    `json:"Spawns"`
	Outputs  []any                    `json:"Outputs"`
	Errors   any                      `json:"Errors"`
	GasUsed  int                      `json:"GasUsed"`
}


func (cu *CU) ReadResult(process string, message string) (*ReadResultResponse, error) {
	resp, err := cu.client.Get(cu.url + "/result/" + message + "?process-id=" + process)
	if err != nil {
		return nil, err
	}
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var readResult ReadResultResponse
	err = json.Unmarshal(res, &readResult)
	if err != nil {
		return nil, err
	}
	return &readResult, nil
}