package handle_connections

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	Port = ":9090"
)

func SingleThreadListeningto(ip string) (net.Conn, error) {
	var connection net.Conn
	listener, err := net.Listen("tcp", ip)
	if err != nil {
		return nil, err
	}
	fmt.Println("Listening to:", ip)
	connection, err = listener.Accept()
	if err != nil {
		return nil, err
	} else {
		fmt.Println("Connection established from:" + connection.RemoteAddr().String())
	}
	return connection, nil
}
func MultiCoreListeningTo(connectionch *chan net.Conn, canstarttosend *chan bool, errorch *chan error, ip string) {
	var connection net.Conn
	listener, err := net.Listen("tcp", ip)
	if err != nil {
		*errorch <- err
	}
	fmt.Println("Listening to:", ip)
	*canstarttosend <- true //client can start to listen

	connection, err = listener.Accept()
	if err != nil {
		*errorch <- err
	} else {
		fmt.Println("Connection established from:" + connection.RemoteAddr().String())
	}
	*connectionch <- connection

}
func SendStringTo(conn net.Conn, stringtoSend string) (int, error) {
	bytes, err := conn.Write([]byte(stringtoSend))
	if err != nil {
		return 0, err
	}
	return bytes, nil
}

func ReceiveStringFrom(conn net.Conn) (string, error) {
	readFromServer := bufio.NewReader(conn)
	dataReceived, err := readFromServer.ReadString('\n')
	if err != nil {
		return "", err
	}
	strings.TrimSuffix(dataReceived, "\n")
	return dataReceived, nil
}

func SetClientServerConnectionTo(ip string) (server net.Conn, client net.Conn, err error) {

	conch := make(chan net.Conn)
	canstart := make(chan bool)
	errorch := make(chan error)

	go MultiCoreListeningTo(&conch, &canstart, &errorch, ip)
	serverClosed := false
	for {
		select {
		case server = <-conch:
			serverClosed = true
			break
		case isServerStarted := <-canstart:
			if isServerStarted {
				client, err = ConnectTo(ip)
				if err != nil {
					return nil, nil, err
				}
			}
			break
		case error := <-errorch:
			return nil, nil, error

		}
		if serverClosed {
			break
		}
	}
	return server, client, nil
}
func GetLocalIp04() (string, error) {
	host, err := os.Hostname()
	if err != nil {
		return "", err
	}
	addrs, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}
	ip4 := addrs[len(addrs)-1]

	return ip4.To4().String(), nil
}
