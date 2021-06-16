package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/thearyanahmed/nlogx/utils/env"
	"strings"
)

func BootAndListen() {
	conf := getReadConfig()

	reader := kafka.NewReader(conf)

	for {
		message, err := reader.ReadMessage(context.Background())

		if err != nil {
			fmt.Println("Error ",err.Error())
			continue
		}


		fmt.Printf("[+] msg: %s\n",string(message.Value))
	}
}

func getReadConfig() kafka.ReaderConfig{
	brokerAddresses := env.Get("KAFKA_BROKER_ADDRESS")
	brokers 		:= strings.Split(brokerAddresses,",")

	topic 	:= env.Get("KAFKA_TOPIC")
	groupID := env.Get("KAFKA_GROUP_ID")


	return kafka.ReaderConfig{
		Brokers:                brokers,
		GroupID:                groupID,
		GroupTopics:            nil,
		Topic:                  topic,
		Partition:              0,
		MaxBytes: 				10,
	}
}
