package config

const (
	RabbitMQDriver = "rabbitmq"
)

type BusConfig struct {
	Driver   string `mapstructure:"driver"`
	Exchange string `mapstructure:"exchange"`
}
