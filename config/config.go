package config

import "github.com/spf13/viper"

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Pvtbc    PvtbcConfig    `mapstructure:"pvtbc"`
	Server   ServerConfig   `mapstructure:"server"`
}

func LoadConfig(path string) (*TickenConfig, error) {
	tickenConfig := new(TickenConfig)

	env, err := loadEnv(path, ".env")
	if err != nil {
		return nil, err
	}

	config, err := loadConfig(path, "config")
	if err != nil {
		return nil, err
	}

	tickenConfig.Env = env
	tickenConfig.Config = config

	return tickenConfig, nil
}

func loadEnv(path string, filename string) (*Env, error) {
	viper.AddConfigPath(path)

	viper.SetConfigName(filename)
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	env := new(Env)
	err = viper.Unmarshal(env)
	if err != nil {
		return nil, err
	}

	return env, nil
}

func loadConfig(path string, filename string) (*Config, error) {
	viper.AddConfigPath(path)

	viper.SetConfigName(filename)
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := new(Config)
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
