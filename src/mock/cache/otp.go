package cache

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

func (o *OtpMock) Cache(ctx context.Context, data *dto.SendOtpReq) {}

func (o *OtpMock) FindByEmail(ctx context.Context, email string) *dto.SendOtpReq {
	arguments := o.Mock.Called(ctx, email)

	return arguments.Get(0).(*dto.SendOtpReq)
}

func (o *OtpMock) DeleteByEmail(ctx context.Context, email string) {}
