package device

import (
	"context"
	"database/sql"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/common/logger"
	"gopkg.in/olivere/elastic.v6"
	"time"
)

type DeviceProcessor struct {
	ctx         context.Context
	db          *sql.DB
	esClient    *elastic.Client
	appHandlers map[string]*DeviceAppHandler
	callbacker  *Callbacker
}

func NewDeviceProcessor(context context.Context, database *sql.DB, client *elastic.Client) *DeviceProcessor {
	return &DeviceProcessor{
		ctx:         context,
		db:          database,
		esClient:    client,
		appHandlers: make(map[string]*DeviceAppHandler),
		callbacker:  NewCallbacker(context, database),
	}
}

func (processor *DeviceProcessor) Start() {
	processor.startNewAppHandler()
	ticker := time.NewTicker(time.Minute * 1)
runningLoop:
	for {
		select {
		case <-processor.ctx.Done():
			logger.MatchLogger.Info("Device Handler stop")
			break runningLoop
		case <-ticker.C:
			processor.startNewAppHandler()
			ticker = time.NewTicker(time.Minute * 1)
		}
	}
}

func (processor *DeviceProcessor) startNewAppHandler() {
	appChannels, err := processor.scanAppChannels()
	if err != nil {
		logger.MatchLogger.Error("Select app_channel from table app_channel error")
		return
	}

	for _, appChannel := range appChannels {
		if _, ok := processor.appHandlers[appChannel.AppKey]; !ok {
			// New appKey, start new DeviceAppHandler
			handler := newDeviceAppHandler(processor.ctx, appChannel.AppKey, processor.esClient, processor.callbacker)
			processor.appHandlers[appChannel.AppKey] = handler
			go handler.start()
			logger.MatchLogger.With(logger.Fields{
				"appKey": appChannel.AppKey,
			}).Info("Start appHandler.")
		}

		handler := processor.appHandlers[appChannel.AppKey]
		if _, ok := handler.matchers[appChannel.ChannelType]; !ok {
			// New app channel_type for this appKey
			instantiateMatcherFunc := matcherMappings[appChannel.ChannelType]
			matcher := instantiateMatcherFunc(processor.esClient)
			handler.matchers[appChannel.ChannelType] = matcher
			logger.MatchLogger.With(logger.Fields{
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
