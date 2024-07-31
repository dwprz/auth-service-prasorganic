package client

import (
	"context"
)

type RabbitMQMock struct{}

func NewRabbitMQMock() *RabbitMQMock {
	return &RabbitMQMock{}
}

func (r *RabbitMQMock) Publish(ctx context.Context, exchange string, key string, message any) {}

func (r *RabbitMQMock) Close() {}
