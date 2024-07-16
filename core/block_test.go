package core

import (
	"github.com/stretchr/testify/assert"
	"sharkchain/crypto"
	"sharkchain/types"
	"testing"
	"time"
)

func randomBlock(height uint32) *Block {
	//privKey := crypto.GeneratePrivateKey()
	//tx := randomTxWithSignature(t)

	header := &Header{
		Version:       1,
		PrevBlockHash: types.RandomHash(),
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}

	tx := Transaction{
		Data: []byte("foo"),
	}

	b, err := NewBlock(header, []*Transaction{&tx})
	if err != nil {
		return nil
	}

	return b
}

func TestSignBlock(t *testing.T) {
	b := randomBlock(0)
	privKey := crypto.GeneratePrivateKey()

	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(0)

	assert.Nil(t, b.Sign(privKey))
	assert.Nil(t, b.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivKey.PublicKey()

	assert.NotNil(t, b.Verify())
}
