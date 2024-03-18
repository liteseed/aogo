package transaction

import (
	"crypto/sha256"

	"github.com/everFinance/goar/utils"
)

func hash(data []byte) ([]byte, error) {
	h := sha256.New()
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}
	r := h.Sum(nil)
	return r, nil
}

func verify(data []byte, signature []byte, owner string) (bool, error) {
	publicKey, err := utils.OwnerToPubKey(owner)
	if err != nil {
		return false, err
	}
	err = utils.Verify(data, publicKey, signature)
	if err != nil {
		return false, err
	}
	return true, err
}
