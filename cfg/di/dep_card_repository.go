package di

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/cfg/di"

	"github.com/VirgilSecurity/virgil-services-cards/src/dao"
)

//
// Dependency name.
//
const (
	DefCardRepository = "CardRepository"
)

//
// registerCardRepository dependency registrar.
//
func (c *Container) registerCardRepository() error {

	return c.RegisterDependency(
		DefCardRepository,
		func(ctx di.Context) (interface{}, error) {

			return dao.NewCardRepository(
				c.GetCassandraClient(),
			), nil
		},
		nil,
	)
}

//
// GetCardRepository dependency retriever.
//
func (c *Container) GetCardRepository() dao.CardRepositoryProvider {

	return c.Container.Get(DefCardRepository).(dao.CardRepositoryProvider)
}
