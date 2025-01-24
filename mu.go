package aogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/tag"
	"github.com/liteseed/goar/transaction/data_item"
)

type IMU interface {
	SendMessage(process string, data string, tags []tag.Tag, s *signer.Signer) (string, error)
	SpawnProcess(data string, tags []tag.Tag, s *signer.Signer) (string, error)

	Monitor()
}
type MU struct {
	client *http.Client
	url    string
}

func newMU(url string) MU {
	return MU{
		client: http.DefaultClient,
		url:    url,
	}
}

type SendMessageResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

type SpawnProcessResponse struct {
	ID string `json:"id"`
}

func (mu *MU) SendMessage(process string, data string, tags *[]tag.Tag, anchor string, s *signer.Signer) (string, error) {
	if tags == nil {
		tags = &[]tag.Tag{}
	}
	*tags = append(*tags, tag.Tag{Name: "Data-Protocol", Value: "ao"},
		tag.Tag{Name: "Variant", Value: "ao.TN.1"},
		tag.Tag{Name: "Type", Value: "Message"},
		tag.Tag{Name: "SDK", Value: SDK})

	dataItem := data_item.New([]byte(data), process, anchor, tags)
	err := dataItem.Sign(s)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", mu.url, bytes.NewBuffer(dataItem.Raw))
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "application/octet-stream")
	req.Header.Set("accept", "application/json")

	resp, err := mu.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return "", fmt.Errorf("message failed: %s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var res SendMessageResponse
	err = json.Unmarshal(b, &res)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return res.ID, nil
}

func (mu *MU) SpawnProcess(module string, data []byte, tags []tag.Tag, s *signer.Signer) (string, error) {
	if data == nil {
		data = []byte("1984")
	}

	// Initialize newTags with the base tags
	newTags := []tag.Tag{
		{Name: "Data-Protocol", Value: "ao"},
		{Name: "Variant", Value: "ao.TN.1"},
		{Name: "Type", Value: "Process"},
		{Name: "Scheduler", Value: SCHEDULER},
		{Name: "Module", Value: module},
		{Name: "SDK", Value: SDK},
	}

	newTags = append(newTags, tags...)

	dataItem := data_item.New(data, "", "", &newTags)
	err := dataItem.Sign(s)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", mu.url, bytes.NewBuffer(dataItem.Raw))
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "application/octet-stream")
	req.Header.Set("accept", "application/json")

	resp, err := mu.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return "", fmt.Errorf("request failed: %s", resp.Status)
	}
	var res SpawnProcessResponse
	err = json.Unmarshal(b, &res)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return res.ID, nil
}
