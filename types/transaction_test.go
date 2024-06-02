package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	proto "github.com/zacksfF/gRPC-P2P-UTXO-Blocker/Proto"
	"github.com/zacksfF/gRPC-P2P-UTXO-Blocker/encrypted"
	"github.com/zacksfF/gRPC-P2P-UTXO-Blocker/util"
)

func TestNewTransaction(t *testing.T) {
	fromPrivateKey := encrypted.GeneratePrivateKey()
	fromAddress := fromPrivateKey.Public().Address().Bytes()
	toPrivateKey := encrypted.GeneratePrivateKey()
	toAddress := toPrivateKey.Public().Address().Bytes()

	input := &proto.TxInput{
		PrevTxHash:   util.RandomHash(),
		PrevOutIndex: 0,
		PublicKey:    fromPrivateKey.Public().Bytes(),
	}

	output1 := &proto.TxOutput{
		Amount:  5,
		Address: toAddress,
	}

	output2 := &proto.TxOutput{
		Amount:  95,
		Address: fromAddress,
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  []*proto.TxInput{input},
		Outputs: []*proto.TxOutput{output1, output2},
	}

	signature := SignTransaction(fromPrivateKey, tx)
	input.Signature = signature.Bytes()

	assert.True(t, VerifyTransaction(tx))
	fmt.Printf("%v\n", tx)
}
