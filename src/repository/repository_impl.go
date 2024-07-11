package repository

import (
	"context"

	"github.com/dwprz/prasorganic-auth-service/src/common/model/entity"
	"gorm.io/gorm"
)

type AuthRepositoryImpl struct {
	Db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &AuthRepositoryImpl{
		Db: db,
	}
}

func (r *AuthRepositoryImpl) FindCredentialByEmail(ctx context.Context, email string) *entity.Credential {
	credential := &entity.Credential{}

	ra := r.Db.WithContext(ctx).Raw("SELECT * FROM credentials WHERE email = ?;", email).Scan(credential).RowsAffected
	if ra == 0 {
		return nil
	}
	
	return credential
}
