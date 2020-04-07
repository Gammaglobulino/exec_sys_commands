package main

import (
	"os/exec"
	"strings"
)

type Commander interface {
	Execute() (string, error)
}
type Command struct {
	name   string
	inFace string
	mac    string
	args   []string
}

func (c *Command) Execute() (string, error) {
	cmd := exec.Command(c.name, c.args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	output := string(out)
	return output, nil
}
func (c *Command) AddArgument(arg string) *Command {
	c.args = append(c.args, arg)
	return c
}
func NewIPCommand() *Command {
	c := Command{
		name:   "sudo",
		inFace: "",
		mac:    "",
		args:   []string{"ip"},
	}
	return &c
}

func ExecuteSystemCommand(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	output := string(out)
	return output, nil
}

func InterfaceExists(iface string) (error, bool, int) {
	out, err := ExecuteSystemCommand("ip", []string{"a"})
	if err != nil {
		return err, false, 0
	}
	outputCommands := strings.Split(out, " ")
	isthere := false
	index := 0
	for i, v := range outputCommands {
		if v == iface {
			isthere = true
			index = i
			break
		}
	}
	return nil, isthere, index
}
func RetrieveCurrentIP() (string, error) {
	out, err := ExecuteSystemCommand("ip", []string{"a"})
	if err != nil {
		return "", err
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
			IP = strings.Split(NextEth0[i+1], "/")[0]
		}
	}

	return IP, nil

}

func main() {

}
