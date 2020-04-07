package Execute_systems_commands

// in order to works, you should have sudo permission to runs all the ip commands
// https://www.cyberciti.biz/faq/linux-unix-running-sudo-command-without-a-password/

import (
	"../Execute_systems_commands"
	"fmt"
	"strings"
	"testing"
)

func TestExecuteSystemCommandWithFunc(t *testing.T) {
	out, err := Execute_systems_commands.ExecuteSystemCommand("ip", []string{"address"})
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	result := strings.Split(out, " ")[2:3][0]
	if result != "<LOOPBACK,UP,LOWER_UP>" {
		t.Fatalf("Output string should contain <LOOPBACK,UP,LOWER_UP>, actual %s", result)
	}

}

func TestCommand_Execute_RetrieveETH0StringIndexinterface(t *testing.T) {
	out, err := Execute_systems_commands.ExecuteSystemCommand("ip", []string{"a"})
	if err != nil {
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
	if index != 56 {
		t.Fatalf("Expected index should be 56, actual %v", index)
	}
}

func TestCommand_Execute_ETH0InterfaceExists(t *testing.T) {
	err, yes, index := Execute_systems_commands.InterfaceExists("eth0:")
	if err != nil {
		t.Fail()
	}
	if !yes {
		t.Fatalf("Interface should exist")
	}
	if index != 56 {
		t.Fatalf("Expected index should be 56, actual %v", index)
	}
}

func TestRetrieveIPaddr(t *testing.T) {
	out, err := Execute_systems_commands.ExecuteSystemCommand("ip", []string{"a"})
	if err != nil {
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
	NextEth0 := outputCommands[index:]
	for i, v := range NextEth0 {
		if v == "inet" {
			IP = NextEth0[i+1]
		}
	}
	IP = strings.Split(IP, "/")[0]
	if IP == "" {
		t.Fatal("Invalid IP or not present")
	}
}

func TestRetrieveCurrentIP(t *testing.T) {
	ip, err := Execute_systems_commands.RetrieveCurrentIP()
	if err != nil {
		t.Fatal(err)
	}
	if ip == "" {
		t.Fatal("Invalid ip address or not present")
	}
}

func TestExecuteSystemCommandWithStruct(t *testing.T) {
	cmd := Execute_systems_commands.NewIPCommand().AddArgument("address")
	out, err := cmd.Execute()
	if err != nil {
		t.Fail()
	}
	result := strings.Split(out, " ")[2:3][0]
	if result != "<LOOPBACK,UP,LOWER_UP>" {
		t.Fatalf("Expected <LOOPBACK,UP,LOWER_UP>, actual:%s", result)
	}
}

func TestExecuteSystemCommandQueryInterface(t *testing.T) {
	args := []string{"link", "ls", "eth0"}
	out, err := Execute_systems_commands.ExecuteSystemCommand("ip", args)
	if err != nil {
		t.Fail()
	}
	result := strings.Split(out, " ")[2:3][0]
	if result != "<BROADCAST,MULTICAST,UP,LOWER_UP>" {
		t.Fatalf("Expected <BROADCAST,MULTICAST,UP,LOWER_UP>, actual:%s", result)
	}
}
func TestExecuteSystemCommandWrongArg(t *testing.T) {
	cmd := Execute_systems_commands.NewIPCommand().AddArgument("wrong")
	_, err := cmd.Execute()
	if err == nil {
		t.Fatal(err)
	}
}

func TestExecuteSystemCommandInterfaceOff(t *testing.T) {
	//ip link set eth0 down
	args := []string{"ip", "link", "set", "eth0", "down"}
	_, err := Execute_systems_commands.ExecuteSystemCommand("sudo", args)
	if err != nil {
		t.Fail()
	}

}
func TestExecuteSystemCommandInterfaceOn(t *testing.T) {
	//ip link set eth0 up
	args := []string{"ip", "link", "set", "eth0", "up"}
	_, err := Execute_systems_commands.ExecuteSystemCommand("sudo", args)
	if err != nil {
		t.Fail()
	}
}

func TestExecuteSystemCommandChangingMAC(t *testing.T) {
	//ip link set eth0 address 02:01:02:03:04:08
	iFace := "eth0"
	newMACaddr := "00:11:22:33:44:55"

	args := []string{"ip", "link", "set", iFace, "down"}
	_, err := Execute_systems_commands.ExecuteSystemCommand("sudo", args)
	if err != nil {
		t.Fail()
	}

	args = []string{"ip", "link", "set", "eth0", "address", newMACaddr}
	_, err = Execute_systems_commands.ExecuteSystemCommand("sudo", args)
	if err != nil {
		t.Fail()
	}

	args = []string{"ip", "link", "set", iFace, "up"}
	_, err = Execute_systems_commands.ExecuteSystemCommand("sudo", args)
	if err != nil {
		t.Fail()
	}
}
