package nginx

import (
	"fmt"
	"github.com/google/gopacket"
)

type dataSource struct {
	gopacket.PacketDataSource
}

func (ds *dataSource) ReadPacketData(data []byte, ci gopacket.CaptureInfo, err error) {

}

func Run()  {
	fmt.Printf("running nginx collector\n")

}