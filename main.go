package main

import (
	"flag"
	"fmt"
	"github.com/thearyanahmed/logflow/rpc/client"
	"github.com/thearyanahmed/logflow/rpc/server"
	"github.com/thearyanahmed/logflow/utils/env"
)

var (
	action string
)

func main()  {
	env.LoadEnv()

	flag.StringVar(&action, "action", "", "start rpc server. available option: server, client")

	flag.Parse()

	switch action {
	case "serve":
		server.Run()
	case "client":
		client.Run()
	default:
		fmt.Printf("You did not select any proper option.quiting\n")
	}

}