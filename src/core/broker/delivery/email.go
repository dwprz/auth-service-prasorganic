package delivery

import (
	"encoding/json"

	"github.com/dwprz/prasorganic-auth-service/src/common/log"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-auth-service/src/interface/delivery"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type EmailBrokerImpl struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewEmailBroker() delivery.EmailBroker {
	conn, err := amqp.Dial(config.Conf.RabbitMQEmailService.DSN)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "delivery.NewEmailBroker", "section": "amqp.Dial"}).Fatal(err)
	}

	chann, err := conn.Channel()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "delivery.NewEmailBroker", "section": "conn.Channel"}).Error(err)
	}

	return &EmailBrokerImpl{
		connection: conn,
		channel:    chann,
	}
}

func (e *EmailBrokerImpl) Publish(exchange string, key string, message any) {
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "delivery.EmailBrokerImpl/Publish", "section": "json.Marshal"}).Error(err)
		return
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(jsonData),
	}

	if err := e.channel.Publish(exchange, key, false, false, msg); err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "delivery.EmailBrokerImpl/Publish", "section": "channel.PublishWithContext"}).Error(err)
	}
}

func (e *EmailBrokerImpl) Close() {
	if err := e.channel.Close(); err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "delivery.EmailBrokerImpl/Close", "section": "channel.Close"}).Error(err)
	}

	if err := e.connection.Close(); err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "delivery.EmailBrokerImpl/Close", "section": "connection.Close"}).Error(err)
	}
}
