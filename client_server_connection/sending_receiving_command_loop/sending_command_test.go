package sending_receiving_command_loop

import (
	"."
	"../core/handle_connections"
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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
	printHello := sending_receiving_command_loop.NewCommand("Print Hello", command)
	out, _ := printHello.Execute()
	assert.EqualValues(t, "ok", out)
}

func TestCommandPipe_AddCommand(t *testing.T) {
	command := func() (string, error) {
		println("Hello world")
		return "ok", nil
	}
	printHello := sending_receiving_command_loop.NewCommand("Print Hello", command)

	command = func() (string, error) {
		fmt.Println("Exec dir")
		return "ok", nil
	}
	execDir := sending_receiving_command_loop.NewCommand("Exec dir", command)

	pipe := &sending_receiving_command_loop.CommandPipe{}
	pipe.AddCommand(printHello).AddCommand(execDir)
	assert.EqualValues(t, "Print Hello", pipe.Pipe[0].(*sending_receiving_command_loop.RemoteCommand).Name)
	pipe.Execute()
}

func TestTCPSendingReceiveHybridPipe(t *testing.T) {
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

	pipe := sending_receiving_command_loop.CommandPipe{Name: "Andy pipe"}
	receivingPipe := sending_receiving_command_loop.CommandPipe{}

	dataByte := []byte("Ciao Andrea come stai?")

	newDRC, err := sending_receiving_command_loop.NewCryptedDRC("test", dataByte)
	newDRC.Filename = "Andrea.txt"

	assert.Nil(t, err)
	assert.NotEmpty(t, newDRC)
	pipe.AddCommand(newDRC).AddCommand(&sending_receiving_command_loop.RemoteCommand{
		Name:            "Dir command",
		CodeToExecute:   nil,
		CommandString:   "dir",
		CommandOutput:   "",
		OutError:        "",
		TargetOSMachine: "windows",
		Args:            nil,
	})

	log.Println(pipe.Pipe[0].(*sending_receiving_command_loop.DataRemoteCommand).DataBytes)

	gobencoder := gob.NewEncoder(clientconn)
	gob.Register(&sending_receiving_command_loop.DataRemoteCommand{})
	gob.Register(&sending_receiving_command_loop.RemoteCommand{})

	err = gobencoder.Encode(pipe)
	if err != nil {
		t.Fatal(err)
	}
	gobDencoder := gob.NewDecoder(serverconn)
	err = gobDencoder.Decode(&receivingPipe)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, pipe.Pipe[0].(*sending_receiving_command_loop.DataRemoteCommand).DataBytes, receivingPipe.Pipe[0].(*sending_receiving_command_loop.DataRemoteCommand).DataBytes)
	receivingPipe.Execute()

	log.Println(string(receivingPipe.Pipe[0].(*sending_receiving_command_loop.DataRemoteCommand).DataBytes))
	log.Println(string(receivingPipe.Pipe[1].(*sending_receiving_command_loop.RemoteCommand).CommandOutput))

}

func TestTCPSendingReceiveEncodedPipe(t *testing.T) {
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

	pipe := sending_receiving_command_loop.CommandPipe{Name: "Andy pipe"}
	receivingPipe := sending_receiving_command_loop.CommandPipe{}

	dataByte := []byte("Ciao Andrea come stai?")

	newDRC, err := sending_receiving_command_loop.NewCryptedDRC("test", dataByte)
	newDRC.Filename = "Andrea.txt"

	assert.Nil(t, err)
	assert.NotEmpty(t, newDRC)

	pipe.AddCommand(newDRC)
	log.Println(pipe.Pipe[0].(*sending_receiving_command_loop.DataRemoteCommand).DataBytes)

	gobencoder := gob.NewEncoder(clientconn)
	gob.Register(&sending_receiving_command_loop.DataRemoteCommand{})

	err = gobencoder.Encode(pipe)
	if err != nil {
		t.Fatal(err)
	}
	gobDencoder := gob.NewDecoder(serverconn)
	err = gobDencoder.Decode(&receivingPipe)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, pipe.Pipe[0].(*sending_receiving_command_loop.DataRemoteCommand).DataBytes, receivingPipe.Pipe[0].(*sending_receiving_command_loop.DataRemoteCommand).DataBytes)
	receivingPipe.Execute()
	log.Println(string(receivingPipe.Pipe[0].(*sending_receiving_command_loop.DataRemoteCommand).DataBytes))

}

