package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnect(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)

	assert.Equal(t, tra.peers[trb.addr], trb)
	assert.Equal(t, trb.peers[tra.addr], tra)
}

func TestLocalTransport_SendMessage(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)

	msg := []byte("hello world")
	assert.Nil(t, tra.SendMessage(trb.addr, msg))

	rpc := <-trb.Consume()
	assert.Equal(t, rpc.Payload, msg)
}
