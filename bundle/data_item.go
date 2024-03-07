package bundle

import (
	"encoding/base64"
	"encoding/binary"
	"errors"

	"github.com/liteseed/argo/signer"
)

// Create a Data Item
func NewDataItem(data []byte, s signer.Signer, target string, anchor string, tags []Tag) (*DataItem, error) {
	ownerBytes := s.S.PubKey.N.Bytes()
	targetBytes := []byte(target)
	anchorBytes := []byte(anchor)

	tagsBytes, err := encodeTags(&tags)
	if err != nil {
		return nil, err
	}
	chunks := []interface{}{
		[]byte(base64.URLEncoding.EncodeToString([]byte("dataitem"))),
		[]byte(base64.URLEncoding.EncodeToString([]byte("1"))),
		[]byte(base64.URLEncoding.EncodeToString(ownerBytes)),
		[]byte(base64.URLEncoding.EncodeToString(tagsBytes)),
		[]byte(base64.URLEncoding.EncodeToString(anchorBytes)),
		[]byte(base64.URLEncoding.EncodeToString(tagsBytes)),
		[]byte(base64.URLEncoding.EncodeToString(data)),
	}
	signature := DeepHash(chunks)
	signatureBytes := []byte(signature[:])
	rawData := make([]byte, 2)
	binary.LittleEndian.PutUint16(rawData, uint16(1))

	rawData = append(rawData, signatureBytes...)
	rawData = append(rawData, ownerBytes...)
	rawData = append(rawData, targetBytes...)
	rawData = append(rawData, anchorBytes...)
	rawData = append(rawData, tagsBytes...)
	rawData = append(rawData, data...)

	return &DataItem{
		SignatureType: 1,
		Owner:         s.S.Owner(),
		Target:        base64.URLEncoding.EncodeToString([]byte(target)),
		Anchor:        base64.URLEncoding.EncodeToString([]byte(anchor)),
		Tags:          tags,
		Data:          base64.URLEncoding.EncodeToString(data),
		Raw:           rawData,
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
		Data:          base64.URLEncoding.EncodeToString(data),
	}, nil
}
