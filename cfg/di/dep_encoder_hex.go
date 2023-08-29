package di

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/cfg/di"

	"github.com/VirgilSecurity/virgil-services-cards/src/dep/crypto/encoder"
)

//
// Dependency name.
//
const (
	DefEncoderHex = "EncoderHex"
)

//
// registerEncoderHex dependency registrar.
//
func (c *Container) registerEncoderHex() error {

	return c.RegisterDependency(
		DefEncoderHex,
		func(ctx di.Context) (interface{}, error) {

			return encoder.NewHex(), nil
		},
		nil,
	)
}

//
// GetEncoderHex dependency retriever.
//
func (c *Container) GetEncoderHex() encoder.Provider {

	return c.Container.Get(DefEncoderHex).(encoder.Provider)
}
