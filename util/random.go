package util

import (
	randc "crypto/rand"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/zacksfF/gRPC-P2P-UTXO-Blocker/proto"
)

// RandomHash generates a random hash of 32 bytes.
func RandomHash() ([]byte, error) {
	hash := make([]byte, 32)
	if _, err := io.ReadFull(randc.Reader, hash); err != nil {
		return nil, fmt.Errorf("could not generate random hash: %w", err)
	}
	return hash, nil
}

// RandomBlock generates a random block with a random hash.
func RandomBlock() *proto.Block {
	prevHash, err := RandomHash()
	if err != nil {
		panic(err)
	}

	// merkleRoot, err := RandomHash()
	// if err != nil {
	// 	return nil, fmt.Errorf("could not generate random hash: %w", err)
	// }

	return &proto.Block{
		Header: &proto.Header{
			Version:      1,
			Height:       uint32(rand.Intn(1000)),
			PreviousHash: prevHash,
			MerkleRoot:   prevHash,
			Timestamp:    time.Now().UnixNano(),
		},
	}
}
