package config

type MockInfo struct {
	DisablePVTBCMock bool `mapstructure:"disable_pvtbc_mock"`
	DisableBusMock   bool `mapstructure:"disable_bus_mock"`
	DisableAuthMock  bool `mapstructure:"disable_auth_mock"`
}

type DevConfig struct {
	Mock          MockInfo `mapstructure:"mock"`
	JWTPublicKey  string   `mapstructure:"jwt_public_key"`
	JWTPrivateKey string   `mapstructure:"jwt_private_key"`
}
