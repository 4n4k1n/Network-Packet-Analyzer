package main

import (
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
	var stats_data Stats_data

	parse_data := parse()

	handle, err = pcap.OpenLive(*parse_data.device, 1024, true, time.Microsecond*10)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	filterInput(parse_data, handle)

	printHeaderLine()

	pack_src := gopacket.NewPacketSource(handle, handle.LinkType())
	for pack := range pack_src.Packets() {
		count++
		data = getData(pack)
		printPacketData(data)
		if time.Since(startTime) > time.Duration(*parse_data.duration)*time.Second {
			break
		}
	}
	printStats(stats_data)
}
