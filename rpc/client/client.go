package client

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

type RpcClientInterface interface {
	Send(data string) error
	Terminate() error
	Wait()
	Add()
}

type RpcClient struct {
	RpcClientInterface
	lock *sync.Mutex
	wg *sync.WaitGroup
	r packet.LogService_StreamLogClient
}

func NewRpcClient() (*RpcClient,error) {
	host := env.Get("RPC_HOST") + ":" + env.Get("RPC_PORT")

	fmt.Printf("trying to connect %v\n",host)

	conn, err := grpc.Dial(host, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		fmt.Printf("error occured %v %v\n",host,err.Error())
		return nil, err
	}

	fmt.Printf("connected to %v\n",host)

	c := packet.NewLogServiceClient(conn)

	ctx := context.Background()

	r , err := c.StreamLog(ctx)

	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	rpcClient := RpcClient{
		lock:               &mutex,
		wg:                 &wg,
		r:                  r,
	}

	return &rpcClient,nil
}

func (rc *RpcClient) Terminate() (*packet.LogResponse, error) {
	fmt.Printf("terminating response 1")
	fmt.Printf("terminating response 2")
	rc.lock.Lock()
	defer rc.wg.Done()
	defer rc.lock.Unlock()

	return rc.r.CloseAndRecv()
}

// Wait rethink these
func (rc *RpcClient) Wait() {
	rc.wg.Wait()
}

func (rc *RpcClient) Add() {
	rc.wg.Add(1)
}

func (rc *RpcClient) Send(data string) error {
	rc.lock.Lock()

	defer rc.wg.Done()
	defer rc.lock.Unlock()

	json, err := json2.Marshal(data)

	logRequest := packet.LogRequest{
		Topics: []string{"hello_world"}, // todo handle topic
		Payload: json,
	}

	if err != nil {
		return err
	}

	return rc.r.Send(&logRequest)

}





func Run()  {
	fmt.Printf("running client\n")

	host := env.Get("RPC_HOST") + ":" + env.Get("RPC_PORT")

 	conn, err := grpc.Dial(host, grpc.WithInsecure(), grpc.WithBlock())

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

	defer file.Close()

	scanner := bufio.NewScanner(file)

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	for scanner.Scan() {

		fmt.Printf("%v\n",scanner.Text())

		wg.Add(1)

		go func(mt *sync.Mutex, wg *sync.WaitGroup) {
			mt.Lock()

			defer wg.Done()
			defer mt.Unlock()

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


