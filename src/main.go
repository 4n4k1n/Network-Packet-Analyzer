package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

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

func parse() Parse_data {
	var parse_data Parse_data

	parse_data.duration = flag.Int("time", 30, "Duration of the program in seconds!")
	parse_data.device = flag.String("device", "wlan0", "Get the device name!")
	parse_data.filter_items[0] = *flag.String("ip", "", "ip")
	parse_data.filter_items[1] = *flag.String("protocol", "", "protocol")
	parse_data.filter_items[2] = *flag.String("port", "", "port")
	flag.Parse()
	return parse_data
}

func main() {
	var err error
	var handle *pcap.Handle
	var count int32 = 0
	startTime := time.Now()
	var data Data

	parse_data := parse()

	handle, err = pcap.OpenLive(*parse_data.device, 1024, true, time.Microsecond*10)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// handle.SetBPFFilter()

	pack_src := gopacket.NewPacketSource(handle, handle.LinkType())
	for pack := range pack_src.Packets() {
		count++
		data = getData(pack)
		fmt.Printf("%-15s  %-15s  %-8s  %-10d  %-10d\n", data.src_ip, data.dst_ip, data.protocol, data.src_port, data.dst_port)
		// fmt.Printf("srcIP: %s, DstIP: %s, Prot: %s\n", data.src_ip, data.dst_ip, data.protocol)
		if time.Since(startTime) > time.Duration(*parse_data.duration)*time.Second {
			break
		}
	}

	// fmt.Printf("Network interface: %s\nPackets captured: %d\n", device, count)
}
