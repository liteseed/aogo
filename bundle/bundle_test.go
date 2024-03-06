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
