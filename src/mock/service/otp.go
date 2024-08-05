package service

import (
	"context"

	"github.com/dwprz/prasorganic-auth-service/src/model/dto"
	"github.com/stretchr/testify/mock"
)

type OtpMock struct {
	mock.Mock
}

func NewOtpMock() *OtpMock {
	return &OtpMock{
		Mock: mock.Mock{},
	}
}

func (o *OtpMock) Send(ctx context.Context, email string) error {
	arguments := o.Mock.Called(ctx, email)

	return arguments.Error(0)
}

func (o *OtpMock) Verify(ctx context.Context, data *dto.VerifyOtpReq) error {
	arguments := o.Mock.Called(ctx, data)

	return arguments.Error(0)
}