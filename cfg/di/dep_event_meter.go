package di

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/cfg/di"
	"github.com/VirgilSecurity/virgil-services-core-kit/metrics"

	"github.com/VirgilSecurity/virgil-services-cards/src/events"
)

//
// Dependency name.
//
const (
	DefEventMeter = "EventMeter"
)

//
// registerEventMeter dependency registrar.
//
func (c *Container) registerEventMeter() error {

	return c.RegisterDependency(
		DefEventMeter,
		func(ctx di.Context) (interface{}, error) {

			logger := c.GetLogger()

			return events.NewEventMeter(
				metrics.NewManager(
					c.config.GetEventsAddress(),
					c.config.GetEventsPushPeriod(),
					c.GetLogger(),
				), logger,
			), nil
		},
		nil,
	)
}

//
// GetEventMeter dependency retriever.
//
func (c *Container) GetEventMeter() events.EventProvider {

	return c.Container.Get(DefEventMeter).(events.EventProvider)
}
