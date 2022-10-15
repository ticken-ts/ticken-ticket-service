package config

const (
	MongoDriver = "mongo"
)

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	Name   string `mapstructure:"name"`
}
