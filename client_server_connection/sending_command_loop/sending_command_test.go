package sending_command_loop

import (
	"../core/handle_connections"
	"../sending_command_loop"
	"encoding/gob"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"os"
	"testing"
)

func TestNewCommand(t *testing.T) {
	command := func() (string, error) {
		println("Hello world")
		return "ok", nil
	}
	printHello := sending_command_loop.NewCommand("Print Hello", command)
	out, _ := printHello.Execute()
	assert.EqualValues(t, "ok", out)
}

func TestCommandPipe_AddCommand(t *testing.T) {
	command := func() (string, error) {
		println("Hello world")
		return "ok", nil
	}
	printHello := sending_command_loop.NewCommand("Print Hello", command)

	command = func() (string, error) {
		fmt.Println("Exec dir")
		return "ok", nil
	}
	execDir := sending_command_loop.NewCommand("Exec dir", command)

	pipe := &sending_command_loop.CommandPipe{}
	pipe.AddCommand(*printHello).AddCommand(*execDir)
	assert.EqualValues(t, "Print Hello", pipe.Pipe[0].Name)
	pipe.Execute()
}

func TestRemoteCommand_Execute(t *testing.T) {
	dirCommand := sending_command_loop.NewCommandString("exec dir", "dir", "windows")
	out, err := dirCommand.Execute()
	if err != nil {
		log.Fatal(err)
	}
	assert.NotEmpty(t, out)
	log.Println(out)
}

func TestRemoteCommandOSNotRecognized(t *testing.T) {
	dirCommand := sending_command_loop.NewCommandString("exec dir", "dir", "OS2")
	out, err := dirCommand.Execute()
	assert.NotNil(t, err)
	assert.Empty(t, out)
	log.Println(err)
}

func TestRemoteCommandAddingParameters(t *testing.T) {
	dirCommand := sending_command_loop.NewCommandString("exec dir", "dir", "windows")
	assert.Nil(t, dirCommand.Args)
	dirCommand.AddParam("c:/")
	assert.EqualValues(t, 1, len(dirCommand.Args))
	out, err := dirCommand.Execute()
	if err != nil {
		log.Fatal(err)
	}
	assert.NotEmpty(t, out)
	log.Println(out)
}

func TestBuildStringCommandPipe(t *testing.T) {
	ipConfigCommand := sending_command_loop.NewCommandString("exec ipconfig", "ipconfig", "windows")
	assert.NotEmpty(t, ipConfigCommand)

	dirCommand := sending_command_loop.NewCommandString("exec dir", "dir", "windows")
	assert.NotEmpty(t, ipConfigCommand)

	pipe := sending_command_loop.CommandPipe{}
	pipe.AddCommand(ipConfigCommand).AddCommand(dirCommand)
	pipe.Execute()

	log.Println(t, pipe.Pipe[0].CommandOutput)
	log.Println(t, pipe.Pipe[1].CommandOutput)

	assert.NotEmpty(t, pipe.Pipe[0].CommandOutput)
	assert.NotEmpty(t, pipe.Pipe[1].CommandOutput)
}

func TestTCPCommandPipeUsingGOB(t *testing.T) {

	//not working: func cannot be serialized with Gob and in general with any static compiled language
	//client/server communication loop created

	var serverconn net.Conn
	var clientconn net.Conn

	t.Run("Test Client/Server connection", func(t *testing.T) {
		localip, err := handle_connections.GetLocalIp04Str()
		assert.Nil(t, err)
		assert.Contains(t, localip, "")
		assert.Nil(t, err)
		localip = localip + handle_connections.Port
		serverconn, clientconn, err = handle_connections.SetClientServerConnectionTo(localip)
		assert.Nil(t, err)
		assert.NotNil(t, serverconn)
		assert.NotNil(t, clientconn)

	})
	defer serverconn.Close()
	defer clientconn.Close()

	//command pipe creation

	command := func() (string, error) {
		println("Hello world")
		return "ok", nil
	}
	printHello := sending_command_loop.NewCommand("Print Hello", command)

	command = func() (string, error) {
		fmt.Println("Exec dir")
		return "ok", nil
	}
	execDir := sending_command_loop.NewCommand("Exec dir", command)

	pipe := &sending_command_loop.CommandPipe{}
	pipe.AddCommand(*printHello).AddCommand(*execDir)

	receiveddata := &sending_command_loop.CommandPipe{Pipe: nil}
	gobencoder := gob.NewEncoder(clientconn)
	err := gobencoder.Encode(pipe)
	if err != nil {
		t.Fail()
	}
	log.Println(pipe)
	gobdecoder := gob.NewDecoder(serverconn)
	err = gobdecoder.Decode(receiveddata)
	if err != nil {
		t.Fail()
	}
	log.Println(receiveddata)

}

