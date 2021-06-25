package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
	"time"
)

type WriterInterface interface {
	Writer(brokers []string,topic string) *kafka.Writer
	Produce(w *kafka.Writer, ctx context.Context, key, msg string) error
}

func ReaderConfig(brokers []string,kafkaTopic ,kafkaClientId string) *kafka.ReaderConfig {
	// make a new reader that consumes from topic-A
	return &kafka.ReaderConfig{
		Brokers:         brokers,
		GroupID:         kafkaClientId,
		Topic:           kafkaTopic,
		MinBytes:        10e3,            // 10KB
		MaxBytes:        10e6,            // 10MB
		MaxWait:         1 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: -1,
	}
}

func Consume(config kafka.ReaderConfig) /*(kafka.Message,error)*/ {
	reader := kafka.NewReader(config)

	defer func(reader *kafka.Reader) {
		_ = reader.Close()
	}(reader)

	for {
		m, err := reader.ReadMessage(context.Background())
		fmt.Printf("tpc: %v | prtn: %v | ofst: %v |: %s\n", m.Topic, m.Partition, m.Offset, string(m.Value))

		if err != nil {
			fmt.Printf("error not nil, %v\n",err.Error())
		}
		//return m, err
	}
}

// GetTopics TODO update
func GetTopics() {

	fmt.Printf("action: get topics \n")
	conn, err := kafka.Dial("tcp", "kafka:9092")

	if err != nil {
		fmt.Printf("erorr getting topic connecton, %v\n",err.Error())
		panic(err.Error())
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()

	if err != nil {
		fmt.Printf("erorr reading topic partition, %v\n",err.Error())
		panic(err.Error())
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}

	count := 0
	for k := range m {
		fmt.Println(k)
		count++
	}

	if count == 0 {
		CreateATopic()
	}

	fmt.Printf("done\n")
}

func CreateATopic() {
	// to create topics when auto.create.topics.enable='false'
	topic := "hello_world"

	conn, err := kafka.Dial("tcp", "kafka:9092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn

	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()


	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}

	GetTopics()
}