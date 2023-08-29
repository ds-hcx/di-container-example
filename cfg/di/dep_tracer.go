package di

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/cfg/di"
	"github.com/VirgilSecurity/virgil-services-core-kit/tracer"
)

//
// Dependency name.
//
const (
	DefTracer = "Tracer"
)

//
// registerTracer dependency registrar.
//
func (c *Container) registerTracer() error {

	return c.RegisterDependency(
		DefTracer,
		func(ctx di.Context) (interface{}, error) {
			return tracer.NewTracer(
				c.GetLogger(),
				tracer.SetDisable(c.GetConfig().GetTracerDisabled()),
				tracer.SetServiceName("virgil-cards-service"),
				tracer.SetAgentHostPort(c.GetConfig().GetTracerAgentAddress()),
				tracer.SetSamplingServerURL(c.GetConfig().GetTracerSamplerManagerAddress()),
				tracer.SetSamplingType(c.GetConfig().GetTracerSamplerType()),
				tracer.SetSamplingParam(c.GetConfig().GetTracerSamplerParam()),
			)
		},
		nil,
	)
}

//
// GetTracer dependency retriever.
//
func (c *Container) GetTracer() tracer.Tracer {

	return c.Container.Get(DefTracer).(tracer.Provider).GetTracer()
}