func TestTCPRemotePipeExecutionUsingGOBSerialization(t *testing.T) {

	//GOB force you to serialize data not pointers
	//you can pass a pointer to a slice of struct, not pointers

	var serverconn net.Conn
	var clientconn net.Conn

	t.Run("Test Client/Server connection", func(t *testing.T) {
		localip, err := handle_connections.GetLocalIp04Str()
		assert.Nil(t, err)
		assert.Contains(t, localip, "")
		assert.Nil(t, err)
		localip = localip + handle_connections.Port
		serverconn, clientconn, err = handle_connections.SetClientServerConnectionTo(localip)
		assert.Nil(t, err)
		assert.NotNil(t, serverconn)
		assert.NotNil(t, clientconn)

	})
	defer serverconn.Close()
	defer clientconn.Close()

	receivedData := &sending_command_loop.CommandPipe{}

	t.Run("Sent complete command pipe to server ", func(t *testing.T) {
		pipe := &sending_command_loop.CommandPipe{Name: "Andy pipe"}

		dirCommand := sending_command_loop.NewCommandString("Exec remote dir", "dir", "windows")
		ipconfigCommand := sending_command_loop.NewCommandString("Exec remote ipconfig", "ipconfig", "windows")

		pipe.AddCommand(dirCommand).AddCommand(ipconfigCommand)

		gobencoder := gob.NewEncoder(clientconn)
		err := gobencoder.Encode(pipe)
		if err != nil {
			log.Fatal(err)
		}
		gobdecoder := gob.NewDecoder(serverconn)
		err = gobdecoder.Decode(receivedData)
		if err != nil {
			log.Fatal(err)
		}
		assert.NotEmpty(t, receivedData.Pipe[0].Name)

		log.Printf("Pipe sent to server %s contains %s", clientconn.RemoteAddr().String(), pipe.Name)
		log.Printf("Pipe received from client %s contains %s", serverconn.RemoteAddr().String(), receivedData.Name)

	})

	t.Run("Sent executed command pipe back to the client", func(t *testing.T) {
		_, err := receivedData.Execute()
		if err != nil {
			log.Fatal(err)
		}

		gobencoder := gob.NewEncoder(serverconn)
		err = gobencoder.Encode(receivedData)
		if err != nil {
			log.Fatal(err)
		}
		pipe := &sending_command_loop.CommandPipe{}
		gobdecoder := gob.NewDecoder(clientconn)
		err = gobdecoder.Decode(pipe)
		if err != nil {
			log.Fatal(err)
		}
		assert.EqualValues(t, pipe.Pipe[0].Name, receivedData.Pipe[0].Name)

		log.Println(pipe.Pipe[0].CommandOutput)
		log.Println(pipe.Pipe[1].CommandOutput)

	})

}

func TestSendingCommandTo(t *testing.T) {
	localIp, err := handle_connections.GetLocalIp04Str()
	assert.Nil(t, err)
	localIp = localIp + handle_connections.Port
	serverConn, clientConn, err := handle_connections.SetClientServerConnectionTo(localIp)
	assert.Nil(t, err)
	assert.NotNil(t, serverConn)
	assert.NotNil(t, clientConn)
	defer clientConn.Close()
	defer serverConn.Close()
}

func TestFileSystemNavigationCommandPipe(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(workingDir)
}
