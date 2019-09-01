package common

import "strings"

type OsType string

const (
	IOS     OsType = "ios"
	Android OsType = "android"
	Unknown OsType = "unknown"
)

func (ot OsType) String() string {
	switch ot {
	case IOS:
		return "ios"
	case Android:
		return "android"
	default:
		return ""
	}
}

func ParseOsType(ot string) OsType {
	switch strings.ToLower(ot) {
	case "ios":
		return IOS
	case "android":
		return Android
	default:
		return Unknown

	}
}
