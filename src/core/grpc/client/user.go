package client

import (
	"context"
	"fmt"
	"log"

	"github.com/dwprz/prasorganic-auth-service/src/interface/client"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/interceptor"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
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

func (u *UserGrpcImpl) FindByEmail(ctx context.Context, data *pb.Email) (*pb.FindUserResponse, error) {
	res, err := u.cbreaker.Execute(func() (any, error) {
		res, err := u.client.FindByEmail(ctx, data)
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

func (u *UserGrpcImpl) Create(ctx context.Context, data *pb.RegisterRequest) error {
	_, err := u.cbreaker.Execute(func() (any, error) {
		_, err := u.client.Create(ctx, data)
		return nil, err
	})

	return err
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