func TestGOBEncodedPipe(t *testing.T) {
	pipe := sending_receiving_command_loop.CommandPipe{Name: "Andy pipe"}
	receivingPipe := sending_receiving_command_loop.CommandPipe{}

	dataByte := []byte("Ciao Andrea come stai?")

	newDRC, err := sending_receiving_command_loop.NewCryptedDRC("test", dataByte)
	assert.Nil(t, err)
	assert.NotEmpty(t, newDRC)

	pipe.AddCommand(newDRC)
	log.Println(pipe.Pipe[0].(*sending_receiving_command_loop.DataRemoteCommand).DataBytes)

	dataOut := bytes.NewBuffer(make([]byte, 0, 1024))

	gobencoder := gob.NewEncoder(dataOut)
	gob.Register(&sending_receiving_command_loop.DataRemoteCommand{})

	err = gobencoder.Encode(pipe)
	if err != nil {
		t.Fatal(err)
	}
	gobDencoder := gob.NewDecoder(dataOut)
	err = gobDencoder.Decode(&receivingPipe)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, pipe.Pipe[0].(*sending_receiving_command_loop.DataRemoteCommand).DataBytes, receivingPipe.Pipe[0].(*sending_receiving_command_loop.DataRemoteCommand).DataBytes)
	receivingPipe.Execute()
	log.Println(string(receivingPipe.Pipe[0].(*sending_receiving_command_loop.DataRemoteCommand).DataBytes))

}

func TestRemoteCommand_Execute(t *testing.T) {
	dirCommand := sending_receiving_command_loop.NewCommandString("exec dir", "dir", "windows")
	out, err := dirCommand.Execute()
	if err != nil {
		log.Fatal(err)
	}
	assert.NotEmpty(t, out)
	log.Println(out)
}

func TestRemoteCommandOSNotRecognized(t *testing.T) {
	dirCommand := sending_receiving_command_loop.NewCommandString("exec dir", "dir", "OS2")
	out, err := dirCommand.Execute()
	assert.NotNil(t, err)
	assert.Empty(t, out)
	log.Println(err)
}

func TestRemoteCommandAddingParameters(t *testing.T) {
	dirCommand := sending_receiving_command_loop.NewCommandString("exec dir", "dir", "windows")
	assert.Nil(t, dirCommand.(*sending_receiving_command_loop.RemoteCommand).Args)
	dirCommand.(*sending_receiving_command_loop.RemoteCommand).AddParam("c:/")
	assert.EqualValues(t, 1, len(dirCommand.(*sending_receiving_command_loop.RemoteCommand).Args))
	out, err := dirCommand.Execute()
	if err != nil {
		log.Fatal(err)
	}
	assert.NotEmpty(t, out)
	log.Println(out)
}

func TestAddWrongParam(t *testing.T) {
	dirCommand := sending_receiving_command_loop.NewCommandString("exec dir", "dir", "windows")
	assert.Nil(t, dirCommand.(*sending_receiving_command_loop.RemoteCommand).Args)
	dirCommand.(*sending_receiving_command_loop.RemoteCommand).AddParam("f:/")
	assert.EqualValues(t, 1, len(dirCommand.(*sending_receiving_command_loop.RemoteCommand).Args))
	_, err := dirCommand.Execute()
	assert.NotNil(t, err)
	fmt.Println(dirCommand.(*sending_receiving_command_loop.RemoteCommand).OutError)
}

