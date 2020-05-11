package cockroachlab

import (
	"../cockroachlab"
	"../core/handle_connections"
	"../sending_receiving_command_loop"
	"encoding/gob"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"testing"
)

func TestCreateNewServerClusterPipe(t *testing.T) {
	localip := "192.168.177"
	pipe := cockroachlab.CreateNewServerClusterPipe(5, localip, 26257)
	assert.EqualValues(t, "cockroach start --insecure --store=node-0 --listen-addr=192.168.177:26257 --http-addr=192.168.177:8080 --join=192.168.177:26257,192.168.177:26258,192.168.177:26259,192.168.177:26260,192.168.177:26261", pipe.Pipe[0].(*sending_receiving_command_loop.RemoteCommand).GetCommandAndArgs())
	assert.EqualValues(t, "cockroach start --insecure --store=node-1 --listen-addr=192.168.177:26258 --http-addr=192.168.177:8081 --join=192.168.177:26257,192.168.177:26258,192.168.177:26259,192.168.177:26260,192.168.177:26261", pipe.Pipe[1].(*sending_receiving_command_loop.RemoteCommand).GetCommandAndArgs())
}

func TestCreateAndStartServerCluster(t *testing.T) {
	localip, err := handle_connections.GetLocalIp04Str()
	assert.Nil(t, err)
	pipe := cockroachlab.CreateNewServerClusterPipe(7, localip, 26257)
	msg, err := pipe.Execute()
	if err != nil {
		log.Println(msg)
		t.Fatal(err)

	}
	log.Println("Cluster started on:" + pipe.Pipe[0].(*sending_receiving_command_loop.RemoteCommand).GetCommandAndArgs())
}
func TestVerifyNodeSTatus(t *testing.T) {
	localip, err := handle_connections.GetLocalIp04Str()
	assert.Nil(t, err)
	host := localip + ":" + "26257"

	out, err := cockroachlab.VerifyNodeStatusOn(host)
	if err != nil {
		fmt.Println(out)
		t.Fatal(err)
	}
	assert.EqualValues(t, "id", out[:2])
}

