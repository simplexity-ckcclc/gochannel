package click

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/common/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v6"
	"strconv"
	"time"
)

type ClickPorter struct {
	db       *sql.DB
	esClient *elastic.Client
}

func NewClickPorter(database *sql.DB) (*ClickPorter, error) {
	client, err := elastic.NewClient(elastic.SetURL(config.GetString(config.EsServer)))
	if err != nil {
		return nil, err
	}

	return &ClickPorter{
		db:       database,
		esClient: client,
	}, nil
}

func (porter ClickPorter) TransferClicks() {
	esDeviceIndex := config.GetString(config.EsClickIndex)
	for {
		devices, err := porter.getClickInfos(config.GetInt(config.EsClickBatchSize))
		if err != nil {
			common.ApiLogger.WithFields(logrus.Fields{}).Error("Get click infos from db error : ", err)
		}

		if len(devices) > 0 {
			porter.putClickIntoEs(devices, esDeviceIndex)
		}

		time.Sleep(10 * time.Second)

	}
}

func (porter ClickPorter) getClickInfos(limit int) ([]ClickInfo, error) {
	rows, err := porter.db.Query(`SELECT id, app_key, channel_id, os_type, device_id, click_time FROM click_info limit ` +
		strconv.Itoa(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clickInfos []ClickInfo
	for rows.Next() {
		click := ClickInfo{}
		if err = rows.Scan(&click.Id, &click.AppKey, &click.ChannelId, &click.OsType, &click.DeviceId, &click.ClickTime); err != nil {
			return nil, err
		}
		clickInfos = append(clickInfos, click)
	}
	err = rows.Err()
	return clickInfos, err
}

func (porter ClickPorter) putClickIntoEs(clickInfos []ClickInfo, index string) {
	bulkRequest := porter.esClient.Bulk()
	for _, click := range clickInfos {
		clickJson, err := json.Marshal(click)
		if err != nil {
			common.ApiLogger.WithFields(logrus.Fields{
				"click": click,
			}).Error("JSON marshal click error : ", err)
			continue
		}

		req := elastic.NewBulkIndexRequest().
			Index(index).
			Type(click.AppKey).
			Doc(string(clickJson))
		bulkRequest.Add(req)
	}

	bulkResponse, err := bulkRequest.Do(context.Background())
	if err != nil {
		common.ApiLogger.WithFields(logrus.Fields{
			"clickInfos": clickInfos,
		}).Error("Bulk put click doc error : ", err)
		return
	}

	failed := bulkResponse.Failed()
	for _, failedResp := range failed {
		common.ApiLogger.WithFields(logrus.Fields{
			"id":       failedResp.Id,
			"errCause": failedResp.Error,
		}).Error("Bulk put click doc error")
	}
}
