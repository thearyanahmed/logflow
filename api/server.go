package main

import (
	"context"
	json2 "encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/thearyanahmed/nlogx/pb/packet"
	"github.com/thearyanahmed/nlogx/utils/env"
	"github.com/thearyanahmed/nlogx/utils/random"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type server struct {
	packet.UnimplementedLogServiceServer
	StreamCount int64
	writerHandler writerHandler
	startedAt time.Time
}

type stream struct {
	grpc.ServerStream
}

func (st *stream) SendAndClose(response *packet.LogResponse) error {
	// send a response
	fmt.Printf("send and close stream. terminate session\n")
	return nil
}

func (st *stream) Recv(response *packet.LogResponse) error { // handle
	fmt.Printf("recv stream \n")
	return nil
}

type writerHandler struct {
	KafkaWriters map[string]*kafka.Writer
	ctx context.Context
	wg *sync.WaitGroup
}

func (wh *writerHandler) SetWriter(brokers []string,topics []string,ctx context.Context, wg *sync.WaitGroup) {

	writers := make(map[string]*kafka.Writer)

	for _,topic := range topics {

		writer := &kafka.Writer{
			Addr:     kafka.TCP(brokers[0]),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}

		writers[topic] = writer
	}

	wh.KafkaWriters = writers

	wh.ctx = ctx
	wh.wg = wg
}

func (wh *writerHandler) Produce(writer *kafka.Writer, key, msg string) error {

	err := writer.WriteMessages(wh.ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(msg),
		},
	)

	wh.wg.Done()

	return err
}


func (s *server) StreamLog(stream packet.LogService_StreamLogServer) error{
	s.StreamCount = 0
	s.startedAt = time.Now()

	for {
		msg, err := stream.Recv()

		if err == io.EOF {

			timeDiff := time.Since(s.startedAt)

			// todo need to close off the writers
			// use channels maybe ?

			return stream.SendAndClose(&packet.LogResponse{
				Success:       true,
				StreamedCount: s.StreamCount,
				Message:       "transmission ended. reason: end of file. duration: " + timeDiff.String(),
			})
		}

		if err != nil {
			return err
		}

		s.StreamCount++

		key := random.Str(3)

		json, err := getJsonDataFromMessage(msg)

		if err != nil {
			fmt.Printf("error marshaling %v\n",err.Error())
			return err
		}

		go func() {

			for _, topic := range msg.GetTopics() {

				if writer, ok := s.writerHandler.KafkaWriters[topic]; ok {
					s.writerHandler.wg.Add(1)

					err = s.writerHandler.Produce(writer,key,string(json))

					if err != nil {
						fmt.Printf("error producing : %v on topic %v\n",err,topic)
						return
					}
				}

			}

		}()
	}
}

func getJsonDataFromMessage(request *packet.LogRequest) ([]byte, error){
	var data []interface{}

	data = append(data,request.GetPacket())
	data = append(data,request.GetLog())
	data = append(data,request.GetHeaders())
	data = append(data,request.GetTopics())

	return json2.Marshal(data)
}

func main()  {
	fmt.Printf("running server\n")
	env.LoadEnv()

	lis, err := net.Listen("tcp", ":5053")

	if err != nil {
		fmt.Printf("error opening tcp server %v\n",err.Error())
	}

	s := grpc.NewServer()

	brokers := strings.Split(env.Get("KAFKA_BROKER_ADDRESS"), ",")
	topics 	:= strings.Split(env.Get("KAFKA_TOPICS"), ",")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wh := writerHandler{}

	var wg sync.WaitGroup

	wh.SetWriter(brokers,topics,ctx,&wg)

	server := server{
		writerHandler: wh,
	}

	packet.RegisterLogServiceServer(s,&server)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}