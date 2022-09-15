package utils

import "github.com/spf13/viper"

type TickenConfig struct {
	Config *Config
	Env    *Env
}

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
}

type Env struct {
	TickenEnv string `mapstructure:"TICKEN_ENV"`
	MongoUri  string `mapstructure:"MONGO_URI"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	Name   string `mapstructure:"name"`
}

const (
	devEnv  = "dev"
	prodEnv = "prod"
	testEnv = "test"
)

const (
	MongoDriver = "mongo"
)

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

func (config *TickenConfig) IsDev() bool {
	return config.Env.TickenEnv == devEnv
}

func (config *TickenConfig) IsProd() bool {
	return config.Env.TickenEnv == prodEnv
}

func (config *TickenConfig) IsTest() bool {
	return config.Env.TickenEnv == testEnv
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
