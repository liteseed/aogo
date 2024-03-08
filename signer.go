package argo

import "github.com/everFinance/goar"

type Signer struct {
	S *goar.Signer
}

func NewSigner(path string) (*Signer, error) {
	s := &Signer{}
	goarSigner, err := goar.NewSignerFromPath(path)
	if err != nil {
		return nil, err
	}
	s.S = goarSigner
	return s, nil
}
