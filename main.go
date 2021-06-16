package main

import (
	"context"
	"fmt"
	"github.com/thearyanahmed/nlogx/kafka"
	"github.com/thearyanahmed/nlogx/utils/env"
)

func main() {
	fmt.Println("startx")

	env.LoadEnv()

	ctx := context.Background()

	kafka.Consume(ctx)

	fmt.Printf("End of the line ")
}