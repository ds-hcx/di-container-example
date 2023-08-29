package di

import "github.com/VirgilSecurity/virgil-services-cards/src/cfg/config"

//
// GetConfig dependency retriever.
//
func (c *Container) GetConfig() *config.Config {

	return c.config
}
