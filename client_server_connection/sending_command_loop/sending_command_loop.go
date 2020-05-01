package sending_command_loop

import (
	"../../Execute_systems_commands"
	"../../tools/secure_crypting_AES"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"io/ioutil"
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

type DataRemoteCommand struct {
	//A remote command with an embedded file or byte data, kind of jacket command
	Command       RemoteCommand
	DataSignature string
	DataBytes     []byte
	Filename      string
	Crypted       bool
}

func (drc *DataRemoteCommand) Execute() (string, error) {
	var err error
	if drc.Crypted {
		drc.DataBytes, err = secure_crypting_AES.AESDecryptFromBytesToByteArray([]byte(drc.DataSignature), drc.DataBytes)
		if err != nil {
			drc.Command.OutError = err.Error()
			return "", err
		}
	}
	ioutil.WriteFile(drc.Filename, drc.DataBytes, 0666)
	drc.Command.Execute()
	return "ok", nil
}

func NewAESCryptedDRC(password string, data []byte) (*DataRemoteCommand, error) {

	//MD5 is not secure - but you need 32byte code for AES. MD5 plus AES is a kind of secure method.
	drc := &DataRemoteCommand{Crypted: true}

	MD5 := md5.Sum([]byte(password))
	drc.DataSignature = base64.StdEncoding.EncodeToString(MD5[:])

	var err error
	drc.DataBytes, err = secure_crypting_AES.AESCryptFromDataBytesToByteArray([]byte(drc.DataSignature), data)

	if err != nil {
		return nil, err
	}
	return drc, nil
}

type CommandPipe struct {
	Name string
	Pipe []Execute_systems_commands.Commander
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
		if rc.CmdPlusArgs() != "" {
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
		} else {
			return "", nil
		}

	}
}
func (p *CommandPipe) Execute() (string, error) {
	for i, _ := range p.Pipe {
		p.Pipe[i].Execute()
	}
	return "ok", nil
}

func (p *CommandPipe) AddCommand(command Execute_systems_commands.Commander) *CommandPipe {
	p.Pipe = append(p.Pipe, command)
	return p
}

func NewCommand(commandName string, f func() (string, error)) *RemoteCommand {
	return &RemoteCommand{
		Name:          commandName,
		CodeToExecute: f,
	}
}
func NewCommandString(commandName string, commandString string, os string) Execute_systems_commands.Commander {
	return &RemoteCommand{
		Name:            commandName,
		CodeToExecute:   nil,
		CommandString:   commandString,
		TargetOSMachine: os,
	}
}
