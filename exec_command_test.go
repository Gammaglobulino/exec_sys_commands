package main

// in order to works, you should have sudo permission to runs all the ip commands
// https://www.cyberciti.biz/faq/linux-unix-running-sudo-command-without-a-password/

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestExecuteSystemCommandWithFunc(t *testing.T) {
	out, err := ExecuteSystemCommand("ip", []string{"address"})
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(strings.Split(out, " ")[2:3])
	assert.EqualValues(t, "<LOOPBACK,UP,LOWER_UP>", strings.Split(out, " ")[2:3][0])
	fmt.Println(out)
}
func TestExecuteSystemCommandWithStruct(t *testing.T) {
	cmd := NewIPCommand().AddArgument("address")
	out, err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(strings.Split(out, " ")[2:3])
	assert.EqualValues(t, "<LOOPBACK,UP,LOWER_UP>", strings.Split(out, " ")[2:3][0])
	fmt.Println(out)
}

func TestExecuteSystemCommandQueryInterface(t *testing.T) {
	args := []string{"link", "ls", "eth0"}
	out, err := ExecuteSystemCommand("ip", args)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(strings.Split(out, " ")[2:3])
	assert.EqualValues(t, "<BROADCAST,MULTICAST,UP,LOWER_UP>", strings.Split(out, " ")[2:3][0])
	fmt.Println(out)

}
func TestExecuteSystemCommandWrongArg(t *testing.T) {
	cmd := NewIPCommand().AddArgument("wrong")
	_, err := cmd.Execute()
	assert.NotNil(t,err,"Wrong argument")
}

func TestExecuteSystemCommandInterfaceOff(t *testing.T) {
	//ip link set eth0 down
	args := []string{"ip", "link", "set", "eth0", "down"}
	_, err := ExecuteSystemCommand("sudo", args)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

}
func TestExecuteSystemCommandInterfaceOn(t *testing.T) {
	//ip link set eth0 up
	args := []string{"ip", "link", "set", "eth0", "up"}
	_, err := ExecuteSystemCommand("sudo", args)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestExecuteSystemCommandChangingMAC(t *testing.T) {
	//ip link set eth0 address 02:01:02:03:04:08
	iFace := "eth0"
	newMACaddr := "00:11:22:33:44:55"

	args := []string{"ip", "link", "set", iFace, "down"}
	_, err := ExecuteSystemCommand("sudo", args)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	args = []string{"ip", "link", "set", "eth0", "address", newMACaddr}
	_, err = ExecuteSystemCommand("sudo", args)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	args = []string{"ip", "link", "set", iFace, "up"}
	_, err = ExecuteSystemCommand("sudo", args)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