func TestTCPStart5NodeCluster(t *testing.T) {
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
	cockroachNode1.(*sending_receiving_command_loop.RemoteCommand).AddParam("--join=" + localip + ":26257," + localip + ":26258," + localip + ":26259," + localip + ":26260," + localip + ":26261")

	fmt.Println(cockroachNode1.(*sending_receiving_command_loop.RemoteCommand).GetCommandAndArgs())

	cockroachNode2 := sending_receiving_command_loop.NewCommandString("start cockroach node2", "cockroach start", "windows")
	cockroachNode2.(*sending_receiving_command_loop.RemoteCommand).AddParam("--insecure")
	cockroachNode2.(*sending_receiving_command_loop.RemoteCommand).AddParam("--store=node2")
	cockroachNode2.(*sending_receiving_command_loop.RemoteCommand).AddParam("--listen-addr=" + localip + ":26258")
	cockroachNode2.(*sending_receiving_command_loop.RemoteCommand).AddParam("--http-addr=" + localip + ":8081")
	cockroachNode2.(*sending_receiving_command_loop.RemoteCommand).AddParam("--join=" + localip + ":26257," + localip + ":26258," + localip + ":26259," + localip + ":26260," + localip + ":26261")

	cockroachNode3 := sending_receiving_command_loop.NewCommandString("start cockroach node3", "cockroach start", "windows")
	cockroachNode3.(*sending_receiving_command_loop.RemoteCommand).AddParam("--insecure")
	cockroachNode3.(*sending_receiving_command_loop.RemoteCommand).AddParam("--store=node3")
	cockroachNode3.(*sending_receiving_command_loop.RemoteCommand).AddParam("--listen-addr=" + localip + ":26259")
	cockroachNode3.(*sending_receiving_command_loop.RemoteCommand).AddParam("--http-addr=" + localip + ":8082")
	cockroachNode3.(*sending_receiving_command_loop.RemoteCommand).AddParam("--join=" + localip + ":26257," + localip + ":26258," + localip + ":26259," + localip + ":26260," + localip + ":26261")

	cockroachNode4 := sending_receiving_command_loop.NewCommandString("start cockroach node4", "cockroach start", "windows")
	cockroachNode4.(*sending_receiving_command_loop.RemoteCommand).AddParam("--insecure")
	cockroachNode4.(*sending_receiving_command_loop.RemoteCommand).AddParam("--store=node4")
	cockroachNode4.(*sending_receiving_command_loop.RemoteCommand).AddParam("--listen-addr=" + localip + ":26260")
	cockroachNode4.(*sending_receiving_command_loop.RemoteCommand).AddParam("--http-addr=" + localip + ":8083")
	cockroachNode4.(*sending_receiving_command_loop.RemoteCommand).AddParam("--join=" + localip + ":26257," + localip + ":26258," + localip + ":26259," + localip + ":26260," + localip + ":26261")

	cockroachNode5 := sending_receiving_command_loop.NewCommandString("start cockroach node5", "cockroach start", "windows")
	cockroachNode5.(*sending_receiving_command_loop.RemoteCommand).AddParam("--insecure")
	cockroachNode5.(*sending_receiving_command_loop.RemoteCommand).AddParam("--store=node5")
	cockroachNode5.(*sending_receiving_command_loop.RemoteCommand).AddParam("--listen-addr=" + localip + ":26261")
	cockroachNode5.(*sending_receiving_command_loop.RemoteCommand).AddParam("--http-addr=" + localip + ":8084")
	cockroachNode5.(*sending_receiving_command_loop.RemoteCommand).AddParam("--join=" + localip + ":26257," + localip + ":26258," + localip + ":26259," + localip + ":26260," + localip + ":26261")

	initCluster := sending_receiving_command_loop.NewCommandString("init cockroach cluster", "cockroach init", "windows")
	initCluster.(*sending_receiving_command_loop.RemoteCommand).AddParam("--insecure")
	initCluster.(*sending_receiving_command_loop.RemoteCommand).AddParam("--host=" + localip + ":26257")

	pipe.AddCommand(cockroachNode1).AddCommand(cockroachNode2).AddCommand(cockroachNode3).AddCommand(cockroachNode4).AddCommand(cockroachNode5).AddCommand(initCluster)

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
	msg, err := receivingPipe.Execute()
	if err != nil {
		log.Println(msg)
		t.Fatal(err)
	}

}

