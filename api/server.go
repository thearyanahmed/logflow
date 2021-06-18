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
)

type server struct {
	packet.UnimplementedLogServiceServer
	StreamCount int64
	writerHandler writerHandler
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
	KafkaWriter *kafka.Writer
	ctx context.Context
	wg *sync.WaitGroup
}

func (wh *writerHandler) SetWriter(brokers []string,topic string,ctx context.Context, wg *sync.WaitGroup) {

	wh.KafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(brokers[0]),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	wh.ctx = ctx
	wh.wg = wg
}

func (wh *writerHandler) Produce(key, msg string) error {

	err := wh.KafkaWriter.WriteMessages(wh.ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(msg),
		},
	)

	wh.wg.Done()

	return err
}


func (s *server) StreamLog(stream packet.LogService_StreamLogServer) error{

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&packet.LogResponse{
				Success:       true,
				StreamedCount: s.StreamCount,
				Message:       "transmission ended. reason: end of file",
			})
		}

		if err != nil {
			return err
		}

		fmt.Printf("received stream msg: %v\n push data to kafka\n",msg.GetPacket().GetAgentIp())

		s.StreamCount++

		s.writerHandler.wg.Add(1)


		go func() {

			key := random.Str(3)

			var data []interface{}

			data = append(data,msg.GetPacket())
			data = append(data,msg.GetLog())
			data = append(data,msg.GetHeaders())


			fmt.Printf("data %v\n",data)

			json, err := json2.Marshal(data)

			if err != nil {
				fmt.Printf("error marshaling %v\n",err.Error())
				return
			}

			err = s.writerHandler.Produce(key,string(json))

			if err != nil {
				fmt.Printf("error producing : %v\n",err.Error())
				return
			}
		}()
	}
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
	topic 	:= env.Get("KAFKA_TOPIC")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wh := writerHandler{}

	var wg sync.WaitGroup

	wh.SetWriter(brokers,topic,ctx,&wg)

	server := server{
		StreamCount:   0,
		writerHandler: wh,
	}

	packet.RegisterLogServiceServer(s,&server)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}