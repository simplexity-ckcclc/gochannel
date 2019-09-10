package match

import (
	"github.com/bsm/sarama-cluster"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/simplexity-ckcclc/gochannel/common"
	"github.com/simplexity-ckcclc/gochannel/common/config"
	"github.com/simplexity-ckcclc/gochannel/common/logger"
	pb "github.com/simplexity-ckcclc/gochannel/protos"
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
	consumerConfig := cluster.NewConfig()
	consumerConfig.Consumer.Return.Errors = true
	consumerConfig.Group.Return.Notifications = true

	// init consumer
	brokers := config.GetStringSlice(config.KafkaBootstrapServer)
	topics := config.GetStringSlice(config.KafkaTopic)
	groupId := config.GetString(config.KafkaGroupId)
	consumer, err := cluster.NewConsumer(brokers, groupId, topics, consumerConfig)
	if err != nil {
		panic(err)
	}

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			logger.MatchLogger.Error("Kafka consumer error : ", err)
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			logger.MatchLogger.Debug("Kafka rebalanced notifications : ", ntf)
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
		logger.MatchLogger.Error("Parse device error : ", err)
	} else {
		if !validate(device) {
			logger.MatchLogger.With(logger.Fields{
				"Device ": device,
			}).Debug("Invalid sdk message.")
			return
		}

		if err := insertIntoDB(device); err != nil {
			logger.MatchLogger.With(logger.Fields{
				"Device ": device,
			}).Error("Insert into DB error.", err)
		} else {
			logger.MatchLogger.With(logger.Fields{
				"Device ": device,
			}).Debug("Insert into DB.")
		}
	}

}

func validate(device *pb.SdkDeviceReport) bool {
	return device != nil && device.GetAppKey() != "" && device.GetActivateTime() != nil
}

func insertIntoDB(device *pb.SdkDeviceReport) error {
	rtime, err := ptypes.Timestamp(device.ActivateTime)
	if err != nil {
		return nil
	}

	_, err = common.DB.Exec("INSERT INTO sdk_device_report (imei, idfa, app_key, channel_id, resolution, "+
		"language, os_type, os_version, activate_time, source_ip) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		device.Imei, device.Idfa, device.AppKey, device.Channel, device.Resolution,
		device.Language, device.OsType.String(), device.OsVersion, common.TimeToMillis(rtime),
		device.SourceIp)
	return err
}
