package main

import (
	"bytes"
	"math/rand"
	"net"
	"sharkchain/core"
	"sharkchain/crypto"
	"sharkchain/network"
	"sharkchain/util"
	"strconv"
	"time"
)

func main() {
	validatorPrivKey := crypto.GeneratePrivateKey()

	// isValidator
	localNode := util.MakeServer("LOCAL_NODE", &validatorPrivKey, ":3000", []string{":4000"}, ":9000")
	go localNode.Start()

	// test gossip mechanism
	//remoteNode := makeServer("REMOTE_NODE", nil, ":4000", []string{":3000"}, "")
	//go remoteNode.Start()

	//remoteNodeB := makeServer("REMOTE_NODE_B", nil, ":5000", nil, "")
	//go remoteNodeB.Start()
	//
	//go func() {
	//	time.Sleep(11 * time.Second)
	//
	//	lateNode := makeServer("LATE_NODE", nil, ":6000", []string{":4000"}, "")
	//	go lateNode.Start()
	//}()

	time.Sleep(1 * time.Second)

	// if err := sendTransaction(validatorPrivKey); err != nil {
	// 	panic(err)
	// }

	// collectionOwnerPrivKey := crypto.GeneratePrivateKey()
	// collectionHash := createCollectionTx(collectionOwnerPrivKey)

	// txSendTicker := time.NewTicker(1 * time.Second)
	// go func() {
	// 	for i := 0; i < 20; i++ {
	// 		nftMinter(collectionOwnerPrivKey, collectionHash)

	// 		<-txSendTicker.C
	// 	}
	// }()

	select {}
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
