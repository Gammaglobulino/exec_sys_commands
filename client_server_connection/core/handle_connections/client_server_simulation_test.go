package handle_connections

import (
	"../handle_connections"
	"encoding/gob"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"testing"
)

type SentData struct {
	Name string
	Id   int
	Age  int
}
type ReceiveData struct {
	Name string
	Id   int
}

func TestClientServerCommunicationLab(t *testing.T) {
	localip, err := handle_connections.GetLocalIp04Str()
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
	localip := "0.0.0.0:0"
	_, _, err := handle_connections.SetClientServerConnectionTo(localip)
	assert.NotNil(t, err)

}

func TestGetLocalIp04(t *testing.T) {
	localip, err := handle_connections.GetLocalIp04Str()
	assert.Nil(t, err)
	assert.Contains(t, localip, "")
	log.Println(localip)
}

func TestSendStreamedDatatoClient(t *testing.T) {
	localip, err := handle_connections.GetLocalIp04Str()
	assert.Nil(t, err)
	assert.Contains(t, localip, "")

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

	data2send = "another string sent\n"
	numberofbytes, err = handle_connections.SendStringTo(clientconn, data2send)
	assert.EqualValues(t, 20, numberofbytes)

	dataReceived, err = handle_connections.ReceiveStringFrom(serverconn)
	assert.EqualValues(t, "another string sent\n", dataReceived)

}

func TestSendGobDataOverTheWire(t *testing.T) {
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

	sentdata := &SentData{
		Name: "Andrea\nMazzanti",
		Id:   100,
		Age:  500,
	}

	receiveddata := &ReceiveData{}
	gobencoder := gob.NewEncoder(clientconn)
	err := gobencoder.Encode(sentdata)
	if err != nil {
		t.Fail()
	}
	log.Println(sentdata)
	gobdecoder := gob.NewDecoder(serverconn)
	err = gobdecoder.Decode(receiveddata)
	if err != nil {
		t.Fail()
	}
	log.Println(receiveddata)
	assert.EqualValues(t, "Andrea\nMazzanti", receiveddata.Name)

}
