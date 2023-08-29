package di

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/cfg/di"
	"github.com/VirgilSecurity/virgil-services-core-kit/db/cassandra"
)

//
// Dependency name.
//
const (
	DefCassandraConnection = "CassandraConnection"
)

//
// registerCassandraClient dependency registrar.
//
func (c *Container) registerCassandraClient() error {

	return c.RegisterDependency(
		DefCassandraConnection,
		func(ctx di.Context) (interface{}, error) {

			return cassandra.NewConnection(
				c.GetLogger(),
				c.GetConfig().GetCassandraConnectionInfo(),
				c.GetConfig().GetMetricPrefix(),
			)
		},
		nil,
	)
}

//
// GetCassandraClient dependency retriever.
//
func (c *Container) GetCassandraClient() *cassandra.Connection {

	return c.Container.Get(DefCassandraConnection).(*cassandra.Connection)
}
