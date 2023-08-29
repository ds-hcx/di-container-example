package config

import "time"

//
// GetEventsAddress returns an address to push statistics.
//
func (c *Config) GetEventsAddress() string {

	return c.config.GetString(ConfEventsAddress)
}

//
// GetEventsPushPeriod returns an statistics pushing period.
//
func (c *Config) GetEventsPushPeriod() time.Duration {

	return c.config.GetDuration(ConfEventsPushPeriod)
}
