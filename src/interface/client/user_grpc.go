package client

import (
	"context"
	pb "github.com/dwprz/prasorganic-proto/protogen/user"
)

type UserGrpc interface {
	Create(ctx context.Context, data *pb.RegisterRequest) error
	FindByEmail(ctx context.Context, email string) (*pb.FindUserResponse, error)
	FindByRefreshToken(ctx context.Context, data *pb.RefreshToken) (*pb.FindUserResponse, error)
	Upsert(ctx context.Context, data *pb.LoginWithGoogleRequest) (*pb.User, error)
	AddRefreshToken(ctx context.Context, data *pb.AddRefreshToken) error
	SetNullRefreshToken(ctx context.Context, refreshToken string) error 
}
