package net_packet_capturing

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"strings"
	"time"
)

var (
	iface    = `\Device\NPF_{74F64D39-4E0C-40D4-B041-B17D82CF0536}`
	DevFound = false
	SnapLen  = int32(65535)
	Timeout  = -1 * time.Second
	Filter   = "tcp and port 80"
	Promisc  = false
)

func ListenToTCPPackets() {

	// you must have permission to capture packet- try using sudo su
	//trying entering username and password to http://diptera.myspecies.info/

	for {

		handle, err := pcap.OpenLive(iface, SnapLen, Promisc, Timeout)
		if err != nil {
			log.Panic(err)
		}

		err = handle.SetBPFFilter(Filter)
		if err != nil {
			log.Panicln(err)
		}

		src := gopacket.NewPacketSource(handle, handle.LinkType())
		search_array := []string{"name", "username", "pass"}

		for pkt := range src.Packets() {
			applayer := pkt.ApplicationLayer()
			if applayer == nil {
				continue
			}
			payload := applayer.Payload()
			for _, s := range search_array {
				index := strings.Index(string(payload), s)
				if index != -1 {
					fmt.Println(string(payload[index : index+100]))
				}

			}
		}
		handle.Close()

	}

}
