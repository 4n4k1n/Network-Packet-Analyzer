package main

import (
	"fmt"

	"github.com/dariubs/percent"
)

func prinTrafficSize(stats Stats_data) {

	if stats.traffic_size > GB {
		fmt.Printf("Traffic size            : %d GB\n\n", stats.traffic_size/GB)
	} else if stats.traffic_size > MB {
		fmt.Printf("Traffic size            : %d MB\n\n", stats.traffic_size/MB)
	} else if stats.traffic_size > KB {
		fmt.Printf("Traffic size            : %d KB\n\n", stats.traffic_size/KB)
	} else {
		fmt.Printf("Traffic size            : %d B\n\n", stats.traffic_size)
	}
}

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
	fmt.Printf("Most active destination : %s (%d packets)\n", key, stats_data.dst_ip_counts[key])
	prinTrafficSize(stats_data)
}

// print the header line
func printHeaderLine() {
	fmt.Printf("%-15s  %-15s  %-8s  %-10s  %-10s  %-10s  %-20s  %-20s\n\n", "src IP", "dst IP", "protocol", "src port", "dest port", "bytes", "src host", "dst host")
}

// print the packet data
func sprintPacketData(data Pack_data, size int, stats *Stats_data) string {
	log := fmt.Sprintf("%-15s  %-15s  %-8s  %-10d  %-10d  %-10d  %-20s  %-20s\n", data.src_ip, data.dst_ip, data.protocol, data.src_port, data.dst_port, size, reverseDNS(data.src_ip, stats), reverseDNS(data.dst_ip, stats))
	return log
}
