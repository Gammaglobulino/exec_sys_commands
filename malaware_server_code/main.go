package main

import (
	"../Execute_systems_commands"
	"fmt"
	"log"
)

const (
	Port = ":9090"
)

func main() {
	//var connection net.Conn
	localip, err := Execute_systems_commands.RetrieveCurrentIP()
	if err != nil {
		log.Fatal(err)
	}
	localip = localip + Port
	fmt.Println(localip)

}
