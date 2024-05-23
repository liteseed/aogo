package aogo

import (
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

type Message struct {
	ID     string      `json:"Id"`
	Target string      `json:"Target"`
	From   string      `json:"From"`
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

func (ao *AO) LoadResult(process string, message string) (*LoadResultResponse, error) {
	return ao.cu.LoadResult(process, message)
}

func (ao *AO) DryRun(message Message) (*DryRunResponse, error) {
	return ao.cu.DryRun(message)
}
