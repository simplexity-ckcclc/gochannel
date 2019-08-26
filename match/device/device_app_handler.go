package device

import (
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/sirupsen/logrus"
)

type DeviceAppHandler struct {
	appKey   string
	stopChan chan bool
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
			//todo start process
		}
	}
}

func (handler DeviceAppHandler) stop() {
	handler.stopChan <- true
}
