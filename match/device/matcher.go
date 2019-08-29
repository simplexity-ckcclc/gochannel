package device

type Matcher interface {
	match(device Device) (bool, *MatchedDevice, error)
}
