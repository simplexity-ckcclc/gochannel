package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	pb "github.com/simplexity-ckcclc/gochannel/match/proto"
	"time"
)

var (
	bootstrapServers = []string{"localhost:9092"}
	topic            = "gochannel"
)

func ProduceDeviceProto() {
	producer, err := sarama.NewSyncProducer(bootstrapServers, nil)
	if err != nil {
		panic("Error creating the sync producer")
	}

	defer func() {
		err := producer.Close()
		if err != nil {
			fmt.Println("Error closing producer: ", err)
			return
		}
		fmt.Println("Producer closed")
	}()

	value, err := marshalDevice()
	if err != nil {
		panic(err)
	}

	message := sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(value)}
	partition, offset, err := producer.SendMessage(&message)
	if err != nil {
		fmt.Println("Error sending message: ", err)
	} else {
		fmt.Printf("Sent message value='%s' at partition = %d, offset = %d\n", value, partition, offset)
	}
}

func marshalDevice() (message []byte, err error) {
	var ptime *timestamp.Timestamp
	ptime, err = ptypes.TimestampProto(time.Unix(1562238598, 132000000))
	device := &pb.SdkDeviceReport{
		OsType:       pb.SdkDeviceReport_ANDROID,
		Imei:         "1234567890123456",
		AppKey:       "appkeyA",
		Channel:      "channelA",
		ActivateTime: ptime,
	}

	message, err = proto.Marshal(device)
	return
}
