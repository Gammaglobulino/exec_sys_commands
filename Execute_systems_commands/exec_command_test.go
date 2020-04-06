package main

// in order to works, you should have sudo permission to runs all the ip commands
// https://www.cyberciti.biz/faq/linux-unix-running-sudo-command-without-a-password/

import (
	"../Execute_systems_commands"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestExecuteSystemCommandWithFunc(t *testing.T) {
	out, err := main.ExecuteSystemCommand("ip", []string{"address"})
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(strings.Split(out, " ")[2:3])
	assert.EqualValues(t, "<LOOPBACK,UP,LOWER_UP>", strings.Split(out, " ")[2:3][0])
	fmt.Println(out)
}

func TestCommand_Execute_RetrieveETH0StringIndexinterface(t *testing.T) {
	out, err := main.ExecuteSystemCommand("ip", []string{"a"})
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	outputCommands := strings.Split(out, " ")
	var index int
	for i, v := range outputCommands {
		if v == "eth0:" {
			index = i
			break
		}
	}
	assert.EqualValues(t, 56, index)
}

func TestCommand_Execute_ETH0InterfaceExists(t *testing.T) {
	err, yes, index := main.InterfaceExists("eth0:")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	assert.True(t, yes)
	assert.EqualValues(t, 56, index)
}

func TestRetrieveIPaddr(t *testing.T) {
	out, err := main.ExecuteSystemCommand("ip", []string{"a"})
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	outputCommands := strings.Split(out, " ")
	var index int
	for i, v := range outputCommands {
		if v == "eth0:" {
			index = i
			break
		}
	}
	IP := ""
	NextEth0 := outputCommands[index:len(outputCommands)]
	fmt.Println(NextEth0)
	for i, v := range NextEth0 {
		if v == "inet" {
			IP = NextEth0[i+1]
		}
	}
	assert.NotContains(t, "", IP)
}

func TestRetrieveCurrentIP(t *testing.T) {
	ip, err := main.RetrieveCurrentIP()
	assert.Nil(t, err)
	assert.NotContains(t, "", ip)
}

func TestExecuteSystemCommandWithStruct(t *testing.T) {
	cmd := main.NewIPCommand().AddArgument("address")
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
	out, err := main.ExecuteSystemCommand("ip", args)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(strings.Split(out, " ")[2:3])
	assert.EqualValues(t, "<BROADCAST,MULTICAST,UP,LOWER_UP>", strings.Split(out, " ")[2:3][0])
	fmt.Println(out)

}
func TestExecuteSystemCommandWrongArg(t *testing.T) {
	cmd := main.NewIPCommand().AddArgument("wrong")
	_, err := cmd.Execute()
	assert.NotNil(t, err, "Wrong argument")
}

func TestExecuteSystemCommandInterfaceOff(t *testing.T) {
	//ip link set eth0 down
	args := []string{"ip", "link", "set", "eth0", "down"}
	_, err := main.ExecuteSystemCommand("sudo", args)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

}
func TestExecuteSystemCommandInterfaceOn(t *testing.T) {
	//ip link set eth0 up
	args := []string{"ip", "link", "set", "eth0", "up"}
	_, err := main.ExecuteSystemCommand("sudo", args)
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
	_, err := main.ExecuteSystemCommand("sudo", args)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	args = []string{"ip", "link", "set", "eth0", "address", newMACaddr}
	_, err = main.ExecuteSystemCommand("sudo", args)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	args = []string{"ip", "link", "set", iFace, "up"}
	_, err = main.ExecuteSystemCommand("sudo", args)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
