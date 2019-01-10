package api

import "database/sql"

type clickInfo struct {
	AppKey   string `form:"appKey" binding:"required"`
	DeviceId string `form:"deviceId" binding:"required"`
}

func insertClickInfo(db *sql.DB, click clickInfo) error {
	stmt, err := db.Prepare("INSERT INTO click_info (app_key, device_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(click.AppKey, click.DeviceId)
	return err
}
