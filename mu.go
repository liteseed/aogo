package aogo

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
)

type IMU interface {
	SendMessage(process string, data string, tags []types.Tag, s *goar.ItemSigner) (string, error)
	SpawnProcess(data string, tags []types.Tag, s *goar.ItemSigner) (string, error)

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

func (mu MU) SendMessage(process string, data string, tags []types.Tag, anchor string, s *goar.ItemSigner) (string, error) {
	tags = append(tags, types.Tag{Name: "Data-Protocol", Value: "ao"})
	tags = append(tags, types.Tag{Name: "Variant", Value: "ao.TN.1"})
	tags = append(tags, types.Tag{Name: "Type", Value: "Message"})
	tags = append(tags, types.Tag{Name: "SDK", Value: SDK})
	dataItem, err := s.CreateAndSignItem([]byte(data), process, anchor, tags)
	if err != nil {
		return "", err
	}
	resp, err := mu.client.Post(mu.url, "application/octet-stream", bytes.NewBuffer(dataItem.ItemBinary))
	if err != nil {
		return "", err
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return dataItem.Id, nil
}

func (mu MU) SpawnProcess(data []byte, tags []types.Tag, s *goar.ItemSigner) (string, error) {
	tags = append(tags, types.Tag{Name: "Data-Protocol", Value: "ao"})
	tags = append(tags, types.Tag{Name: "Variant", Value: "ao.TN.1"})
	tags = append(tags, types.Tag{Name: "Type", Value: "Process"})
	tags = append(tags, types.Tag{Name: "Scheduler", Value: SCHEDULER})
	tags = append(tags, types.Tag{Name: "Module", Value: MODULE})
	tags = append(tags, types.Tag{Name: "SDK", Value: SDK})

	dataItem, err := s.CreateAndSignItem(data, "", "", tags)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", mu.url, bytes.NewBuffer(dataItem.ItemBinary))
	req.Header.Set("content-type", "application/octet-stream")
	req.Header.Set("accept", "application/json")
	if err != nil {
		return "", err
	}
	resp, err := mu.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		return "", errors.New(resp.Status)
	}
	return dataItem.Id, nil
}
