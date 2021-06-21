package main

import (
	"flag"
	"fmt"
	"github.com/thearyanahmed/logflow/collectors"
	"github.com/thearyanahmed/logflow/collectors/file"
	"github.com/thearyanahmed/logflow/rpc/client"
	"github.com/thearyanahmed/logflow/rpc/server"
	"github.com/thearyanahmed/logflow/utils/env"
	"log"
	"os"
	"sync"
)

var (
	action string
)

func main()  {
	env.LoadEnv()

	flag.StringVar(&action, "action", "help", "start rpc server. available option: server, client")

	flag.Parse()

	switch action {
	case "serve":
		server.Run()
	case "client":

		rootDir, err := os.Getwd()

		if err != nil {
			log.Fatalf("could not open directory. %v\n",err.Error())
		}

		rootDir = rootDir + "/" + env.Get("TEST_DATA_FILE")

		// create collectorOptions for creating collector
		collOpts := collectors.CollectorOptions{FilePath: rootDir}

		collector, ch , err := file.NewCollector(collOpts)

		if err != nil {
			log.Fatalf("error creating collector %v\n",err.Error())
		}

		rpcClient , err := client.NewRpcClient()

		if err != nil {
			log.Fatalf("%v\n",err.Error())
		}

		wg := sync.WaitGroup{}
		wg.Add(1)

		chOpen := true

		go collector.Read(ch,&wg)

		for chOpen {
			select {
			case line, open := <- ch:

				fmt.Printf("%v\n%v\n",line,open)

				if open == false {
					chOpen = false
					fmt.Printf("channel closed.\n")
					break
				}

				rpcClient.Add()

				go func() {
					err := rpcClient.Send(line)

					if err != nil {
						fmt.Printf("send error %v.\n",err.Error())
						return
					}

					fmt.Printf("request sent\n")
				}()
			}
		}

		fmt.Printf("before wait")

		rpcClient.Wait()

		fmt.Printf("waiting on rpc")

		wg.Wait()

		resp, err := rpcClient.Terminate()

		if err != nil {
			fmt.Printf("err terminating response %v\n",err.Error())
			return
		}

		fmt.Printf("response terminated \n stream count %v\n msg %v\n",resp.GetStreamedCount(),resp.GetMessage())

	case "help":
		printHelp()
	default:
		fmt.Printf("You did not select any proper option.quiting\n")
	}
}

func printHelp() {
	fmt.Printf(`
	
About : Logflow reads nginx log data and sends to kafka using grpc. 
	
Dependencies: Logflow requires kafka to be install and running on your machine. Make sure .env file is setup correctly.

Available commands    
------------------------------------------------------------------------------------------------
serve				: Will star the gRPC server
client			 	: Will start a gRPC client that reads from nginx log (specified in the .env) and streams to gRPC server. Everything is subject to change.

help				: Displayes help menu.	

How to use : For now, you can use it like
"go run main.go --action serve"
`)
}