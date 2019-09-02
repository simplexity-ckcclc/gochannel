package device

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/common/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v6"
	"strings"
	"time"
)

type DeviceStatus int

const Processed DeviceStatus = iota

type DeviceAppHandler struct {
	appKey   string
	esClient *elastic.Client
	stopChan chan bool
	// channelType(ios, android, gdt, etc.) - matcher
	matchers   map[common.ChannelType]*Matcher
	callbacker *Callbacker
}

func (handler *DeviceAppHandler) start() {
runningLoop:
	for {
		select {
		case <-handler.stopChan:
			common.MatchLogger.WithFields(logrus.Fields{
				"appKey": handler.appKey,
			}).Info("Device App Handler stop")
			break runningLoop
		default:
			latestProcessTime := getLatestProcessTime(handler.appKey)
			if latestProcessTime == 0 {
				// New appKey, process the last 3 days
				latestProcessTime = common.TimeToMillis(time.Now().AddDate(0, 0, -3))
			}
			processEndTime := latestProcessTime + config.GetInt64(config.ProcessPeriod)

			var matchedDevices []*Device
			devices := handler.getDevices(latestProcessTime, processEndTime, 5)
			if len(devices) > 0 {
				for _, device := range devices {
					device.ResetMatchInfo()
					for _, matcher := range handler.matchers {
						if err := (*matcher).Match(device); err != nil {
							common.MatchLogger.WithFields(logrus.Fields{
								"device": device,
							}).Error("match device error.", err)
						}
					}

					if device.MatchInfo.IsMatched {
						matchedDevices = append(matchedDevices, device)
					}
				}

				if len(matchedDevices) > 0 {
					if err := handler.callbacker.preHandle(matchedDevices); err != nil {
						common.MatchLogger.Error("PreHandle matched devices error.", err)
						continue
					}
				}

				if err := updateLatestProcessTime(handler.appKey, processEndTime); err != nil {
					common.MatchLogger.WithFields(logrus.Fields{
						"appKey":      handler.appKey,
						"processTime": processEndTime,
					}).Error("Update process time error.", err)
				}
			} else {
				common.MatchLogger.WithFields(logrus.Fields{
					"startTime": latestProcessTime,
					"endTime":   processEndTime,
					"appKey":    handler.appKey,
				}).Info("No sdk device activation.")

				if err := updateLatestProcessTime(handler.appKey, processEndTime); err != nil {
					common.MatchLogger.WithFields(logrus.Fields{
						"appKey":      handler.appKey,
						"processTime": processEndTime,
					}).Error("Update process time error.", err)
				}
				time.Sleep(30 * time.Second)
			}
		}
	}
}

func (handler *DeviceAppHandler) stop() {
	handler.stopChan <- true
}

func getLatestProcessTime(appKey string) (processTime int64) {
	if err := common.DB.QueryRow("SELECT process_time FROM app_process_info WHERE app_key = ?", appKey).
		Scan(&processTime); err != nil && err != sql.ErrNoRows {
		common.MatchLogger.WithFields(logrus.Fields{
			"appKey": appKey,
		}).Error("Get app_key last process time error")
	}
	return
}

func updateLatestProcessTime(appKey string, processTime int64) (err error) {
	_, err = common.DB.Exec("INSERT INTO app_process_info (app_key, process_time) VALUES (?, ?) ON DUPLICATE KEY UPDATE process_time = ?",
		appKey, processTime, processTime)
	return
}

func (handler *DeviceAppHandler) getDevices(startTime int64, endTime int64, batchSize int) (devices []*Device) {
	index := config.GetString(config.EsDeviceIndex)
	query := elastic.NewBoolQuery().
		MustNot(elastic.NewTermQuery("status", Processed)).
		Must(elastic.NewTermQuery("app_key.keyword", handler.appKey)).
		Must(elastic.NewRangeQuery("activate_time").Gte(startTime).Lt(endTime))

	esResponse, err := handler.esClient.Search().
		Index(index).
		Type(handler.appKey).
		Query(query).
		Sort("activate_time", true).
		Size(batchSize).
		Do(context.Background())
	if err != nil {
		common.MatchLogger.Error("Search es device error.", err)
		return
	}

	if esResponse.Hits.TotalHits <= 0 {
		return
	}

	var device *Device
	for _, value := range esResponse.Hits.Hits {
		if err := json.Unmarshal(*value.Source, &device); err != nil {
			common.MatchLogger.WithFields(logrus.Fields{
				"value": value.Source,
			}).Error("Construct device from es error.", err)
		}
		device.OsType = common.ParseOsType(strings.ToLower(string(device.OsType)))
		devices = append(devices, device)
	}
	return devices
}
