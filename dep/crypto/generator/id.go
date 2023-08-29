package generator

import (
	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto/encoder"
	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto/hasher"
)

//
// IDProvider calculates IDs for Virgil Instances like Public Keys, Virgil Card, etc.
//
type IDProvider interface {
	//
	// VirgilCardID returns an ID for a Virgil Card by its content snapshot.
	//
	VirgilCardID([]byte) string

	//
	// PublicKeyID returns a Public Key ID by its content.
	//
	PublicKeyID([]byte) string
}

//
// IDs constants.
//
const (
	PublicKeyBytesInID  = 8
	VirgilCardBytesInID = 32
)

//
// IDGenerator is a default implementation for the Virgil Security IDs calculation.
//
type IDGenerator struct {
	hasher  hasher.Provider
	encoder encoder.Provider
}

//
// NewID returns a default ID generator instance.
//
func NewID(hasher hasher.Provider, encoder encoder.Provider) *IDGenerator {

	return &IDGenerator{hasher: hasher, encoder: encoder}
}

//
// VirgilCardID returns an ID for a Virgil Card by its content snapshot.
//
func (g *IDGenerator) VirgilCardID(snapshot []byte) string {

	h := g.hasher.Hash(snapshot)
	data := h[:VirgilCardBytesInID]

	return g.encoder.EncodeToString(data)
}

//
// PublicKeyID returns a Public Key ID by its content.
//
func (g *IDGenerator) PublicKeyID(key []byte) string {

	h := g.hasher.Hash(key)
	data := h[:PublicKeyBytesInID]

	return g.encoder.EncodeToString(data)
}
