package restful

import (
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/handler"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/middleware"
	"github.com/dwprz/prasorganic-auth-service/src/core/restful/server"
	"github.com/dwprz/prasorganic-auth-service/src/interface/service"
)

func InitServer(as service.Auth) *server.Restful {
	authHandler := handler.NewAuthRESTful(as)
	middleware := middleware.New()

	restfulServer := server.NewRestful(authHandler, middleware)
	return restfulServer
}
