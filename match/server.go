package match

import (
    "github.com/bsm/sarama-cluster"
    "log"
)

func Serve() {
    messages := make(chan []byte, 5)
    handler := MatchHandler{}
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
    brokers := []string{"127.0.0.1:9092"}
    topics := []string{"foo"}
    consumer, err := cluster.NewConsumer(brokers, "my-consumer-group", topics, config)
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
        consumer.MarkOffset(msg, "")	// mark message as processed
    }
    return nil
}

func consumeMessage(messages <-chan []byte, handler messageHandler) {
    for msg := range messages {
        handler.handle(msg)
    }
}