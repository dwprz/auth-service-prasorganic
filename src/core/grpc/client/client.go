package client

import (
	"github.com/dwprz/prasorganic-auth-service/src/common/log"
	"github.com/dwprz/prasorganic-auth-service/src/interface/delivery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// this main grpc client
type Grpc struct {
	User     delivery.UserGrpc
	userConn *grpc.ClientConn
}

func NewGrpc(ugd delivery.UserGrpc, userConn *grpc.ClientConn) *Grpc {

	return &Grpc{
		User:     ugd,
		userConn: userConn,
	}
}

func (g *Grpc) Close() {
	if err := g.userConn.Close(); err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "client.Grpc/Close", "section": "userConn.Close"}).Error(err)
	}
}
