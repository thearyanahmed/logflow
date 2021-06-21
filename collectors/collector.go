package collectors

import "sync"

type CollectorOptions struct {
	FilePath string
}

type CollectorInterface interface {
	Read(ch chan string, wg *sync.WaitGroup)

	FKill(reason string) error
}
