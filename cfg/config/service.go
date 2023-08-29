package config

//
// GetServicePrivateKey returns service Private Key.
//
func (c *Config) GetServicePrivateKey() []byte {
	// TODO we should move this data to the Vault.
	return c.config.GetBase64String(ConfServicePrivateKey)
}

//
// GetServicePrivateKeyPassword returns service Private Key password.
//
func (c *Config) GetServicePrivateKeyPassword() []byte {
	// TODO we should move this data to the Vault.
	return c.config.GetBase64String(ConfServicePrivateKeyPassword)
}
