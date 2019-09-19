package device

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/common/config"
	"github.com/simplexity-ckcclc/gochannel/common/logger"
	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v6"
	"strings"
	"time"
)

type DeviceAppHandler struct {
	appKey   string
	esClient *elastic.Client
	stopChan chan bool
	// channelType(ios, android, gdt, etc.) - matcher
	matchers   map[common.ChannelType]Matcher
	callbacker *Callbacker
}

func newDeviceAppHandler(appkey string, client *elastic.Client, cb *Callbacker) *DeviceAppHandler {
	return &DeviceAppHandler{
		appKey:     appkey,
		esClient:   client,
		stopChan:   make(chan bool, 1),
		matchers:   make(map[common.ChannelType]Matcher),
		callbacker: cb,
	}
}

func (handler *DeviceAppHandler) start() {
runningLoop:
	for {
		select {
		case <-handler.stopChan:
			logger.MatchLogger.With(logger.Fields{
				"appKey": handler.appKey,
			}).Info("Device App Handler stop")
			break runningLoop
		default:
			latestProcessTime := getLatestProcessTime(handler.appKey)
			if latestProcessTime == 0 {
				// New appKey, process the last 3 days
				latestProcessTime = common.TimeToMillis(time.Now().AddDate(0, 0, -3))
			}
			processEndTime := latestProcessTime + config.GetInt64(config.ProcessPeriod) - common.TimeDurationToMillis(2*time.Second)
			// ES segment flush duration, for insert new sdk device into es
			if processEndTime < latestProcessTime {
				time.Sleep(2 * time.Second)
				continue
			}

			var matchedDevices []*Device
			devices, newLatestProcessTime := handler.getDevices(latestProcessTime, processEndTime, 5)
			if len(devices) > 0 {
				for _, device := range devices {
					for _, matcher := range handler.matchers {
						if err := matcher.Match(device); err != nil {
							logger.MatchLogger.With(logger.Fields{
								"device": device,
							}).Error("match device error.", err)
						}
					}

					if device.Status == Matched {
						matchedDevices = append(matchedDevices, device)
					} else {
						device.Status = Processed
					}
				}

				if len(matchedDevices) > 0 {
					if err := handler.callbacker.preHandle(matchedDevices); err != nil {
						logger.MatchLogger.Error("PreHandle matched devices error.", err)
						continue
					}
				}

				_ = handler.updateDevices(devices)

				if err := updateLatestProcessTime(handler.appKey, newLatestProcessTime); err != nil {
					logger.MatchLogger.WithFields(logrus.Fields{
						"appKey":      handler.appKey,
						"processTime": processEndTime,
					}).Error("Update process time error.", err)
				}
				time.Sleep(2 * time.Second) // ES segment flush duration, for update device status
			} else {
				logger.MatchLogger.With(logger.Fields{
					"startTime": latestProcessTime,
					"endTime":   processEndTime,
					"appKey":    handler.appKey,
				}).Info("No sdk device activation.")

				if err := updateLatestProcessTime(handler.appKey, processEndTime); err != nil {
					logger.MatchLogger.With(logger.Fields{
						"appKey":      handler.appKey,
						"processTime": processEndTime,
					}).Error("Update process time error.", err)
				}
				time.Sleep(30 * time.Second)
			}
		}
	}
}

// do not use context.Cancel here, because DeviceAppHandler do process device in batch in a time.
// don`t want to interrupt the process in case of unstable state, but wait for the inflight batch complete
func (handler *DeviceAppHandler) stop() {
	handler.stopChan <- true
}

func getLatestProcessTime(appKey string) (processTime int64) {
	if err := common.DB.QueryRow("SELECT process_time FROM app_process_info WHERE app_key = ?", appKey).
		Scan(&processTime); err != nil && err != sql.ErrNoRows {
		logger.MatchLogger.With(logger.Fields{
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

func (handler *DeviceAppHandler) getDevices(startTime int64, endTime int64, batchSize int) (devices []*Device, latestProcessTime int64) {
	index := config.GetString(config.EsDeviceIndex)
	query := elastic.NewBoolQuery().
		//MustNot(elastic.NewTermQuery("status", Matched)).
		Must(elastic.NewTermQuery("status", New)).
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
		logger.MatchLogger.Error("Search es device error.", err)
		return
	}

	if esResponse.Hits.TotalHits <= 0 {
		return
	}

	for _, value := range esResponse.Hits.Hits {
		var device Device
		if err := json.Unmarshal(*value.Source, &device); err != nil {
			logger.MatchLogger.With(logger.Fields{
				"value": value.Source,
			}).Error("Construct device from es error.", err)
		}
		device.OsType = common.ParseOsType(strings.ToLower(string(device.OsType)))
		device.ESId = value.Id
		device.ResetMatchInfo()
		devices = append(devices, &device)

		if latestProcessTime < device.ActivateTime {
			latestProcessTime = device.ActivateTime
		}
	}
	return
}

func (handler *DeviceAppHandler) updateDevices(devices []*Device) error {
	index := config.GetString(config.EsDeviceIndex)
	bulkRequest := handler.esClient.Bulk()
	for _, device := range devices {
		req := elastic.NewBulkUpdateRequest().
			Index(index).
			Type(device.AppKey).
			Id(device.ESId).
			Doc(device).
			DocAsUpsert(true)
		bulkRequest.Add(req)
	}

	bulkResponse, err := bulkRequest.Do(context.Background())
	if err != nil {
		logger.MatchLogger.With(logger.Fields{
			"devices": devices,
		}).Error("Bulk put device doc error : ", err)
		return err
	}

	failed := bulkResponse.Failed()
	for _, failedResp := range failed {
		logger.MatchLogger.With(logger.Fields{
			"id":       failedResp.Id,
			"errCause": failedResp.Error,
		}).Error("Bulk put device doc error : ", err)
	}
	return nil
}
