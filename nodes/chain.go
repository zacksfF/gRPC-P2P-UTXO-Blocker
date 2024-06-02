package nodes

import (
	"bytes"
	"encoding/hex"
	"fmt"

	proto "github.com/zacksfF/gRPC-P2P-UTXO-Blocker/Proto"
	"github.com/zacksfF/gRPC-P2P-UTXO-Blocker/encrypted"
	"github.com/zacksfF/gRPC-P2P-UTXO-Blocker/types"
)

const godSeed = "d12cda4733e2e24377cc161b55bf447a13a615d48838b33ab7634b77531734dc"

type HeaderList struct {
	headers []*proto.Header
}

func NewHeaderList() *HeaderList {
	return &HeaderList{[]*proto.Header{}}
}

func (h *HeaderList) Add(header *proto.Header) {
	h.headers = append(h.headers, header)
}

func (h *HeaderList) Get(index int) *proto.Header {
	if index > h.Height() {
		panic("index out of range")
	}

	return h.headers[index]
}

func (h *HeaderList) Height() int {
	return h.Len() - 1
}

func (h *HeaderList) Len() int {
	return len(h.headers)
}

type UTXO struct {
	Hash     string
	OutIndex int
	Amount   int64
	Spent    bool
}

type Chain struct {
	txStore    TXStorer
	blockStore BlockStorer
	utxoStore  UTXOStorer
	headers    *HeaderList
}

func NewChain(blockStore BlockStorer, txStore TXStorer) *Chain {
	chain := &Chain{
		blockStore: blockStore,
		txStore:    txStore,
		utxoStore:  NewMemoryUTXOStore(),
		headers:    NewHeaderList(),
	}

	err := chain.addBlock(createGenesisBlock())
	if err != nil {
		panic(err)
	}

	return chain
}

func (c *Chain) Height() int {
	return c.headers.Height()
}

func (c *Chain) AddBlock(block *proto.Block) error {
	if err := c.ValidateBlock(block); err != nil {
		return err
	}

	return c.addBlock(block)
}

func (c *Chain) addBlock(block *proto.Block) error {
	c.headers.Add(block.Header)

	for _, tx := range block.Transactions {
		if err := c.txStore.Put(tx); err != nil {
			return err
		}
		hash := hex.EncodeToString(types.HashTransaction(tx))

		for it, output := range tx.Outputs {
			utxo := &UTXO{
				Hash:     hash,
				Amount:   output.Amount,
				OutIndex: it,
				Spent:    false,
			}

			if err := c.utxoStore.Put(utxo); err != nil {
				return err
			}
		}

		for _, input := range tx.Inputs {
			key := fmt.Sprintf("%s_%d", hex.EncodeToString(input.PrevTxHash), input.PrevOutIndex)
			utxo, err := c.utxoStore.Get(key)
			if err != nil {
				return err
			}

			utxo.Spent = true
			if err := c.utxoStore.Put(utxo); err != nil {
				return err
			}
		}
	}

	return c.blockStore.Put(block)
}

func (c *Chain) GetBlockByHash(hash []byte) (*proto.Block, error) {
	hashHex := hex.EncodeToString(hash)

	return c.blockStore.Get(hashHex)
}

func (c *Chain) GetBlockByHeight(height int) (*proto.Block, error) {
	if c.Height() < height {
		return nil, fmt.Errorf("given height (%d) too high - current height (%d)", height, c.Height())
	}

	header := c.headers.Get(height)
	hash := types.HashHeader(header)

	return c.GetBlockByHash(hash)
}

func (c *Chain) ValidateBlock(block *proto.Block) error {
	// validate signature
	if !types.VerifyBlock(block) {
		return fmt.Errorf("invalid block signature")
	}

	// validate prev block hash
	currentBlock, err := c.GetBlockByHeight(c.Height())
	if err != nil {
		return err
	}

	hash := types.HashBlock(currentBlock)
	if !bytes.Equal(hash, block.Header.PrevHash) {
		return fmt.Errorf("prev block hash mismatch")
	}

	for _, tx := range block.Transactions {
		if err := c.ValidateTransaction(tx); err != nil {
			return err
		}
	}

	return nil
}

func (c *Chain) ValidateTransaction(tx *proto.Transaction) error {
	// verify signature
	if !types.VerifyTransaction(tx) {
		return fmt.Errorf("invalid transaction signature")
	}

	// check if inputs are not spent
	var (
		nInputs   = len(tx.Inputs)
		hash      = hex.EncodeToString(types.HashTransaction(tx))
		sumInputs = 0
	)
	for i := 0; i < nInputs; i++ {
		prevHash := hex.EncodeToString(tx.Inputs[i].PrevTxHash)
		key := fmt.Sprintf("%s_%d", prevHash, i)
		utxo, err := c.utxoStore.Get(key)
		if err != nil {
			return err
		}

		sumInputs += int(utxo.Amount)

		if utxo.Spent {
			return fmt.Errorf("input %d of transaction %s is already spent", i, hash)
		}
	}

	sumOutputs := 0
	for _, output := range tx.Outputs {
		sumOutputs += int(output.Amount)
	}
	if sumInputs < sumOutputs {
		return fmt.Errorf("transaction %s has insufficient funds", hash)
	}

	return nil
}

func createGenesisBlock() *proto.Block {
	privateKey := encrypted.NewPrivateKeyFromSeedString(godSeed)

	block := &proto.Block{
		Header: &proto.Header{
			Version: 1,
		},
		Transactions: []*proto.Transaction{},
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  []*proto.TxInput{},
		Outputs: []*proto.TxOutput{
			{
				Amount:  1000,
				Address: privateKey.Public().Address().Bytes(),
			},
		},
	}

	block.Transactions = append(block.Transactions, tx)
	types.SignBlock(privateKey, block)

	return block
}
