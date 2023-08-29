package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/VirgilSecurity/virgil-services-core-kit/models"

	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto/encoder"
	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto/generator"
	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto/hasher"
)

//
// Testing constants.
//
var (
	dataToBeSigned       = []byte("Some data to be encoded")
	invalidSignature     = []byte("Some invalid signature")
	invalidPublicKey     = []byte("Some invalid public key")
	emptyPublicKey       []byte
	emptyContentSnapshot []byte
	emptyExtraSnapshot   []byte
)

//
// Test ValidateSign :: with incorrect data :: returns an error.
//
func TestValidateSignWithIncorrectData(t *testing.T) {

	crypto := getCryptoUnderTest()

	err := crypto.ValidateVirgilCardSignature(dataToBeSigned, emptyExtraSnapshot, invalidPublicKey, invalidSignature)

	assert.Error(t, err)
}

//
// Test ValidateSign :: with an empty public key :: returns an error.
//
func TestValidateSignWithAnEmptyPublicKey(t *testing.T) {

	crypto := getCryptoUnderTest()

	err := crypto.ValidateVirgilCardSignature(dataToBeSigned, emptyExtraSnapshot, emptyPublicKey, invalidSignature)

	assert.Error(t, err)
}

//
// Test ValidateSign :: with an incorrect signature :: returns an error.
//
func TestValidateSignWithAnIncorrectSignature(t *testing.T) {

	crypto := getCryptoUnderTest()
	keyPair, err := crypto.GenerateKeyPair()

	assert.Nil(t, err)

	publicKeyBytes, err := keyPair.PublicKey().Encode()

	assert.Nil(t, err)

	err = crypto.ValidateVirgilCardSignature(dataToBeSigned, emptyExtraSnapshot, publicKeyBytes, invalidSignature)

	assert.Error(t, err)
}

//
// Test ValidateSign :: with all correct data :: passes and returns nil.
//
func TestValidateSignWithAllCorrectData(t *testing.T) {

	crypto := getCryptoUnderTest()
	keyPair, err := crypto.GenerateKeyPair()

	assert.Nil(t, err)

	publicKeyBytes, err := keyPair.PublicKey().Encode()

	assert.Nil(t, err)

	d := hasher.NewSHA512().Hash(dataToBeSigned)
	validSignature, err := crypto.SignVirgilCard(d, []byte{}, keyPair.PrivateKey())

	assert.Nil(t, err)

	err = crypto.ValidateVirgilCardSignature(d, emptyExtraSnapshot, publicKeyBytes, validSignature)

	assert.Nil(t, err)
}

//
// Test Sign :: with all correct data :: passes.
//
func TestSignWithAllCorrectData(t *testing.T) {

	crypto := getCryptoUnderTest()
	keyPair, err := crypto.GenerateKeyPair()

	assert.Nil(t, err)

	publicKeyBytes, err := keyPair.PublicKey().Encode()

	assert.Nil(t, err)

	d := hasher.NewSHA512().Hash(dataToBeSigned)
	signature, err := crypto.SignVirgilCard(d, []byte{}, keyPair.PrivateKey())

	assert.Nil(t, err)
	assert.NotNil(t, signature)
	assert.Empty(t, crypto.ValidateVirgilCardSignature(d, emptyExtraSnapshot, publicKeyBytes, signature))
}

//
// Test CalculateCardID :: with empty Snapshots :: passes.
//
func TestCalculateCardIDWithEmptySnapshots(t *testing.T) {

	crypto := getCryptoUnderTest()

	fingerprint := crypto.CalculateCardID(emptyContentSnapshot)

	assert.NotEmpty(t, fingerprint)
	assert.Len(t, fingerprint, models.IDLength)
}

//
// getCryptoUnderTest returns a Crypto object under test.
//
func getCryptoUnderTest() *Crypto {

	h := hasher.NewSHA512()
	e := encoder.NewHex()

	return NewCrypto(generator.NewID(h, e))

}
