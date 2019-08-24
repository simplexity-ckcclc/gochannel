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
	ReceiveTime int64  `json:"receiveTime"`
	Channel     string `json:"channel"`
	AppKey      string `json:"appKey"`
	OsType      OsType `json:"osType"`
	OsVersion   string `json:"osVersion"`
	Language    string `json:"language"`
	Resolution  string `json:"resolution"`
	SourceIp    string `json:"sourceIp"`
}
