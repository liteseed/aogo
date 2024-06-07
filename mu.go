package aogo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/tx"
	"github.com/liteseed/goar/types"
)

type IMU interface {
	SendMessage(process string, data string, tags []types.Tag, s *signer.Signer) (string, error)
	SpawnProcess(data string, tags []types.Tag, s *signer.Signer) (string, error)

	Monitor()
}
type MU struct {
	client *http.Client
	url    string
}

func newMU(url string) MU {
	return MU{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		url: url,
	}
}

type SendMessageResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

type SpawnProcessResponse struct {
	ID string `json:"id"`
}

func (mu MU) SendMessage(process string, data string, tags []types.Tag, anchor string, s *signer.Signer) (string, error) {
	tags = append(tags, types.Tag{Name: "Data-Protocol", Value: "ao"})
	tags = append(tags, types.Tag{Name: "Variant", Value: "ao.TN.1"})
	tags = append(tags, types.Tag{Name: "Type", Value: "Message"})
	tags = append(tags, types.Tag{Name: "SDK", Value: SDK})

	dataItem, err := tx.NewDataItem([]byte(data), process, anchor, tags)
	if err != nil {
		return "", err
	}

	err = s.SignDataItem(dataItem)
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

	if resp.StatusCode >= 400 {
		return "", errors.New("message failed: " + resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var res SendMessageResponse
	err = json.Unmarshal(b, &res)
	if err != nil {
		return "", err
	}

	return res.ID, nil
}

func (mu MU) SpawnProcess(module string, data string, tags []types.Tag, s *signer.Signer) (string, error) {
	if data == "" {
		data = "1984"
	}
	tags = append(tags, types.Tag{Name: "Data-Protocol", Value: "ao"})
	tags = append(tags, types.Tag{Name: "Variant", Value: "ao.TN.1"})
	tags = append(tags, types.Tag{Name: "Type", Value: "Process"})
	tags = append(tags, types.Tag{Name: "Scheduler", Value: SCHEDULER})
	tags = append(tags, types.Tag{Name: "Module", Value: module})
	tags = append(tags, types.Tag{Name: "SDK", Value: SDK})

	dataItem, err := tx.NewDataItem([]byte(data), "", "", tags)
	if err != nil {
		return "", err
	}

	err = s.SignDataItem(dataItem)
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
	if resp.StatusCode >= 400 {
		return "", errors.New(resp.Status)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var res SpawnProcessResponse
	err = json.Unmarshal(b, &res)
	if err != nil {
		return "", err
	}

	return res.ID, nil
}
