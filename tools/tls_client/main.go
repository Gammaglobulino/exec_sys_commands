package main

import (
	"../../client_server_connection/core/handle_connections"
	"crypto/tls"
	"log"
)

func main() {
	hostname, err := handle_connections.GetLocalIp04Str()
	hostname = hostname + ":9090"

	if err != nil {
		log.Fatal(err)
	}
	messageToSend := "Hello tls server\n"

	certificate, err := tls.LoadX509KeyPair("PEMCertificate", "privatePEM")
	if err != nil {
		log.Fatal(err)
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{certificate},
		ServerName:         hostname,
	}
	connection, err := tls.Dial("tcp", hostname, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	numBytesWritten, err := connection.Write([]byte(messageToSend))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes to the server %s", numBytesWritten, hostname)

	bufferRead := make([]byte, 100)
	numBytesRead, err := connection.Read(bufferRead)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Read %d bytes from the server %s", numBytesRead, hostname)
	log.Println(string(bufferRead))

}
