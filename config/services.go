package config

type ServicesConfig struct {
	Keycloak  string `mapstructure:"keycloak"`
	Validator string `mapstructure:"validator"`
}
