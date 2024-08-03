package client

import (
	"context"
	"fmt"
	"log"

	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/interceptor"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-auth-service/src/interface/client"
	pb "github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserGrpcImpl struct {
	client   pb.UserServiceClient
	cbreaker *gobreaker.CircuitBreaker[any]
}

func NewUserGrpc(cb *gobreaker.CircuitBreaker[any], conf *config.Config, unaryRequest *interceptor.UnaryRequest) (client.UserGrpc, *grpc.ClientConn) {
	var opts []grpc.DialOption
	opts = append(
		opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(unaryRequest.AddBasicAuth),
	)

	conn, err := grpc.NewClient(conf.ApiGateway.BaseUrl, opts...)
	if err != nil {
		log.Fatalf("new user grpc client: %v", err.Error())
	}

	client := pb.NewUserServiceClient(conn)

	return &UserGrpcImpl{
		client:   client,
		cbreaker: cb,
	}, conn
}

func (u *UserGrpcImpl) Create(ctx context.Context, data *pb.RegisterRequest) error {
	_, err := u.cbreaker.Execute(func() (any, error) {
		_, err := u.client.Create(ctx, data)
		return nil, err
	})

	return err
}

func (u *UserGrpcImpl) FindByEmail(ctx context.Context, email string) (*pb.FindUserResponse, error) {
	res, err := u.cbreaker.Execute(func() (any, error) {
		res, err := u.client.FindByEmail(ctx, &pb.Email{Email: email})
		return res, err
	})

	if err != nil {
		return nil, err
	}

	user, ok := res.(*pb.FindUserResponse)
	if !ok {
		return nil, fmt.Errorf("client.UserGrpcImpl/FindByEmail | unexpected type: %T", res)
	}

	return user, err
}

func (u *UserGrpcImpl) FindByRefreshToken(ctx context.Context, data *pb.RefreshToken) (*pb.FindUserResponse, error) {
	res, err := u.cbreaker.Execute(func() (any, error) {
		res, err := u.client.FindByRefreshToken(ctx, &pb.RefreshToken{
			Token: data.Token,
		})
		return res, err
	})

	user, ok := res.(*pb.FindUserResponse)
	if !ok {
		return nil, fmt.Errorf("client.UserGrpcImpl/FindByRefreshToken | unexpected type: %T", res)
	}

	return user, err
}

func (u *UserGrpcImpl) Upsert(ctx context.Context, data *pb.LoginWithGoogleRequest) (*pb.User, error) {
	res, err := u.cbreaker.Execute(func() (any, error) {
		res, err := u.client.Upsert(ctx, data)
		return res, err
	})

	if err != nil {
		return nil, err
	}

	user, ok := res.(*pb.User)
	if !ok {
		return nil, fmt.Errorf("client.UserGrpcImpl/FindByEmail | unexpected type: %T", res)
	}

	return user, err
}

func (u *UserGrpcImpl) AddRefreshToken(ctx context.Context, data *pb.AddRefreshToken) error {
	_, err := u.cbreaker.Execute(func() (any, error) {
		_, err := u.client.AddRefreshToken(ctx, data)
		return nil, err
	})

	return err
}

func (u *UserGrpcImpl) SetNullRefreshToken(ctx context.Context, refreshToken string) error {
	_, err := u.cbreaker.Execute(func() (any, error) {
		_, err := u.client.SetNullRefreshToken(ctx, &pb.RefreshToken{
			Token: refreshToken,
		})
		return nil, err
	})

	return err
}
