package malaware_server_code

import (
	"../Execute_systems_commands"
	"../malaware_server_code"
	"../malaware_victim_client"
	"fmt"
	"net"
	"testing"
)

func TestRetrievingLocalIP(t *testing.T) {
	localip, err := Execute_systems_commands.RetrieveCurrentIP()
	if err != nil {
		t.Fail()
	}
	if localip != "10.0.2.15" {
		t.Fail()
	}
}

func TestStartListeningto(t *testing.T) {

	var clientconn net.Conn
	var serverconn net.Conn

	conch := make(chan net.Conn, 1)
	canstart := make(chan bool)

	localip, err := Execute_systems_commands.RetrieveCurrentIP()
	if err != nil {
		t.Fail()
	}
	localip = localip + malaware_server_code.Port

	go malaware_server_code.AsyncStartListeningTo(&conch, &canstart, localip)

	for {
		select {
		case serverconn = <-conch:
			serverconn.Close()
			return
		case startend := <-canstart:
			if startend {
				clientconn, err = malaware_victim_client.ConnectTo(localip)
				if err != nil {
					t.Fatalf("Connection error %s", err.Error())
					return
				}
				clientconn.Close()
			}

		}
	}

	//TODO remove debug message
	fmt.Println(serverconn.RemoteAddr().String())
	fmt.Println(clientconn.RemoteAddr().String())

}
