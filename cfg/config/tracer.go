package config

import "github.com/VirgilSecurity/virgil-services-core-kit/tracer"

//
// GetTracerDisabled returns the status of tracer - should tracer be enabled or disabled.
//
func (c *Config) GetTracerDisabled() bool {

	return c.config.GetBool(ConfTracerDisabled)
}

//
// GetTracerAgentAddress returns a host where agent is listening on.
//
func (c *Config) GetTracerAgentAddress() string {

	return c.config.GetString(ConfTracerAgentAddress)
}

//
// GetTracerSamplerType returns sampler type.
//
func (c *Config) GetTracerSamplerType() tracer.SamplerType {

	return tracer.SamplerType(c.config.GetString(ConfTracerSamplerType))
}

//
// GetTracerSamplerParam returns sampler param.
//
func (c *Config) GetTracerSamplerParam() float64 {

	return c.config.GetFloat64(ConfTracerSamplerParam)
}

//
// GetTracerSamplerManagerAddress returns a host where sampler manager is listening on.
//
func (c *Config) GetTracerSamplerManagerAddress() string {

	return c.config.GetString(ConfTracerSamplerManagerAddress)
}
