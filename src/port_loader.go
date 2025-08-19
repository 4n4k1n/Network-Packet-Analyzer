package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
)

// Load port mappings from JSON file into service_cache
func loadPortMappings(stats *Stats_data) {
	data, err := ioutil.ReadFile("data/ports.json")
	if err != nil {
		log.Printf("Warning: Could not read ports.json: %v", err)
		return
	}

	var portMap map[string]string
	err = json.Unmarshal(data, &portMap)
	if err != nil {
		log.Printf("Warning: Could not parse ports.json: %v", err)
		return
	}

	for portStr, service := range portMap {
		port, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			log.Printf("Warning: Invalid port number '%s': %v", portStr, err)
			continue
		}
		stats.service_cache[uint16(port)] = service
	}

	log.Printf("Loaded %d port mappings", len(stats.service_cache))
}

// Get service name for a port
func getPortName(port uint16, stats *Stats_data) string {
	if service, exists := stats.service_cache[port]; exists {
		return service
	}
	return "Unknown"
}

func getServiceName(pack *Pack_data, stats *Stats_data) {

	src_service := getPortName(pack.src_port, stats)
	dst_service := getPortName(pack.dst_port, stats)

	if src_service != "Unknown" {
		pack.service = src_service
	} else if dst_service != "Unknown" {
		pack.service = dst_service
	} else {
		pack.service = "Unknown"
	}
}
