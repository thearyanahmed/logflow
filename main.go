package main

import (
	"context"
	"flag"
	"github.com/thearyanahmed/nlogx/kafka"
	"github.com/thearyanahmed/nlogx/utils/env"
	"github.com/thearyanahmed/nlogx/utils/random"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
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

	writer := kafka.Writer(brokers,topic)

	err := kafka.Produce(writer, ctx,"hello", "world")

	if err != nil {
		log.Fatalf("error writing to kafka %v\n",err.Error())
		return
	}
}