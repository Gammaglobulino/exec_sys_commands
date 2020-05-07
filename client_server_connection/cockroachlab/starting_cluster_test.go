package cockroachlab

import (
	"../core/handle_connections"
	"../sending_receiving_command_loop"
	"encoding/gob"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestTCPStartNewCluster(t *testing.T) {
	var serverconn net.Conn
	var clientconn net.Conn
	var localip string
	var err error

	t.Run("Test Client/Server connection", func(t *testing.T) {
		localip, err = handle_connections.GetLocalIp04Str()
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

	localip, _ = handle_connections.GetLocalIp04Str()

	pipe := sending_receiving_command_loop.CommandPipe{Name: "Remote start a local cockroach cluster"}
	receivingPipe := sending_receiving_command_loop.CommandPipe{}

	cockroachNode1 := sending_receiving_command_loop.NewCommandString("start cockroach node1", "cockroach start", "windows")
	cockroachNode1.(*sending_receiving_command_loop.RemoteCommand).AddParam("--insecure")
	cockroachNode1.(*sending_receiving_command_loop.RemoteCommand).AddParam("--store=node1")
	cockroachNode1.(*sending_receiving_command_loop.RemoteCommand).AddParam("--listen-addr=" + localip + ":26257")
	cockroachNode1.(*sending_receiving_command_loop.RemoteCommand).AddParam("--http-addr=" + localip + ":8080")
	cockroachNode1.(*sending_receiving_command_loop.RemoteCommand).AddParam("--join=" + localip + ":26257," + localip + ":26258," + localip + ":26259")

	cockroachNode2 := sending_receiving_command_loop.NewCommandString("start cockroach node2", "cockroach start", "windows")
	cockroachNode2.(*sending_receiving_command_loop.RemoteCommand).AddParam("--insecure")
	cockroachNode2.(*sending_receiving_command_loop.RemoteCommand).AddParam("--store=node2")
	cockroachNode2.(*sending_receiving_command_loop.RemoteCommand).AddParam("--listen-addr=" + localip + ":26258")
	cockroachNode2.(*sending_receiving_command_loop.RemoteCommand).AddParam("--http-addr=" + localip + ":8081")
	cockroachNode2.(*sending_receiving_command_loop.RemoteCommand).AddParam("--join=" + localip + ":26257," + localip + ":26258," + localip + ":26259")

	cockroachNode3 := sending_receiving_command_loop.NewCommandString("start cockroach node3", "cockroach start", "windows")
	cockroachNode3.(*sending_receiving_command_loop.RemoteCommand).AddParam("--insecure")
	cockroachNode3.(*sending_receiving_command_loop.RemoteCommand).AddParam("--store=node3")
	cockroachNode3.(*sending_receiving_command_loop.RemoteCommand).AddParam("--listen-addr=" + localip + ":26259")
	cockroachNode3.(*sending_receiving_command_loop.RemoteCommand).AddParam("--http-addr=" + localip + ":8082")
	cockroachNode3.(*sending_receiving_command_loop.RemoteCommand).AddParam("--join=" + localip + ":26257," + localip + ":26258," + localip + ":26259")

	initCluster := sending_receiving_command_loop.NewCommandString("init cockroach cluster", "cockroach init", "windows")
	initCluster.(*sending_receiving_command_loop.RemoteCommand).AddParam("--insecure")
	initCluster.(*sending_receiving_command_loop.RemoteCommand).AddParam("--host=" + localip + ":26257")

	pipe.AddCommand(cockroachNode1).AddCommand(cockroachNode2).AddCommand(cockroachNode3).AddCommand(initCluster)
	assert.EqualValues(t, "cockroach start --insecure --store=node1 --listen-addr=192.168.36.161:26257 --http-addr=192.168.36.161:8080 --join=192.168.36.161:26257,192.168.36.161:26258,192.168.36.161:26259", pipe.Pipe[0].(*sending_receiving_command_loop.RemoteCommand).GetCommandAndArgs())

	gobencoder := gob.NewEncoder(clientconn)
	gob.Register(&sending_receiving_command_loop.RemoteCommand{})

	err = gobencoder.Encode(pipe)
	if err != nil {
		t.Fatal(err)
	}
	gobDencoder := gob.NewDecoder(serverconn)
	gob.Register(&sending_receiving_command_loop.RemoteCommand{})

	err = gobDencoder.Decode(&receivingPipe)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, pipe.Pipe[0].(*sending_receiving_command_loop.RemoteCommand).Name, receivingPipe.Pipe[0].(*sending_receiving_command_loop.RemoteCommand).Name)
	_, err = receivingPipe.Execute()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("End")
}
