package main

import (
	"../Execute_systems_commands"
	"fmt"
	"log"
	"net"
)

const (
	Port = ":9090"
)

func main() {
	var connection net.Conn
	var listener net.Listener
	localip, err := Execute_systems_commands.RetrieveCurrentIP()
	if err != nil {
		log.Fatal(err)
	}
	localip = localip + Port

	listener, err = net.Listen("tcp", localip)
	if err != nil {
		log.Fatal(err)
	}
	connection, err = listener.Accept()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection established from:" + connection.RemoteAddr().String())
	}

}
