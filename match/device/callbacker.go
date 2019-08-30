package device

import (
	"database/sql"
	"strings"
)

type Callbacker struct {
	db *sql.DB
}

func (cb Callbacker) preHandle(devices []*Device) (err error) {
	sqlStr := "INSERT INTO callback_info(app_key, channel_id, device_id, os_type, click_time, activate_time) VALUES "
	vals := []interface{}{}

	for _, device := range devices {
		sqlStr += "(?, ?, ?),"
		vals = append(vals, device.AppKey, device.MatchInfo.Channel, device.OsType, device.MatchInfo.ClickTime, device.ActivateTime)
	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	stmt, _ := cb.db.Prepare(sqlStr)
	_, err = stmt.Exec(vals...)
	return
}