func TestTCPStart5NodeClusterAndDecomissionOneNode(t *testing.T) {
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
	cockroachNode1.(*sending_receiving_command_loop.RemoteCommand).AddParam("--join=" + localip + ":26257," + localip + ":26258," + localip + ":26259," + localip + ":26260," + localip + ":26261")

	fmt.Println(cockroachNode1.(*sending_receiving_command_loop.RemoteCommand).GetCommandAndArgs())

	cockroachNode2 := sending_receiving_command_loop.NewCommandString("start cockroach node2", "cockroach start", "windows")
	cockroachNode2.(*sending_receiving_command_loop.RemoteCommand).AddParam("--insecure")
	cockroachNode2.(*sending_receiving_command_loop.RemoteCommand).AddParam("--store=node2")
	cockroachNode2.(*sending_receiving_command_loop.RemoteCommand).AddParam("--listen-addr=" + localip + ":26258")
	cockroachNode2.(*sending_receiving_command_loop.RemoteCommand).AddParam("--http-addr=" + localip + ":8081")
	cockroachNode2.(*sending_receiving_command_loop.RemoteCommand).AddParam("--join=" + localip + ":26257," + localip + ":26258," + localip + ":26259," + localip + ":26260," + localip + ":26261")

	cockroachNode3 := sending_receiving_command_loop.NewCommandString("start cockroach node3", "cockroach start", "windows")
	cockroachNode3.(*sending_receiving_command_loop.RemoteCommand).AddParam("--insecure")
	cockroachNode3.(*sending_receiving_command_loop.RemoteCommand).AddParam("--store=node3")
	cockroachNode3.(*sending_receiving_command_loop.RemoteCommand).AddParam("--listen-addr=" + localip + ":26259")
	cockroachNode3.(*sending_receiving_command_loop.RemoteCommand).AddParam("--http-addr=" + localip + ":8082")
	cockroachNode3.(*sending_receiving_command_loop.RemoteCommand).AddParam("--join=" + localip + ":26257," + localip + ":26258," + localip + ":26259," + localip + ":26260," + localip + ":26261")

	cockroachNode4 := sending_receiving_command_loop.NewCommandString("start cockroach node4", "cockroach start", "windows")
	cockroachNode4.(*sending_receiving_command_loop.RemoteCommand).AddParam("--insecure")
	cockroachNode4.(*sending_receiving_command_loop.RemoteCommand).AddParam("--store=node4")
	cockroachNode4.(*sending_receiving_command_loop.RemoteCommand).AddParam("--listen-addr=" + localip + ":26260")
	cockroachNode4.(*sending_receiving_command_loop.RemoteCommand).AddParam("--http-addr=" + localip + ":8083")
	cockroachNode4.(*sending_receiving_command_loop.RemoteCommand).AddParam("--join=" + localip + ":26257," + localip + ":26258," + localip + ":26259," + localip + ":26260," + localip + ":26261")

	cockroachNode5 := sending_receiving_command_loop.NewCommandString("start cockroach node5", "cockroach start", "windows")
	cockroachNode5.(*sending_receiving_command_loop.RemoteCommand).AddParam("--insecure")
	cockroachNode5.(*sending_receiving_command_loop.RemoteCommand).AddParam("--store=node5")
	cockroachNode5.(*sending_receiving_command_loop.RemoteCommand).AddParam("--listen-addr=" + localip + ":26261")
	cockroachNode5.(*sending_receiving_command_loop.RemoteCommand).AddParam("--http-addr=" + localip + ":8084")
	cockroachNode5.(*sending_receiving_command_loop.RemoteCommand).AddParam("--join=" + localip + ":26257," + localip + ":26258," + localip + ":26259," + localip + ":26260," + localip + ":26261")

	initCluster := sending_receiving_command_loop.NewCommandString("init cockroach cluster", "cockroach init", "windows")
	initCluster.(*sending_receiving_command_loop.RemoteCommand).AddParam("--insecure")
	initCluster.(*sending_receiving_command_loop.RemoteCommand).AddParam("--host=" + localip + ":26257")

	pipe.AddCommand(cockroachNode1).AddCommand(cockroachNode2).AddCommand(cockroachNode3).AddCommand(cockroachNode4).AddCommand(cockroachNode5).AddCommand(initCluster)

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
	msg, err := receivingPipe.Execute()
	if err != nil {
		log.Println(msg)
		t.Fatal(err)
	}

	pipe = sending_receiving_command_loop.CommandPipe{Name: "Decomission a node"}
	receivingPipe = sending_receiving_command_loop.CommandPipe{}

	killNode := sending_receiving_command_loop.NewCommandString("stop cockroach node3", "cockroach quit", "windows")
	killNode.(*sending_receiving_command_loop.RemoteCommand).AddParam("--decommission")
	killNode.(*sending_receiving_command_loop.RemoteCommand).AddParam("--insecure")
	killNode.(*sending_receiving_command_loop.RemoteCommand).AddParam("--host=" + localip + ":26257")

	pipe.AddCommand(killNode)

	gobencoder = gob.NewEncoder(clientconn)
	gob.Register(&sending_receiving_command_loop.RemoteCommand{})

	err = gobencoder.Encode(pipe)
	if err != nil {
		t.Fatal(err)
	}
	gobDencoder = gob.NewDecoder(serverconn)
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
	log.Println(receivingPipe.Pipe[0].(*sending_receiving_command_loop.RemoteCommand).GetCommandAndArgs())
}
