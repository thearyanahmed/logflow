package main

import (
	"fmt"
	"github.com/thearyanahmed/nlogx/utils/env"
)

func main() {
	fmt.Println("startx")

	env.LoadEnv()

	fmt.Println(env.Get("APP_ENV"))
}