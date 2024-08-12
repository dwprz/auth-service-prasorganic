package broker

import (
	"github.com/dwprz/prasorganic-auth-service/src/core/broker/client"
	"github.com/dwprz/prasorganic-auth-service/src/core/broker/delivery"
)

func InitClient() *client.RabbitMQ {
	emailBrokerDelivery := delivery.NewEmailBroker()
	rabbitMQClient := client.NewRabbitMQ(emailBrokerDelivery)

	return rabbitMQClient
}
