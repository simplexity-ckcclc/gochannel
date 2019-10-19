package common

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/simplexity-ckcclc/gochannel/common/config"
	"gopkg.in/olivere/elastic.v6"
	"time"
)

const _defaultMaxLifeTimeSecs = 30

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
	DB.SetConnMaxLifetime(_defaultMaxLifeTimeSecs * time.Second)
	return DB, err
}

func InitEsClient() (*elastic.Client, error) {
	var err error
	EsClient, err = elastic.NewClient(elastic.SetURL(config.GetString(config.EsServer)))
	return EsClient, err
}
