package device

import "github.com/simplexity-ckcclc/gochannel/common"

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
	MatchInfo    *MatchInfo    `json:"-"`
}

type MatchInfo struct {
	IsMatched bool
	Channel   string
	ClickTime int64
}

func (device Device) ResetMatchInfo() {
	device.MatchInfo = &MatchInfo{
		IsMatched: false,
		ClickTime: 0,
		Channel:   "",
	}
}
