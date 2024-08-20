package producer

import "github.com/dwprz/prasorganic-auth-service/src/interface/delivery"

type RabbitMQ struct {
	Email delivery.EmailBroker
}

func NewRabbitMQ(email delivery.EmailBroker) *RabbitMQ {
	return &RabbitMQ{
		Email: email,
	}
}
