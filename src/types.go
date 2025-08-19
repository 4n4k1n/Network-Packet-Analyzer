package main

import (
	"github.com/google/gopacket/layers"
)

type ByteSize int64

const (
	B  ByteSize = 1 << (10 * iota) // 1 << 0 = 1
	KB ByteSize = 1 << (10 * iota) // 1 << 10 = 1024
	MB ByteSize = 1 << (10 * iota) // 1 << 20 = 1048576
	GB ByteSize = 1 << (10 * iota) // 1 << 30 = 1073741824
)

// struct for the pasing data
type Parse_data struct {
	device   *string
	duration *int
	// protocol *string
	// port     *string
	// ip       *string
	filter_items []string
}

// general data struct
type Pack_data struct {
	src_ip   string
	dst_ip   string
	protocol layers.IPProtocol
	tcp      *layers.TCP
	udp      *layers.UDP
	src_port uint16
	dst_port uint16
	service  string
}

// struct for the stats data
type Stats_data struct {
	total_packets     int
	captured_duration int
	average_rate      float32
	tcp_packets       int
	udp_packets       int
	other_packets     int
	src_ip_counts     map[string]int
	dst_ip_counts     map[string]int
	traffic_size      ByteSize
	dns_cache         map[string]string
	service_cache     map[uint16]string
}
