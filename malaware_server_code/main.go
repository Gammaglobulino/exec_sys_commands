package malaware_server_code

import (
	"../Execute_systems_commands"
	"fmt"
	"log"
	"net"
)

const (
	Port = ":9090"
)

func StartListeningto(ip string) (net.Conn, error) {
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
func AsyncStartListeningTo(connectionch *chan net.Conn, canstarttosend *chan bool, ip string) {
	var connection net.Conn
	listener, err := net.Listen("tcp", ip)
	if err != nil {
		close(*connectionch)
		close(*canstarttosend)
		log.Fatal(err)
	}
	fmt.Println("Listening to:", ip)
	*canstarttosend <- true //client can start to listen
	connection, err = listener.Accept()

	if err != nil {
		close(*connectionch)
		close(*canstarttosend)
		log.Fatal(err)
	} else {
		fmt.Println("Connection established from:" + connection.RemoteAddr().String())
	}
	*connectionch <- connection
}

func main() {

	localip, err := Execute_systems_commands.RetrieveCurrentIP()
	if err != nil {
		log.Fatal(err)
	}

	localip = localip + Port

	con, err := StartListeningto(localip)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(con.RemoteAddr().String())
}
