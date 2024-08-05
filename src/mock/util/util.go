package util

import "github.com/stretchr/testify/mock"

type UtilMock struct {
	mock.Mock
}

func NewMock() *UtilMock {
	return &UtilMock{
		Mock: mock.Mock{},
	}
}

func (u *UtilMock) GenerateOtp() (string, error) {
	arguments := u.Mock.Called()

	return arguments.String(0), arguments.Error(1)
}