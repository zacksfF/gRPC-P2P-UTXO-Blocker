package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	proto "github.com/zacksfF/gRPC-P2P-UTXO-Blocker/Proto"
	"github.com/zacksfF/gRPC-P2P-UTXO-Blocker/encrypted"
	"github.com/zacksfF/gRPC-P2P-UTXO-Blocker/nodes"
	"github.com/zacksfF/gRPC-P2P-UTXO-Blocker/util"
	"google.golang.org/grpc"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	validatorIndex := rand.Intn(3)

	makeNode("localhost:3000", []string{}, validatorIndex == 0)
	time.Sleep(time.Second)
	makeNode("localhost:3001", []string{"localhost:3000"}, validatorIndex == 1)
	time.Sleep(time.Second)
	makeNode("localhost:3002", []string{"localhost:3001"}, validatorIndex == 2)

	for {
		time.Sleep(time.Second)
		makeTransaction()
	}
}

func makeNode(listenAddr string, bootstrapNodes []string, isValidator bool) *nodes.Node {
	cfg := nodes.ServerConfig{
		Version:    "0.0.1",
		ListenAddr: listenAddr,
	}
	// if isValidator {
	// 	privKey := encrypted.GeneratePrivateKey()
	// 	cfg.PrivateKey = privKey
	// }

	n := nodes.NewNode(cfg)
	go func() {
		log.Fatal(n.Start(listenAddr, bootstrapNodes))
	}()

	return n
}

func makeTransaction() {
	client, err := grpc.Dial("localhost:3000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	c := proto.NewNodeClient(client)
	privKey := encrypted.GeneratePrivateKey()

	tx := &proto.Transaction{
		Version: 1,
		Inputs: []*proto.TxInput{
			{
				PrevTxHash:   util.RandomHash(),
				PrevOutIndex: 0,
				PublicKey:    privKey.Public().Bytes(),
			},
		},
		Outputs: []*proto.TxOutput{
			{
				Amount:  99,
				Address: privKey.Public().Address().Bytes(),
			},
		},
	}

	_, err = c.HandleTransaction(context.TODO(), tx)
	if err != nil {
		log.Fatal(err)
	}
}
