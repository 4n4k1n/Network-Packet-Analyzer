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

type Stats_data struct {
	total_packets     int
	captured_duration int
	average_rate      float32
	unique_src_ids    int
	unique_dst_ids    int
	tcp_packets       int
	udp_packets       int
	other_packets     int
	most_active_src   string
	most_active_dst   string
	src_ip_counts     map[string]int
	dst_ip_counts     map[string]int
}
