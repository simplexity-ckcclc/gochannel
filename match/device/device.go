package device

type OsType string

const (
	IOS     OsType = "ios"
	ANDROID OsType = "android"
)

type Device struct {
	Id           int64  `json:"id"`
	Imei         string `json:"imei"`
	Idfa         string `json:"idfa"`
	ActivateTime int64  `json:"activate_time"`
	Channel      string `json:"channel"`
	AppKey       string `json:"app_key"`
	OsType       OsType `json:"os_type"`
	OsVersion    string `json:"os_version"`
	Language     string `json:"language"`
	Resolution   string `json:"resolution"`
	SourceIp     string `json:"source_ip"`
}

type MatchedDevice struct {
	*Device
	MatchedChannel string
	ClickTime      int64
}
