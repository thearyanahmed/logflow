package file

import (
	"fmt"
	"github.com/nxadm/tail"
)

func Run(filepath string) {

	t, err := tail.TailFile(filepath, tail.Config{Follow: true, ReOpen: true})

	if err != nil {
		fmt.Printf("error creating tail %v\n",err.Error())
		return
	}

	// Print the text of each received line
	for line := range t.Lines {
		fmt.Println(line.Text)
	}
}
