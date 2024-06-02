package nodes

import (
	"encoding/hex"
	"fmt"
	"sync"

	proto "github.com/zacksfF/gRPC-P2P-UTXO-Blocker/Proto"
	"github.com/zacksfF/gRPC-P2P-UTXO-Blocker/types"
)

type UTXOStorer interface {
	Put(*UTXO) error
	Get(string) (*UTXO, error)
}

type MemoryUTXOStore struct {
	lock sync.RWMutex
	data map[string]*UTXO
}

func NewMemoryUTXOStore() *MemoryUTXOStore {
	return &MemoryUTXOStore{
		data: make(map[string]*UTXO),
	}
}

func (m *MemoryUTXOStore) Put(utxo *UTXO) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	key := fmt.Sprintf("%s_%d", utxo.Hash, utxo.OutIndex)
	m.data[key] = utxo

	return nil
}

func (m *MemoryUTXOStore) Get(hash string) (*UTXO, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	utxo, ok := m.data[hash]
	if !ok {
		return nil, fmt.Errorf("utxo with hash [%s] does not exist", hash)
	}

	return utxo, nil
}

type TXStorer interface {
	Put(*proto.Transaction) error
	Get(string) (*proto.Transaction, error)
}

type MemoryTXStore struct {
	lock sync.RWMutex
	txx  map[string]*proto.Transaction
}

func NewMemoryTXStore() *MemoryTXStore {
	return &MemoryTXStore{
		txx: make(map[string]*proto.Transaction),
	}
}

func (m *MemoryTXStore) Put(tx *proto.Transaction) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	hash := hex.EncodeToString(types.HashTransaction(tx))
	m.txx[hash] = tx

	return nil
}

func (m *MemoryTXStore) Get(hash string) (*proto.Transaction, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	tx, ok := m.txx[hash]
	if !ok {
		return nil, fmt.Errorf("transaction with hash [%s] does not exist", hash)
	}

	return tx, nil
}

type BlockStorer interface {
	Put(*proto.Block) error
	Get(string) (*proto.Block, error)
}

type MemoryBlockStore struct {
	lock   sync.RWMutex
	blocks map[string]*proto.Block
}

func NewMemoryBlockStore() *MemoryBlockStore {
	return &MemoryBlockStore{
		blocks: make(map[string]*proto.Block),
	}
}

func (m *MemoryBlockStore) Put(block *proto.Block) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	hash := hex.EncodeToString(types.HashBlock(block))
	m.blocks[hash] = block

	return nil
}

func (m *MemoryBlockStore) Get(hash string) (*proto.Block, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	block, ok := m.blocks[hash]
	if !ok {
		return nil, fmt.Errorf("block with hash [%s] does not exist", hash)
	}

	return block, nil
}
