package bundle

import (
	"encoding/base64"
	"errors"
)


func NewDataItem(data []byte) (*DataItem, error) {
	dataItem := &DataItem{}
	return dataItem, nil
}

// Decode a DataItem from bytes
func DecodeDataItem(data []byte) (*DataItem, error) {
	N := len(data)
	if N < 2 {
		return nil, errors.New("binary too small")
	}

	signatureType, signatureLength, publicKeyLength, err := getSignatureMetadata(data[:2], N)
	if err != nil {
		return nil, err
	}
	signatureStart := 2
	signatureEnd := signatureLength + signatureStart
	signature := base64.URLEncoding.EncodeToString(data[signatureStart:signatureEnd])

	ownerStart := signatureEnd
	ownerEnd := ownerStart + publicKeyLength
	owner := base64.URLEncoding.EncodeToString(data[ownerStart:ownerEnd])

	position := 2 + ownerEnd
	target, position := getTarget(&data, position)
	anchor, position := getAnchor(&data, position)
	tags, position, err := decodeTags(&data, position)
	if err != nil {
		return nil, err
	}

	rawData := data[position:]

	return &DataItem{
		SignatureType: signatureType,
		Signature:     signature,
		Owner:         owner,
		Target:        target,
		Anchor:        anchor,
		Tags:          *tags,
		RawData:       base64.URLEncoding.EncodeToString(rawData),
	}, nil
}

func (d *DataItem) Sign() error {

	return nil
}
