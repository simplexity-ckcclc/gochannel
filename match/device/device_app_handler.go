package device

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/common/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v6"
	"time"
)

type DeviceStatus int

const Processed DeviceStatus = iota

type DeviceAppHandler struct {
	appKey   string
	esClient *elastic.Client
	stopChan chan bool
	matchers []Matcher
}

func (handler DeviceAppHandler) start() {
runningLoop:
	for {
		select {
		case <-handler.stopChan:
			common.MatchLogger.WithFields(logrus.Fields{
				"appKey": handler.appKey,
			}).Info("Device App Handler stop")
			break runningLoop
		default:
			latestProcessTime := latestProcessTime(handler.appKey)
			if latestProcessTime == 0 {
				// New appKey, process the last 3 days
				latestProcessTime = time.Now().Add(-3*24*time.Hour).UnixNano() / int64(time.Millisecond)
			}

			devices := handler.getDevices(latestProcessTime, time.Now().UnixNano()/int64(time.Millisecond), 5)
			fmt.Println(devices)

			//fmt.Println(latestProcessTime)
			//todo start process
		}
	}
}

func (handler DeviceAppHandler) stop() {
	handler.stopChan <- true
}

func latestProcessTime(appKey string) (processTime int64) {
	if err := common.DB.QueryRow("SELECT process_time FROM app_process_info WHERE app_key = ?", appKey).Scan(&processTime); err != nil && err != sql.ErrNoRows {
		common.MatchLogger.WithFields(logrus.Fields{
			"appKey": appKey,
		}).Error("Get app_key last process time error")
	}
	return
}

func updateLatestProcessTime(appKey string, processTime int64) error {
	stmt, err := common.DB.Prepare("INSERT INTO app_process_info (app_key, process_time) VALUES (?, ?) ON DUPLICATE KEY UPDATE process_time = ?")
	defer stmt.Close()
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(appKey, processTime, processTime); err != nil {
		return err
	}
	return nil
}

func (handler DeviceAppHandler) getDevices(startTime int64, endTime int64, batchSize int) []Device {
	index := config.GetString(config.EsDeviceIndex)
	query := elastic.NewBoolQuery().
		MustNot(elastic.NewTermQuery("status", Processed)).
		Must(elastic.NewTermQuery("app_key.keyword", handler.appKey))
		//Must(elastic.NewRangeQuery("activate_time").Gte(startTime).Lt(endTime))

	esResponse, _ := handler.esClient.Search().
		Index(index).
		Type(handler.appKey).
		Query(query).
		Sort("activate_time", true).
		Size(batchSize).
		Do(context.Background())

	var devices []Device
	var device Device
	for _, value := range esResponse.Hits.Hits {
		if err := json.Unmarshal(*value.Source, &device); err != nil {
			common.MatchLogger.WithFields(logrus.Fields{
				"value": value.Source,
			}).Error("Construct device from es error", err)
		}
		devices = append(devices, device)
	}
	return devices
}
