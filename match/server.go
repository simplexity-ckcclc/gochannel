package match

import (
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/match/device"
)

func Serve() {
	msgChan := make(chan []byte, 5)
	receiver := SdkMsgReceiver{
		SdkMsgHandler: PbSdkMsgHandler{},
		SdkMsgChannel: msgChan,
	}
	go receiver.ConsumeMessage()

	// transfer devices from mysql to es
	devicePorter := device.NewDevicePorter(common.DB, common.EsClient)
	go devicePorter.TransferDevices()

	// bulk get devices from es and then run match process
	deviceHandler := device.NewDeviceHandler(common.DB, common.EsClient)
	go deviceHandler.Start()

	common.MatchLogger.Info("Match server started")
}
