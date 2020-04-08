package handle_connections

import (
	"../../../Execute_systems_commands"
	"../handle_connections"
	"fmt"
	"net"
	"testing"
)

func TestClientServerCommunicationLab(t *testing.T) {

	var clientconn net.Conn
	var serverconn net.Conn

	conch := make(chan net.Conn, 1)
	canstart := make(chan bool)

	localip, err := Execute_systems_commands.RetrieveCurrentIP()
	if err != nil {
		t.Fail()
	}
	localip = localip + handle_connections.Port

	go handle_connections.MultiCoreListeningTo(&conch, &canstart, localip)
	serverClosed := false
	for {
		select {
		case serverconn = <-conch:
			serverconn.Close()
			serverClosed = true
			break
		case isServerStarted := <-canstart:
			if isServerStarted {
				clientconn, err = handle_connections.ConnectTo(localip)
				if err != nil {
					t.Fatalf("Connection error %s", err.Error())
					return
				}
				clientconn.Close()
			}
			break
		}
		if serverClosed {
			break
		}
	}

	//TODO remove
	fmt.Println("Bye")

}
