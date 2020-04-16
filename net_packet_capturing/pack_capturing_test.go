package net_packet_capturing

import (
	"../net_packet_capturing"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"strings"
	"testing"
	"time"
)

func TestFindAllLocalDevices(t *testing.T) {
	//getmac /fo csv /v
	err := pcap.LoadWinPCAP()
	if err != nil {
		log.Fatal(err)
	}
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	//80:a5:89:9e:60:87

	for _, dev := range devices {
		fmt.Println("Divice Name:", dev.Name)
		fmt.Println("Device description:", dev.Description)
		fmt.Println("Device address:")

		for _, addrs := range dev.Addresses {
			fmt.Println("\tIP:", addrs.IP)
			fmt.Println("\tNet Mask:", addrs.Netmask)
			fmt.Println("\tBroad:", addrs.Broadaddr)
			fmt.Println("\tP2P:", addrs.P2P)
		}
	}
}

func TestGetDeviceNameFromInterface(t *testing.T) {
	ip, err := net_packet_capturing.GetLocalInterfaceIP4FromName("Wi-Fi")
	assert.Nil(t, err)

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	var deviceName string
	for _, dev := range devices {
		for _, addrs := range dev.Addresses {
			if addrs.IP.String() == ip {
				deviceName = dev.Name
			}
		}
	}
	assert.NotEmpty(t, deviceName)
	log.Println(deviceName)

}

func TestGetDeviceNameFromInterfaceFunc(t *testing.T) {
	deviceName, err := net_packet_capturing.GetDeviceNameFromInterface("Wi-Fi")
	assert.Nil(t, err)
	assert.NotEmpty(t, deviceName)
	log.Println(deviceName)
}

func TestFindDeviceIP6Addr(t *testing.T) {
	localIpv6 := "fe80::9df5:35e5:9971"
	//netsh int ipv6 show addresses
	err := pcap.LoadWinPCAP()
	if err != nil {
		log.Fatal(err)
	}
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	//80:a5:89:9e:60:87
	found := false
	for _, dev := range devices {
		for _, addrs := range dev.Addresses {
			fmt.Println("\tIP:", addrs.IP)
			if strings.Contains(addrs.IP.String(), localIpv6) {
				found = true
			}

		}
	}
	assert.True(t, found)
}

func TestEthernetConnectionDeviceNamePresent(t *testing.T) {
	var ethernetDeviceName string
	isNetworkAdapterThere := func() bool {
		devices, err := pcap.FindAllDevs()
		if err != nil {
			log.Fatal(err)
		}
		for _, dev := range devices {
			if strings.Contains(dev.Description, " Ethernet Connection") {
				ethernetDeviceName = dev.Name
				return true
			}
		}
		return false
	}
	assert.True(t, isNetworkAdapterThere())
	fmt.Println(ethernetDeviceName)

}

func TestRetrieveLocalNetworkDevice(t *testing.T) {
	device, err := net.InterfaceByName("Wi-Fi")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, device)
	addrs, err := device.Addrs()
	assert.Nil(t, err)
	for _, addr := range addrs {
		fmt.Println(strings.Split(addr.String(), "/")[0])
	}
}

func TestGetLocalInterfaceIP4FromName(t *testing.T) {
	ip, err := net_packet_capturing.GetLocalInterfaceIP4FromName("Wi-Fi")
	assert.Nil(t, err)
	log.Println(ip)

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
func TestListenToallPacketsFromDevice(t *testing.T) {
	device := "\\Device\\NPF_{8470239F-5090-4316-A7AE-0427DCB73C4F}"
	snapshotlen := int32(1024)
	promiscuous := false
	time_out := 30 * time.Second
	var handle *pcap.Handle

	handle, err := pcap.OpenLive(device, snapshotlen, promiscuous, time_out)
	assert.Nil(t, err)
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		log.Println(packet)
	}
}

func TestListenToTCPPort80traffic(t *testing.T) {
	device := "\\Device\\NPF_{8470239F-5090-4316-A7AE-0427DCB73C4F}"
	snapshotlen := int32(1024)
	promiscuous := false
	time_out := 30 * time.Second
	var handle *pcap.Handle

	handle, err := pcap.OpenLive(device, snapshotlen, promiscuous, time_out)
	assert.Nil(t, err)
	var filter = "TCP and port 80"
	handle.SetBPFFilter(filter)
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		log.Println(packet)
	}
}

func TestGetOnePacketFromEthernetDevice(t *testing.T) {
	var ethernetDeviceName string
	isNetworkAdapterThere := func() bool {
		devices, err := pcap.FindAllDevs()
		if err != nil {
			log.Fatal(err)
		}
		for _, dev := range devices {
			if strings.Contains(dev.Description, "Ethernet Connection") {
				ethernetDeviceName = dev.Name
				return true
			}
		}
		return false
	}
	assert.True(t, isNetworkAdapterThere())
	// Open Device
	handle, err := pcap.OpenLive(ethernetDeviceName, net_packet_capturing.SnapLen, net_packet_capturing.Promisc, net_packet_capturing.Timeout)
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
func TestSamplePacketsFromDevice(t *testing.T) {
	deviceName, err := net_packet_capturing.GetDeviceNameFromInterface("Wi-Fi")
	assert.Nil(t, err)
	assert.NotEmpty(t, deviceName)
	log.Println(deviceName)
	// Open Device
	handle, err := pcap.OpenLive(deviceName, net_packet_capturing.SnapLen, net_packet_capturing.Promisc, net_packet_capturing.Timeout)
	if err != nil {
		log.Fatal("OpenLive Call: ", err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		fmt.Println("========== Packet Layers ==========")
		for _, layer := range packet.Layers() {
			fmt.Println(layer)
		}
		fmt.Println("==================================")
		// Get IPv4 layer
		ip4Layer := packet.Layer(layers.LayerTypeIPv4)
		if ip4Layer != nil {
			fmt.Println("IPv4 layer detected.")
			ip, _ := ip4Layer.(*layers.IPv4)
			fmt.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
			fmt.Println("Protocol:", ip.Protocol)
			fmt.Println()
		} else {
			fmt.Println("No Ipv4 layer Detected")
		}

		//Get TCP layer
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer != nil {
			fmt.Println("TCP layer detected")
			tcp, _ := tcpLayer.(*layers.TCP)
			fmt.Println("ACK: ", tcp.ACK)
			fmt.Println("SYN: ", tcp.SYN)
			fmt.Println("seq: ", tcp.Seq)
			fmt.Println("DstPort: ", tcp.DstPort)
			fmt.Println("SrcPort: ", tcp.SrcPort)
		} else {
			fmt.Println("No TCP layer detected.")
		}
	}
}
