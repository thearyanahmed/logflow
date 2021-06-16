package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

func BootAndListen() {
	conf := kafka.ReaderConfig{
		Brokers:                []string{"localhost:9092"},
		GroupID:                "g_1",
		GroupTopics:            nil,
		Topic:                  "hello_world",
		Partition:              0,
		MaxBytes: 				10,
	}

	reader := kafka.NewReader(conf)

	for {
		message, err := reader.ReadMessage(context.Background())

		if err != nil {
			fmt.Println("Error ",err.Error())
			continue
		}

		fmt.Printf("message: %s\n",string(message.Value))
	}
	
}
