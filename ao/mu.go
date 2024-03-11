package ao

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/liteseed/argo/signer"
	"github.com/liteseed/argo/transaction"
)

type IMU interface {
	SendMessage(process string, data string, tags []transaction.Tag, s *signer.Signer) (string, error)
	SpawnProcess(data string, tags []transaction.Tag, scheduler string, s *signer.Signer) (string, error)
}
type MU struct {
	client *http.Client
	url    string
}

func NewMU() MU {
	return MU{
		client: http.DefaultClient,
		url:    MU_URL,
	}
}

func (mu MU) SendMessage(process string, data string, tags []transaction.Tag, anchor string, s *signer.Signer) (string, error) {
	tags = append(tags, transaction.Tag{Name: "Data-Protocol", Value: "ao"})
	tags = append(tags, transaction.Tag{Name: "Variant", Value: "ao.TN.1"})
	tags = append(tags, transaction.Tag{Name: "Type", Value: "Message"})
	tags = append(tags, transaction.Tag{Name: "SDK", Value: sdk})
	dataItem, err := transaction.NewDataItem([]byte(data), *s, process, anchor, tags)
	if err != nil {
		return "", err
	}
	resp, err := mu.client.Post(mu.url, "application/octet-stream", bytes.NewBuffer(dataItem.Raw))
	if err != nil {
		return "", err
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return dataItem.ID, nil
}

func (mu MU) SpawnProcess(data string, tags []transaction.Tag, s *signer.Signer) (string, error) {
	tags = append(tags, transaction.Tag{Name: "Data-Protocol", Value: "ao"})
	tags = append(tags, transaction.Tag{Name: "Variant", Value: "ao.TN.1"})
	tags = append(tags, transaction.Tag{Name: "Type", Value: "Process"})
	tags = append(tags, transaction.Tag{Name: "Scheduler", Value: SCHEDULER})
	tags = append(tags, transaction.Tag{Name: "Module", Value: MODULE})
	tags = append(tags, transaction.Tag{Name: "SDK", Value: sdk})
	dataItem, err := transaction.NewDataItem([]byte(data), *s, "", "", tags)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", mu.url, bytes.NewBuffer(dataItem.Raw))
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

	res, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}
	log.Print(string(res))
	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}
	return dataItem.ID, nil
}
