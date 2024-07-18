package main

import (
	"sharkchain/util"
	"time"
)

func main() {
	// test gossip mechanism
	remoteNode := util.MakeServer("REMOTE_NODE", nil, ":4000", []string{":3000"}, "")
	go remoteNode.Start()

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
