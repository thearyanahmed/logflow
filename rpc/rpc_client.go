package main

import (
	"bufio"
	"context"
	json2 "encoding/json"
	"fmt"
	"github.com/thearyanahmed/logflow/pb/packet"
	"github.com/thearyanahmed/logflow/utils/env"
	"google.golang.org/grpc"
	"log"
	"os"
	"sync"
	"time"
)

func init()  {
	env.LoadEnv()
}

func main()  {
	fmt.Printf("running client\n")

	host := env.Get("RPC_HOST") + ":" + env.Get("RPC_PORT")

 	conn, err := grpc.Dial(host, grpc.WithInsecure(), grpc.WithBlock())

 	fmt.Println(conn,err)

	if err != nil {
		fmt.Printf("could not dial grpc\n%v\n",err.Error())
		return
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
	fmt.Printf("checkpoint 2\n")

	r , err := c.StreamLog(ctx)

	if err != nil{
		fmt.Printf("error with stream log\n%v\n",err.Error())
		return
	}

	filePath, _ := os.Getwd()

	filePath = filePath + "/" + env.Get("TEST_DATA_FILE")

	file, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("error with file %v\n",err.Error())
		return
	}
	fmt.Printf("checkpoint 3\n")

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var wg sync.WaitGroup
	var mutex = sync.Mutex{}

	for scanner.Scan() {

		fmt.Printf("%v\n",scanner.Text())

		wg.Add(1)

		go func(mt *sync.Mutex, wg *sync.WaitGroup) {
			defer wg.Done()
			defer mt.Unlock()

			mt.Lock()

			json, err := json2.Marshal(scanner.Text())

			if err != nil {
				fmt.Printf("error marshaling %v\n",err.Error())
				return
			}

			logRequest := packet.LogRequest{
				Topics: []string{"hello_world"},
				Payload: json,
			}

			err = r.Send(&logRequest)

			if err != nil {
				fmt.Printf("error while sending data to grpc server %v\n",err.Error())
				return
			}
		}(&mutex,&wg)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("scanner error %v\n",err.Error())
	}

	wg.Wait()

	reply , err := r.CloseAndRecv()

	if err != nil {
		log.Fatalf("error reply %v\n",err.Error())
	}

	fmt.Printf("\nstream count :%v\nmsg : %v\n",reply.GetStreamedCount(),reply.GetMessage())
}
