package entity

import (
	"database/sql"
	"encoding/json"
)

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

func (clickInfo ClickInfo) String() string {
	jsonStr, err := json.Marshal(clickInfo)
	if err != nil {
		return ""
	}
	return string(jsonStr)
}
