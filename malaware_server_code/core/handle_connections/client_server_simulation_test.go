package handle_connections

import (
	"../../../Execute_systems_commands"
	"../handle_connections"
	"net"
	"testing"
)

func TestStartListeningto(t *testing.T) {

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

	for {
		select {
		case serverconn = <-conch:
			serverconn.Close()
			return
		case startend := <-canstart:
			if startend {
				clientconn, err = handle_connections.ConnectTo(localip)
				if err != nil {
					t.Fatalf("Connection error %s", err.Error())
					return
				}
				clientconn.Close()
			}

		}
	}

}
