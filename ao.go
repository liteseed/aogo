package argo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const (
	MU_URL    = "https://mu.ao-testnet.xyz"
	CU_URL    = "https://cu.ao-testnet.xyz"
	GATEWAY   = "https://arweave.net"
	AO_MODULE = ""
)

func SendMessage(process string, data string, tags []Tag, anchor string, s *Signer) (string, error) {
	dataItem, err := NewDataItem([]byte(data), *s, process, anchor, tags)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(MU_URL+"/monitor/"+process, "application/octet-stream", bytes.NewBuffer(dataItem.Raw))
	if err != nil {
		return "", err
	}
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var messageId string
	err = json.Unmarshal(res, &messageId)
	if err != nil {
		return "", err
	}
	return messageId, nil
}

type ReadResultResponse struct {
	Messages []any
	Spawns   []any
	Outputs  []any
	Errors   any
	GasUsed  int
}

func ReadResult(process string, message string) (*ReadResultResponse, error) {
	resp, err := http.Get(MU_URL + "/result/" + message + "?process-id=" + process)
	if err != nil {
		return nil, err
	}
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var readResult *ReadResultResponse
	err = json.Unmarshal(res, readResult)
	if err != nil {
		return nil, err
	}
	return readResult, nil
}

func SpawnProcess() {

}
