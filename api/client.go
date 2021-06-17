package main

import (
	"context"
	"fmt"
	"github.com/thearyanahmed/nlogx/pb"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main()  {
	fmt.Printf("running client\n")

	conn, err := grpc.Dial("localhost:5053", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewFlowServiceClient(conn)

	ctx , cancel := context.WithTimeout(context.Background(),time.Second)
	defer cancel()

	flow := pb.Flow{Hello: 2312}

	req := pb.FlowRequest{Flow: &flow}

	r , err := c.RecordFlow(ctx,&req)

	if err != nil{
		log.Fatalf("error in request %\v",err.Error())
	}

	fmt.Printf("response : %v\n",r.Status)

}
