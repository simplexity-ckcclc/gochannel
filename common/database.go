package common

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/simplexity-ckcclc/gochannel/common/config"
	"time"
)

var (
	DB *sql.DB // todo should adopt an elegant way to share this var
)

func OpenDB() (*sql.DB, error) {
	var err error
	DB, err = sql.Open("mysql", config.GetString(config.DatabaseDsn))
	if err != nil {
		return nil, err
	}

	DB.SetMaxOpenConns(config.GetInt(config.DatabaseMaxOpenCons))
	DB.SetMaxIdleConns(config.GetInt(config.DatabaseMaxIdleCons))
	DB.SetConnMaxLifetime(1 * time.Minute)
	return DB, err
}
