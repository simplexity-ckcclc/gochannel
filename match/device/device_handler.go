package device

import (
	"database/sql"
	"github.com/simplexity-ckcclc/gochannel/common"
	"gopkg.in/olivere/elastic.v6"
	"time"
)

type DeviceHandler struct {
	db          *sql.DB
	esClient    *elastic.Client
	appHandlers map[string]DeviceAppHandler
	stopChan    chan bool
}

func NewDeviceHandler(database *sql.DB, client *elastic.Client) DeviceHandler {
	return DeviceHandler{
		db:          database,
		esClient:    client,
		appHandlers: make(map[string]DeviceAppHandler),
		stopChan:    make(chan bool, 1),
	}
}

func (handler DeviceHandler) Start() {
	handler.startNewAppHandler()
	ticker := time.NewTicker(time.Second * 3)
runningLoop:
	for {
		select {
		case <-handler.stopChan:
			common.MatchLogger.Info("Device Handler stop")
			break runningLoop
		case <-ticker.C:
			handler.startNewAppHandler()
			ticker = time.NewTicker(time.Second * 3)
		}
	}
}

func (handler DeviceHandler) Stop() {
	handler.stopChan <- true
	for _, appHandler := range handler.appHandlers {
		appHandler.stop()
	}
}

func (handler DeviceHandler) startNewAppHandler() {
	appChannels, err := handler.scanAppChannels()
	if err != nil {
		common.MatchLogger.Error("Select app_key from table app_channel error")
		return
	}

	for _, appChannel := range appChannels {
		if _, ok := handler.appHandlers[appChannel.AppKey]; !ok {
			// New appKey, start new DeviceAppHandler
			var matchers []Matcher
			matcherFunc := matcherMappings[appChannel.ChannelType]
			matcher := matcherFunc(handler.esClient)
			matchers = append(matchers, matcher)

			appHandler := DeviceAppHandler{
				appKey:   appChannel.AppKey,
				esClient: handler.esClient,
				stopChan: make(chan bool, 1),
				matchers: matchers,
			}
			handler.appHandlers[appChannel.AppKey] = appHandler
			go appHandler.start()
		}
	}

	// Deprecate appKey
	for appKey := range handler.appHandlers {
		if !contains(appChannels, appKey) {
			appHandler := handler.appHandlers[appKey]
			appHandler.stop()
			delete(handler.appHandlers, appKey)
		}
	}
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func (handler DeviceHandler) scanAppChannels() ([]common.AppChannel, error) {
	rows, err := handler.db.Query(`SELECT app_key, channel_id, channel_type FROM app_channel`)
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
