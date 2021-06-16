package kafka

import (
	"context"
	"fmt"
)

func Consume(ctx context.Context) {
	fmt.Println("hello world")
	//r := kafka.NewReader(kafka.ReaderConfig{
	//	Brokers: []string{"localhost:9092"},
	//	Topic:   "hello_world",
	//	GroupID: "my-group",
	//})
	//for {
	//	// the `ReadMessage` method blocks until we receive the next event
	//	msg, err := r.ReadMessage(ctx)
	//	if err != nil {
	//		panic("could not read message " + err.Error())
	//	}
	//	// after receiving the message, log its value
	//	fmt.Println("received: ", string(msg.Value))
	//}
}
