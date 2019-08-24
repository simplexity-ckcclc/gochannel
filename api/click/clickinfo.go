package click

import (
	"database/sql"
	"encoding/json"
)

type ClickInfo struct {
	Id        int64  `json:"id"`
	AppKey    string `json:"appKey" form:"appKey" binding:"required"`
	ChannelId string `json:"channelId" form:"channelId" binding:"required"`
	DeviceId  string `json:"deviceId" form:"deviceId" binding:"required"` // idfa or imei
	OsType    string `json:"osType" form:"osType" binding:"required"`     // idfa or imei
	ClickTime int64  `json:"clickTime" form:"clickTime" binding:"required"`
}

func (clickInfo *ClickInfo) InsertDB(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO click_info (app_key, channel_id, os_type, device_id, click_time) VALUES (?, ?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(clickInfo.AppKey, clickInfo.ChannelId, clickInfo.OsType, clickInfo.DeviceId, clickInfo.ClickTime)
	return err
}

func (clickInfo ClickInfo) String() string {
	jsonStr, err := json.Marshal(clickInfo)
	if err != nil {
		return ""
	}
	return string(jsonStr)
}
