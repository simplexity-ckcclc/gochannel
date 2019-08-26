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
	devicePorter, err := device.NewDevicePorter(common.DB)
	if err != nil {
		panic(err)
	}
	go devicePorter.TransferDevices()

	// bulk get devices from es and then run match process
	deviceHandler := device.NewDeviceHandler()
	go deviceHandler.Start()

}
