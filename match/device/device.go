package device

type OsType string

const (
	IOS     OsType = "ios"
	ANDROID OsType = "android"
)

type Device struct {
	Id          int64  `json:"id"`
	Imei        string `json:"imei"`
	Idfa        string `json:"idfa"`
	ReceiveTime int64  `json:"receive_time"`
	Channel     string `json:"channel"`
	AppKey      string `json:"app_key"`
	OsType      OsType `json:"os_type"`
	OsVersion   string `json:"os_version"`
	Language    string `json:"language"`
	Resolution  string `json:"resolution"`
	SourceIp    string `json:"source_ip"`
}
