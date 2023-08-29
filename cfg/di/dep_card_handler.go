package di

import (
	"github.com/VirgilSecurity/virgil-services-cards/src/transport"
	"github.com/VirgilSecurity/virgil-services-core-kit/cfg/di"
)

//
// Dependency name.
//
const (
	DefCardsHandler = "CardsHandler"
)

//
// registerCardsHandler dependency registrar.
//
func (c *Container) registerCardsHandler() error {

	return c.RegisterDependency(
		DefCardsHandler,
		func(ctx di.Context) (interface{}, error) {

			return transport.NewCardsHandler(
				c.GetCardController(),
				c.GetEventMeter(),
			), nil
		},
		nil,
	)
}

//
// GetCardsHandler dependency retriever.
//
func (c *Container) GetCardsHandler() *transport.CardsHandler {

	return c.Container.Get(DefCardsHandler).(*transport.CardsHandler)
}
