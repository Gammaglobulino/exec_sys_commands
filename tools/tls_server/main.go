package main

import (
	"../../malaware_server_code/core/handle_connections"
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

func HandleConnectionFor(clientConnection net.Conn) {
	defer clientConnection.Close()
	socketReader := bufio.NewReader(clientConnection)
	for {
		message, err := socketReader.ReadString('\n')
		if err != nil {
			return
		}
		fmt.Println(message)
		message = message + "<sent from:" + clientConnection.LocalAddr().String() + ">"
		numBytesWritten, err := clientConnection.Write([]byte(message))
		if err != nil {
			log.Println("Error writing data back to the client", err)
			return
		}
		fmt.Printf("Wrote %d bytes back to client. \n", numBytesWritten)
	}
}

func StartTLSServerTo(ip string) {

	serverCert, err := tls.LoadX509KeyPair("PEMCertificate", "privatePEM")
	if err != nil {
		log.Fatal(err)
	}
	serverConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}
	if err != nil {
		log.Fatal(err)
	}
	listener, err := tls.Listen("tcp", ip, serverConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		clientConnection, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting client connection")
			continue
		}
		go HandleConnectionFor(clientConnection)
	}
}

func main() {
	localip, err := handle_connections.GetLocalIp04Str()
	if err != nil {
		log.Fatal(err)
	}
	localip = localip + ":9090"
	StartTLSServerTo(localip)
}
