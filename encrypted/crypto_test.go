package encrypted

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	privateKey := GeneratePrivateKey()
	assert.Equal(t, PrivateKeyLen, len(privateKey.Bytes()))

	publicKey := privateKey.Public()
	assert.Equal(t, PublicKeyLen, len(publicKey.Bytes()))
}

func TestPrivateKeySign(t *testing.T) {
	privateKey := GeneratePrivateKey()
	publicKey := privateKey.Public()
	signature := privateKey.Sign([]byte("hello world"))

	assert.True(t, signature.Verify(publicKey, []byte("hello world")))
	assert.False(t, signature.Verify(publicKey, []byte("hello world!"))) // different message

	invalidPrivateKey := GeneratePrivateKey()
	invalidPublicKey := invalidPrivateKey.Public()
	assert.False(t, signature.Verify(invalidPublicKey, []byte("hello world"))) // different public key
}

func TestPublicKeyAddress(t *testing.T) {
	privateKey := GeneratePrivateKey()
	publicKey := privateKey.Public()
	address := publicKey.Address()

	assert.Equal(t, AddressLen, len(address.Bytes()))
}

func TestNewPrivateKeyFromString(t *testing.T) {
	var (
		seed       = "c96e14d8abd284946d6f8bd3fabe5d4e4a22fd63013382b040b247ec1a471060"
		privateKey = NewPrivateKeyFromString(seed)
		addressStr = "5955d2b0bbad576748b1dc3d499b57b517e2c144"
	)

	assert.Equal(t, PrivateKeyLen, len(privateKey.Bytes()))
	address := privateKey.Public().Address()
	assert.Equal(t, addressStr, address.String())
}
