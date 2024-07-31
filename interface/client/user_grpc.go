package client

import (
	"context"
	"github.com/dwprz/prasorganic-proto/protogen/user"
)

type UserGrpc interface {
	FindByEmail(ctx context.Context, data *user.Email) (*user.FindUserResponse, error)
	Create(ctx context.Context, data *user.RegisterRequest) error
}
