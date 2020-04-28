package ssh_client

import (
	"../../client_server_connection/core/handle_connections"
	"../ssh_client"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
	"log"
	"testing"
)

func TestSSHClientWithUNPWAuth(t *testing.T) {
	un := "gamma"
	pw := "pan"

	localip, err := handle_connections.GetLocalIp04Str()
	assert.Nil(t, err)
	localip = localip + ":2200"

	clientConfig := &ssh.ClientConfig{
		Config:            ssh.Config{},
		User:              un,
		Auth:              []ssh.AuthMethod{ssh.Password(pw)},
		HostKeyCallback:   ssh.InsecureIgnoreHostKey(),
		BannerCallback:    nil,
		ClientVersion:     "",
		HostKeyAlgorithms: nil,
		Timeout:           0,
	}
	client, err := ssh.Dial("tcp", localip, clientConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(client.ClientVersion()))
}
func TestSSHClientWithKeyPairAuth(t *testing.T) {
	un := "gamma"

	localip, err := handle_connections.GetLocalIp04Str()
	assert.Nil(t, err)
	localip = localip + ":2200"

	privateKey, err := ssh_client.GetKeySignerFrom("gammaserver")
	if err != nil {
		t.Fail()
	}
	clientConfig := &ssh.ClientConfig{
		Config:            ssh.Config{},
		User:              un,
		Auth:              []ssh.AuthMethod{ssh.PublicKeys(privateKey)},
		HostKeyCallback:   ssh.FixedHostKey(privateKey.PublicKey()),
		BannerCallback:    nil,
		ClientVersion:     "",
		HostKeyAlgorithms: nil,
		Timeout:           0,
	}
	client, err := ssh.Dial("tcp", localip, clientConfig)
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}
	log.Println(string(client.ClientVersion()))

}

func TestSSHVerifyRemoteHostKeyPair(t *testing.T) {
	un := "gamma"

	localip, err := handle_connections.GetLocalIp04Str()
	assert.Nil(t, err)
	localip = localip + ":2200"

	privateKey, err := ssh_client.GetKeySignerFrom("gammaserver")
	if err != nil {
		t.Fail()
	}
	log.Println(privateKey.PublicKey())
	serverPublicKey, err := ssh_client.LoadServerPublicKeyFrom("c:/Users/Drako/.ssh/known_hosts")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}
	log.Println(serverPublicKey)

	clientConfig := &ssh.ClientConfig{
		Config:          ssh.Config{},
		User:            un,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(privateKey)},
		HostKeyCallback: ssh.FixedHostKey(serverPublicKey),
		BannerCallback:  nil,
		ClientVersion:   "",
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoDSA,
			ssh.KeyAlgoECDSA256,
			ssh.KeyAlgoECDSA384,
			ssh.KeyAlgoECDSA521,
			ssh.KeyAlgoED25519,
		},
		Timeout: 0,
	}
	client, err := ssh.Dial("tcp", localip, clientConfig)
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}
	log.Println(string(client.ClientVersion()))

}

func TestVerifyRemoteSSHServerConnectionUsingJeyPairAuth(t *testing.T) {
	userName := "gamma"
	pkFileName := "gammaserver"
	knownHostsFilename := "c:/Users/Drako/.ssh/known_hosts"

	localip, err := handle_connections.GetLocalIp04Str()
	assert.Nil(t, err)
	localip = localip + ":2200"

	client, err := ssh_client.ConnectToSSHServerWithKPAuthOn(pkFileName, knownHostsFilename, localip, userName)
	assert.Nil(t, err)
	assert.NotEmpty(t, client)
	assert.EqualValues(t, "SSH-2.0-Go", string(client.ClientVersion()))

}

func TestExecuteCommandOnRemoteSSHServer(t *testing.T) {
	userName := "gamma"
	pkFileName := "gammaserver"
	knownHostsFilename := "c:/Users/Drako/.ssh/known_hosts"

	localip, err := handle_connections.GetLocalIp04Str()
	assert.Nil(t, err)
	localip = localip + ":2200"

	client, err := ssh_client.ConnectToSSHServerWithKPAuthOn(pkFileName, knownHostsFilename, localip, userName)
	assert.Nil(t, err)
	assert.NotEmpty(t, client)
	assert.EqualValues(t, "SSH-2.0-Go", string(client.ClientVersion()))
	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	session.Run("hostname")
	if err != nil {
		log.Fatal(err)
	}

}
