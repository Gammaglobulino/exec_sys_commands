package ssh_server

import (
	"../../malaware_server_code/core/handle_connections"
	"../storing_secure_password"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"log"
	"net"
	"os"
	"testing"
)

func TestStartingNewSSHServerWithUN_PWAuth(t *testing.T) {

	authCallback := func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {

		hashedPasword := storing_secure_password.HashPassword(string(pass))

		if c.User() == "gamma" && hashedPasword == "39a38b570c1cfa4f169a220b12d79fcb8ce936ab35ebe8f445f506d33f537f75" {
			return nil, nil
		}
		return nil, fmt.Errorf("Not a vaild password for %s provided", c.User())
	}
	serverConfig := &ssh.ServerConfig{
		Config:                      ssh.Config{},
		NoClientAuth:                false,
		MaxAuthTries:                0,
		PasswordCallback:            authCallback,
		PublicKeyCallback:           nil,
		KeyboardInteractiveCallback: nil,
		AuthLogCallback:             nil,
		ServerVersion:               "",
		BannerCallback:              nil,
		GSSAPIWithMICConfig:         nil,
	}
	privatePEMBytes, err := ioutil.ReadFile("gammaserver")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(privatePEMBytes)

	pk, err := ssh.ParsePrivateKey(privatePEMBytes)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	serverConfig.AddHostKey(pk)
	localip, err := handle_connections.GetLocalIp04Str()
	assert.Nil(t, err)
	localip = localip + ":2200"

	listener, err := net.Listen("tcp", localip)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	log.Printf("Listening on %s \n", localip)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to handshake (%s)\n", err)
			continue
		}
		sshConn, chans, reqs, err := ssh.NewServerConn(conn, serverConfig)
		if err != nil {
			log.Println(err)
			continue
		}
		go ssh.DiscardRequests(reqs)
		// Service the incoming Channel channel.
		log.Printf("logged in %s", sshConn.User())

		for newChannel := range chans {

			if newChannel.ChannelType() != "session" {
				newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
				continue
			}
			channel, requests, err := newChannel.Accept()
			if err != nil {
				log.Fatalf("Could not accept channel: %v", err)
			}

			go func(in <-chan *ssh.Request) {
				for req := range in {
					switch req.Type {
					case "shell":
						// We only accept the default shell
						// (i.e. no command in the Payload)
						if len(req.Payload) == 0 {
							req.Reply(true, nil)
						}
					case "pty-req":
						req.Reply(true, nil)
					}
				}
			}(requests)

			term := terminal.NewTerminal(channel, "> ")

			go func() {
				defer channel.Close()
				for {
					line, err := term.ReadLine()
					if err != nil {
						break
					}
					if line == "exit" {
						_, _ = term.Write([]byte("Hello\n"))
						fmt.Println("Bye!")
						os.Exit(2)
					}
					fmt.Println(line)
					nbytes, err := term.Write([]byte("Hello\n"))
					if err != nil {
						log.Fatal(err)
					}
					fmt.Printf("Sent %d bytes back to client\n", nbytes)
				}
			}()
		}

	}

}

func TestStartingNewSSHServerWithPrivateKeyAuth(t *testing.T) {
	authorizedKeysBytes, err := ioutil.ReadFile("authorized_keys")
	if err != nil {
		log.Fatalf("Failed to load authorized_keys, err: %v", err)
	}

	authorizedKeysMap := map[string]bool{}
	for len(authorizedKeysBytes) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil {
			log.Fatal(err)
		}
		authorizedKeysMap[string(pubKey.Marshal())] = true
		authorizedKeysBytes = rest
	}

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	serverConfig := &ssh.ServerConfig{
		// Remove to disable public key auth.
		PublicKeyCallback: func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			if authorizedKeysMap[string(pubKey.Marshal())] {
				return &ssh.Permissions{
					// Record the public key used for authentication.
					Extensions: map[string]string{
						"pubkey-fp": ssh.FingerprintSHA256(pubKey),
					},
				}, nil
			}
			return nil, fmt.Errorf("unknown public key for %q", c.User())
		},
	}

	privatePEMBytes, err := ioutil.ReadFile("gammaserver")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(privatePEMBytes)

	pk, err := ssh.ParsePrivateKey(privatePEMBytes)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	serverConfig.AddHostKey(pk)
	localip, err := handle_connections.GetLocalIp04Str()
	assert.Nil(t, err)
	localip = localip + ":2200"

	listener, err := net.Listen("tcp", localip)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	log.Printf("Listening on %s \n", localip)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to handshake (%s)\n", err)
			continue
		}
		sshConn, chans, reqs, err := ssh.NewServerConn(conn, serverConfig)
		if err != nil {
			log.Println(err)
			continue
		}
		go ssh.DiscardRequests(reqs)
		// Service the incoming Channel channel.
		log.Printf("logged in %s", sshConn.User())

		for newChannel := range chans {

			if newChannel.ChannelType() != "session" {
				newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
				continue
			}
			channel, requests, err := newChannel.Accept()
			if err != nil {
				log.Fatalf("Could not accept channel: %v", err)
			}

			go func(in <-chan *ssh.Request) {
				for req := range in {
					switch req.Type {
					case "shell":
						// We only accept the default shell
						// (i.e. no command in the Payload)
						if len(req.Payload) == 0 {
							log.Println(req.Payload)
							req.Reply(true, nil)
						}
					case "pty-req":
						req.Reply(true, nil)
					}
				}
			}(requests)

			term := terminal.NewTerminal(channel, "> ")

			go func() {
				defer channel.Close()
				for {
					line, err := term.ReadLine()
					if err != nil {
						break
					}
					if line == "exit" {
						_, _ = term.Write([]byte("Hello\n"))
						fmt.Println("Bye!")
						os.Exit(2)
					}
					fmt.Println(line)
					nbytes, err := term.Write([]byte("Hello\n"))
					if err != nil {
						log.Fatal(err)
					}
					fmt.Printf("Sent %d bytes back to client\n", nbytes)
				}
			}()
		}

	}

}
