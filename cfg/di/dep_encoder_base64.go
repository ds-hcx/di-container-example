package di

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/cfg/di"

	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto/encoder"
)

//
// Dependency name.
//
const (
	DefEncoderBase64 = "EncoderBase64"
)

//
// registerEncoderBase64 dependency registrar.
//
func (c *Container) registerEncoderBase64() error {

	return c.RegisterDependency(
		DefEncoderBase64,
		func(ctx di.Context) (interface{}, error) {

			return encoder.NewBase64(), nil
		},
		nil,
	)
}

//
// GetEncoderBase64 dependency retriever.
//
func (c *Container) GetEncoderBase64() encoder.Provider {

	return c.Container.Get(DefEncoderBase64).(encoder.Provider)
}
