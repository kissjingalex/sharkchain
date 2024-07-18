package main

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"log"
	"math/rand"
	"net"
	"sharkchain/core"
	"sharkchain/crypto"
	"sharkchain/network"
	"strconv"
	"time"
)

var LocalAddr = &net.TCPAddr{
	IP:   net.ParseIP("127.0.0.1"),
	Port: 8080,
	Zone: "",
}

var RemoteAddr = &net.TCPAddr{
	IP:   net.ParseIP("127.0.0.1"),
	Port: 8081,
	Zone: "",
}

func main() {
	validatorPrivKey := crypto.GeneratePrivateKey()

	trLocal := network.NewLocalTransport(LocalAddr)
	trRemote := network.NewLocalTransport(RemoteAddr)

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			//trRemote.SendMessage(trLocal.Addr(), []byte("hello world"))
			if err := sendTransaction(trRemote, trLocal.Addr()); err != nil {
				logrus.Errorf("failed to send transaction from remote, %+v\n", err)
			}

			time.Sleep(1 * time.Second)
		}
	}()

	opts := network.ServerOpts{
		PrivateKey: &validatorPrivKey,
		ID:         "Local",
		Transports: []network.Transport{trLocal},
	}

	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	s.Start()
}

func sendTransaction(tr network.Transport, to net.Addr) error {
	privKey := crypto.GeneratePrivateKey()
	data := []byte(strconv.FormatInt(int64(rand.Intn(1000)), 10))
	tx := core.NewTransaction(data)
	tx.Sign(privKey)

	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())
	tr.SendMessage(to, msg.Bytes())
	//fmt.Printf("send transaction = %+v\n", tx)
	return nil
}
