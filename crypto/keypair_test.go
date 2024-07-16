package crypto

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratePrivateKey()

	pubKey := privKey.PublicKey()
	address := pubKey.Address()

	fmt.Println(address)
}

func TestPrivateKey_Sign(t *testing.T) {
	privKey := GeneratePrivateKey()

	pubKey := privKey.PublicKey()

	msg := []byte("hello world")
	sig, err := privKey.Sign(msg)
	assert.Nil(t, err)

	b := sig.Verify(pubKey, msg)
	assert.True(t, b)
}
