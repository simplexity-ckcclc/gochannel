package match

import (
	"github.com/bsm/sarama-cluster"
	"github.com/simplexity-ckcclc/gochannel/common"
	"log"
)

func Serve() {
	messages := make(chan []byte, 5)
	handler := SdkMessageReceiver{}
	go consumeMessage(messages, handler)

	err := startKafkaConsumer(messages)
	if err != nil {
		panic(err)
	}

}

func startKafkaConsumer(messages chan<- []byte) error {
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
		return err
	}

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	// consume messages, watch signals
	for msg := range consumer.Messages() {
		messages <- msg.Value
		consumer.MarkOffset(msg, "") // mark message as processed
	}
	return nil
}

func consumeMessage(messages <-chan []byte, handler messageReceiver) {
	for msg := range messages {
		handler.receive(msg)
	}
}
