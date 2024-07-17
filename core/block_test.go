package core

import (
	"github.com/stretchr/testify/assert"
	"sharkchain/crypto"
	"sharkchain/types"
	"testing"
	"time"
)

func randomZeroBlock(t *testing.T) *Block {
	return randomBlock(t, 0, types.Hash{})
}

func randomBlock(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privKey := crypto.GeneratePrivateKey()
	tx := randomTxWithSignature(t)

	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}

	b, err := NewBlock(header, []*Transaction{tx})
	if err != nil {
		return nil
	}

	if err1 := b.Sign(privKey); err1 != nil {
		return nil
	}

	return b
}

func TestSignBlock(t *testing.T) {
	b := randomZeroBlock(t)
	privKey := crypto.GeneratePrivateKey()

	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomZeroBlock(t)

	assert.Nil(t, b.Sign(privKey))
	assert.Nil(t, b.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivKey.PublicKey()

	assert.NotNil(t, b.Verify())
}
