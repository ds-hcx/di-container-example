package di

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/cfg/di"

	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto/hasher"
)

//
// Dependency name.
//
const (
	DefCryptoHasher = "CryptoHasher"
)

//
// registerCryptoHasher dependency registrar.
//
func (c *Container) registerCryptoHasher() error {

	return c.RegisterDependency(
		DefCryptoHasher,
		func(ctx di.Context) (interface{}, error) {

			return hasher.NewSHA512(), nil
		},
		nil,
	)
}

//
// GetCryptoHasher dependency retriever.
//
func (c *Container) GetCryptoHasher() hasher.Provider {

	return c.Container.Get(DefCryptoHasher).(hasher.Provider)
}
