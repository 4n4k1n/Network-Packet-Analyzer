package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func main() {
	var err error
	var handle *pcap.Handle
	var count int32 = 0
	startTime := time.Now()
	var data Data

	parse_data := parse()

	handle, err = pcap.OpenLive(*parse_data.device, 1024, true, time.Microsecond*10)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	filterInput(parse_data, handle)

	fmt.Printf("%-15s  %-15s  %-8s  %-10s  %-10s\n", "src IP", "test IP", "protocol", "src port", "dest port")

	pack_src := gopacket.NewPacketSource(handle, handle.LinkType())
	for pack := range pack_src.Packets() {
		count++
		data = getData(pack)
		fmt.Printf("%-15s  %-15s  %-8s  %-10d  %-10d\n", data.src_ip, data.dst_ip, data.protocol, data.src_port, data.dst_port)
		// fmt.Printf("srcIP: %s, DstIP: %s, Prot: %s\n", data.src_ip, data.dst_ip, data.protocol)
		if time.Since(startTime) > time.Duration(*parse_data.duration)*time.Second {
			break
		}
	}

	// fmt.Printf("Network interface: %s\nPackets captured: %d\n", device, count)
}
