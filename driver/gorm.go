package driver

import (
	_ "github.com/lib/pq"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/assyatier21/learn-distributed-trx/config"
)

func NewGormDatabase(cfg config.DBConfig) (db *gorm.DB) {
	dsn := cfg.GetDSN()
	var dialector gorm.Dialector

	if cfg.Driver == "mysql" {
		dialector = mysql.Open(dsn)
	} else if cfg.Driver == "postgres" {
		dialector = postgres.Open(dsn)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("db connection failed")
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)

	if cfg.DebugMode {
		db = db.Debug()
	}

	return
}
