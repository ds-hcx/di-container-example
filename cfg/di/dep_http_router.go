package di

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/cfg/di"
	"github.com/VirgilSecurity/virgil-services-core-kit/http"

	"github.com/VirgilSecurity/virgil-services-cards/src/router/routes"
)

//
// Dependency name.
//
const (
	DefHTTPRouter = "HTTPRouter"
)

//
// registerHTTPRouter dependency registrar.
//
func (c *Container) registerHTTPRouter() error {

	return c.RegisterDependency(
		DefHTTPRouter,
		func(ctx di.Context) (interface{}, error) {

			r := http.NewRouter(
				c.GetLogger(),
				c.GetConfig().GetMetricPrefix(),
				http.SetupHealthDependencyList(
					c.GetCassandraClient(),
				),
			)

			// Cards endpoints.
			routes.InitCardsRouteList(c.GetTracer(), r, c.GetCardsHandler())

			return r, nil

		}, nil,
	)
}

//
// GetHTTPRouter dependency retriever.
//
func (c *Container) GetHTTPRouter() http.RouterProvider {

	return c.Container.Get(DefHTTPRouter).(http.RouterProvider)
}
