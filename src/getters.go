package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// get protocol (TCP/UDP) and get the in/out port
func getLayerPortData(pack gopacket.Packet, data *Data) {

	switch data.protocol {
	case layers.IPProtocolTCP:
		tcp_layer := pack.Layer(layers.LayerTypeTCP)
		data.tcp = tcp_layer.(*layers.TCP)
		data.src_port = uint16(data.tcp.SrcPort)
		data.dst_port = uint16(data.tcp.DstPort)
	case layers.IPProtocolUDP:
		udp_layer := pack.Layer(layers.LayerTypeUDP)
		data.udp = udp_layer.(*layers.UDP)
		data.src_port = uint16(data.udp.SrcPort)
		data.dst_port = uint16(data.udp.DstPort)
	}
}

// get the data of the packet (IPv4 layer, ports, protocol, IPs)
func getData(pack gopacket.Packet) Data {
	var data Data

	ip_layer := pack.Layer(layers.LayerTypeIPv4)
	if ip_layer == nil {
		data.src_ip = "N/A"
		data.dst_ip = "N/A"
		data.protocol = 0
		return data
	}

	ipv4 := ip_layer.(*layers.IPv4)
	data.src_ip = ipv4.SrcIP.String()
	data.dst_ip = ipv4.DstIP.String()
	data.protocol = ipv4.Protocol
	getLayerPortData(pack, &data)

	return data
}
