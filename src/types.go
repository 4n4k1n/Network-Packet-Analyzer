package main

import "github.com/google/gopacket/layers"

type Parse_data struct {
	device   *string
	duration *int
	// protocol *string
	// port     *string
	// ip       *string
	filter_items []string
}

type Data struct {
	src_ip   string
	dst_ip   string
	protocol layers.IPProtocol
	tcp      *layers.TCP
	udp      *layers.UDP
	src_port uint16
	dst_port uint16
}
