package bundle

import (
	"encoding/base64"
	"encoding/binary"
	"errors"

	"github.com/liteseed/argo/signer"
)

// Create a Data Item
func NewDataItem(rawData []byte, s signer.Signer, target string, anchor string, tags []Tag) (*DataItem, error) {
	rawOwner := s.S.PubKey.N.Bytes()
	rawTarget := []byte(target)
	rawAnchor := []byte(anchor)

	tagsBytes, err := encodeTags(&tags)
	if err != nil {
		return nil, err
	}

	chunks := []interface{}{
		[]byte(base64.URLEncoding.EncodeToString([]byte("dataitem"))),
		[]byte(base64.URLEncoding.EncodeToString([]byte("1"))),
		[]byte(base64.URLEncoding.EncodeToString(rawOwner)),
		[]byte(base64.URLEncoding.EncodeToString(tagsBytes)),
		[]byte(base64.URLEncoding.EncodeToString(rawAnchor)),
		[]byte(base64.URLEncoding.EncodeToString(tagsBytes)),
		[]byte(base64.URLEncoding.EncodeToString(rawData)),
	}
	signature := DeepHash(chunks)

	rawSignature := []byte(signature[:])

	raw := make([]byte, 2)
	binary.LittleEndian.PutUint16(raw, uint16(1))
	raw = append(raw, rawSignature...)
	raw = append(raw, rawOwner...)
	raw = append(raw, rawTarget...)
	raw = append(raw, rawAnchor...)
	raw = append(raw, tagsBytes...)
	raw = append(raw, rawData...)

	return &DataItem{
		SignatureType: 1,
		Signature:     base64.URLEncoding.EncodeToString(rawSignature),
		Owner:         base64.URLEncoding.EncodeToString(rawOwner),
		Target:        base64.URLEncoding.EncodeToString(rawTarget),
		Anchor:        base64.URLEncoding.EncodeToString(rawAnchor),
		Tags:          tags,
		Data:          base64.URLEncoding.EncodeToString(rawData),
		Raw:           raw,
	}, nil
}

// Decode a DataItem from bytes
func DecodeDataItem(raw []byte) (*DataItem, error) {
	N := len(raw)
	if N < 2 {
		return nil, errors.New("binary too small")
	}

	signatureType, signatureLength, publicKeyLength, err := getSignatureMetadata(raw[:2], N)
	if err != nil {
		return nil, err
	}
	signatureStart := 2
	signatureEnd := signatureLength + signatureStart
	signature := base64.URLEncoding.EncodeToString(raw[signatureStart:signatureEnd])
	idBytes, err := hash(raw[signatureStart:signatureEnd])
	if err != nil {
		return nil, err
	}
	id := base64.URLEncoding.EncodeToString(idBytes)
	ownerStart := signatureEnd
	ownerEnd := ownerStart + publicKeyLength
	owner := base64.URLEncoding.EncodeToString(raw[ownerStart:ownerEnd])

	position := 2 + ownerEnd
	target, position := getTarget(&raw, position)
	anchor, position := getAnchor(&raw, position)
	tags, position, err := decodeTags(&raw, position)
	if err != nil {
		return nil, err
	}
	data := raw[position:]
	return &DataItem{
		ID:            id,
		SignatureType: signatureType,
		Signature:     signature,
		Owner:         owner,
		Target:        target,
		Anchor:        anchor,
		Tags:          *tags,
		Data:          base64.URLEncoding.EncodeToString(data),
		Raw:           raw,
	}, nil
}
