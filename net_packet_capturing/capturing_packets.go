package net_packet_capturing

import (
	"fmt"
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
