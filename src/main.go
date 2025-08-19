package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	var err error
	var handle *pcap.Handle
	startTime := time.Now()
	var data Pack_data
	var stats_data Stats_data = Stats_data{
		src_ip_counts: make(map[string]int),
		dst_ip_counts: make(map[string]int),
		dns_cache:     make(map[string]string),
		service_cache: make(map[uint16]string)}

	parse_data := parse()

	loadPortMappings(&stats_data)

	handle, err = pcap.OpenLive(*parse_data.device, 1024, true, time.Microsecond*10)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	file, err := os.OpenFile("logs/packets.log", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	filterInput(parse_data, handle)

	printHeaderLine()

	stats_data.traffic_size = 0
	pack_src := gopacket.NewPacketSource(handle, handle.LinkType())
	for pack := range pack_src.Packets() {
		stats_data.total_packets++
		data = getData(pack)
		stats_data.traffic_size += ByteSize(pack.Metadata().Length)
		getServiceName(&data, &stats_data)
		log := sprintPacketData(data, pack.Metadata().Length, &stats_data)
		fmt.Printf("%s", log)
		file.WriteString(log)
		switch data.protocol {
		case layers.IPProtocolTCP:
			stats_data.tcp_packets++
		case layers.IPProtocolUDP:
			stats_data.udp_packets++
		default:
			stats_data.other_packets++
		}
		stats_data.src_ip_counts[data.src_ip]++
		stats_data.dst_ip_counts[data.dst_ip]++
		if time.Since(startTime) > time.Duration(*parse_data.duration)*time.Second {
			break
		}
	}
	stats_data.captured_duration = *parse_data.duration
	stats_data.average_rate = float32(stats_data.total_packets) / float32(stats_data.captured_duration)
	printStats(stats_data)
}
