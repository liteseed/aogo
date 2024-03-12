package ao

import (
	"errors"

	"github.com/liteseed/argo/signer"
	"github.com/liteseed/argo/transaction"
)

const (
	MU_URL    = "https://mu.ao-testnet.xyz"
	CU_URL    = "https://cu.ao-testnet.xyz"
	SCHEDULER = "_GQ33BkPtZrqxA84vM8Zk-N2aO0toNNu_C-l-rawrBA"
	GATEWAY   = "https://arweave.net"
	MODULE    = "9afQ1PLf2mrshqCTZEzzJTR2gWaC9zNPnYgYEqg1Pt4"

	SDK = "argo"
)

type AO struct {
	mu MU
	cu CU
}

func New() *AO {
	return &AO{
		mu: NewMU(),
		cu: NewCU(),
	}
}

// MU Functions

func (ao *AO) SpawnProcess(data string, tags []transaction.Tag, s *signer.Signer) (string, error) {
	return ao.mu.SpawnProcess(data, tags, s)
}

func (ao *AO) SendMessage(process string, data string, tags []transaction.Tag, anchor string, s *signer.Signer) (string, error) {
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
