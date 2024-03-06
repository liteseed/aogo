package bundle

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/hamba/avro"
)

func getSignatureMetadata(data []byte, N int) (SignatureType int, SignatureLength int, PublicKeyLength int, err error) {
	SignatureType = int(binary.LittleEndian.Uint16(data))
	signatureMeta, ok := SignatureConfig[SignatureType]
	if !ok {
		return -1, -1, -1, fmt.Errorf("unsupported signature type:%d", SignatureType)
	}
	if N < signatureMeta.SignatureLength+2 {
		return -1, -1, -1, errors.New("dataItem longer than expected signature length")
	}
	SignatureLength = signatureMeta.SignatureLength
	PublicKeyLength = signatureMeta.PublicKeyLength
	err = nil
	return
}

func getTarget(data *[]byte, startAt int) (string, int) {
	target := ""
	position := startAt
	if (*data)[startAt] == 1 {
		target = base64.URLEncoding.EncodeToString((*data)[startAt+1 : startAt+1+32])
		position += 32
	}
	return target, position
}

func getAnchor(data *[]byte, startAt int) (string, int) {
	anchor := ""
	position := startAt
	if (*data)[startAt] == 1 {
		anchor = base64.URLEncoding.EncodeToString((*data)[position+1 : position+1+32])
		position += 32
	}
	return anchor, position
}

func decodeTags(data *[]byte, startAt int) (*[]Tag, int, error) {
	tags := &[]Tag{}
	tagsEnd := startAt + 8

	numberOfTags := int(binary.LittleEndian.Uint16((*data)[startAt : startAt+8]))

	if numberOfTags > 0 {

		numberOfTagBytesStart := startAt + 8
		numberOfTagBytesEnd := numberOfTagBytesStart + 8
		numberOfTagBytes := int(binary.LittleEndian.Uint16((*data)[numberOfTagBytesStart:numberOfTagBytesEnd]))

		bytesDataStart := numberOfTagBytesEnd
		bytesDataEnd := numberOfTagBytesEnd + numberOfTagBytes
		bytesData := (*data)[bytesDataStart:bytesDataEnd]

		tags, err := decodeAvro(bytesData)
		if err != nil {
			return nil, tagsEnd, err
		}
		println(tags)
		tagsEnd = bytesDataEnd
		return tags, tagsEnd, nil
	}
	return tags, tagsEnd, nil
}

func decodeAvro(data []byte) (*[]Tag, error) {
	codec, err := avro.Parse(`
	{
		"type": "array",
		"items": {
			"type": "record",
			"name": "Tag",
			"fields": [
				{ "name": "name", "type": "bytes" },
				{ "name": "value", "type": "bytes" }
			]
		}
	}`)
	if err != nil {
		panic(err)
	}
	avroTags := &[]map[string]any{}
	err = avro.Unmarshal(codec, data, avroTags)
	if err != nil {
		return nil, err
	}

	tags := []Tag{}
	for _, v := range *avroTags {
		tags = append(tags, Tag{Name: string(v["name"].([]byte)), Value: string(v["value"].([]byte))})
	}
	return &tags, err
}

type bundleHeader struct {
	id   int
	size int
}

func decodeBundleHeader(data *[]byte) (*[]bundleHeader, int) {
	N := int(binary.LittleEndian.Uint16((*data)[:32]))
	headers := []bundleHeader{}
	for i := 32; i < 32+64*N; i += 64 {
		size := int(binary.LittleEndian.Uint16((*data)[i : i+32]))
		id := int(binary.LittleEndian.Uint16((*data)[i+32 : i+64]))
		headers = append(headers, bundleHeader{id: id, size: size})
	}
	return &headers, N
}
