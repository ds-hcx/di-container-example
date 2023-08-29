package di

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/cfg/di"

	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto"
)

//
// Dependency name.
//
const (
	DefCrypto = "Crypto"
)

//
// registerCrypto dependency registrar.
//
func (c *Container) registerCrypto() error {

	return c.RegisterDependency(
		DefCrypto,
		func(ctx di.Context) (interface{}, error) {

			return crypto.NewCrypto(c.GetCryptoIDGenerator()), nil
		},
		nil,
	)
}

//
// GetCrypto dependency retriever.
//
func (c *Container) GetCrypto() crypto.Provider {

	return c.Container.Get(DefCrypto).(crypto.Provider)
}
