package config

type ServerConfig struct {
	Host           string `mapstructure:"host"`
	Port           string `mapstructure:"port"`
	APIPrefix      string `mapstructure:"api_prefix"`
	ClientID       string `mapstructure:"client_id"`
	IdentityIssuer string `mapstructure:"identity_issuer"`
}

func (c *ServerConfig) GetServerURL() string {
	return c.Host + ":" + c.Port
}
