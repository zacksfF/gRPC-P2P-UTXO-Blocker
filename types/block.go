package types

import (
	"bytes"
	"crypto/sha256"

	"github.com/cbergoon/merkletree"
	proto "github.com/zacksfF/gRPC-P2P-UTXO-Blocker/Proto"
	"github.com/zacksfF/gRPC-P2P-UTXO-Blocker/encrypted"
	pb "google.golang.org/protobuf/proto"
)

type TxHash struct {
	hash []byte
}

func NewTxHash(hash []byte) TxHash {
	return TxHash{hash: hash}
}

func (t TxHash) CalculateHash() ([]byte, error) {
	return t.hash, nil
}

func (t TxHash) Equals(other merkletree.Content) (bool, error) {
	equals := bytes.Equal(t.hash, other.(TxHash).hash)
	return equals, nil
}

func SignBlock(pk *encrypted.PrivateKey, block *proto.Block) *encrypted.Signature {
	if len(block.Transactions) > 0 {
		tree, err := GetMerkleTree(block)
		if err != nil {
			panic(err)
		}

		block.Header.RootHash = tree.MerkleRoot()
	}

	hash := HashBlock(block)
	signature := pk.Sign(hash)
	block.PublicKey = pk.Public().Bytes()
	block.Signature = signature.Bytes()

	return signature
}

// HashBlock returns SHA256 of the header.
func HashBlock(block *proto.Block) []byte {
	return HashHeader(block.Header)
}

func HashHeader(header *proto.Header) []byte {
	b, err := pb.Marshal(header)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)

	return hash[:]
}

func VerifyBlock(block *proto.Block) bool {
	if len(block.Transactions) > 0 {
		if !VerifyRootHash(block) {
			return false
		}
	}

	if len(block.PublicKey) != encrypted.PublicKeyLen {
		return false
	}
	if len(block.Signature) != encrypted.SignatureLen {
		return false
	}

	var (
		signature = encrypted.SignatureFromBytes(block.Signature)
		publicKey = encrypted.PublicKeyFromBytes(block.PublicKey)
		hash      = HashBlock(block)
	)

	return signature.Verify(publicKey, hash)
}

func VerifyRootHash(block *proto.Block) bool {
	tree, err := GetMerkleTree(block)
	if err != nil {
		return false
	}

	valid, err := tree.VerifyTree()
	if err != nil {
		return false
	}

	if !valid {
		return false
	}

	return bytes.Equal(block.Header.RootHash, tree.MerkleRoot())
}

func GetMerkleTree(block *proto.Block) (*merkletree.MerkleTree, error) {
	list := make([]merkletree.Content, len(block.Transactions))
	for i := 0; i < len(block.Transactions); i++ {
		list[i] = NewTxHash(HashTransaction(block.Transactions[i]))
	}

	tree, err := merkletree.NewTree(list)
	if err != nil {
		return nil, err
	}

	return tree, nil
}
