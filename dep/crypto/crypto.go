package crypto

import (
	"gopkg.in/virgil.v4/virgilcrypto"
	"gopkg.in/virgilsecurity/virgil-crypto-go.v4"

	"github.com/VirgilSecurity/virgil-services-core-kit/errors"

	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto/generator"
)

//
// Provider is an interface for all cryptographic operations.
//
type Provider interface {
	//
	// ValidateVirgilCardSignature validates the signature for the Virgil Card CSR data using provided public key and
	// returns an error if the validation fails.
	//
	ValidateVirgilCardSignature(csr, extraCSR []byte, publicKey []byte, signature []byte) error

	//
	// Sign signs the data with provided private key and returns the signature value.
	//
	Sign(data []byte, key PrivateKey) ([]byte, error)

	//
	// CalculateCardID returns the card ID for the content snapshot and extra snapshot values.
	//
	CalculateCardID(contentSnapshot []byte) string

	//
	// CalculatePublicKeyID returns a public key id for key content.
	//
	CalculatePublicKeyID(key []byte) string

	//
	// SignVirgilCard signs a Virgil Card.
	//
	SignVirgilCard(csr, extraCSR []byte, key PrivateKey) ([]byte, error)

	//
	// ImportPrivateKey returns private key instance based on private key value and password.
	//
	ImportPrivateKey(privateKey []byte, password string) (PrivateKey, error)
}

//
// PrivateKey is a private key interface.
//
type PrivateKey virgilcrypto.PrivateKey

//
// PublicKey is a public key interface.
//
type PublicKey virgilcrypto.PublicKey

//
// KeyPair is a key pair interface.
//
type KeyPair virgilcrypto.Keypair

//
// Crypto is an adapter object for the Virgil Crypto SDK.
//
type Crypto struct {
	crypto      *virgil_crypto_go.NativeCrypto
	idGenerator generator.IDProvider
}

//
// NewCrypto returns a new crypto adapter instance.
//
func NewCrypto(idGenerator generator.IDProvider) *Crypto {

	return &Crypto{
		crypto:      &virgil_crypto_go.NativeCrypto{},
		idGenerator: idGenerator,
	}
}

//
// Sign signs the data with provided private key and its password and returns the signature value.
//
func (c *Crypto) Sign(data []byte, privateKey PrivateKey) ([]byte, error) {

	signature, err := c.crypto.Sign(data, privateKey)
	if nil != err {
		return nil, errors.WithMessage(err, `signing error for data (%s)`, data)
	}

	return signature, nil
}

//
// CalculateCardID returns a Virgil Card ID for the content snapshot.
//
func (c *Crypto) CalculateCardID(snapshot []byte) string {

	return c.idGenerator.VirgilCardID(snapshot)
}

//
// CalculatePublicKeyID returns a public key id for key content.
//
func (c *Crypto) CalculatePublicKeyID(key []byte) string {

	return c.idGenerator.PublicKeyID(key)
}

//
// SignVirgilCard signs a Virgil Card content snapshot.
//
func (c *Crypto) SignVirgilCard(csr, extraCSR []byte, privateKey PrivateKey) ([]byte, error) {

	return c.crypto.SignSHA512(append(csr, extraCSR...), privateKey)
}

//
// ValidateVirgilCardSignature performs a signature validation.
//
func (c *Crypto) ValidateVirgilCardSignature(csr, extraCSR []byte, publicKey []byte, signature []byte) error {

	key, err := c.crypto.ImportPublicKey(publicKey)
	if err != nil {
		return errors.WithMessage(err, `public key (%s) import error`, publicKey)
	}

	data := append(csr, extraCSR...)
	ok, err := c.crypto.Verify(data, signature, key)
	if nil != err {
		return errors.WithMessage(err, `signature (%s) verification error for data (%s) and key (%s)`,
			signature, csr, publicKey)
	}
	if !ok {
		return errors.New(`signature (%s) is incorrect fot csr (%s) and key (%s)`, signature, csr, publicKey)
	}

	return nil
}

//
// ImportPrivateKey returns private key instance based on private key value and password.
//
func (c *Crypto) ImportPrivateKey(privateKey []byte, password string) (PrivateKey, error) {

	return c.crypto.ImportPrivateKey(privateKey, password)
}

//
// GenerateKeyPair returns a new key-pair instance.
//
func (c *Crypto) GenerateKeyPair() (KeyPair, error) {

	return c.crypto.GenerateKeypair()
}
