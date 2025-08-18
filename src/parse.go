package main

import (
	"flag"
	"strings"

	"github.com/google/gopacket/pcap"
)

// pase the analyzer args using the flag libary
func parse() Parse_data {
	var parse_data Parse_data

	parse_data.duration = flag.Int("time", 30, "Duration of the program in seconds!")
	parse_data.device = flag.String("device", "wlan0", "Get the device name!")

	ip := flag.String("ip", "", "ip")
	protocol := flag.String("protocol", "", "protocol")
	port := flag.String("port", "", "port")
	flag.Parse()

	if *ip != "" {
		parse_data.filter_items = append(parse_data.filter_items, "host "+*ip)
	}
	if *protocol != "" {
		parse_data.filter_items = append(parse_data.filter_items, *protocol)
	}
	if *port != "" {
		parse_data.filter_items = append(parse_data.filter_items, "port "+*port)
	}

	return parse_data
}

// join the args and filter the packets using BPF
func filterInput(parse_data Parse_data, handle *pcap.Handle) {

	filter_str := strings.Join(parse_data.filter_items, " and ")

	if filter_str != "" {
		handle.SetBPFFilter(filter_str)
	}
}
