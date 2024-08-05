package client

type RabbitMQ interface {
	Publish(exchange string, key string, message any)
	Close()
}
