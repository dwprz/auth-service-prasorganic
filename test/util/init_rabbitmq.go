package util

import (
	"github.com/dwprz/prasorganic-auth-service/src/core/broker/producer"
	"github.com/dwprz/prasorganic-auth-service/src/mock/delivery"
)

func InitRabbitMQ() (*producer.RabbitMQ, *delivery.EmailBrokerMock) {
	emailBrokerDelivery := delivery.NewEmailBrokerMock()
	rabbitMQClient := producer.NewRabbitMQ(emailBrokerDelivery)

	return rabbitMQClient, emailBrokerDelivery
}
