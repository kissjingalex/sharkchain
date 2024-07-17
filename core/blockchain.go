package core

import (
	"fmt"
	"github.com/go-kit/log"
	"sync"
)

type Blockchain struct {
	logger log.Logger
	store  Storage

	lock sync.RWMutex

	headers   []*Header
	blocks    []*Block
	validator Validator
}

func NewBlockchain(l log.Logger, genesis *Block) (*Blockchain, error) {
	// We should create all states inside the scope of the newblockchain.

	bc := &Blockchain{
		logger:    l,
		store:     NewMemoryStore(),
		headers:   []*Header{},
		validator: nil,
	}

	bc.validator = NewBlockValidator(bc)

	err := bc.addBlockWithoutValidation(genesis)
	return bc, err
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.headers[height], nil
}

func (bc *Blockchain) GetBlock(height uint32) (*Block, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.blocks[height], nil
}

func (bc *Blockchain) AddBlock(b *Block) error {
	// validate
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.headers = append(bc.headers, b.Header)

	bc.logger.Log(
		"msg", "new block",
		"hash", b.Hash(BlockHasher{}),
		"height", b.Height,
		"transactions", len(b.Transactions),
	)

	return bc.store.Put(b)
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

// [0, 1, 2] : height=2
func (bc *Blockchain) Height() uint32 {
	if len(bc.headers) == 0 {
		return 0
	}

	return uint32(len(bc.headers) - 1) // maybe -1
}
