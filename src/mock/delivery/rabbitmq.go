package delivery

type RabbitMQMock struct{}

func NewRabbitMQMock() *RabbitMQMock {
	return &RabbitMQMock{}
}

func (r *RabbitMQMock) Publish(exchange string, key string, message any) {}

func (r *RabbitMQMock) Close() {}
