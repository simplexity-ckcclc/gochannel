package device

import (
	"github.com/simplexity-ckcclc/gochannel/common"
	"gopkg.in/olivere/elastic.v6"
)

type instantiateMatcherFunc func(*elastic.Client) Matcher

var matcherMappings = make(map[common.ChannelType]instantiateMatcherFunc)

func init() {
	matcherMappings[common.IOSChannelType] = NewIdfaMatcher
}

type Matcher interface {
	Match(device *Device) error
}
