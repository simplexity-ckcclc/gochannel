package common

import "strings"

type ChannelType string

const (
	UnknownChannelType ChannelType = "unknown"
	IOSChannelType     ChannelType = "ios"
	AndroidChannelType ChannelType = "android"
)

func (ct ChannelType) String() string {
	return string(ct)
}

func ParseChannelType(ct string) ChannelType {
	switch strings.ToLower(ct) {
	case "ios":
		return IOSChannelType
	case "android":
		return AndroidChannelType
	default:
		return UnknownChannelType

	}
}

type AppChannel struct {
	AppKey      string      `json:"app_key"`
	ChannelId   string      `json:"channel"`
	ChannelType ChannelType `json:"channel_type"`
	PublicKey   string      `json:"pub_key"`
	PrivateKey  string      `json:"pri_key"`
	CallbackUrl string      `json:"-"`
}
