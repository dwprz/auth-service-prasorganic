package grpc

import (
	"github.com/dwprz/prasorganic-auth-service/src/interface/client"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// this main grpc client
type Client struct {
	User      client.UserGrpc
	userConn  *grpc.ClientConn
	logger    *logrus.Logger
}

func NewClient(ugc client.UserGrpc, userConn *grpc.ClientConn, l *logrus.Logger) *Client {

	return &Client{
		User:      ugc,
		userConn:  userConn,
		logger:    l,
	}
}

func (g *Client) Close() {
	if err := g.userConn.Close(); err != nil {
		g.logger.WithFields(logrus.Fields{"location": "grpc.Client/Close", "section": "userConn.Close"}).Errorf(err.Error())
	}
}
