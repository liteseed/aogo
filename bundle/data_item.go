package bundle

import (
	"encoding/base64"
	"errors"
	"log"

	"github.com/liteseed/argo/signer"
)

// Create a Data Item
func NewDataItem(data []byte, s signer.Signer, target string, anchor string, tags []Tag) (*DataItem, error) {

	targetBytes := []byte(target)
	anchorBytes := []byte(anchor)
	tagsBytes, err := encodeTags(&tags)

	if err != nil {
		return nil, err
	}

	print(targetBytes, anchorBytes, tagsBytes)
	rawData := []byte{uint8(1)}
	log.Println(rawData)

	return &DataItem{
		SignatureType: 1,
		Owner:         s.S.Owner(),
		Target:        target,
		Anchor:        anchor,
		Tags:          tags,
		RawData:       base64.URLEncoding.EncodeToString(rawData),
	}, nil
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
	idBytes, err := hash(data[signatureStart:signatureEnd])
	if err != nil {
		return nil, err
	}
	id := base64.URLEncoding.EncodeToString(idBytes)
	ownerStart := signatureEnd
	ownerEnd := ownerStart + publicKeyLength
	owner := base64.URLEncoding.EncodeToString(data[ownerStart:ownerEnd])

	position := 2 + ownerEnd
	target, position := getTarget(&data, position)
	anchor, position := getAnchor(&data, position)
	tags, _, err := decodeTags(&data, position)
	if err != nil {
		return nil, err
	}

	return &DataItem{
		ID:            id,
		SignatureType: signatureType,
		Signature:     signature,
		Owner:         owner,
		Target:        target,
		Anchor:        anchor,
		Tags:          *tags,
		RawData:       base64.URLEncoding.EncodeToString(data),
	}, nil
}
