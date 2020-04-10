package handle_connections

import (
	"fmt"
	"net"
)

func ConnectTo(ip string) (net.Conn, error) {
	connection, err := net.Dial("tcp", ip)
	if err != nil {
		return nil, err
	}
	fmt.Println(connection.RemoteAddr().String())
	return connection, nil
}
