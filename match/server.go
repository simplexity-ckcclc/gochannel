package match

import (
	"context"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/common/logger"
	"github.com/simplexity-ckcclc/gochannel/match/device"
)

var cancelProc context.CancelFunc

func Serve() {
	msgChan := make(chan []byte, 5)
	receiver := SdkMsgReceiver{
		SdkMsgHandler: PbSdkMsgHandler{},
		SdkMsgChannel: msgChan,
	}
	go receiver.ConsumeMessage()
	logger.MatchLogger.Info("Starting consume kafka message")

	// transfer devices from mysql to es
	devicePorter := device.NewDevicePorter(common.DB, common.EsClient)
	go devicePorter.TransferDevices()
	logger.MatchLogger.Info("Starting transfer data from mysql to es")

	// bulk get devices from es and then run match process
	var ctx context.Context
	ctx, cancelProc = context.WithCancel(context.Background())
	processor := device.NewDeviceProcessor(ctx, common.DB, common.EsClient)
	go processor.Start()
	logger.MatchLogger.Info("Starting process data")

	logger.MatchLogger.Info("Match server started")
}

func Stop() {
	cancelProc()
}
