package ssh

import (
	"github.com/golang/crypto/ssh"
)

type Config struct {
	ssh.ClientConfig
}

func Connection(user,host, pass string) (*ssh.Client, *ssh.Session, error) {

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},

	}

	signer, err := ssh.ParsePrivateKey([]byte(pass))

	if err == nil {

		sshConfig.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer) }
	}

	//sshConfig.SetDefaults()

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()

	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}
