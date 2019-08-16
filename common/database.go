package common

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	DB *sql.DB // todo should adopt an elegant way to share this var
)

func OpenDB(conf ConfYaml) (*sql.DB, error) {
	var err error
	DB, err = sql.Open("mysql", conf.Core.Database.DSN)
	DB.SetMaxOpenConns(conf.Core.Database.MaxOpenConns)
	DB.SetMaxIdleConns(conf.Core.Database.MaxIdleConns)
	DB.SetConnMaxLifetime(1 * time.Minute)
	return DB, err
}
