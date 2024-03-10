package ao

import (
	"github.com/liteseed/argo/signer"
	"github.com/liteseed/argo/transaction"
)

const (
	MU_URL    = "https://mu.ao-testnet.xyz"
	CU_URL    = "https://cu.ao-testnet.xyz"
	SU_URL    = "https://g8way.io/1SafZGlZT4TLI8xoc0QEQ4MylHhuyQUblxD8xLKvEKI"
	GATEWAY   = "https://arweave.net"
	AO_MODULE = ""
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

func (ao *AO) SpawnProcess(data string, tags []transaction.Tag, s *signer.Signer) (string, error) {
	return ao.mu.SpawnProcess(data, tags, s)
}

func (ao *AO) SendMessage(process string, data string, tags []transaction.Tag, anchor string, s *signer.Signer) (string, error) {
	return ao.mu.SendMessage(process, data, tags, anchor, s)
}

func (ao *AO) ReadResult(process string, message string) (*ReadResultResponse, error) {
	return ao.cu.ReadResult(process, message)
}
