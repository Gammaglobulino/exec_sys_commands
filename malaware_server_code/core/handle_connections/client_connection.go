package handle_connections

import (
	"fmt"
	"net"
)

const (
	ServerIP = "10.0.2.15"
	Port     = ":9090"
)

func ConnectTo(ip string) (net.Conn, error) {
	LocalIPAddr := ServerIP + Port
	connection, err := net.Dial("tcp", LocalIPAddr)
	if err != nil {
		return nil, err
	}
	fmt.Println(connection.RemoteAddr().String())
	return connection, nil
}
