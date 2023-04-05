package config

const (
	RabbitMQDriver = "rabbitmq"
	DevBusDriver   = "devbus"
)

type BusConfig struct {
	Driver      string   `mapstructure:"driver"`
	Exchange    string   `mapstructure:"exchange"`
	SendQueues  []string `mapstructure:"send_queues"`
	ListenQueue string   `mapstructure:"listen_queue"`
}
