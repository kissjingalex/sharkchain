package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math/big"
	"sharkchain/types"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		R: r,
		S: s,
	}, nil
}

func GeneratePrivateKey() PrivateKey {
	return NewPrivateKeyFromReader(rand.Reader)
}

func NewPrivateKeyFromReader(r io.Reader) PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), r)
	if err != nil {
		panic(err)
	}

	return PrivateKey{
		key: key,
	}
}

func (k PrivateKey) PublicKey() PublicKey {
	return elliptic.MarshalCompressed(k.key.PublicKey, k.key.PublicKey.X, k.key.PublicKey.Y)
}

// PublicKey change PublicKey to type []byte due to failed to encode transaction
//
//	type PublicKey struct {
//		key *ecdsa.PublicKey
//	}
type PublicKey []byte

//func (k PublicKey) ToSlice() []byte {
//	b := elliptic.MarshalCompressed(k.key, k.key.X, k.key.Y)
//	return b
//}

func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k)

	return types.AddressFromBytes(h[len(h)-20:])
}

type Signature struct {
	S, R *big.Int
}

func (sig Signature) String() string {
	b := append(sig.S.Bytes(), sig.R.Bytes()...)
	return hex.EncodeToString(b)
}

func (sig Signature) Verify(pubKey PublicKey, data []byte) bool {
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), pubKey)
	key := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	return ecdsa.Verify(key, data, sig.R, sig.S)
}
