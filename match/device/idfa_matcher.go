package device

import (
	"context"
	"encoding/json"
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
		Must(elastic.NewTermQuery("device_id.keyword", device.Idfa)).
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

	var click ClickInfo
	recentHit := searchHits.Hits[0]
	if err = json.Unmarshal(*recentHit.Source, &click); err != nil {
		return
	}

	clickTime := click.ClickTime
	if device.ActivateTime-clickTime > config.GetInt64(config.ActivateValidPeriod) {
		return
	}

	device.MatchInfo.IsMatched = true
	device.MatchInfo.ClickTime = clickTime
	device.MatchInfo.Channel = click.ChannelId
	return
}

type ClickInfo struct {
	AppKey    string `json:"app_key"`
	ChannelId string `json:"channel_id"`
	DeviceId  string `json:"device_id"` // idfa or imei
	ClickTime int64  `json:"click_time" form:"clickTime" binding:"required"`
}
