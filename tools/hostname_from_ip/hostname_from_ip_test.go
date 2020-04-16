package hostname_from_ip

import (
	"../../malaware_server_code/core/handle_connections"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"testing"
)

func TestDetectingHostnameFromIP(t *testing.T) {
	localip, err := handle_connections.GetLocalIp04()
	assert.Nil(t, err)
	log.Println(localip)
	ip := net.ParseIP(localip)
	assert.NotNil(t, ip)
	log.Println(ip)
	hostnames, err := net.LookupAddr(localip)
	assert.Nil(t, err)
	for _, hostname := range hostnames {
		log.Println(hostname)
	}
}

func TestDetectingIPsfromHostname(t *testing.T) {
	localip, err := handle_connections.GetLocalIp04()
	assert.Nil(t, err)
	log.Println(localip)
	ip := net.ParseIP(localip)
	assert.NotNil(t, ip)
	log.Println(ip)
	hostnames, err := net.LookupAddr(localip)
	localHost := hostnames[0]
	log.Println(localHost)
	ips, err := net.LookupHost(localHost)
	assert.Nil(t, err)
	for _, ip := range ips {
		log.Println(ip)
	}
}
