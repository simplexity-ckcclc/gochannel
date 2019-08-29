package device

import (
	"context"
	"github.com/simplexity-ckcclc/gochannel/common/config"
	"gopkg.in/olivere/elastic.v6"
)

type IdfaMatcher struct {
	esClient *elastic.Client
}

func (matcher IdfaMatcher) match(device Device) (matched bool, matchedDevice *MatchedDevice, err error) {
	index := config.GetString(config.EsClickIndex)
	query := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("app_key.keyword", device.AppKey)).
		Must(elastic.NewTermQuery("device_id.keyword", device.Imei)).
		Must(elastic.NewRangeQuery("click_time").Lt(device.ActivateTime))

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

	matched = false
	searchHits := esResponse.Hits
	if searchHits.TotalHits <= 0 {
		return
	}

	recentHit := searchHits.Hits[0]
	clickTime := recentHit.Fields["click_time"].(int64)
	if device.ActivateTime-clickTime > config.GetInt64(config.ActivateValidPeriod) {
		return
	}

	matchedDevice = &MatchedDevice{
		Device:         &device,
		MatchedChannel: recentHit.Fields["channel"].(string),
		ClickTime:      clickTime,
	}
	return
}
