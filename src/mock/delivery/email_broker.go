package delivery

import "github.com/stretchr/testify/mock"

type EmailBrokerMock struct {
	mock.Mock
}

func NewEmailBrokerMock() *EmailBrokerMock {
	return &EmailBrokerMock{
		Mock: mock.Mock{},
	}
}

func (o *EmailBrokerMock) Publish(exchange string, key string, message any) {}

func (o *EmailBrokerMock) Close() {}
