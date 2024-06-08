package aogo

import (
	"github.com/liteseed/goar/signer"
	"github.com/liteseed/goar/types"
)

const (
	MU_URL    = "https://mu.ao-testnet.xyz"
	CU_URL    = "https://cu.ao-testnet.xyz"
	SCHEDULER = "_GQ33BkPtZrqxA84vM8Zk-N2aO0toNNu_C-l-rawrBA"
	GATEWAY   = "https://arweave.net"

	SDK = "aogo"
)

type AO struct {
	mu IMU
	cu ICU
}

type Message struct {
	ID     string      `json:"Id"`
	Target string      `json:"Target"`
	Owner  string      `json:"Owner"`
	Data   any         `json:"Data"`
	Tags   []types.Tag `json:"Tags"`
}

func New(options ...func(*AO)) (*AO, error) {
	ao := &AO{cu: newCU(CU_URL), mu: newMU(MU_URL)}
	for _, o := range options {
		o(ao)
	}
	return ao, nil
}

func WithMU(url string) func(*AO) {
	return func(ao *AO) {
		ao.mu = newMU(url)
	}
}

func WithCU(url string) func(*AO) {
	return func(ao *AO) {
		ao.cu = newCU(url)
	}
}

// MU Functions

func (ao *AO) SpawnProcess(module string, data string, tags []types.Tag, s *signer.Signer) (string, error) {
	return ao.mu.SpawnProcess(module, data, tags, s)
}

func (ao *AO) SendMessage(process string, data string, tags []types.Tag, anchor string, s *signer.Signer) (string, error) {
	return ao.mu.SendMessage(process, data, tags, anchor, s)
}

// CU Functions

func (ao *AO) LoadResult(process string, message string) (*Response, error) {
	return ao.cu.LoadResult(process, message)
}

func (ao *AO) DryRun(message Message) (*Response, error) {
	return ao.cu.DryRun(message)
}
