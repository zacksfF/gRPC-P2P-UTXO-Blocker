package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	proto "github.com/zacksfF/gRPC-P2P-UTXO-Blocker/Proto"
	"github.com/zacksfF/gRPC-P2P-UTXO-Blocker/encrypted"
	"github.com/zacksfF/gRPC-P2P-UTXO-Blocker/util"
)

func TestVerifyBlock(t *testing.T) {
	var (
		block      = util.RandomBlock()
		privateKey = encrypted.GeneratePrivateKey()
		publicKey  = privateKey.Public()
	)

	signature := SignBlock(privateKey, block)

	assert.Equal(t, 64, len(signature.Bytes()))
	assert.True(t, signature.Verify(publicKey, HashBlock(block)))

	assert.Equal(t, block.PublicKey, publicKey.Bytes())
	assert.Equal(t, block.Signature, signature.Bytes())

	assert.True(t, VerifyBlock(block))

	invalidPrivKey := encrypted.GeneratePrivateKey()
	block.PublicKey = invalidPrivKey.Public().Bytes()

	assert.False(t, VerifyBlock(block))
}

func TestHashBlock(t *testing.T) {
	block := util.RandomBlock()
	hash := HashBlock(block)

	assert.Equal(t, 32, len(hash))
}

func TestCalculateRootHash(t *testing.T) {
	privateKey := encrypted.GeneratePrivateKey()
	block := util.RandomBlock()
	tx := &proto.Transaction{
		Version: 1,
	}
	block.Transactions = append(block.Transactions, tx)
	SignBlock(privateKey, block)

	assert.True(t, VerifyRootHash(block))
	assert.Equal(t, 32, len(block.Header.RootHash))
}
