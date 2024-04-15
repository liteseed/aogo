package aogo

import (
	"errors"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
)

const (
	MU_URL    = "https://mu.ao-testnet.xyz"
	CU_URL    = "https://cu.ao-testnet.xyz"
	SCHEDULER = "_GQ33BkPtZrqxA84vM8Zk-N2aO0toNNu_C-l-rawrBA"
	GATEWAY   = "https://arweave.net"

	SDK = "aogo"
)

type AO struct {
	mu MU
	cu CU
}

func New(options ...func(*AO)) (*AO, error) {
	ao := &AO{cu: newCU(CU_URL), mu: newMU(MU_URL)}
	for _, o := range options {
		o(ao)
	}
	return ao, nil
}

func WthMU(url string) func(*AO) {
	return func(ao *AO) {
		ao.mu = newMU(url)
	}
}

func WthCU(url string) func(*AO) {
	return func(ao *AO) {
		ao.cu = newCU(url)
	}
}

// MU Functions

func (ao *AO) SpawnProcess(module string, data string, tags []types.Tag, s *goar.ItemSigner) (string, error) {
	return ao.mu.SpawnProcess(module, data, tags, s)
}

func (ao *AO) SendMessage(process string, data string, tags []types.Tag, anchor string, s *goar.ItemSigner) (string, error) {
	return ao.mu.SendMessage(process, data, tags, anchor, s)
}

// CU Functions

func (ao *AO) ReadResult(process string, message string) (*ReadResultResponse, error) {
	return ao.cu.ReadResult(process, message)
}

func (ao *AO) ReadResults(process string, message string) error {
	return errors.New("Unimplemented")
}

func (ao *AO) DryRun() error {
	return errors.New("Unimplemented")
}
