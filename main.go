package main

import (
	"context"
	"flag"
	"fmt"
	kafkaReader "github.com/thearyanahmed/nlogx/kafka"
	"github.com/thearyanahmed/nlogx/utils/env"
	"github.com/thearyanahmed/nlogx/utils/random"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var (
	kafkaBrokerUrl     string
	kafkaTopic         string
	kafkaClientId      string
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	env.LoadEnv()

	flag.StringVar(&kafkaBrokerUrl, "kafka-brokers", "localhost:9092", "Kafka brokers in comma separated value")
	flag.StringVar(&kafkaTopic, "kafka-topic", "hello_world", "Kafka topic. Only one topic per worker.")
	//flag.StringVar(&kafkaConsumerGroup, "kafka-consumer-group", "consumer-group", "Kafka consumer group")
	flag.StringVar(&kafkaClientId, "kafka-client-id", "kafka_client_"+random.Str(3), "Kafka client id")

	flag.Parse()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		select {
		case <-signalChan: // first signal, cancel context
			cancel()
		case <-ctx.Done():
		}
		<-signalChan // second signal, hard exit
		os.Exit(1)
	}()

	brokers := strings.Split(kafkaBrokerUrl, ",")
	topic 	:= env.Get("KAFKA_TOPIC")

	var wg sync.WaitGroup
	wg.Add(1)

	readerConfig := kafkaReader.ReaderConfig(brokers,topic,kafkaClientId)

	kafkaReader.Consume(*readerConfig)

	fmt.Println("end of the line from client")
	wg.Wait()

}