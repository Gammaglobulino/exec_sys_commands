package malaware_victim_client

import (
	"fmt"
	"log"
	"net"
)

const (
	ServerIP = "10.0.2.15"
	Port     = ":9090"
)

func main() {

	LocalIPAddr := ServerIP + Port
	connection, err := net.Dial("tcp", LocalIPAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(connection.RemoteAddr().String())
}
