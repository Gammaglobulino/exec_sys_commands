package net_packet_capturing

import (
	"../net_packet_capturing"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"testing"
)

func TestFindAllLocalDevices(t *testing.T) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for index, dev := range devices {
		fmt.Println(index, "Divice Name:", dev.Name)
		for _, addrs := range dev.Addresses {
			fmt.Println("\tIP:", addrs.IP)
			fmt.Println("\tNet Mask:", addrs.Netmask)
			fmt.Println("\tBroad:", addrs.Broadaddr)

		}
	}
}

func TestRetrieveLocalNetworkDevice(t *testing.T) {
	device, err := net.InterfaceByName("Ethernet")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, device)
	fmt.Println(device.Name)
}

func TestListenToTCPPackets(t *testing.T) {
	net_packet_capturing.ListenToTCPPackets()
}

func TestLocalNetworkInterfacesNames(t *testing.T) {

	ifaces, err := net.Interfaces()
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range ifaces {
		fmt.Println(i)

	}
}
func TestGetPacketFromDevice(t *testing.T) {
	// Open Device
	handle, err := pcap.OpenLive(`\Device\NPF_{74F64D39-4E0C-40D4-B041-B17D82CF0536}`, net_packet_capturing.SnapLen, net_packet_capturing.Promisc, net_packet_capturing.Timeout)
	if err != nil {
		log.Fatal("OpenLive Call: ", err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Get the next packet
	packet, err := packetSource.NextPacket()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(packet)
}
