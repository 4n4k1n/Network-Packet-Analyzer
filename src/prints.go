package main

import (
	"fmt"

	"github.com/dariubs/percent"
)

func printStats(stats_data Stats_data) {
	fmt.Printf("\n=== Capture Summary ===\n")
	fmt.Printf("Total packets captured  : %d\n", stats_data.total_packets)
	fmt.Printf("Capture duration        : %d\n", stats_data.captured_duration)
	fmt.Printf("Average rate            : %.1f\n", stats_data.average_rate)
	fmt.Printf("Unique source IPs       : %d\n", len(stats_data.src_ip_counts))
	fmt.Printf("Unique destination IPs  : %d\n", len(stats_data.dst_ip_counts))
	fmt.Printf("Protocol breakdown\n")
	fmt.Printf("   TCP : %d packets (%.1f)\n", stats_data.tcp_packets, percent.PercentOf(stats_data.tcp_packets, stats_data.total_packets))
	fmt.Printf("   UDP : %d packets (%.1f)\n", stats_data.udp_packets, percent.PercentOf(stats_data.udp_packets, stats_data.total_packets))
	fmt.Printf("Most active source      : %s (%d packets)\n", stats_data.most_active_src, 0)
	fmt.Printf("Most active destination : %s (%d packets)\n\n", stats_data.most_active_dst, 0)
}

func printHeaderLine() {
	fmt.Printf("%-15s  %-15s  %-8s  %-10s  %-10s\n", "src IP", "dest IP", "protocol", "src port", "dest port")
}

func printPacketData(data Data) {
	fmt.Printf("%-15s  %-15s  %-8s  %-10d  %-10d\n", data.src_ip, data.dst_ip, data.protocol, data.src_port, data.dst_port)
}
