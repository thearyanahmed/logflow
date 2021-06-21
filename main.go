package main

import (
	"flag"
	"fmt"
	"github.com/thearyanahmed/logflow/collectors"
	"github.com/thearyanahmed/logflow/collectors/file"
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

		wg := sync.WaitGroup{}
		wg.Add(1)

		chOpen := true

		go collector.Read(ch,&wg)

		for chOpen {
			select {
			case xLine, open := <- ch:
				fmt.Printf("line : %v\n",xLine)

				fmt.Printf("open %v\n",open)


				if open == false {
					chOpen = false
				}
			}
		}


		fmt.Printf("\ndone\n")
		wg.Wait()
		// create a client instance
		// when creating an instance, it will create a grpc dial
		// which will hold the connect
		// a lock and possibly an wg

		// client should have a function


		//client.Run()

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