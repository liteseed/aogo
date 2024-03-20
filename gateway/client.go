package gateway

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	Data "github.com/liteseed/argo/data"
)

type Gateway struct {
	client *http.Client
	url    string
}

func New(url string) *Gateway {
	httpClient := http.DefaultClient
	return &Gateway{client: httpClient, url: url}
}

func (g *Gateway) GetTransactionStatus(id string) (*TransactionStatus, error) {
	resp, err := g.client.Get(fmt.Sprintf("%s/tx/%s/status", g.url, id))
	if err != nil {
		return nil, err
	}

	status := &TransactionStatus{}
	if resp.StatusCode != 200 {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, status)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (g *Gateway) GetTransactionAnchor() (string, error) {
	resp, err := g.client.Get(fmt.Sprintf("%s/tx_anchor", g.url))
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (g *Gateway) SendTransaction(t *Data.Transaction) error {
	resp, err := g.client.Post(fmt.Sprintf("%s/tx", g.url), "application/json", bytes.NewBuffer(t.Raw))
	if err != nil {
		return err
	}
	if resp.StatusCode == 208 {
		return errors.New("transaction already posted")
	}
	if resp.StatusCode != 200 {
		return errors.New("something went wrong")
	}
	return nil
}
