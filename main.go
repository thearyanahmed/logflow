package main

import (
	"github.com/thearyanahmed/nlogx/kafka"
	"github.com/thearyanahmed/nlogx/utils/env"
	"sync"
)

func main() {
	env.LoadEnv()

	var wg sync.WaitGroup

	wg.Add(1)

	go kafka.BootAndListen()

	wg.Wait()

}