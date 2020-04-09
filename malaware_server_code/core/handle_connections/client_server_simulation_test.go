package handle_connections

import (
	"../../../Execute_systems_commands"
	"../handle_connections"
	"encoding/gob"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

type Data struct {
	Name string
	Id   int
	Age  int
}

func TestClientServerCommunicationLab(t *testing.T) {
	localip, err := Execute_systems_commands.RetrieveCurrentIP()
	assert.Nil(t, err)
	localip = localip + handle_connections.Port
	serverconn, clientconn, err := handle_connections.SetClientServerConnectionTo(localip)
	assert.Nil(t, err)
	assert.NotNil(t, serverconn)
	assert.NotNil(t, clientconn)
	defer clientconn.Close()
	defer serverconn.Close()
}

func TestClientServerCommunicationNotValidIp(t *testing.T) {
	localip := "1.1.1.1:111"
	_, _, err := handle_connections.SetClientServerConnectionTo(localip)
	assert.NotNil(t, err)
}

func TestSendStreamedDatatoClient(t *testing.T) {
	localip, err := Execute_systems_commands.RetrieveCurrentIP()
	assert.Nil(t, err)
	localip = localip + handle_connections.Port
	serverconn, clientconn, err := handle_connections.SetClientServerConnectionTo(localip)
	assert.Nil(t, err)
	assert.NotNil(t, serverconn)
	assert.NotNil(t, clientconn)
	defer clientconn.Close()
	defer serverconn.Close()
	data2send := "This is a string message\n"
	numberofbytes, err := handle_connections.SendStringTo(clientconn, data2send)
	assert.EqualValues(t, 25, numberofbytes)
	dataReceived, err := handle_connections.ReceiveStringFrom(serverconn)
	assert.EqualValues(t, "This is a string message\n", dataReceived)
}

func TestSendGobDataOverTheWire(t *testing.T) {
	var serverconn net.Conn
	var clientconn net.Conn

	t.Run("Test Client/Server connection", func(t *testing.T) {
		localip, err := Execute_systems_commands.RetrieveCurrentIP()
		assert.Nil(t, err)
		localip = localip + handle_connections.Port
		serverconn, clientconn, err = handle_connections.SetClientServerConnectionTo(localip)
		assert.Nil(t, err)
		assert.NotNil(t, serverconn)
		assert.NotNil(t, clientconn)

	})
	defer serverconn.Close()
	defer clientconn.Close()

	sentdata := &Data{
		Name: "Andrea\nMazzanti",
		Id:   100,
		Age:  500,
	}

	receiveddata := &Data{}
	gobencoder := gob.NewEncoder(clientconn)
	err := gobencoder.Encode(sentdata)
	if err != nil {
		t.Fail()
	}
	gobdecoder := gob.NewDecoder(serverconn)
	err = gobdecoder.Decode(receiveddata)
	if err != nil {
		t.Fail()
	}
	assert.EqualValues(t, "Andrea\nMazzanti", receiveddata.Name)

}
