package device

import (
	"database/sql"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v6"
	"time"
)

type DeviceProcessor struct {
	db          *sql.DB
	esClient    *elastic.Client
	appHandlers map[string]*DeviceAppHandler
	callbacker  *Callbacker
	stopChan    chan bool
}

func NewDeviceProcessor(database *sql.DB, client *elastic.Client) *DeviceProcessor {
	return &DeviceProcessor{
		db:          database,
		esClient:    client,
		appHandlers: make(map[string]*DeviceAppHandler),
		stopChan:    make(chan bool, 1),
		callbacker:  NewCallbacker(database),
	}
}

func (processor *DeviceProcessor) Start() {
	processor.startNewAppHandler()
	ticker := time.NewTicker(time.Minute * 1)
runningLoop:
	for {
		select {
		case <-processor.stopChan:
			common.MatchLogger.Info("Device Handler stop")
			break runningLoop
		case <-ticker.C:
			processor.startNewAppHandler()
			ticker = time.NewTicker(time.Minute * 1)
		}
	}
}

func (processor *DeviceProcessor) Stop() {
	processor.stopChan <- true
	for _, appHandler := range processor.appHandlers {
		appHandler.stop()
	}
	processor.callbacker.stop()
}

func (processor *DeviceProcessor) startNewAppHandler() {
	appChannels, err := processor.scanAppChannels()
	if err != nil {
		common.MatchLogger.Error("Select app_channel from table app_channel error")
		return
	}

	for _, appChannel := range appChannels {
		if _, ok := processor.appHandlers[appChannel.AppKey]; !ok {
			// New appKey, start new DeviceAppHandler
			handler := &DeviceAppHandler{
				appKey:     appChannel.AppKey,
				esClient:   processor.esClient,
				stopChan:   make(chan bool, 1),
				matchers:   make(map[common.ChannelType]Matcher),
				callbacker: processor.callbacker,
			}
			processor.appHandlers[appChannel.AppKey] = handler
			go handler.start()
			common.MatchLogger.WithFields(logrus.Fields{
				"appKey": appChannel.AppKey,
			}).Info("Start appHandler.")
		}

		handler := processor.appHandlers[appChannel.AppKey]
		if _, ok := handler.matchers[appChannel.ChannelType]; !ok {
			// New app channel_type for this appKey
			instantiateMatcherFunc := matcherMappings[appChannel.ChannelType]
			matcher := instantiateMatcherFunc(processor.esClient)
			handler.matchers[appChannel.ChannelType] = matcher
			common.MatchLogger.WithFields(logrus.Fields{
				"appKey":      appChannel.AppKey,
				"channelType": appChannel.ChannelType,
			}).Info("Instantiate and add matcher to appHandler.")
		}
	}

	// todo Remove deprecate appKey and appChannelTypeMatcher
}

func (processor *DeviceProcessor) scanAppChannels() ([]common.AppChannel, error) {
	rows, err := processor.db.Query(`SELECT app_key, channel_id, channel_type FROM app_channel`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appChannels []common.AppChannel
	var appChannel common.AppChannel
	var appKey, channelId, channelType string
	for rows.Next() {
		if err = rows.Scan(&appKey, &channelId, &channelType); err != nil {
			return nil, err
		}

		appChannel = common.AppChannel{
			AppKey:      appKey,
			ChannelId:   channelId,
			ChannelType: common.ParseChannelType(channelType),
		}
		appChannels = append(appChannels, appChannel)
	}
	err = rows.Err()
	return appChannels, err
}
