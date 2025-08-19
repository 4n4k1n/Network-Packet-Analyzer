package main

import "net"

func reverseDNS(ip string, stats *Stats_data) string {

	if host, exists := stats.dns_cache[ip]; exists {
		return host
	}

	names, err := net.LookupAddr(ip)
	if err != nil || len(names) == 0 {
		stats.dns_cache[ip] = names[0]
		return ip
	}

	stats.dns_cache[ip] = names[0]
	return names[0]
}
