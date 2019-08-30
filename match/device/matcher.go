package device

import (
	"github.com/simplexity-ckcclc/gochannel/common"
	"gopkg.in/olivere/elastic.v6"
)

type newMatcherFunc func(*elastic.Client) Matcher

var matcherMappings map[common.ChannelType]newMatcherFunc

func init() {
	matcherMappings[common.IOSChannelType] = NewIdfaMatcher
}

type Matcher interface {
	match(device Device) (bool, *MatchedDevice, error)
}
