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

	c := &collector{
		tail: t,
	}

	ch := make(chan string)

	return c, ch, err
}

func (c *collector) Read(ch chan<- string, wg *sync.WaitGroup) {

	defer c.tail.Done()
	defer close(ch)
	defer wg.Done()

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
