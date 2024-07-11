package database

import (
	"github.com/dwprz/prasorganic-auth-service/src/common/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	lg "gorm.io/gorm/logger"
)

func NewPostgres(conf *config.Config, logger *logrus.Logger) *gorm.DB {

	db, err := gorm.Open(postgres.Open(conf.Postgres.Dsn), &gorm.Config{
		Logger: lg.Default.LogMode(lg.Info),
	})
	if err != nil {
		logger.Errorf("database (gorm): %s", err.Error())
	}

	return db
}
