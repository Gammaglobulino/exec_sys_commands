package main

import (
	"fmt"
	"github.com/google/gopacket/pcap"
	"log"
)

func main() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for index, dev := range devices {
		fmt.Println(index, " ", dev.Name)
		for _, addrs := range dev.Addresses {
			fmt.Println("\tIP:", addrs.IP)
			fmt.Println("\tNet Mask:", addrs.Netmask)
		}
	}
}
