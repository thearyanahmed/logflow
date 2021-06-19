package main

import (
	"context"
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

	defer conn.Close()

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
			pckt := &packet.Packet{
				Version:    i,
				IvVersion:  nil,
				AgentIp:    2288332211,
				SubAgentId: 55772288,
			}

			nginxLog := &packet.NginxLog{
				BytesSent:               762,
				Connection:              "connection",
				ConnectionRequests:      666,
				Status:                  1,
				Host:                    "my host",
				NginxVersion:            "nginx version",
				ProxyProtocolAddr:       "proxy protocol addr",
				ProxyProtocolPort:       1234,
				ProxyProtocolServerAddr: "proxy protocol server addr",
				ProxyProtocolServerPort: 2345,
				RemoteAddr:              "remote addr",
				RemotePort:              3456,
				RemoteUser:              "remote user",
				RequestMethod:           "get",
				ServerAddr:              "server addr",
				ServerName:              "server name",
				ServerPort:              5678,
				Endpoint:                "/end/point",
				HttpVersion:             "http1.1",
				UserAgent:               "safari",
			}

			var headers []*packet.Header

			header := &packet.Header{
				Key:   "hello",
				Value: "world",
			}

			headers = append(headers,header)

			logRequest := packet.LogRequest{
				Packet:  pckt,
				Log:     nginxLog,
				Headers: headers,
				Topics: []string{"hello_world","banana_world"},
			}

			err := r.Send(&logRequest)

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
