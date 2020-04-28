package ssh_client

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

func GetKeySignerFrom(privateKeyFile string) (ssh.Signer, error) {
	privateKeyData, err := ioutil.ReadFile("gammaserver")
	if err != nil {
		return nil, err
	}
	privateKey, err := ssh.ParsePrivateKey(privateKeyData)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func LoadServerPublicKeyFrom(knownHostFile string) (ssh.PublicKey, error) {
	publicKeyData, err := ioutil.ReadFile(knownHostFile)
	if err != nil {
		return nil, err
	}
	_, _, publicKey, _, _, err := ssh.ParseKnownHosts(publicKeyData)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

func ConnectToSSHServerWithKPAuthOn(privateKeyFile string, knownHostFile string, ip string, username string) (*ssh.Client, error) {

	privateKey, err := GetKeySignerFrom(privateKeyFile)
	if err != nil {
		return nil, err
	}
	serverPublicKey, err := LoadServerPublicKeyFrom(knownHostFile)
	if err != nil {
		return nil, err
	}

	clientConfig := &ssh.ClientConfig{
		Config:          ssh.Config{},
		User:            username,
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
	client, err := ssh.Dial("tcp", ip, clientConfig)
	if err != nil {
		return nil, nil
	}
	return client, nil
}
