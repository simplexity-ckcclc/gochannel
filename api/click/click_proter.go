package click

import (
	"context"
	"database/sql"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/common/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v6"
	"strconv"
	"strings"
	"time"
)

// dump mysql to es. Can be replace by third-party open-source tool [go-mysql-elasticsearch](https://github.com/siddontang/go-mysql-elasticsearch)
type ClickPorter struct {
	db       *sql.DB
	esClient *elastic.Client
}

func NewClickPorter(database *sql.DB, client *elastic.Client) *ClickPorter {
	return &ClickPorter{
		db:       database,
		esClient: client,
	}
}

func (porter ClickPorter) TransferClicks() {
	esClickIndex := config.GetString(config.EsClickIndex)
	for {
		clicks, err := porter.getClickInfos(config.GetInt(config.EsClickBatchSize))
		if err != nil {
			common.ApiLogger.WithFields(logrus.Fields{}).Error("Get click infos from db error : ", err)
		}

		if len(clicks) > 0 {
			if err = porter.putClickIntoEs(clicks, esClickIndex); err == nil {
				if err = porter.deleteClicks(clicks); err != nil {
					common.ApiLogger.Error("Delete click infos error. ", err)
				}
			}
		}

		time.Sleep(10 * time.Second)

	}
}

func (porter ClickPorter) getClickInfos(limit int) ([]ClickInfo, error) {
	rows, err := porter.db.Query(`SELECT id, app_key, channel_id, os_type, device_id, click_time FROM click_info limit ?`,
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

func (porter ClickPorter) putClickIntoEs(clickInfos []ClickInfo, index string) error {
	bulkRequest := porter.esClient.Bulk()
	for _, click := range clickInfos {
		req := elastic.NewBulkIndexRequest().
			Index(index).
			Type(click.AppKey).
			Doc(click)
		bulkRequest.Add(req)
	}

	bulkResponse, err := bulkRequest.Do(context.Background())
	if err != nil {
		common.ApiLogger.WithFields(logrus.Fields{
			"clickInfos": clickInfos,
		}).Error("Bulk put click doc error : ", err)
		return err
	}

	failed := bulkResponse.Failed()
	for _, failedResp := range failed {
		common.ApiLogger.WithFields(logrus.Fields{
			"id":       failedResp.Id,
			"errCause": failedResp.Error,
		}).Error("Bulk put click doc error")
	}
	return nil
}

func (porter ClickPorter) deleteClicks(clicks []ClickInfo) error {
	var ids []string
	for _, click := range clicks {
		ids = append(ids, strconv.Itoa(int(click.Id)))
	}

	s := strings.Join(ids, ",")
	_, err := porter.db.Exec(`DELETE FROM click_info WHERE id in (` + s + `)`)
	return err
}
