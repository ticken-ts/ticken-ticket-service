package config

import "github.com/spf13/viper"

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Pvtbc    PvtbcConfig    `mapstructure:"pvtbc"`
	Server   ServerConfig   `mapstructure:"server"`
}

func Load(path string) (*Config, error) {
	config := new(Config)

	viper.AddConfigPath(path)

	viper.SetConfigName("config")
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
