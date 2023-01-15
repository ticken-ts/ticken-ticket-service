package config

type PublicBlockchainConfig struct {
	AddressPK string `mapstructure:"address_pk"`
	ChainURL  string `mapstructure:"chain_url"`
	ChainID   int64  `mapstructure:"chain_id"`
}
