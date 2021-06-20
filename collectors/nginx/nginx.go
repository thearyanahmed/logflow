package main

import (
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

type dataSource struct {
	gopacket.PacketDataSource
}

type httpStream struct {
	net 		gopacket.Flow
	transport   gopacket.Flow
	r 			tcpreader.ReaderStream
}

var (

	iface 		  = flag.String("i", "eth0", "Interface to get packets from")
	fname 		  = flag.String("r", "", "Filename to read from, overrides -i")
	snaplen 	  = flag.Int("s", 1600, "SnapLen for pcap packet capture")
	filter 	      = flag.String("f", "tcp and dst port 80", "BPF filter for pcap")
	logAllPackets = flag.Bool("v", false, "Logs every packet in great detail")
)

type httpStreamFactory struct{}
//
//type httpStream struct{
//	gopacket.StreamPool
//}

func main()  {
	flag.Parse()

	fmt.Printf("running nginx collector \n args: iface %v\t fname %v\t snaplen %v\t filter %v\t logAllPackets %v\n",*iface,*fname,*snaplen,*filter,*logAllPackets)

	var handle *pcap.Handle
	var err error

	handle, err = pcap.OpenLive(*iface, int32(*snaplen), true, pcap.BlockForever)

	if err != nil {
		fmt.Printf("could not open connection %v\n",err.Error())
		return
	}

	if err = handle.SetBPFFilter(*filter); err != nil {
		fmt.Printf("could not set BPF Filter: %v\n",err.Error())
		return
	}

	// Set up assembly
	//streamFactory := &httpStreamFactory{}
	//streamPool := tcpassembly.NewStreamPool(streamFactory)
	//assembler := tcpassembly.NewAssembler(streamPool)


	packetSource := gopacket.NewPacketSource(handle,handle.LinkType())

	for packet := range packetSource.Packets() {
		//fmt.Printf("%v\n",packet.)
		if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			fmt.Printf("unstable packet. layer : %v\n",packet.LinkLayer())
			continue
		}

		if app := packet.ApplicationLayer(); app != nil {
			fmt.Printf("application layer : %v\n",app.Payload())
		}

		fmt.Printf("packet data: %v\n packet metadata timestamp: %v\n",string(packet.Data()[:]),packet.Metadata().CaptureInfo.Timestamp)
	}



	//for {
	//	select {
	//	case packet := <-packets:
	//		if packet == nil {
	//			fmt.Printf("nil packet,end of pcap file, if we are using it\n")
	//			return
	//		}
	//	}
	//}
}