func TestBuildStringCommandPipe(t *testing.T) {
	ipConfigCommand := sending_receiving_command_loop.NewCommandString("exec ipconfig", "ipconfig", "windows")
	assert.NotEmpty(t, ipConfigCommand)

	execDirCommand := sending_receiving_command_loop.NewCommandString("exec Dir", "dir", "windows")
	assert.NotEmpty(t, ipConfigCommand)

	execDirCommand.(*sending_receiving_command_loop.RemoteCommand).AddParam("c:")

	pipe := sending_receiving_command_loop.CommandPipe{}
	pipe.AddCommand(ipConfigCommand).AddCommand(execDirCommand)
	_, err := pipe.Execute()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(pipe.Pipe[0].(*sending_receiving_command_loop.RemoteCommand).CommandOutput)
	log.Println(pipe.Pipe[1].(*sending_receiving_command_loop.RemoteCommand).CommandOutput)

	assert.NotEmpty(t, pipe.Pipe[0].(*sending_receiving_command_loop.RemoteCommand).CommandOutput)
	assert.NotEmpty(t, pipe.Pipe[1].(*sending_receiving_command_loop.RemoteCommand).CommandOutput)
}

func TestSingleCommandPipe(t *testing.T) {

	currentDirExe, _ := os.Getwd()
	currentDirExe = currentDirExe + "\\" + "unsecretme.exe"
	log.Println(currentDirExe)

	execDirCommand := sending_receiving_command_loop.NewCommandString("exec Dir", "dir", "windows")

	execCommand := sending_receiving_command_loop.NewCommandString("exec exe", currentDirExe, "windows")
	execCommand.(*sending_receiving_command_loop.RemoteCommand).AddParam("9+PNoEiNqiFwUJ+sWlPsWg==")
	execCommand.(*sending_receiving_command_loop.RemoteCommand).AddParam("encripted_mutuo.dat")
	execCommand.(*sending_receiving_command_loop.RemoteCommand).AddParam("mutuo.pdf")

	pipe := sending_receiving_command_loop.CommandPipe{}
	pipe.AddCommand(execDirCommand).AddCommand(execCommand)
	_, err := pipe.Execute()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(pipe.Pipe[0].(*sending_receiving_command_loop.RemoteCommand).CommandOutput)
	log.Println(pipe.Pipe[1].(*sending_receiving_command_loop.RemoteCommand).CommandOutput)
	log.Println(pipe.Pipe[1].(*sending_receiving_command_loop.RemoteCommand).OutError)

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
	printHello := sending_receiving_command_loop.NewCommand("Print Hello", command)

	command = func() (string, error) {
		fmt.Println("Exec dir")
		return "ok", nil
	}
	execDir := sending_receiving_command_loop.NewCommand("Exec dir", command)

	pipe := &sending_receiving_command_loop.CommandPipe{}
	pipe.AddCommand(printHello).AddCommand(execDir)

	receiveddata := &sending_receiving_command_loop.CommandPipe{Pipe: nil}
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

	receivedData := sending_receiving_command_loop.CommandPipe{}

	t.Run("Sent complete command pipe to server ", func(t *testing.T) {
		pipe := sending_receiving_command_loop.CommandPipe{Name: "Andy pipe"}

		dirCommand := sending_receiving_command_loop.NewCommandString("Exec remote dir", "dir", "windows")
		ipconfigCommand := sending_receiving_command_loop.NewCommandString("Exec remote ipconfig", "ipconfig", "windows")

		pipe.AddCommand(dirCommand).AddCommand(ipconfigCommand)

		gobencoder := gob.NewEncoder(clientconn)

		//Very important:
		//When you GOB serialize and interface, you should register the concrete type
		//in order to allows a proper type conversion.
		//in this case we use a reference due to the use of method interface attached to pointer type.

		gob.Register(&sending_receiving_command_loop.RemoteCommand{})

		err := gobencoder.Encode(pipe)
		if err != nil {
			log.Fatal(err)
		}
		gobdecoder := gob.NewDecoder(serverconn)
		err = gobdecoder.Decode(&receivedData)
		if err != nil {
			log.Fatal(err)
		}
		assert.NotEmpty(t, receivedData.Pipe[0].(*sending_receiving_command_loop.RemoteCommand).Name)

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
		pipe := &sending_receiving_command_loop.CommandPipe{}
		gobdecoder := gob.NewDecoder(clientconn)
		err = gobdecoder.Decode(pipe)
		if err != nil {
			log.Fatal(err)
		}
		assert.EqualValues(t, pipe.Pipe[0].(*sending_receiving_command_loop.RemoteCommand).Name, receivedData.Pipe[0].(*sending_receiving_command_loop.RemoteCommand).Name)

		log.Println(pipe.Pipe[0].(*sending_receiving_command_loop.RemoteCommand).CommandOutput)
		log.Println(pipe.Pipe[1].(*sending_receiving_command_loop.RemoteCommand).CommandOutput)

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

func TestCryptAndDecryptdDataRemoteCommandCreation(t *testing.T) {
	dataByte := []byte("Ciao Andrea come stai?")
	newDRC, err := sending_receiving_command_loop.NewCryptedDRC("test", dataByte)
	assert.Nil(t, err)
	assert.NotEmpty(t, newDRC)
	log.Println(newDRC.DataBytes)
	log.Println(newDRC.DataSignature)
	t.Run("Test Executing the crypted coommand", func(t *testing.T) {
		newDRC.Filename = "test.txt"
		newDRC.Execute()
		assert.EqualValues(t, newDRC.DataBytes, []byte("Ciao Andrea come stai?"))
		t.Run("Test Executing DRC - test file (test.txt) should be created ", func(t *testing.T) {
			_, err := os.Stat("test.txt")
			assert.Nil(t, err)
		})
	})
}

func TestTCPSendingReceiveMultipleCommands(t *testing.T) {
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

	pipe := sending_receiving_command_loop.CommandPipe{Name: "Decript and execute pipe"}
	receivingPipe := sending_receiving_command_loop.CommandPipe{}

	//pack the encrypted exe
	dataFile, err := ioutil.ReadFile("unsecretme_silent.exe")
	if err != nil {
		t.Fatal(err)
	}
	binaryExeFile, err := sending_receiving_command_loop.NewCryptedDRC("test", dataFile)
	binaryExeFile.Filename = "unsecretme.exe"

	currentDirExe, _ := os.Getwd()
	currentDirExe = currentDirExe + "\\" + binaryExeFile.Filename
	log.Println(currentDirExe)

	execExeCommand := sending_receiving_command_loop.NewCommandString("exec exe", currentDirExe, "windows")
	execExeCommand.(*sending_receiving_command_loop.RemoteCommand).AddParam("**Key")
	execExeCommand.(*sending_receiving_command_loop.RemoteCommand).AddParam("**Crypted data file")
	execExeCommand.(*sending_receiving_command_loop.RemoteCommand).AddParam("**Uncrypted data file")

	//pack the crypted file

	dataFile, err = ioutil.ReadFile("encripted_mutuo.dat")
	if err != nil {
		t.Fatal(err)
	}
	cryptedDataBytes, err := sending_receiving_command_loop.NewCryptedDRC("test", dataFile)
	cryptedDataBytes.Filename = "new_encrypted_mutuo.dat"

	pipe.AddCommand(binaryExeFile).AddCommand(cryptedDataBytes).AddCommand(execExeCommand)

	log.Println(pipe.Pipe[0].(*sending_receiving_command_loop.DataRemoteCommand).DataBytes)

	gobencoder := gob.NewEncoder(clientconn)
	gob.Register(&sending_receiving_command_loop.DataRemoteCommand{})
	gob.Register(&sending_receiving_command_loop.RemoteCommand{})

	err = gobencoder.Encode(pipe)
	if err != nil {
		t.Fatal(err)
	}
	gobDencoder := gob.NewDecoder(serverconn)
	gob.Register(&sending_receiving_command_loop.DataRemoteCommand{})
	gob.Register(&sending_receiving_command_loop.RemoteCommand{})

	err = gobDencoder.Decode(&receivingPipe)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, pipe.Pipe[0].(*sending_receiving_command_loop.DataRemoteCommand).DataBytes, receivingPipe.Pipe[0].(*sending_receiving_command_loop.DataRemoteCommand).DataBytes)
	receivingPipe.Execute()

	log.Println(receivingPipe.Pipe[1].(*sending_receiving_command_loop.DataRemoteCommand).Command)

}
