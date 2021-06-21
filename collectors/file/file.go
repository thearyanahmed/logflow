package file

import (
	"errors"
	"github.com/nxadm/tail"
	"github.com/thearyanahmed/logflow/collectors"
	"sync"
)

type collector struct {
	tail *tail.Tail
	collectors.CollectorInterface
	lock *sync.Mutex
}

func NewCollector(options collectors.CollectorOptions) (*collector,chan string,error) {
	t, err := tail.TailFile(options.FilePath, tail.Config{
		ReOpen:      true,
		MustExist:   true,
		Poll:        false,
		Pipe:        false,
		Follow:      true,
		MaxLineSize: 0,
		RateLimiter: nil,
		Logger:      nil,
	})

	if err != nil {
		return nil, nil,err
	}

	lock := sync.Mutex{}

	c := &collector{
		tail: t,
		lock: &lock,
	}

	ch := make(chan string)

	return c, ch, nil
}

func (c *collector) Read(ch chan<- string, wg *sync.WaitGroup) {
	c.lock.Lock()

	defer c.lock.Unlock()
	defer wg.Done()
	defer close(ch)

	defer c.tail.Done()
	defer c.tail.Cleanup()

	//var err error

	for line := range c.tail.Lines {
		ch <- line.Text

		//if err = c.tail.Err(); err != nil {
			//fmt.Println("err",err.Error())
			//return
		//}
	}
}


func (c *collector) FKill(reason string) error {
	c.tail.Kill(errors.New(reason))

	return nil
}
