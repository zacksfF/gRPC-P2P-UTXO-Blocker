package types

import (
	"crypto/sha256"

	proto "github.com/zacksfF/gRPC-P2P-UTXO-Blocker/Proto"
	"github.com/zacksfF/gRPC-P2P-UTXO-Blocker/encrypted"
	pb "google.golang.org/protobuf/proto"
)

func SignTransaction(pk *encrypted.PrivateKey, tx *proto.Transaction) *encrypted.Signature {
	return pk.Sign(HashTransaction(tx))
}

func HashTransaction(tx *proto.Transaction) []byte {
	b, err := pb.Marshal(tx)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]
}

func VerifyTransaction(tx *proto.Transaction) bool{
	for _, input := range tx.Inputs{
		if len(input.Signature) == 0{
			panic("Transaction with no signature")
		}
		var(
			signature = encrypted.SignatureFromBytes(input.Signature)
			publicKey = encrypted.PublicKeyFromBytes(input.PublicKey)
		)

		input.Signature = nil
		if !signature.Verify(publicKey, HashTransaction(tx)){
			return false
		}
	}
	return true
}