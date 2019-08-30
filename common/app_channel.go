package common

import "strings"

type ChannelType int

const (
	UnknownChannelType ChannelType = iota
	IOSChannelType
	AndroidChannelType
)

func (ct ChannelType) String() string {
	switch ct {
	case IOSChannelType:
		return "ios"
	case AndroidChannelType:
		return "android"
	default:
		return ""
	}
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
	AppKey      string
	ChannelId   string
	ChannelType ChannelType
	PublicKey   string
	PrivateKey  string
}
