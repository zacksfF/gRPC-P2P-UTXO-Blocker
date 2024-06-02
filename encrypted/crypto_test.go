package encrypted

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratePrivateKey()

	assert.Equal(t, len(privKey.Bytes()), PrivKeyLen)
	pubKey := privKey.Public()
	assert.Equal(t, len(pubKey.Bytes()), PubKeyLen)
}

func TestPrivateKeySign(t *testing.T) {
	privkey := GeneratePrivateKey()
	pubKey := privkey.Public()
	msg := []byte("Zack Block")

	sign := privkey.Sign(msg)
	assert.True(t, sign.Verify(pubKey, msg))

	//test with invalid msg
	assert.False(t, sign.Verify(pubKey, []byte("zakck")))

	//test with invalid pubkey
	invalidPrivateKey := GeneratePrivateKey()
	invalidPubkey := invalidPrivateKey.Public()
	assert.False(t, sign.Verify(invalidPubkey, msg))
}

func TestPublicKeyAddress(t *testing.T) {
	privkey := GeneratePrivateKey()
	pubKey := privkey.Public()
	address := pubKey.Address()
	assert.Equal(t, addressLen, len(address.Bytes()))
	fmt.Println(address)
}

func TestNewPrivateKeyFromString(t *testing.T) {
	var (
		seed          = "ba547efd3869cd3f6c74bc7fab1178499da44ba1ba10e7b9063c386defe8921c"
		privKey       = newPrivateKeyFromString(seed)
		addressString = "b9b15db4d715ab58ac628d55fb7264571203ddb39d817ac29844d75781c6"
	)

	// seed := make([]byte, 32)
	// io.ReadFull(rand.Reader, seed)
	// fmt.Println(hex.EncodeToString(seed))
	assert.Equal(t, PrivKeyLen, len(privKey.Bytes()))
	Address := privKey.Public().Address()
	assert.Equal(t, addressString, Address.String())

}
