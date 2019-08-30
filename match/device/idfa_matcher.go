package device

import (
	"context"
	"github.com/simplexity-ckcclc/gochannel/common/config"
	"gopkg.in/olivere/elastic.v6"
)

type IdfaMatcher struct {
	esClient *elastic.Client
}

func NewIdfaMatcher(client *elastic.Client) Matcher {
	return &IdfaMatcher{
		esClient: client,
	}
}

func (matcher IdfaMatcher) Match(device *Device) (err error) {
	index := config.GetString(config.EsClickIndex)
	query := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("app_key.keyword", device.AppKey)).
		Must(elastic.NewTermQuery("device_id.keyword", device.Imei)).
		Must(elastic.NewRangeQuery("click_time").Lt(device.ActivateTime).Gt(device.MatchInfo.ClickTime))

	esResponse, esErr := matcher.esClient.Search().
		Index(index).
		Type(device.AppKey).
		Query(query).
		Sort("click_time", false).
		Do(context.Background())
	if esErr != nil {
		err = esErr
		return
	}

	searchHits := esResponse.Hits
	if searchHits.TotalHits <= 0 {
		return
	}

	recentHit := searchHits.Hits[0]
	clickTime := recentHit.Fields["click_time"].(int64)
	if device.ActivateTime-clickTime > config.GetInt64(config.ActivateValidPeriod) {
		return
	}
	channel := recentHit.Fields["channel_id"].(string)

	device.MatchInfo.IsMatched = true
	device.MatchInfo.ClickTime = clickTime
	device.MatchInfo.Channel = channel
	return
}
