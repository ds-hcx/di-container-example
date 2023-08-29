package di

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/cfg/di"
	"github.com/VirgilSecurity/virgil-services-core-kit/errors"
	"github.com/VirgilSecurity/virgil-services-core-kit/log"

	"github.com/VirgilSecurity/virgil-services-cards/src/cfg/config"
)

//
// Container is a dependency resolver object.
//
type Container struct {
	logger log.Logger
	config *config.Config
	*di.Container
}

//
// NewContainer returns an instance of the DI DIContainer.
//
func NewContainer(c *config.Config, l log.Logger) (*Container, error) {

	container, err := di.NewContainer()
	if err != nil {
		return nil, errors.WithMessage(err, `di container instantiating error`)
	}

	return &Container{
		config:    c,
		logger:    l,
		Container: container,
	}, nil
}

//
// Build builds the application dependencies at once.
//
func (c *Container) Build() error {

	for _, dep := range []func() error{
		c.registerCrypto,
		c.registerCryptoIDGenerator,
		c.registerCryptoHasher,
		c.registerHTTPRouter,
		c.registerLogger,
		c.registerCassandraClient,
		c.registerEventMeter,
		c.registerCardsHandler,
		c.registerCardController,
		c.registerCardRepository,
		c.registerCardSigner,
		c.registerTracer,
		c.registerEncoderBase64,
		c.registerEncoderHex,
		c.registerValidatorCSR,
		c.registerValidatorCSRStamps,
		c.registerValidatorCreateCard,
		c.registerValidatorSearchCard,
		c.registerValidatorDeleteCard,
	} {
		if err := dep(); err != nil {
			return err
		}
	}

	c.Container.Build()

	return nil
}
