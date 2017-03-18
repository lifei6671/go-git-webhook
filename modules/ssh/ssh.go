package ssh

import (
	"github.com/golang/crypto/ssh"
	"unicode/utf8"
)

type Config struct {
	ssh.ClientConfig
}

func Connection(user,host, pass string) (*ssh.Client, *ssh.Session, error) {

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},

	}

	if count := utf8.RuneCountInString(pass); count > 50 {
		//signer, err := ssh.ParsePrivateKey([]byte(pass))

		//ssh.PublicKeys()
		//sshConfig.Auth = []ssh.AuthMethod{  }
	}

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
