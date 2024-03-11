package ao

import (
	"bytes"
	"io"
	"net/http"

	"github.com/liteseed/argo/signer"
	"github.com/liteseed/argo/transaction"
)

type IMU interface {
	SendMessage(process string, data string, tags []transaction.Tag, s *signer.Signer) (string, error)
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

func (mu MU) SpawnProcess(data string, tags []transaction.Tag, anchor string, scheduler string, s *signer.Signer) (string, error) {
	tags = append(tags, transaction.Tag{Name: "Data-Protocol", Value: "ao"})
	tags = append(tags, transaction.Tag{Name: "Variant", Value: "ao.TN.1"})
	tags = append(tags, transaction.Tag{Name: "Type", Value: "Message"})
	tags = append(tags, transaction.Tag{Name: "SDK", Value: sdk})
	tags = append(tags, transaction.Tag{Name: "Scheduler", Value: scheduler})
	dataItem, err := transaction.NewDataItem([]byte(data), *s, "", anchor, tags)
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
