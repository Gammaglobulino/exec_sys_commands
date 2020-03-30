package main

import (
	"os/exec"
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
func NewIPCommand(n string) *Command {
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

func main() {

}
