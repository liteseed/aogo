package aogo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/everFinance/goar/types"
)

type ICU interface {
	ReadResult(process string, message string) (*LoadResultResponse, error)
}

type CU struct {
	client *http.Client
	url    string
}

func newCU(url string) CU {
	return CU{
		client: http.DefaultClient,
		url:    url,
	}
}

type DryRunResponse struct {
	Messages []map[string]interface{} `json:"Messages"`
	Spawns   []any                    `json:"Spawns"`
	Outputs  []any                    `json:"Outputs"`
	Error    string                   `json:"Error"`
	GasUsed  int                      `json:"GasUsed"`
}

type LoadResultResponse struct {
	Messages []map[string]interface{} `json:"Messages"`
	Spawns   []any                    `json:"Spawns"`
	Outputs  []any                    `json:"Outputs"`
	Error    string                   `json:"Error"`
	GasUsed  int                      `json:"GasUsed"`
}

func (cu *CU) LoadResult(process string, message string) (*LoadResultResponse, error) {
	resp, err := cu.client.Get(cu.url + "/result/" + message + "?process-id=" + process)
	if err != nil {
		return nil, err
	}
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var readResult LoadResultResponse
	err = json.Unmarshal(res, &readResult)
	if err != nil {
		return nil, err
	}
	return &readResult, nil
}

func (cu *CU) DryRun(message Message) (*DryRunResponse, error) {
	message.Tags = append(
		message.Tags,
		[]types.Tag{{Name: "Data-Protocol", Value: "ao"}, {Name: "Type", Value: "Message"}, {Name: "Variant", Value: "ao.TN.1"}}...,
	)

	body, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", cu.url+"/dry-run?process-id="+message.Target, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")

	resp, err := cu.client.Do(req)
	if err != nil {
		return nil, err
	}
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var dryRun DryRunResponse
	err = json.Unmarshal(res, &dryRun)
	if err != nil {
		return nil, err
	}
	return &dryRun, nil
}
