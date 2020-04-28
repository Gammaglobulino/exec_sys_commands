package sending_command_loop

import (
	"bytes"
	"errors"
	"os/exec"
)

const (
	win = "windows"
	lin = "linux"
)

type RemoteCommand struct {
	Name            string
	CodeToExecute   func() (string, error)
	CommandString   string
	CommandOutput   string
	OutError        string
	TargetOSMachine string //"windows" or "linux"
	Args            []string
}

type CommandPipe struct {
	Name string
	Pipe []RemoteCommand
}

func (rc *RemoteCommand) CmdPlusArgs() string {
	outString := rc.CommandString
	for _, v := range rc.Args {
		outString += " " + v
	}
	return outString
}

func (rc *RemoteCommand) AddParam(param string) {
	rc.Args = append(rc.Args, param)
}
func (rc *RemoteCommand) Execute() (string, error) {

	if rc.CodeToExecute != nil {
		result, err := rc.CodeToExecute()
		if err != nil {
			return "", err
		}
		return result, err
	} else {
		var cmd *exec.Cmd
		switch rc.TargetOSMachine {
		case win:
			cmd = exec.Command("powershell.exe", "/C", rc.CmdPlusArgs())
		case lin:
			cmd = exec.Command(rc.CmdPlusArgs())
		default:
			return "", errors.New("OS not recognized")
		}

		cmdOut := bytes.NewBuffer(make([]byte, 0, 1024))
		cmdErr := bytes.NewBuffer(make([]byte, 0, 1024))

		cmd.Stdout = cmdOut
		cmd.Stderr = cmdErr

		err := cmd.Run()
		if err != nil {
			return "", err
		}
		rc.CommandOutput = cmdOut.String()
		rc.OutError = cmdErr.String()
		return rc.CommandOutput, nil

	}
}
func (p *CommandPipe) Execute() (string, error) {
	for i, _ := range p.Pipe {
		p.Pipe[i].Execute()
	}
	return "ok", nil
}

func (p *CommandPipe) AddCommand(command RemoteCommand) *CommandPipe {
	p.Pipe = append(p.Pipe, command)
	return p
}

func NewCommand(commandName string, f func() (string, error)) *RemoteCommand {
	return &RemoteCommand{
		Name:          commandName,
		CodeToExecute: f,
	}
}
func NewCommandString(commandName string, commandString string, os string) RemoteCommand {
	return RemoteCommand{
		Name:            commandName,
		CodeToExecute:   nil,
		CommandString:   commandString,
		TargetOSMachine: os,
	}
}
