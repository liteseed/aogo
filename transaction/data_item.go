package transaction

import (
	"encoding/base64"
	"encoding/binary"
	"errors"

	"github.com/liteseed/argo/signer"
)

func NewDataItem(rawData []byte, s signer.Signer, target string, anchor string, tags []Tag) (*DataItem, error) {

	rawOwner := []byte(s.S.PubKey.N.Bytes())
	rawTarget, err := base64.RawURLEncoding.DecodeString(target)
	if err != nil {
		return nil, err
	}
	rawAnchor := []byte(anchor)

	tagsBytes, err := encodeTags(&tags)
	if err != nil {
		return nil, err
	}

	chunks := []interface{}{
		[]byte("dataitem"),
		[]byte("1"),
		[]byte("1"),
		rawOwner,
		rawTarget,
		rawAnchor,
		tagsBytes,
		rawData,
	}
	signatureData := DeepHash(chunks)

	rawSignature, err := s.S.SignMsg(signatureData[:])
	if err != nil {
		return nil, err
	}
	raw := make([]byte, 0)
	raw = binary.LittleEndian.AppendUint16(raw, uint16(1))
	raw = append(raw, rawSignature...)
	raw = append(raw, rawOwner...)

	if target == "" {
		raw = append(raw, 0)
	} else {
		raw = append(raw, 1)
	}
	raw = append(raw, rawTarget...)

	if anchor == "" {
		raw = append(raw, 0)
	} else {
		raw = append(raw, 1)
	}
	raw = append(raw, rawAnchor...)

	numberOfTags := make([]byte, 8)
	binary.LittleEndian.PutUint16(numberOfTags, uint16(len(tags)))
	raw = append(raw, numberOfTags...)

	tagsLength := make([]byte, 8)
	binary.LittleEndian.PutUint16(tagsLength, uint16(len(tagsBytes)))
	raw = append(raw, tagsLength...)
	raw = append(raw, tagsBytes...)

	raw = append(raw, rawData...)

	rawID, err := hash(rawSignature)
	if err != nil {
		return nil, err
	}
	return &DataItem{
		SignatureType: 1,
		Signature:     base64.RawURLEncoding.EncodeToString(rawSignature),
		ID:            base64.RawURLEncoding.EncodeToString(rawID),
		Owner:         s.S.Owner(),
		Target:        target,
		Anchor:        anchor,
		Tags:          tags,
		Data:          string(rawData),
		Raw:           raw,
	}, nil
}

// Decode a DataItem from bytes
func DecodeDataItem(raw []byte) (*DataItem, error) {
	N := len(raw)
	if N < 2 {
		return nil, errors.New("binary too small")
	}

	signatureType, signatureLength, publicKeyLength, err := getSignatureMetadata(raw[:2])
	if err != nil {
		return nil, err
	}

	signatureStart := 2
	signatureEnd := signatureLength + signatureStart
	signature := base64.RawURLEncoding.EncodeToString(raw[signatureStart:signatureEnd])
	rawId, err := hash(raw[signatureStart:signatureEnd])
	if err != nil {
		return nil, err
	}
	id := base64.RawURLEncoding.EncodeToString(rawId)
	ownerStart := signatureEnd
	ownerEnd := ownerStart + publicKeyLength
	owner := base64.RawURLEncoding.EncodeToString(raw[ownerStart:ownerEnd])

	position := ownerEnd
	target, position := getTarget(&raw, position)
	anchor, position := getAnchor(&raw, position)
	tags, position, err := decodeTags(&raw, position)
	if err != nil {
		return nil, err
	}
	data := base64.RawURLEncoding.EncodeToString(raw[position:])

	return &DataItem{
		ID:            id,
		SignatureType: signatureType,
		Signature:     signature,
		Owner:         owner,
		Target:        target,
		Anchor:        anchor,
		Tags:          *tags,
		Data:          data,
		Raw:           raw,
	}, nil
}
