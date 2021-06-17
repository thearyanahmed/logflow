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

	r , err := c.StreamFlow(ctx)

	if err != nil{
		log.Fatalf("error in request %v\n",err.Error())
	}

	var i int32

	for i = 0 ; i < 10; i++ {
		flow := pb.Flow{Hello: i}

		err := r.Send(&pb.FlowRequest{
			Flow: &flow,
		})


		if err != nil {
			return
		}
	}

	reply , err := r.CloseAndRecv()

	if err != nil {
		log.Fatalf("error reply %v\n",err.Error())
	}


	fmt.Printf("response : %v\n",reply.GetStatus())

}
