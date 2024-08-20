package broker

import (
	"github.com/dwprz/prasorganic-auth-service/src/core/broker/delivery"
	"github.com/dwprz/prasorganic-auth-service/src/core/broker/producer"
)

func InitClient() *producer.RabbitMQ {
	emailBrokerDelivery := delivery.NewEmailBroker()
	rabbitMQClient := producer.NewRabbitMQ(emailBrokerDelivery)

	return rabbitMQClient
}
