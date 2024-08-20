package delivery

type EmailBroker interface {
	Publish(exchange string, key string, message any)
	Close()
}
