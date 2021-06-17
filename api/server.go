package main

import (
	"fmt"
	"github.com/thearyanahmed/nlogx/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedFlowServiceServer
}

type stream struct {
	grpc.ServerStream
}

func (st *stream) SendAndClose(*pb.FlowResponse) error {
	return nil
}

func (st *stream) Recv(*pb.FlowResponse) error {
	return nil
}

func (s *server) StreamFlow(stream pb.FlowService_StreamFlowServer) error{

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&pb.FlowResponse{Status: true})
		}

		if err != nil {
			return err
		}

		fmt.Printf("received stream msg: %v\n push data to kafka",msg.GetFlow().GetHello())
	}
}

func main()  {
	fmt.Printf("running server\n")

	lis, err := net.Listen("tcp", ":5053")

	if err != nil {
		fmt.Printf("error opening tcp server %v\n",err.Error())
	}

	s := grpc.NewServer()
	server := server{}
	pb.RegisterFlowServiceServer(s,&server)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}