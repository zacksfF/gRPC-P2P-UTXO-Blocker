package nodes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	proto "github.com/zacksfF/gRPC-P2P-UTXO-Blocker/Proto"
	"github.com/zacksfF/gRPC-P2P-UTXO-Blocker/encrypted"
	"github.com/zacksfF/gRPC-P2P-UTXO-Blocker/types"
	"github.com/zacksfF/gRPC-P2P-UTXO-Blocker/util"
)

func randomBlock(t *testing.T, chain *Chain) *proto.Block {
	privKey := encrypted.GeneratePrivateKey()
	block := util.RandomBlock()
	prevBlock, err := chain.GetBlockByHeight(chain.Height())
	require.Nil(t, err)

	block.Header.PrevHash = types.HashBlock(prevBlock)
	types.SignBlock(privKey, block)

	return block
}

func TestNewChain(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
	require.Equal(t, 0, chain.Height())

	_, err := chain.GetBlockByHeight(0)
	require.Nil(t, err)
}

func TestChainHeight(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())

	for i := 0; i < 10; i++ {
		block := randomBlock(t, chain)

		require.Nil(t, chain.AddBlock(block))
		require.Equal(t, i+1, chain.Height())
	}
}

func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())

	for i := 0; i < 100; i++ {
		block := randomBlock(t, chain)
		blockHash := types.HashBlock(block)

		require.Nil(t, chain.AddBlock(block))

		fetchedBlock, err := chain.GetBlockByHash(blockHash)
		require.Nil(t, err)
		assert.Equal(t, block, fetchedBlock)

		fetchedBlockByHeight, err := chain.GetBlockByHeight(i + 1)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlockByHeight)
	}
}

func TestAddBlockWithTX(t *testing.T) {
	var (
		chain      = NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
		block      = randomBlock(t, chain)
		privateKey = encrypted.NewPrivateKeyFromSeedString(godSeed)
		recipient  = encrypted.GeneratePrivateKey().Public().Address().Bytes()
	)

	prevTx, err := chain.txStore.Get("737cb341c1dd4a34f78faeb9ae8ffc07b5e9cf987a1ef0bfb65d98542f48aacb")
	assert.Nil(t, err)

	inputs := []*proto.TxInput{
		{
			PrevTxHash:   types.HashTransaction(prevTx),
			PrevOutIndex: 0,
			PublicKey:    privateKey.Public().Bytes(),
		},
	}
	outputs := []*proto.TxOutput{
		{
			Amount:  100,
			Address: recipient,
		},
		{
			Amount:  900,
			Address: privateKey.Public().Address().Bytes(),
		},
	}
	tx := &proto.Transaction{
		Version: 1,
		Inputs:  inputs,
		Outputs: outputs,
	}

	signature := types.SignTransaction(privateKey, tx)
	tx.Inputs[0].Signature = signature.Bytes()

	block.Transactions = append(block.Transactions, tx)
	types.SignBlock(privateKey, block)
	require.Nil(t, chain.AddBlock(block))
}

func TestAddBlockWithInsufficientFunds(t *testing.T) {
	var (
		chain      = NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
		block      = randomBlock(t, chain)
		privateKey = encrypted.NewPrivateKeyFromSeedString(godSeed)
		recipient  = encrypted.GeneratePrivateKey().Public().Address().Bytes()
	)

	prevTx, err := chain.txStore.Get("737cb341c1dd4a34f78faeb9ae8ffc07b5e9cf987a1ef0bfb65d98542f48aacb")
	assert.Nil(t, err)

	inputs := []*proto.TxInput{
		{
			PrevTxHash:   types.HashTransaction(prevTx),
			PrevOutIndex: 0,
			PublicKey:    privateKey.Public().Bytes(),
		},
	}
	outputs := []*proto.TxOutput{
		{
			Amount:  1001,
			Address: recipient,
		},
	}
	tx := &proto.Transaction{
		Version: 1,
		Inputs:  inputs,
		Outputs: outputs,
	}

	signature := types.SignTransaction(privateKey, tx)
	tx.Inputs[0].Signature = signature.Bytes()

	block.Transactions = append(block.Transactions, tx)
	require.NotNil(t, chain.AddBlock(block))
}
