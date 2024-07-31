package broker

import (
	"context"
	"encoding/json"
	"github.com/dwprz/prasorganic-auth-service/interface/client"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type RabbitMQClientImpl struct {
	conf       *config.Config
	connection *amqp.Connection
	logger     *logrus.Logger
}

func NewRabbitMQClient(conf *config.Config, logger *logrus.Logger) client.RabbitMQ {
	conn, err := amqp.Dial(conf.RabbitMQEmailService.DSN)
	if err != nil {
		logger.WithFields(logrus.Fields{"location": "broker.NewRabbitMQClient", "section": "amqp.Dial"}).Fatal(err)
	}

	return &RabbitMQClientImpl{
		conf:       conf,
		connection: conn,
	}
}

func (r *RabbitMQClientImpl) Publish(ctx context.Context, exchange string, key string, message any) {
	jsonData, err := json.Marshal(message)
	if err != nil {
		r.logger.WithFields(logrus.Fields{"location": "broker.RabbitMQClientImpl/Publish", "section": "json.Marshal"}).Error(err)
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(jsonData),
	}

	channel, err := r.connection.Channel()
	if err != nil {
		r.logger.WithFields(logrus.Fields{"location": "broker.RabbitMQClientImpl/Publish", "section": "connection.Channel"}).Error(err)
	}

	defer channel.Close()

	if err := channel.PublishWithContext(ctx, exchange, key, false, false, msg); err != nil {
		r.logger.WithFields(logrus.Fields{"location": "broker.RabbitMQClientImpl/Publish", "section": "channel.PublishWithContext"}).Error(err)
	}
}

func (r *RabbitMQClientImpl) Close() {
	if err := r.connection.Close(); err != nil {
		r.logger.WithFields(logrus.Fields{"location": "broker.RabbitMQClientImpl/Close", "section": "connection.Close"}).Error(err)
	}
}
