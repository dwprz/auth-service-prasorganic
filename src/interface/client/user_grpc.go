package client

import (
	"context"
	pb "github.com/dwprz/prasorganic-proto/protogen/user"
)

type UserGrpc interface {
	FindByEmail(ctx context.Context, email string) (*pb.FindUserResponse, error)
	Create(ctx context.Context, data *pb.RegisterRequest) (error)
	Upsert(ctx context.Context, data *pb.LoginWithGoogleRequest) (*pb.User, error)
	UpdateRefreshToken(ctx context.Context, data *pb.RefreshToken) (error)
}
