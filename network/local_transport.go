package network

import (
	"fmt"
	"sync"
)

type LocalTransport struct {
	addr      NetAddr
	peers     map[NetAddr]*LocalTransport
	consumeCh chan RPC
	lock      sync.RWMutex
}

func NewLocalTransport(addr NetAddr) Transport {
	return &LocalTransport{
		addr:      addr,
		peers:     make(map[NetAddr]*LocalTransport),
		consumeCh: make(chan RPC),
	}
}

func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeCh
}

func (t *LocalTransport) Connect(tr Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.peers[tr.Addr()] = tr.(*LocalTransport)

	return nil
}

func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s: could not send message to unknown peer %s", t.addr, to)
	}

	peer.consumeCh <- RPC{
		From:    t.addr,
		Payload: payload,
	}

	return nil
}

func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}
