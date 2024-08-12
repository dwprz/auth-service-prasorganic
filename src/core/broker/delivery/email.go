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
}

func NewEmailBroker() delivery.EmailBroker {
	conn, err := amqp.Dial(config.Conf.RabbitMQEmailService.DSN)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "delivery.NewRabbitMQClient", "section": "amqp.Dial"}).Fatal(err)
	}

	return &EmailBrokerImpl{
		connection: conn,
	}
}

func (r *EmailBrokerImpl) Publish(exchange string, key string, message any) {
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "delivery.EmailBrokerImpl/Publish", "section": "json.Marshal"}).Error(err)
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(jsonData),
	}

	channel, err := r.connection.Channel()
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "delivery.EmailBrokerImpl/Publish", "section": "connection.Channel"}).Error(err)
	}

	defer channel.Close()

	if err := channel.Publish(exchange, key, false, false, msg); err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "delivery.EmailBrokerImpl/Publish", "section": "channel.PublishWithContext"}).Error(err)
	}
}

func (r *EmailBrokerImpl) Close() {
	if err := r.connection.Close(); err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "delivery.EmailBrokerImpl/Close", "section": "connection.Close"}).Error(err)
	}
}
