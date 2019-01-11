package entity

import "database/sql"

type ClickInfo struct {
	AppKey   string `form:"appKey" binding:"required"`
	DeviceId string `form:"deviceId" binding:"required"`
}

func InsertClickInfo(db *sql.DB, click ClickInfo) error {
	stmt, err := db.Prepare("INSERT INTO click_info (app_key, device_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(click.AppKey, click.DeviceId)
	return err
}
