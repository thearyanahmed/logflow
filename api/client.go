package main

import (
	"context"
	json2 "encoding/json"
	"fmt"
	"github.com/thearyanahmed/nlogx/pb/packet"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

func main()  {
	fmt.Printf("running client\n")

	conn, err := grpc.Dial("localhost:5053", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("connection error:%v\n",err.Error())
		}
	}(conn)

	c := packet.NewLogServiceClient(conn)

	ctx , cancel := context.WithTimeout(context.Background(),time.Second * 20)
	defer cancel()

	r , err := c.StreamLog(ctx)

	if err != nil{
		log.Fatalf("error in request %v\n",err.Error())
	}

	var i int32


	var wg sync.WaitGroup
	var mutex = sync.Mutex{}


	for i = 0 ; i < 10000000; i++ {
		wg.Add(1)

		go func( i int32, mt *sync.Mutex, wg *sync.WaitGroup) {
			defer wg.Done()
			defer mt.Unlock()

			mt.Lock()

			json, err := json2.Marshal("hello world")

			if err != nil {
				fmt.Printf("error marshaling %v\n",err.Error())
				return
			}

			var headers []*packet.Header

			header := &packet.Header{
				Key:   "hello",
				Value: "world",
			}

			headers = append(headers,header)

			logRequest := packet.LogRequest{
				Headers: headers,
				Topics: []string{"hello_world","banana_world"},
				Payload: json,
			}

			err = r.Send(&logRequest)

			if err != nil {
				return
			}
		}(i,&mutex,&wg)
	}

	wg.Wait()

	reply , err := r.CloseAndRecv()

	if err != nil {
		log.Fatalf("error reply %v\n",err.Error())
	}


	fmt.Printf("\nstream count :%v\nmsg : %v\n",reply.GetStreamedCount(),reply.GetMessage())

}
