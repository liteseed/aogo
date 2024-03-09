package argo

import (
	"encoding/hex"
	"log"
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestDecodeBundleHeader(t *testing.T) {
	data, err := os.ReadFile("./test/stubs/bundleHeader")
	if err != nil {
		log.Fatal(err)
	}
	headers, N := decodeBundleHeader(&data)
	assert.Equal(t, N, 1)
	assert.Equal(t, (*headers)[0].size, 1115)
	assert.Equal(t, (*headers)[0].id, 39234)
}

func TestGenerateBundleHeader(t *testing.T) {
	data, err := os.ReadFile("./test/stubs/1115BDataItem")
	assert.NilError(t, err)

	dataItem, err := DecodeDataItem(data)
	assert.NilError(t, err)
	headers, err := generateBundleHeader(&[]DataItem{*dataItem})

	assert.NilError(t, err)
	assert.Equal(t, (*headers)[0].size, 1115)
	assert.Equal(t, (*headers)[0].id, 39234)
}

func TestEncodeAvro(t *testing.T) {
	data := []byte{6, 24, 67, 111, 110, 116, 101, 110, 116, 45, 84, 121, 112, 101, 20, 116, 101, 120, 116, 47, 112, 108, 97, 105, 110, 16, 65, 112, 112, 45, 78, 97, 109, 101, 22, 65, 114, 68, 114, 105, 118, 101, 45, 67, 76, 73, 22, 65, 112, 112, 45, 86, 101, 114, 115, 105, 111, 110, 12, 49, 46, 50, 49, 46, 48, 0}
	tags := []Tag{
		{Name: "Content-Type", Value: "text/plain"},
		{Name: "App-Name", Value: "ArDrive-CLI"},
		{Name: "App-Version", Value: "1.21.0"},
	}

	rawTags, err := encodeAvro(&tags)
	assert.NilError(t, err)
	assert.DeepEqual(t, data, rawTags)
}

func TestDecodeAvro(t *testing.T) {
	data := []byte{6, 24, 67, 111, 110, 116, 101, 110, 116, 45, 84, 121, 112, 101, 20, 116, 101, 120, 116, 47, 112, 108, 97, 105, 110, 16, 65, 112, 112, 45, 78, 97, 109, 101, 22, 65, 114, 68, 114, 105, 118, 101, 45, 67, 76, 73, 22, 65, 112, 112, 45, 86, 101, 114, 115, 105, 111, 110, 12, 49, 46, 50, 49, 46, 48, 0}
	tags, err := decodeAvro(data)
	assert.NilError(t, err)
	assert.DeepEqual(t, *tags, []Tag{
		{Name: "Content-Type", Value: "text/plain"},
		{Name: "App-Name", Value: "ArDrive-CLI"},
		{Name: "App-Version", Value: "1.21.0"},
	})
}

func TestDeepHash0(t *testing.T) {
	data := []byte{1, 2, 3}
	r := DeepHash(data)
	h, err := hex.DecodeString("41300af79285f856e833164518c7ec4974f5869ec77ca3458113fe6c587680d050f9f6864fd77f9eb62bd4e2faea9ae8")
	assert.NilError(t, err)
	assert.DeepEqual(t, r[:], h)
}

func TestDeepHash1(t *testing.T) {
	data := []byte{}
	r := DeepHash(data)
	h, err := hex.DecodeString("fbf00cc444f5fea9dc3bedf62a13fba8ae87e7445fc910567a23bec4eb82fadb1143c433069314d8362983dc3c2e4a38")
	assert.NilError(t, err)
	assert.DeepEqual(t, r[:], h)
}

func TestDeepHash2(t *testing.T) {
	data := [][]byte{{1, 2, 3}, {4, 5, 6}}
	r := DeepHash(data)
	h, err := hex.DecodeString("4dacdcc81acd09f38c77a07a2a7ae81f77c61e6b97ee5cc7b92f3a7f258e8d5ba69d14d7d66070797b083873717c9896")
	assert.NilError(t, err)
	assert.DeepEqual(t, r[:], h)
}
