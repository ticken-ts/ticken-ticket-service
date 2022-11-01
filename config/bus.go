package config

const (
	RabbitMQDriver = "rabbitmq"
	DevBusDriver   = "devbus"
)

type BusConfig struct {
	Driver   string `mapstructure:"driver"`
	Exchange string `mapstructure:"exchange"`
}
