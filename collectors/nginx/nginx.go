package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"io"
	"log"
	"net/http"
)

var (

	iface 		  = flag.String("i", "eth0", "Interface to get packets from")
	fname 		  = flag.String("r", "", "Filename to read from, overrides -i")
	snaplen 	  = flag.Int("s", 1600, "SnapLen for pcap packet capture")
	filter 	      = flag.String("f", "tcp and dst port 80", "BPF filter for pcap")
	logAllPackets = flag.Bool("v", false, "Logs every packet in great detail")
)

type httpStream struct {
	net 		gopacket.Flow
	transport   gopacket.Flow
	r 			tcpreader.ReaderStream
}

// httpStreamFactory implements tcpassembly.StreamFactory
type httpStreamFactory struct{}


func (h *httpStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hstream := &httpStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}
	go hstream.run() // Important... we must guarantee that data from the reader stream is read.

	// ReaderStream implements tcpassembly.Stream, so we can return a pointer to it.
	return &hstream.r
}

func (h *httpStream) run() {
	buf := bufio.NewReader(&h.r)

	fmt.Printf("transport :%v\n\n",h.transport.String())

	for {

		req, err := http.ReadRequest(buf)

		if err == io.EOF {
			// We must read until we see an EOF... very important!
			fmt.Printf("endof line.")
			return
		} else if err != nil {
			log.Println("Error reading stream", h.net, h.transport, ":", err)
		} else {
			bodyBytes := tcpreader.DiscardBytesToEOF(req.Body)
			req.Body.Close()
			fmt.Printf("Received request from stream \n net: %v\n transport %v\n request %v body bytes %v\n",
				h.net, h.transport,req, bodyBytes)
		}
	}
}

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
	streamFactory := &httpStreamFactory{}
	streamPool    := tcpassembly.NewStreamPool(streamFactory)
	assembler 	  := tcpassembly.NewAssembler(streamPool)

	packetSource := gopacket.NewPacketSource(handle,handle.LinkType())

	for packet := range packetSource.Packets() {
		//fmt.Printf("%v\n",packet.)
		if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			fmt.Printf("unstable packet. layer : %v\n",packet.LinkLayer())
			continue
		}

		tcp := packet.TransportLayer().(*layers.TCP)

		assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
		//
		//for i, l := range packet.Layers()  {
		//	fmt.Printf( "- Layer %d (%02d bytes) = %s\n", i+1, len(l.LayerContents()), gopacket.LayerString(l))
		//}

		//fmt.Printf("packet %v\n\n\n",packet.String())
		//
		//if appLayer := packet.ApplicationLayer(); appLayer != nil {
		//	fmt.Printf("applicaiton layer %v\n\n\n",appLayer.LayerPayload())
		//}
	}
}