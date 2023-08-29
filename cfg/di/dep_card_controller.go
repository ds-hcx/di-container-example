package di

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/cfg/di"

	"github.com/VirgilSecurity/virgil-services-cards/src/app/controller"
)

//
// Dependency name.
//
const (
	DefCardController = "CardController"
)

//
// registerCardController dependency registrar.
//
func (c *Container) registerCardController() error {

	return c.RegisterDependency(
		DefCardController,
		func(ctx di.Context) (interface{}, error) {

			return controller.New(
				c.GetCardSigner(),
				c.GetCardRepository(),
				c.GetValidatorCreateCard(),
				c.GetValidatorSearchCard(),
				c.GetValidatorDeleteCard(),
			), nil
		},
		nil,
	)
}

//
// GetCardController dependency retriever.
//
func (c *Container) GetCardController() controller.Provider {

	return c.Container.Get(DefCardController).(controller.Provider)
}
