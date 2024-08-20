package delivery

import (
	"context"
	"fmt"

	"github.com/dwprz/prasorganic-auth-service/src/common/log"
	"github.com/dwprz/prasorganic-auth-service/src/core/grpc/interceptor"
	"github.com/dwprz/prasorganic-auth-service/src/infrastructure/config"
	"github.com/dwprz/prasorganic-auth-service/src/interface/delivery"
	pb "github.com/dwprz/prasorganic-proto/protogen/user"
	"github.com/sirupsen/logrus"
	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserGrpcImpl struct {
	client   pb.UserServiceClient
	cbreaker *gobreaker.CircuitBreaker[any]
}

func NewUserGrpc(CBreakerUserGrpc *gobreaker.CircuitBreaker[any], unaryRequest *interceptor.UnaryRequest) (delivery.UserGrpc, *grpc.ClientConn) {
	var opts []grpc.DialOption
	opts = append(
		opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(unaryRequest.AddBasicAuth),
	)

	conn, err := grpc.NewClient(config.Conf.ApiGateway.BaseUrl, opts...)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "delivery.NewUserGrpc", "section": "grpc.NewClient"}).Fatal(err)
	}

	client := pb.NewUserServiceClient(conn)

	return &UserGrpcImpl{
		client:   client,
		cbreaker: CBreakerUserGrpc,
	}, conn
}

func (u *UserGrpcImpl) Create(ctx context.Context, data *pb.RegisterReq) error {
	_, err := u.cbreaker.Execute(func() (any, error) {
		_, err := u.client.Create(ctx, data)
		return nil, err
	})

	return err
}

func (u *UserGrpcImpl) FindByEmail(ctx context.Context, email string) (*pb.FindUserRes, error) {
	res, err := u.cbreaker.Execute(func() (any, error) {
		res, err := u.client.FindByEmail(ctx, &pb.Email{Email: email})
		return res, err
	})

	if err != nil {
		return nil, err
	}

	user, ok := res.(*pb.FindUserRes)
	if !ok {
		return nil, fmt.Errorf("delivery.UserGrpcImpl/FindByEmail | unexpected type: %T", res)
	}

	return user, err
}

func (u *UserGrpcImpl) FindByRefreshToken(ctx context.Context, data *pb.RefreshToken) (*pb.FindUserRes, error) {
	res, err := u.cbreaker.Execute(func() (any, error) {
		res, err := u.client.FindByRefreshToken(ctx, &pb.RefreshToken{
			Token: data.Token,
		})
		return res, err
	})

	user, ok := res.(*pb.FindUserRes)
	if !ok {
		return nil, fmt.Errorf("delivery.UserGrpcImpl/FindByRefreshToken | unexpected type: %T", res)
	}

	return user, err
}

func (u *UserGrpcImpl) Upsert(ctx context.Context, data *pb.LoginWithGoogleReq) (*pb.User, error) {
	res, err := u.cbreaker.Execute(func() (any, error) {
		res, err := u.client.Upsert(ctx, data)
		return res, err
	})

	if err != nil {
		return nil, err
	}

	user, ok := res.(*pb.User)
	if !ok {
		return nil, fmt.Errorf("delivery.UserGrpcImpl/Upsert | unexpected type: %T", res)
	}

	return user, err
}

func (u *UserGrpcImpl) AddRefreshToken(ctx context.Context, data *pb.AddRefreshTokenReq) error {
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
