package common

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	DB *sql.DB
)

func OpenDB(DSN string) (*sql.DB, error) {
	var err error
	DB, err = sql.Open("mysql", DSN)
	return DB, err
}
