package device

import (
	"database/sql"
	"github.com/simplexity-ckcclc/gochannel/common"
	"time"
)

type DeviceHandler struct {
	db          *sql.DB
	appHandlers map[string]DeviceAppHandler
	stopChan    chan bool
}

func NewDeviceHandler() DeviceHandler {
	return DeviceHandler{
		db:          common.DB,
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
	appKeys, err := handler.scanAppKeys()
	if err != nil {
		common.MatchLogger.Error("Select app_key from table channel_sig error")
	}

	for _, appKey := range appKeys {
		if _, ok := handler.appHandlers[appKey]; !ok {
			// New appKey, start new DeviceAppHandler
			appHandler := DeviceAppHandler{
				appKey:   appKey,
				stopChan: make(chan bool, 1),
			}
			handler.appHandlers[appKey] = appHandler
			go appHandler.start()
		}
	}

	// Deprecate appKey
	for appKey := range handler.appHandlers {
		if !contains(appKeys, appKey) {
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

func (handler DeviceHandler) scanAppKeys() ([]string, error) {
	rows, err := handler.db.Query(`SELECT app_key FROM channel_sig`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appKeys []string
	var appKey string
	for rows.Next() {
		if err = rows.Scan(&appKey); err != nil {
			return nil, err
		}
		appKeys = append(appKeys, appKey)
	}
	err = rows.Err()
	return appKeys, err
}
