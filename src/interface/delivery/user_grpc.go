package delivery

import (
	"context"
	pb "github.com/dwprz/prasorganic-proto/protogen/user"
)

type UserGrpc interface {
	Create(ctx context.Context, data *pb.RegisterReq) error
	FindByEmail(ctx context.Context, email string) (*pb.FindUserRes, error)
	FindByRefreshToken(ctx context.Context, data *pb.RefreshToken) (*pb.FindUserRes, error)
	Upsert(ctx context.Context, data *pb.LoginWithGoogleReq) (*pb.User, error)
	AddRefreshToken(ctx context.Context, data *pb.AddRefreshTokenReq) error
	SetNullRefreshToken(ctx context.Context, refreshToken string) error 
}
