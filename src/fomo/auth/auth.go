package auth

import (
	"errors"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// Agent ...
func Agent() (ssh.AuthMethod, error) {
	addr := os.Getenv("SSH_AUTH_SOCK")
	if addr == "" {
		return nil, errors.New("no agent found")
	}

	c, err := net.Dial("unix", addr)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeysCallback(agent.NewClient(c).Signers), nil
}
