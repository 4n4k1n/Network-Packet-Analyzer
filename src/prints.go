package main

import (
	"fmt"

	"github.com/dariubs/percent"
)

// function to print the stats at the end of the program
func printStats(stats_data Stats_data) {
	var key string

	fmt.Printf("\n=== Capture Summary ===\n")
	fmt.Printf("Total packets captured  : %d\n", stats_data.total_packets)
	fmt.Printf("Capture duration        : %d\n", stats_data.captured_duration)
	fmt.Printf("Average rate            : %.1f\n", stats_data.average_rate)
	fmt.Printf("Unique source IPs       : %d\n", len(stats_data.src_ip_counts))
	fmt.Printf("Unique destination IPs  : %d\n", len(stats_data.dst_ip_counts))
	fmt.Printf("Protocol breakdown\n")
	fmt.Printf("   TCP   : %d packets (%.1f)\n", stats_data.tcp_packets, percent.PercentOf(stats_data.tcp_packets, stats_data.total_packets))
	fmt.Printf("   UDP   : %d packets (%.1f)\n", stats_data.udp_packets, percent.PercentOf(stats_data.udp_packets, stats_data.total_packets))
	fmt.Printf("   other : %d packets (%.1f)\n", stats_data.other_packets, percent.PercentOf(stats_data.other_packets, stats_data.total_packets))
	key = getMaxOfMap(stats_data.src_ip_counts)
	fmt.Printf("Most active source      : %s (%d packets)\n", key, stats_data.src_ip_counts[key])
	key = getMaxOfMap(stats_data.dst_ip_counts)
	fmt.Printf("Most active destination : %s (%d packets)\n\n", key, stats_data.dst_ip_counts[key])
}

// print the header line
func printHeaderLine() {
	fmt.Printf("%-15s  %-15s  %-8s  %-10s  %-10s\n", "src IP", "dest IP", "protocol", "src port", "dest port")
}

// print the packet data
func printPacketData(data Data) {
	fmt.Printf("%-15s  %-15s  %-8s  %-10d  %-10d\n", data.src_ip, data.dst_ip, data.protocol, data.src_port, data.dst_port)
}
