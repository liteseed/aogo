package transaction

import (
	"github.com/liteseed/argo/signer"
)

func New(data []byte, anchor string, tags []Tag, winston string, s *signer.Signer) (*Transaction, error) {
	transaction := &Transaction{
		Format: 2,
		LastTx: anchor,
		Tags:   tags,
		Owner:  s.S.Owner(),
	}

	return transaction, nil
}
