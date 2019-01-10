package match

import (
	"encoding/json"
	"fmt"
)

type messageHandler interface {
	handle(message []byte)
}

type MatchHandler struct {
}

func (handler MatchHandler) handle(message []byte) {
	fmt.Println(string(message))
	var device deviceInfo
	err := json.Unmarshal(message, &device)
	if err == nil {
		fmt.Println(device.DeviceId)
	} else {
		fmt.Println(err)
	}
}
