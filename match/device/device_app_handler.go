package device

import (
	"database/sql"
	"fmt"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/sirupsen/logrus"
	"time"
)

type DeviceAppHandler struct {
	appKey   string
	stopChan chan bool
}

func (handler DeviceAppHandler) start() {
	latestProcessTime, err := latestProcessTime(handler.appKey)
	if err == sql.ErrNoRows {
		latestProcessTime = time.Now().Add(-3 * 24 * time.Hour).Unix()
	} else if err != nil {
		common.MatchLogger.WithFields(logrus.Fields{
			"appKey": handler.appKey,
		}).Error("Get app_key last process time error")
	}

runningLoop:
	for {
		select {
		case <-handler.stopChan:
			common.MatchLogger.WithFields(logrus.Fields{
				"appKey": handler.appKey,
			}).Info("Device App Handler stop")
			break runningLoop
		default:
			fmt.Println(latestProcessTime)
			//todo start process
		}
	}
}

func (handler DeviceAppHandler) stop() {
	handler.stopChan <- true
}

func latestProcessTime(appKey string) (processTime int64, err error) {
	err = common.DB.QueryRow("SELECT process_time FROM app_key_process_info WHERE app_key = ?", appKey).Scan(&processTime)
	return
}
