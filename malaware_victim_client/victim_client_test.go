package malaware_victim_client

import (
	"../malaware_victim_client"
	"fmt"
	"net"
	"testing"
)

func TestConnectionToServer(t *testing.T) {
	LocalIPAddr := malaware_victim_client.ServerIP + malaware_victim_client.Port
	connection, err := net.Dial("tcp", LocalIPAddr)
	if err != nil {
		t.Fatalf("Connection error %s", err.Error())
	}
	fmt.Println(connection.RemoteAddr().String())
}
