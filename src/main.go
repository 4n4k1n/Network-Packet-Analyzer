package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func getIPs(pack gopacket.Packet, src_ip *string, dst_ip *string, protocol *layers.IPProtocol) {
	ip_layer := pack.Layer(layers.LayerTypeIPv4)
	if ip_layer == nil {
		*src_ip = "N/A"
		*dst_ip = "N/A"
		*protocol = 0
		return
	}
	ipv4 := ip_layer.(*layers.IPv4)
	*src_ip = ipv4.SrcIP.String()
	*dst_ip = ipv4.DstIP.String()
	*protocol = ipv4.Protocol
}

func main() {
	var err error
	var handle *pcap.Handle
	var count int32 = 0
	startTime := time.Now()
	var src_ip string
	var dst_ip string
	var Protocol layers.IPProtocol

	handle, err = pcap.OpenLive("wlan0", 1024, true, time.Microsecond*10)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	pack_src := gopacket.NewPacketSource(handle, handle.LinkType())
	for pack := range pack_src.Packets() {
		count++
		getIPs(pack, &src_ip, &dst_ip, &Protocol)
		fmt.Printf("srcIP: %s, DstIP: %s, Prot: %s\n", src_ip, dst_ip, Protocol)
		if time.Since(startTime) > 30*time.Second {
			break
		}
	}

	// fmt.Printf("Network interface: %s\nPackets captured: %d\n", device, count)
}
