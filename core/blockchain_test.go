package core

import (
	"fmt"
	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"sharkchain/types"
	"testing"
)

func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(log.NewNopLogger(), randomZeroBlock(t))
	assert.Nil(t, err)
	return bc
}

func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)
	return BlockHasher{}.Hash(prevHeader)
}

func TestBlockchain(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.NotNil(t, bc)
	assert.Equal(t, bc.Height(), uint32(0))
	fmt.Println(bc.Height())
}

func TestAddBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	lenBlocks := 1000
	for i := 0; i < lenBlocks; i++ {
		height := uint32(i + 1)
		prevBlockHash := getPrevBlockHash(t, bc, height)
		block := randomBlock(t, height, prevBlockHash)
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, bc.Height(), uint32(lenBlocks))
	assert.Equal(t, len(bc.headers), lenBlocks+1)

	assert.NotNil(t, bc.AddBlock(randomBlock(t, 89, types.Hash{})))
}

func TestAddBlockTooHigh(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	assert.NotNil(t, bc.AddBlock(randomBlock(t, 3, types.Hash{})))
}
