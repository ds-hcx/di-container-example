package di

import (
	"os"

	"github.com/VirgilSecurity/virgil-services-core-kit/cfg/di"
	"github.com/VirgilSecurity/virgil-services-core-kit/log"
)

//
// Dependency name.
//
const (
	DefLogger = "Logger"
)

//
// registerLogger dependency registrar.
//
func (c *Container) registerLogger() error {

	return c.RegisterDependency(
		DefLogger,
		func(ctx di.Context) (interface{}, error) {

			return log.New(os.Stdout, c.config.GetLogLevel()), nil
		},
		nil,
	)
}

//
// GetLogger dependency retriever.
//
func (c *Container) GetLogger() log.Logger {

	return c.Container.Get(DefLogger).(log.Logger)
}
