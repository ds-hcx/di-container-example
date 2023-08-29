package di

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/cfg/di"

	"github.com/VirgilSecurity/virgil-services-cards/src/app/controller"
)

//
// Dependency name.
//
const (
	DefValidatorCSR        = "ValidatorCSR"
	DefValidatorCSRStamps  = "ValidatorCSRStamps"
	DefValidatorCreateCard = "ValidatorCreateCard"
	DefValidatorSearchCard = "ValidatorSearchCard"
	DefValidatorDeleteCard = "ValidatorDeleteCard"
)

//
// registerValidatorCSR dependency registrar.
//
func (c *Container) registerValidatorCSR() error {

	return c.RegisterDependency(
		DefValidatorCSR,
		func(ctx di.Context) (interface{}, error) {

			return controller.NewCSRValidator(
				c.GetEncoderBase64()), nil
		},
		nil,
	)
}

//
// GetValidatorCSRStamps dependency retriever.
//
func (c *Container) GetValidatorCSRStamps() controller.CSRStampsValidatorProvider {

	return c.Container.Get(DefValidatorCSRStamps).(controller.CSRStampsValidatorProvider)
}

//
// registerValidatorCSRStamps dependency registrar.
//
func (c *Container) registerValidatorCSRStamps() error {

	return c.RegisterDependency(
		DefValidatorCSRStamps,
		func(ctx di.Context) (interface{}, error) {

			return controller.NewCSRStampsValidator(
				c.GetCrypto(),
				c.GetEncoderBase64()), nil
		},
		nil,
	)
}

//
// GetValidatorCSR dependency retriever.
//
func (c *Container) GetValidatorCSR() controller.CSRValidatorProvider {

	return c.Container.Get(DefValidatorCSR).(controller.CSRValidatorProvider)
}

//
// registerValidatorCreateCard dependency registrar.
//
func (c *Container) registerValidatorCreateCard() error {

	return c.RegisterDependency(
		DefValidatorCreateCard,
		func(ctx di.Context) (interface{}, error) {

			return controller.NewCreateCardValidator(
				c.GetCardRepository(),
				c.GetCrypto(),
				c.GetValidatorCSR(),
				c.GetValidatorCSRStamps(),
			), nil
		},
		nil,
	)
}

//
// GetValidatorCreateCard dependency retriever.
//
func (c *Container) GetValidatorCreateCard() controller.CreateCardValidatorProvider {

	return c.Container.Get(DefValidatorCreateCard).(controller.CreateCardValidatorProvider)
}

//
// registerValidatorSearchCard dependency registrar.
//
func (c *Container) registerValidatorSearchCard() error {

	return c.RegisterDependency(
		DefValidatorSearchCard,
		func(ctx di.Context) (interface{}, error) {

			return controller.NewSearchCardValidator(), nil
		},
		nil,
	)
}

//
// GetValidatorSearchCard dependency retriever.
//
func (c *Container) GetValidatorSearchCard() controller.SearchCardValidatorProvider {

	return c.Container.Get(DefValidatorSearchCard).(controller.SearchCardValidatorProvider)
}

//
// registerValidatorDeleteCard dependency registrar.
//
func (c *Container) registerValidatorDeleteCard() error {

	return c.RegisterDependency(
		DefValidatorDeleteCard,
		func(ctx di.Context) (interface{}, error) {

			return controller.NewDeleteCardValidator(
				c.GetCardRepository(),
				c.GetCrypto(),
				c.GetValidatorCSR(),
				c.GetValidatorCSRStamps(),
			), nil
		},
		nil,
	)
}

//
// GetValidatorDeleteCard dependency retriever.
//
func (c *Container) GetValidatorDeleteCard() controller.DeleteCardValidatorProvider {

	return c.Container.Get(DefValidatorDeleteCard).(controller.DeleteCardValidatorProvider)
}
