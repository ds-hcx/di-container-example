package di

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/cfg/di"

	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto/generator"
)

//
// Dependency name.
//
const (
	DefCryptoIDGenerator = "CryptoIDGenerator"
)

//
// registerCryptoIDGenerator dependency registrar.
//
func (c *Container) registerCryptoIDGenerator() error {

	return c.RegisterDependency(
		DefCryptoIDGenerator,
		func(ctx di.Context) (interface{}, error) {

			return generator.NewID(
				c.GetCryptoHasher(),
				c.GetEncoderHex(),
			), nil
		},
		nil,
	)
}

//
// GetCryptoIDGenerator dependency retriever.
//
func (c *Container) GetCryptoIDGenerator() generator.IDProvider {

	return c.Container.Get(DefCryptoIDGenerator).(generator.IDProvider)
}
