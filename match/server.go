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

	devicePorter, err := device.NewDevicePorter(common.DB)
	if err != nil {
		panic(err)
	}
	go devicePorter.TransferDevices()

}
