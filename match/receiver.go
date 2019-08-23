package match

import (
	"github.com/bsm/sarama-cluster"
	"github.com/golang/protobuf/proto"
	"github.com/simplexity-ckcclc/gochannel/common"
	pb "github.com/simplexity-ckcclc/gochannel/match/proto"
	"github.com/sirupsen/logrus"
	"log"
)

type SdkMsgReceiver struct {
	SdkMsgChannel chan []byte
	SdkMsgHandler
}

type SdkMsgHandler interface {
	handle(message []byte)
}

func (receiver SdkMsgReceiver) ConsumeMessage() {
	go receiver.consumeKafkaMsg()

	for msg := range receiver.SdkMsgChannel {
		receiver.handle(msg)
	}
}

func (receiver SdkMsgReceiver) consumeKafkaMsg() {
	// init (custom) config, enable errors and notifications
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	// init consumer
	brokers := common.Conf.Kafka.Consumer.BootstrapServer
	topics := common.Conf.Kafka.Consumer.Topic
	groupId := common.Conf.Kafka.Consumer.GroupId
	consumer, err := cluster.NewConsumer(brokers, groupId, topics, config)
	if err != nil {
		panic(err)
	}

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Println("Error:", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Println("Rebalanced: ", ntf)
		}
	}()

	// consume messages, watch signals
	msgChan := receiver.SdkMsgChannel
	for msg := range consumer.Messages() {
		msgChan <- msg.Value
		consumer.MarkOffset(msg, "") // mark message as processed
	}
}

type PbSdkMsgHandler struct {
}

func (handler PbSdkMsgHandler) handle(message []byte) {
	device := &pb.SdkDeviceReport{}
	if err := proto.Unmarshal(message, device); err != nil {
		common.MatchLogger.Error("Parse device error : ", err)
	} else {
		if !validate(device) {
			common.MatchLogger.WithFields(logrus.Fields{
				"Device ": device,
			}).Debug("Invalid sdk message")
			return
		}

		if err := insertIntoDB(device); err != nil {
			common.MatchLogger.WithFields(logrus.Fields{
				"Device ": device,
			}).Error("Insert into DB error", err)
		} else {
			common.MatchLogger.WithFields(logrus.Fields{
				"Device ": device,
			}).Debug("Insert into DB")
		}
	}

}

func validate(device *pb.SdkDeviceReport) bool {
	return device != nil && device.GetAppKey() != "" && device.GetReceiveTime() != nil
}

func insertIntoDB(device *pb.SdkDeviceReport) error {
	stmt, err := common.DB.Prepare("INSERT INTO sdk_device_report (imei, idfa, app_key, channel_id, resolution, " +
		"language, os_type, os_version, receive_time, source_ip) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(device.Imei, device.Idfa, device.AppKey, device.Channel, device.Resolution,
		device.Language, device.OsType, device.OsVersion, device.ReceiveTime, device.SourceIp)
	return err
}
