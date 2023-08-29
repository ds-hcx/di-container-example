package config

import "github.com/VirgilSecurity/virgil-services-core-kit/cfg/config"

//
// GetCassandraConnectionInfo returns a cassandra connection string.
//
func (c *Config) GetCassandraConnectionInfo() config.CassandraConnectionInfoProvider {

	return c.config.GetCassandraConnectionInfo(ConfCassandra)
}
