package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func main() {
	var device string = "wlan0"
	var snaplen int32 = 1024
	var promisc bool = true
	var err error
	var timeout time.Duration = time.Microsecond * 100
	var handle *pcap.Handle
	var count int32 = 0
	startTime := time.Now()

	handle, err = pcap.OpenLive(device, snaplen, promisc, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	pack_src := gopacket.NewPacketSource(handle, handle.LinkType())
	for pack := range pack_src.Packets() {
		count++
		fmt.Println(pack)
		if time.Since(startTime) > 30*time.Second {
			break
		}
	}

	fmt.Printf("Network interface: %s\nPackets captured: %d\n", device, count)
}
