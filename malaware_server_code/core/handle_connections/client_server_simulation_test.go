package handle_connections

import (
	"../../../Execute_systems_commands"
	"../handle_connections"
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
)

func TestClientServerCommunicationLab(t *testing.T) {

	var clientconn net.Conn
	var serverconn net.Conn

	conch := make(chan net.Conn)
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
			serverClosed = true
			break
		case isServerStarted := <-canstart:
			if isServerStarted {
				clientconn, err = handle_connections.ConnectTo(localip)
				if err != nil {
					t.Fatalf("Connection error %s", err.Error())
					return
				}
			}
			break
		}
		if serverClosed {
			break
		}
	}
	clientconn.Close()
	serverconn.Close()
}
func TestSendStreamedDatatoClient(t *testing.T) {

	var clientconn net.Conn
	var serverconn net.Conn

	conch := make(chan net.Conn)
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
			serverClosed = true
			break
		case isServerStarted := <-canstart:
			if isServerStarted {
				clientconn, err = handle_connections.ConnectTo(localip)
				if err != nil {
					t.Fatalf("Connection error %s", err.Error())
					return
				}
			}
			break
		}
		if serverClosed {
			break
		}
	}

	defer clientconn.Close()
	defer serverconn.Close()

	var (
		data2send    string
		dataReceived string
	)
	var numberofbytes int

	data2send = "Manu è culo\n"

	numberofbytes, err = clientconn.Write([]byte(data2send))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Number of bytes written %v\n", numberofbytes)

	readFromServer := bufio.NewReader(serverconn)
	dataReceived, err = readFromServer.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	strings.TrimSuffix(dataReceived, "\n")
	fmt.Printf("Data recieved from the server: %s", dataReceived)

}
