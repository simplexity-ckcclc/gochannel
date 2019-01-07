package api

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func open(DSN string) (*sql.DB, error) {
	return sql.Open("mysql", DSN)
}

func insertClickInfo(db *sql.DB, click clickInfo) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO click_info (app_key, device_id) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(click.AppKey, click.DeviceId)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
