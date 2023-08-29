package config

import (
	"time"

	"github.com/VirgilSecurity/virgil-services-core-kit/cfg/config"
	"github.com/VirgilSecurity/virgil-services-core-kit/errors"
)

//
// Configuration parameter names.
//
const (
	ConfCassandra                   = "CARDS5_CASSANDRA"
	ConfServerHTTPAddress           = "CARDS5_SERVER_ADDRESS"
	ConfServerReadTimeout           = "CARDS5_SERVER_READ_TIMEOUT"
	ConfServerWriteTimeout          = "CARDS5_SERVER_WRITE_TIMEOUT"
	ConfLogLevel                    = "CARDS5_LOG_LEVEL"
	ConfEventsAddress               = "CARDS5_EVENTS_ADDRESS"
	ConfEventsPushPeriod            = "CARDS5_EVENTS_PUSH_PERIOD"
	ConfServicePrivateKey           = "CARDS5_PRIVATE_KEY"
	ConfServicePrivateKeyPassword   = "CARDS5_PRIVATE_KEY_PASSWORD"
	ConfTracerDisabled              = "CARDS5_TRACER_DISABLED"
	ConfTracerAgentAddress          = "CARDS5_TRACER_AGENT_ADDRESS"
	ConfTracerSamplerType           = "CARDS5_TRACER_SAMPLER_TYPE"
	ConfTracerSamplerParam          = "CARDS5_TRACER_SAMPLER_PARAM"
	ConfTracerSamplerManagerAddress = "CARDS5_TRACER_SAMPLER_MANAGER_ADDRESS"
)

//
// General constants.
//
const (
	MetricPrefix = "virgil_cards5"
	ServiceName  = "cards"
)

//
// Config is an application config object.
//
type Config struct {
	config *config.Config
}

//
// New returns a new Config instance.
//
func New() (*Config, error) {

	c := config.New()
	c.RegisterParameters(
		config.NewString(
			ConfServerHTTPAddress,
			"HTTP server address for binding.",
			":8080",
		),
		config.NewDuration(
			ConfServerReadTimeout,
			"HTTP server read timeout.",
			5*time.Second,
		),
		config.NewDuration(
			ConfServerWriteTimeout,
			"HTTP server write timeout.",
			5*time.Second,
		),

		config.NewLoggerLevel(
			ConfLogLevel,
			"Logging level",
		),

		config.NewCassandraInfo(
			ConfCassandra,
			"Cassandra db connection string.",
		),

		config.NewString(
			ConfEventsAddress,
			"Business events listener.",
			"",
		),
		config.NewDuration(
			ConfEventsPushPeriod,
			"Events push period.",
			5*time.Second,
		),

		config.NewBase64String(
			ConfServicePrivateKey,
			"Cards Service Private Key.",
			"",
		),
		config.NewBase64String(
			ConfServicePrivateKeyPassword,
			"Cards Service Private Key Password.",
			"",
		),

		config.NewBool(
			ConfTracerDisabled,
			"Enable/Disable request tracing..",
			true,
		),
		config.NewString(
			ConfTracerAgentAddress,
			"Address where tracer agent is listening on.",
			"",
		),
		config.NewString(
			ConfTracerSamplerType,
			"Sampler type. Allowed values are: remote, const, probabilistic, rateLimiting.",
			"probabilistic",
		),
		config.NewFloat64(
			ConfTracerSamplerParam,
			"The sampler parameter (number)",
			0.1,
		),
		config.NewString(
			ConfTracerSamplerManagerAddress,
			"Address of remote sampler manager.",
			"",
		),
	)

	if err := c.Parse(); nil != err {
		return nil, err
	}

	// TODO I want to setup this check as a custom validator for given parameter.
	privateKey := c.GetBase64String(ConfServicePrivateKey)
	if len(privateKey) == 0 {
		return nil, errors.New("config parameter (%s) was not set", ConfServicePrivateKey)
	}

	if len(privateKey) != 0 && len(c.GetBase64String(ConfServicePrivateKeyPassword)) == 0 {
		return nil, errors.New("config parameter (%s) was not set", ConfServicePrivateKeyPassword)
	}

	return &Config{
		config: c,
	}, nil
}
