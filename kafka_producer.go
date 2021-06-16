package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/thearyanahmed/nlogx/utils/env"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	numberOfMessages = 0
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main()  {
	env.LoadEnv()

	appEnv := env.Get("APP_ENV")

	if appEnv != "local" {
		fmt.Println("[+] can run only on local environment. Quiting...")
		os.Exit(1)
	}

	produceRandomStuff()

	fmt.Printf("%v messages sent!\n",numberOfMessages)
}


func produceRandomStuff() {
	brokerAddresses := env.Get("KAFKA_BROKER_ADDRESS")
	brokers 		:= strings.Split(brokerAddresses,",")

	topic 	:= env.Get("KAFKA_TOPIC")

	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers[0]),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go writeMessage(w,&wg)
	}

	wg.Wait()

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	fmt.Printf("end of the line \n")
}

func writeMessage(w *kafka.Writer, group *sync.WaitGroup) {
	key := randStringRunes(4)
	msg := randStringRunes(10)

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(msg),
		},
	)

	group.Done()

	if err != nil {
		_ = fmt.Errorf("failed to write messages: %v\n", err.Error())
		return
	}

	numberOfMessages++

}


var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}