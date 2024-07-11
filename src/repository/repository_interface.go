package repository

import (
	"context"
	"github.com/dwprz/prasorganic-auth-service/src/common/model/entity"
)

type AuthRepository interface {
	FindCredentialByEmail(ctx context.Context, email string) *entity.Credential
}
