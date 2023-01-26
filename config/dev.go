package config

type DevUser struct {
	Email     string `mapstructure:"email"`
	UserID    string `mapstructure:"user_id"`
	Username  string `mapstructure:"username"`
	Firstname string `mapstructure:"firstname"`
	Lastname  string `mapstructure:"lastname"`
}

type MockInfo struct {
	DisablePVTBCMock bool `mapstructure:"disable_pvtbc_mock"`
	DisableBusMock   bool `mapstructure:"disable_bus_mock"`
}

type DevConfig struct {
	User          DevUser  `mapstructure:"user"`
	Mock          MockInfo `mapstructure:"mock"`
	JWTPublicKey  string   `mapstructure:"jwt_public_key"`
	JWTPrivateKey string   `mapstructure:"jwt_private_key"`
}
