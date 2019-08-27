package common

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/simplexity-ckcclc/gochannel/common/config"
	"gopkg.in/olivere/elastic.v6"
	"time"
)

var (
	DB       *sql.DB // todo should adopt an elegant way to share this var
	EsClient *elastic.Client
)

func InitSqlClient() (*sql.DB, error) {
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

func InitEsClient() (*elastic.Client, error) {
	var err error
	EsClient, err = elastic.NewClient(elastic.SetURL(config.GetString(config.EsServer)))
	return EsClient, err
}
