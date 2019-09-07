package device

import "github.com/simplexity-ckcclc/gochannel/common"

type DeviceStatus int

const (
	New DeviceStatus = iota
	Processed
	Matched
)

type Device struct {
	Id           int64         `json:"id,omitempty"`
	Imei         string        `json:"imei,omitempty"`
	Idfa         string        `json:"idfa,omitempty"`
	ActivateTime int64         `json:"activate_time,omitempty"`
	Channel      string        `json:"channel,omitempty"`
	AppKey       string        `json:"app_key,omitempty"`
	OsType       common.OsType `json:"os_type,omitempty"`
	OsVersion    string        `json:"os_version,omitempty"`
	Language     string        `json:"language,omitempty"`
	Resolution   string        `json:"resolution,omitempty"`
	SourceIp     string        `json:"source_ip,omitempty"`
	MatchInfo    *MatchInfo    `json:"match_info"`
	ESId         string        `json:"es_id"`
	Status       DeviceStatus  `json:"status"`
}

type MatchInfo struct {
	Channel   string `json:"channel,omitempty"`
	ClickTime int64  `json:"click_time,omitempty"`
}

func (device *Device) ResetMatchInfo() {
	device.MatchInfo = &MatchInfo{
		ClickTime: 0,
		Channel:   "",
	}
}
