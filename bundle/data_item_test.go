package bundle

import (
	"log"
	"os"
	"testing"

	"github.com/liteseed/argo/signer"
	"gotest.tools/v3/assert"
)

func TestDecodeDataItem(t *testing.T) {
	data, err := os.ReadFile("../test/stubs/1115BDataItem")
	if err != nil {
		log.Fatal(err)
	}
	dataItem, err := DecodeDataItem(data)
	assert.NilError(t, err)
	assert.Equal(t, dataItem.ID, "QpmY8mZmFEC8RxNsgbxSV6e36OF6quIYaPRKzvUco0o=")
	assert.Equal(t, dataItem.Signature, "wUIlPaBflf54QyfiCkLnQcfakgcS5B4Pld-hlOJKyALY82xpAivoc0fxBJWjoeg3zy9aXz8WwCs_0t0MaepMBz2bQljRrVXnsyWUN-CYYfKv0RRglOl-kCmTiy45Ox13LPMATeJADFqkBoQKnGhyyxW81YfuPnVlogFWSz1XHQgHxrFMAeTe9epvBK8OCnYqDjch4pwyYUFrk48JFjHM3-I2kcQnm2dAFzFTfO-nnkdQ7ulP3eoAUr-W-KAGtPfWdJKFFgWFCkr_FuNyHYQScQo-FVOwIsvj_PVWEU179NwiqfkZtnN8VoBgCSxbL1Wmh4NYL-GsRbKz_94hpcj5RiIgq0_H5dzAp-bIb49M4SP-DcuIJ5oT2v2AfPWvznokDDVTeikQJxCD2n9usBOJRpLw_P724Yurbl30eNow0U-Jmrl8S6N64cjwKVLI-hBUfcpviksKEF5_I4XCyciW0TvZj1GxK6ET9lx0s6jFMBf27-GrFx6ZDJUBncX6w8nDvuL6A8TG_ILGNQU_EDoW7iil6NcHn5w11yS_yLkqG6dw_zuC1Vkg1tbcKY3703tmbF-jMEZUvJ6oN8vRwwodinJjzGdj7bxmkUPThwVWedCc8wCR3Ak4OkIGASLMUahSiOkYmELbmwq5II-1Txp2gDPjCpAf9gT6Iu0heAaXhjk=")
	assert.Equal(t, dataItem.Owner, "0zBGbs8Y4wvdS58cAVyxp7mDffScOkbjh50ZrqnWKR_5NGwjezT6J40ejIg5cm1KnuDnw9OhvA7zO6sv1hEE6IaGNnNJWiXFecRMxCl7iw78frrT8xJvhBgtD4fBCV7eIvydqLoMl8K47sacTUxEGseaLfUdYVJ5CSock5SktEEdqqoe3MAso7x4ZsB5CGrbumNcCTifr2mMsrBytocSoHuiCEi7-Nwv4CqzB6oqymBtEECmKYWdINnNQHVyKK1l0XP1hzByHv_WmhouTPos9Y77sgewZrvLF-dGPNWSc6LaYGy5IphCnq9ACFrEbwkiCRgZHnKsRFH0dfGaCgGb3GZE-uspmICJokJ9CwDPDJoxkCBEF0tcLSIA9_ofiJXaZXbrZzu3TUXWU3LQiTqYr4j5gj_7uTclewbyZSsY-msfbFQlaACc02nQkEkr4pMdpEOdAXjWP6qu7AJqoBPNtDPBqWbdfsLXgyK90NbYmf3x4giAmk8L9REy7SGYugG4VyqG39pNQy_hdpXdcfyE0ftCr5tSHVpMreJ0ni7v3IDCbjZFcvcHp0H6f6WPfNCoHg1BM6rHUqkXWd84gdHUzo9LTGq9-7wSBCizpcc_12_I-6yvZsROJvdfYOmjPnd5llefa_X3X1dVm5FPYFIabydGlh1Vs656rRu4dzeEQwc=")
	assert.Equal(t, dataItem.Target, "")
	assert.Equal(t, dataItem.Anchor, "")
	assert.DeepEqual(
		t,
		dataItem.Tags,
		[]Tag{
			{Name: "Content-Type", Value: "text/plain"},
			{Name: "App-Name", Value: "ArDrive-CLI"},
			{Name: "App-Version", Value: "1.21.0"},
		},
	)
	assert.Equal(t, dataItem.RawData, "AQDBQiU9oF-V_nhDJ-IKQudBx9qSBxLkHg-V36GU4krIAtjzbGkCK-hzR_EElaOh6DfPL1pfPxbAKz_S3Qxp6kwHPZtCWNGtVeezJZQ34Jhh8q_RFGCU6X6QKZOLLjk7HXcs8wBN4kAMWqQGhAqcaHLLFbzVh-4-dWWiAVZLPVcdCAfGsUwB5N716m8Erw4KdioONyHinDJhQWuTjwkWMczf4jaRxCebZ0AXMVN876eeR1Du6U_d6gBSv5b4oAa099Z0koUWBYUKSv8W43IdhBJxCj4VU7Aiy-P89VYRTXv03CKp-Rm2c3xWgGAJLFsvVaaHg1gv4axFsrP_3iGlyPlGIiCrT8fl3MCn5shvj0zhI_4Ny4gnmhPa_YB89a_OeiQMNVN6KRAnEIPaf26wE4lGkvD8_vbhi6tuXfR42jDRT4mauXxLo3rhyPApUsj6EFR9ym-KSwoQXn8jhcLJyJbRO9mPUbEroRP2XHSzqMUwF_bv4asXHpkMlQGdxfrDycO-4voDxMb8gsY1BT8QOhbuKKXo1wefnDXXJL_IuSobp3D_O4LVWSDW1twpjfvTe2ZsX6MwRlS8nqg3y9HDCh2KcmPMZ2PtvGaRQ9OHBVZ50JzzAJHcCTg6QgYBIsxRqFKI6RiYQtubCrkgj7VPGnaAM-MKkB_2BPoi7SF4BpeGOdMwRm7PGOML3UufHAFcsae5g330nDpG44edGa6p1ikf-TRsI3s0-ieNHoyIOXJtSp7g58PTobwO8zurL9YRBOiGhjZzSVolxXnETMQpe4sO_H660_MSb4QYLQ-HwQle3iL8nai6DJfCuO7GnE1MRBrHmi31HWFSeQkqHJOUpLRBHaqqHtzALKO8eGbAeQhq27pjXAk4n69pjLKwcraHEqB7oghIu_jcL-AqsweqKspgbRBApimFnSDZzUB1ciitZdFz9Ycwch7_1poaLkz6LPWO-7IHsGa7yxfnRjzVknOi2mBsuSKYQp6vQAhaxG8JIgkYGR5yrERR9HXxmgoBm9xmRPrrKZiAiaJCfQsAzwyaMZAgRBdLXC0iAPf6H4iV2mV262c7t01F1lNy0Ik6mK-I-YI_-7k3JXsG8mUrGPprH2xUJWgAnNNp0JBJK-KTHaRDnQF41j-qruwCaqATzbQzwalm3X7C14MivdDW2Jn98eIIgJpPC_URMu0hmLoBuFcqht_aTUMv4XaV3XH8hNH7Qq-bUh1aTK3idJ4u79yAwm42RXL3B6dB-n-lj3zQqB4NQTOqx1KpF1nfOIHR1M6PS0xqvfu8EgQos6XHP9dvyPusr2bETib3X2Dpoz53eZZXn2v1919XVZuRT2BSGm8nRpYdVbOueq0buHc3hEMHAAADAAAAAAAAAEIAAAAAAAAABhhDb250ZW50LVR5cGUUdGV4dC9wbGFpbhBBcHAtTmFtZRZBckRyaXZlLUNMSRZBcHAtVmVyc2lvbgwxLjIxLjAANTY3MAo=")
}

func TestNewDataItem(t *testing.T) {
	data := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{};':\",./<>?`~"
	tags := []Tag{
		{Name: "tag1", Value: "value1"},
		{Name: "tag2", Value: "value2"},
	}
	anchor := "thisSentenceIs32BytesLongTrustMe"
	target := "OXcT1sVRSA5eGwt2k6Yuz8-3e3g9WJi5uSE99CWqsBs"

	s, err := signer.New("../data/wallet.json")
	assert.NilError(t, err)
	dataItem, err := NewDataItem([]byte(data), *s, target, anchor, tags)
	assert.NilError(t, err)
	assert.Equal(t, dataItem.Target, target)
}
