package core

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sharkchain/types"
	"testing"
	"time"
)

func TestHeader_DecodeBinary(t *testing.T) {
	h := &Header{
		Version:   1,
		PrevBlock: types.RandomHash(),
		Timestamp: time.Now().UnixNano(),
		Height:    10,
		Nonce:     98319381,
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, h.EncodeBinary(buf))
	fmt.Printf("buf : %+v\n", buf.Bytes())

	hDecode := &Header{}
	assert.Nil(t, hDecode.DecodeBinary(buf))
	fmt.Printf("Header : %+v\n", hDecode)
	assert.Equal(t, h, hDecode)
}

func TestBlock_DecodeBinary(t *testing.T) {
	b := &Block{
		Header: Header{
			Version:   1,
			PrevBlock: types.RandomHash(),
			Timestamp: time.Now().UnixNano(),
			Height:    10,
			Nonce:     98319381,
		},
		Transactions: nil,
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, b.EncodeBinary(buf))

	bDecode := &Block{}
	assert.Nil(t, bDecode.DecodeBinary(buf))
	assert.Equal(t, b, bDecode)
}

func TestBlock_Hash(t *testing.T) {
	b := &Block{
		Header: Header{
			Version:   1,
			PrevBlock: types.RandomHash(),
			Timestamp: time.Now().UnixNano(),
			Height:    10,
			Nonce:     98319381,
		},
		Transactions: nil,
	}

	hash := b.Hash()
	assert.False(t, hash.IsZero())
}
