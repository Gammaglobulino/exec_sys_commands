package ssh_server

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func StartSSHServerWithKPairAuthTo(hostname string, privateKeyFilename string) error {

	//use ssh -i gammaserver gamma@172.17.211.113 -p 2200 to connect with a standard client
	//you need to create the authorized_keys file on the same server folder, adding the auth keys

	authorizedKeysBytes, err := ioutil.ReadFile("authorized_keys")
	if err != nil {
		log.Fatalf("Failed to load authorized_keys, err: %v", err)
		return err
	}

	authorizedKeysMap := map[string]bool{}
	for len(authorizedKeysBytes) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil {
			log.Fatal(err)
			return err
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

	privatePEMBytes, err := ioutil.ReadFile(privateKeyFilename)
	if err != nil {
		log.Fatal(err)
		return err
	}

	pk, err := ssh.ParsePrivateKey(privatePEMBytes)
	if err != nil {
		return err
	}
	serverConfig.AddHostKey(pk)

	listener, err := net.Listen("tcp", hostname)
	if err != nil {
		return err
	}
	log.Printf("Listening on %s \n", hostname)
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
						req.Reply(true, nil)
					case "pty-req":
						req.Reply(true, nil)
					case "env":
						req.Reply(true, nil)
					case "exec":
						log.Println(string(req.Payload))
						req.Reply(true, []byte("ok\n"))
					default:
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
						os.Exit(0)
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
