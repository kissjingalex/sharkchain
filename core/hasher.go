package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"sharkchain/types"
)

// Hasher define hash behavior
type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct{}

// TODO only hash header?
func (BlockHasher) Hash(b *Header) types.Hash {
	h := sha256.Sum256(b.Bytes())
	return h
}

type TxHasher struct{}

// Hash will hash the whole bytes of the TX no exception.
func (TxHasher) Hash(tx *Transaction) types.Hash {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, tx.Data)
	//binary.Write(buf, binary.LittleEndian, tx.To)
	//binary.Write(buf, binary.LittleEndian, tx.Value)
	//binary.Write(buf, binary.LittleEndian, tx.From)
	//binary.Write(buf, binary.LittleEndian, tx.Nonce)

	return types.Hash(sha256.Sum256(buf.Bytes()))
}
