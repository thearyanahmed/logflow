package main

import (
	"context"
	"fmt"
	"github.com/thearyanahmed/nlogx/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedFlowServiceServer
}

func (s *server) RecordFlow (_ context.Context, in *pb.FlowRequest) (*pb.FlowResponse, error) {

	fmt.Printf("got data from request %v\n",in.GetFlow().Hello)

	res := pb.FlowResponse{Status: true}

	return &res, nil
}

func main()  {
	fmt.Printf("running server\n")

	lis, err := net.Listen("tcp", ":5053")

	if err != nil {
		fmt.Printf("error opening tcp server %v\n",err.Error())
	}

	s := grpc.NewServer()
	pb.RegisterFlowServiceServer(s,&server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}