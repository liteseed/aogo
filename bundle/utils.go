package bundle

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"

	"github.com/linkedin/goavro/v2"
)

const avroTagSchema = `
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
}`

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
		tagsEnd = bytesDataEnd
		return tags, tagsEnd, nil
	}
	return tags, tagsEnd, nil
}

func decodeAvro(data []byte) (*[]Tag, error) {
	codec, err := goavro.NewCodec(avroTagSchema)
	if err != nil {
		return nil, err
	}

	avroTags, _, err := codec.NativeFromBinary(data)
	if err != nil {
		return nil, err
	}

	tags := &[]Tag{}

	for _, v := range avroTags.([]interface{}) {
		tag := v.(map[string]any)
		*tags = append(*tags, Tag{Name: string(tag["name"].([]byte)), Value: string(tag["value"].([]byte))})
	}
	return tags, err
}

func encodeTags(tags *[]Tag) ([]byte, error) {
	data, err := encodeAvro(tags)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func encodeAvro(tags *[]Tag) ([]byte, error) {
	codec, err := goavro.NewCodec(avroTagSchema)
	if err != nil {
		return nil, err
	}

	avroTags := []map[string]any{}

	for _, tag := range *tags {
		m := map[string]any{"name": []byte(tag.Name), "value": []byte(tag.Value)}
		avroTags = append(avroTags, m)
	}
	data, err := codec.BinaryFromNative(nil, avroTags)
	if err != nil {
		return nil, err
	}

	return data, err
}

type Header struct {
	id   int
	size int
	raw  []byte
}

func generateBundleHeader(d *[]DataItem) (*[]Header, error) {
	headers := []Header{}

	for _, dataItem := range *d {
		idBytes, err := base64.URLEncoding.DecodeString(dataItem.ID)
		if err != nil {
			return nil, err
		}
		rawData, err := base64.URLEncoding.DecodeString(dataItem.Data)
		if err != nil {
			return nil, err
		}
		id := int(binary.LittleEndian.Uint16(idBytes))
		size := len(rawData)
		raw := make([]byte, 64)
		binary.LittleEndian.PutUint16(raw, uint16(size))
		binary.LittleEndian.AppendUint16(raw, uint16(id))
		headers = append(headers, Header{id: id, size: size, raw: raw})
	}
	return &headers, nil
}

func decodeBundleHeader(data *[]byte) (*[]Header, int) {
	N := int(binary.LittleEndian.Uint32((*data)[:32]))
	headers := []Header{}
	for i := 32; i < 32+64*N; i += 64 {
		size := int(binary.LittleEndian.Uint16((*data)[i : i+32]))
		id := int(binary.LittleEndian.Uint16((*data)[i+32 : i+64]))
		headers = append(headers, Header{id: id, size: size})
	}
	return &headers, N
}

func hash(data []byte) ([]byte, error) {
	h := sha256.New()
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}
	r := h.Sum(nil)
	return r, nil
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}
func DeepHash(data any) [48]byte {
	if typeof(data) == "[]uint8" {
		tag := append([]byte("blob"), []byte(fmt.Sprintf("%d", len(data.([]byte))))...)
		tagHashed := sha512.Sum384(tag)
		return tagHashed
	} else {
		tag := append([]byte("list"), []byte(fmt.Sprintf("%d", len(data.([]interface{}))))...)
		tagHashed := sha512.Sum384(tag)
		return deepHashChunk(data.([]interface{}), tagHashed)
	}
}
func deepHashChunk(data []interface{}, acc [48]byte) [48]byte {
	if len(data) < 1 {
		return acc
	}
	var dHash [48]byte
	if typeof(data[0]) == "[]uint8" {
		dHash = DeepHash(data[0].([]byte))
	} else {
		dHash = DeepHash(data[0].([]interface{}))

	}
	hashPair := append(acc[:], dHash[:]...)
	newAcc := sha512.Sum384(hashPair)
	return deepHashChunk(data[1:], newAcc)
}
