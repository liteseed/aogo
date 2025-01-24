package aogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/liteseed/goar/tag"
)

type ICU interface {
	LoadResult(process string, message string) (*Response, error)
	DryRun(message Message) (*Response, error)
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

type Response struct {
	Messages []map[string]any `json:"Messages"`
	Spawns   []any            `json:"Spawns"`
	Outputs  []any            `json:"Outputs"`
	Error    string           `json:"Error"`
	GasUsed  int              `json:"GasUsed"`
}

func (cu *CU) LoadResult(process string, message string) (*Response, error) {
	resp, err := cu.client.Get(fmt.Sprintf("%s/result/%s?process-id=%s", cu.url, message, process))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("cu request failed with status: %s, code: %d, server: %s", resp.Status, resp.StatusCode, resp.Request.Host)
	}
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var readResult Response
	err = json.Unmarshal(res, &readResult)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	return &readResult, nil
}

func (cu *CU) DryRun(message Message) (*Response, error) {
	if message.Tags == nil {
		message.Tags = &[]tag.Tag{}
	}
	*message.Tags = append(*message.Tags,
		tag.Tag{Name: "Data-Protocol", Value: "ao"},
		tag.Tag{Name: "Type", Value: "Message"},
		tag.Tag{Name: "Variant", Value: "ao.TN.1"},
	)
	if message.Data == "" {
		message.Data = "1984"
	}
	body, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/dry-run?process-id=%s", cu.url, message.Target), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	resp, err := cu.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("dry-run request failed with status: %s, code: %d, server: %s", resp.Status, resp.StatusCode, req.URL.Host)
	}
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var dryRun Response
	err = json.Unmarshal(res, &dryRun)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dry-run response: %v", err)
	}
	return &dryRun, nil
}
