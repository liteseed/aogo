package bundle

import (
	"log"
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestDecodeBundleHeader(t *testing.T) {
	data, err := os.ReadFile("../test/stubs/bundleHeader")
	if err != nil {
		log.Fatal(err)
	}
	headers, N := decodeBundleHeader(&data)
	assert.Equal(t, N, 1)
	assert.Equal(t, (*headers)[0].size, 1115)
	assert.Equal(t, (*headers)[0].id, 39234)
}

func TestGenerateBundleHeader(t *testing.T) {
	data, err := os.ReadFile("../test/stubs/1115BDataItem")
	if err != nil {
		log.Fatal(err)
	}
	println(data)
	dataItem, err := DecodeDataItem(data)
	assert.NilError(t, err)
	headers, err := generateBundleHeader(&[]DataItem{*dataItem})

	assert.NilError(t, err)
	assert.Equal(t, (*headers)[0].size, 1115)
	assert.Equal(t, (*headers)[0].id, 39234)
}
