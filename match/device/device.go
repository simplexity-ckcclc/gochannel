package device

type OsType string

const (
	IOS     OsType = "ios"
	ANDROID OsType = "android"
)

type Device struct {
	Id          int64
	Imei        string
	Idfa        string
	ReceiveTime int64
	Channel     string
	AppKey      string
	OsType      OsType
	OsVersion   string
	Language    string
	Resolution  string
	SourceIp    string
}
