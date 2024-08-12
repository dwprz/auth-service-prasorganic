package util

import (
	"github.com/dwprz/prasorganic-auth-service/src/core/broker/client"
	"github.com/dwprz/prasorganic-auth-service/src/mock/delivery"
)

func InitRabbitMQ() (*client.RabbitMQ, *delivery.RabbitMQMock) {
	emailBrokerDelivery := delivery.NewRabbitMQMock()
	rabbitMQClient := client.NewRabbitMQ(emailBrokerDelivery)

	return rabbitMQClient, emailBrokerDelivery
}
