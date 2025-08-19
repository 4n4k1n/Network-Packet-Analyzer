package main

import (
	"fmt"

	"github.com/dariubs/percent"
)

func prinTrafficSize(stats Stats_data) {

	if stats.traffic_size > GB {
		fmt.Printf("Traffic size            : %.2f GB\n\n", float64(stats.traffic_size)/float64(GB))
	} else if stats.traffic_size > MB {
		fmt.Printf("Traffic size            : %.2f MB\n\n", float64(stats.traffic_size)/float64(MB))
	} else if stats.traffic_size > KB {
		fmt.Printf("Traffic size            : %.2f KB\n\n", float64(stats.traffic_size)/float64(KB))
	} else {
		fmt.Printf("Traffic size            : %d B\n\n", stats.traffic_size)
	}
}

// function to print the stats at the end of the program
func printStats(stats_data Stats_data) {

	fmt.Printf("\n=== Capture Summary ===\n")
	fmt.Printf("Total packets captured  : %d\n", stats_data.total_packets)
	fmt.Printf("Capture duration        : %d\n", stats_data.captured_duration)
	fmt.Printf("Average rate            : %.1f packets/sec\n", stats_data.average_rate)
	fmt.Printf("Unique source IPs       : %d\n", len(stats_data.src_ip_counts))
	fmt.Printf("Unique destination IPs  : %d\n", len(stats_data.dst_ip_counts))
	fmt.Printf("Protocol breakdown\n")
	fmt.Printf("   TCP   : %d packets (%.1f%%)\n", stats_data.tcp_packets, percent.PercentOf(stats_data.tcp_packets, stats_data.total_packets))
	fmt.Printf("   UDP   : %d packets (%.1f%%)\n", stats_data.udp_packets, percent.PercentOf(stats_data.udp_packets, stats_data.total_packets))
	fmt.Printf("   other : %d packets (%.1f%%)\n", stats_data.other_packets, percent.PercentOf(stats_data.other_packets, stats_data.total_packets))
	fmt.Printf("Top active sources      :\n")
	topSrc := getTopNFromMap(stats_data.src_ip_counts, 3)
	for i, entry := range topSrc {
		fmt.Printf("  %d. %s (%d packets)\n", i+1, entry.Key, entry.Value)
	}

	fmt.Printf("Top active destinations :\n")
	topDst := getTopNFromMap(stats_data.dst_ip_counts, 3)
	for i, entry := range topDst {
		fmt.Printf("  %d. %s (%d packets)\n", i+1, entry.Key, entry.Value)
	}

	prinTrafficSize(stats_data)
}

// print the header line
func printHeaderLine() {
	fmt.Printf("%-15s  %-15s  %-8s  %-10s  %-10s  %-8s  %-15s\n\n", "src IP", "dst IP", "protocol", "src port", "dest port", "bytes", "service")
}

// print the packet data
func sprintPacketData(data Pack_data, size int, _ *Stats_data) string {
	// Format size with appropriate unit
	var sizeStr string
	if size > 1024 {
		sizeStr = fmt.Sprintf("%.1fK", float64(size)/1024.0)
	} else {
		sizeStr = fmt.Sprintf("%dB", size)
	}

	log := fmt.Sprintf("%-15s  %-15s  %-8s  %-10d  %-10d  %-8s  %-15s\n", data.src_ip, data.dst_ip, data.protocol, data.src_port, data.dst_port, sizeStr, data.service)
	return log
}
