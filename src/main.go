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
		stats_data.total_packets++
		data = getData(pack)
		printPacketData(data)
		if time.Since(startTime) > time.Duration(*parse_data.duration)*time.Second {
			break
		}
	}
	stats_data.captured_duration = *parse_data.duration
	stats_data.average_rate = float32(stats_data.total_packets) / float32(stats_data.captured_duration)
	printStats(stats_data)
}
