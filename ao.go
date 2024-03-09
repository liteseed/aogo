package argo

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const (
	MU_URL    = "https://mu.ao-testnet.xyz"
	CU_URL    = "https://cu.ao-testnet.xyz"
	SU_URL    = "https://g8way.io/1SafZGlZT4TLI8xoc0QEQ4MylHhuyQUblxD8xLKvEKI"
	GATEWAY   = "https://arweave.net"
	AO_MODULE = ""
)

type ReadResultResponse struct {
	Messages []map[string]interface{} `json:"Messages"`
	Spawns   []any                    `json:"Spawns"`
	Outputs  []any                    `json:"Outputs"`
	Errors   any                      `json:"Errors"`
	GasUsed  int                      `json:"GasUsed"`
}

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

func SpawnProcess() {

}

type IMU interface {
	SendMessage(process string, data string, tags []Tag, anchor string, s *Signer) (string, error)
}
type MU struct {
	client *http.Client
	url    string
}

func NewMU(URL string) MU {
	return MU{
		client: http.DefaultClient,
		url:    URL,
	}
}

func (mu MU) SendMessage(process string, data string, tags []Tag, anchor string, s *Signer) (string, error) {
	tags = append(tags, Tag{Name: "Data-Protocol", Value: "ao"})
	tags = append(tags, Tag{Name: "Variant", Value: "ao.TN.1"})
	tags = append(tags, Tag{Name: "Type", Value: "Message"})
	tags = append(tags, Tag{Name: "SDK", Value: "argo"})
	dataItem, err := NewDataItem([]byte(data), *s, process, anchor, tags)
	if err != nil {
		return "", err
	}
	resp, err := mu.client.Post(mu.url, "application/octet-stream", bytes.NewBuffer(dataItem.Raw))
	if err != nil {
		return "", err
	}
	log.Println(resp.Status)
	log.Println(err)
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return dataItem.ID, nil
}
