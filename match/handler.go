package match

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	pb "github.com/simplexity-ckcclc/gochannel/match/proto"
)

type messageHandler interface {
	handle(message []byte)
}

type MatchHandler struct {
}

func (handler MatchHandler) handle(message []byte) {
	fmt.Println("Receive string : ", string(message))
	device := &pb.Device{}
	if err := proto.Unmarshal(message, device); err != nil {
		fmt.Println("Failed to parse device :", err)
	} else {
		fmt.Println("Device : ", device)
	}

	// json
	//var device pb.Device
	//if err := json.Unmarshal(message, &device); err != nil {
	//	fmt.Println("Failed to parse device :", err)
	//} else {
	//	fmt.Println(device)
	//}

}
