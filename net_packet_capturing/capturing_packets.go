package net_packet_capturing

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	"strings"
	"time"
)

var (
	DevFound = false
	SnapLen  = int32(65535)
	Timeout  = -1 * time.Second
	Filter   = "tcp and port 80"
	Promisc  = false
)

func GetLocalInterfaceIP4FromName(inface string) (string, error) {
	device, err := net.InterfaceByName(inface)
	if err != nil {
		return "", err
	}
	addrs, err := device.Addrs()
	if err != nil {
		return "", err
	}
	ip4Name := strings.Split(addrs[1].String(), "/")[0]
	return ip4Name, nil
}

type deviceNotFoundErr struct {
	Message string
}

func (de *deviceNotFoundErr) Error() string {
	return de.Message
}

func GetDeviceNameFromInterface(inface string) (string, error) {

	ip, err := GetLocalInterfaceIP4FromName(inface)

	if err != nil {
		return "", err
	}

	devices, err := pcap.FindAllDevs()
	if err != nil {
		return "", err
	}
	var deviceName string
	for _, dev := range devices {
		for _, addrs := range dev.Addresses {
			if addrs.IP.String() == ip {
				deviceName = dev.Name
			}
		}
	}
	if deviceName == "" {
		error := &deviceNotFoundErr{fmt.Sprintf("Device:%s not found \n", inface)}
		return "", error
	}
	log.Println(deviceName)
	return deviceName, nil
}

func PrintPacketInfo(packet gopacket.Packet) {
	//Ethernet layer
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		fmt.Println("**************Ethernet Layer Detected**************")
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		fmt.Println("Source MAC:", ethernetPacket.SrcMAC)
		fmt.Println("Destination MAC:", ethernetPacket.SrcMAC)
		fmt.Println("Ethernet Type:", ethernetPacket.EthernetType)
		fmt.Println("***************************************************")
	}
	//IP Layer
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		fmt.Println("**************IPV4 layer detected******************")
		ip, _ := ipLayer.(*layers.IPv4)
		fmt.Println("Source IP:", ip.SrcIP)
		fmt.Println("Destination IP:", ip.DstIP)
		fmt.Println("Protocol: ", ip.Protocol)
	}
	//TCP Layer
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		fmt.Println("**************TCP Layer Detected******************")
		fmt.Println("From port:", tcp.SrcPort)
		fmt.Println("From port:", tcp.DstPort)
		fmt.Println("Sequence n;", tcp.Seq)
	}

}
func FindPayloadStringMatchFor(searchStringMatch string, packet gopacket.Packet, silentMode bool) {
	for _, layer := range packet.Layers() {
		if silentMode != true {
			fmt.Println("***************All layers********************")
			fmt.Println("Layer Type:", layer.LayerType())
		}
		applicationLayer := packet.ApplicationLayer()
		if applicationLayer != nil {
			if silentMode != true {
				fmt.Println("Application layer/Payload found.")
				fmt.Printf("%s\n", string(applicationLayer.Payload()))
			}
			if strings.Contains(string(applicationLayer.Payload()), searchStringMatch) {
				fmt.Println("Match found for string:", searchStringMatch)
			}
		}
		if err := packet.ErrorLayer(); err != nil {
			fmt.Println("Error decoding some part of the packet")
		}
	}
}
