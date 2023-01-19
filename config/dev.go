package config

type DevUser struct {
	Email     string `mapstructure:"email"`
	UserID    string `mapstructure:"user_id"`
	Username  string `mapstructure:"username"`
	Firstname string `mapstructure:"firstname"`
	Lastname  string `mapstructure:"lastname"`
}

type DevConfig struct {
	User          DevUser `mapstructure:"user"`
	JWTPublicKey  string  `mapstructure:"jwt_public_key"`
	JWTPrivateKey string  `mapstructure:"jwt_private_key"`
}
