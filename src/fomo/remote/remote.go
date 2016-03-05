package remote

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"time"

	"fomo/auth"
	"fomo/hosts"

	"golang.org/x/crypto/ssh"
)

const (
	tmpDir = "/tmp"
)

func sessionId() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func send(w io.Writer, name, perm string, data []byte) error {
	if _, err := fmt.Fprintf(w, "C%s %d %s\n", perm, len(data), name); err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}

	_, err := fmt.Fprint(w, "\x00")
	return err
}

func sendAsset(w io.Writer, name, perm string) error {
	b, err := Asset(name)
	if err != nil {
		return err
	}

	return send(w, name, perm, b)
}

func sendFile(w io.Writer, name, perm string) error {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	return send(w, filepath.Base(name), perm, b)
}

func makeDir(w io.Writer, name, perm string) error {
	if _, err := fmt.Fprintf(w, "D%s 0 %s\n", perm, name); err != nil {
		return err
	}
	return nil
}

func setup(c *ssh.Client, id, src string) error {
	sess, err := c.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	w, err := sess.StdinPipe()
	if err != nil {
		return err
	}

	var e error
	go func() {
		defer w.Close()

		if err := makeDir(w, id, "0755"); err != nil {
			e = err
			return
		}

		if err := sendAsset(w, "remote.py", "0755"); err != nil {
			e = err
			return
		}

		if err := sendFile(w, src, "0644"); err != nil {
			e = err
			return
		}
	}()

	if err := sess.Run(fmt.Sprintf("/usr/bin/scp -tr %s", tmpDir)); err != nil {
		return err
	}

	return e
}

func serve(nl net.Listener) {
	for {
		c, err := nl.Accept()
		if err == io.EOF {
			return
		} else if err != nil {
			log.Panic(err)
		}

		log.Println(c)
	}
}

func Run(h *hosts.Host, loc io.WriteCloser, src string) error {
	id, err := sessionId()
	if err != nil {
		return err
	}

	agt, err := auth.Agent()
	if err != nil {
		return err
	}

	s, err := ssh.Dial("tcp", h.Addr(), &ssh.ClientConfig{
		User: h.User,
		Auth: []ssh.AuthMethod{
			agt,
		}})
	if err != nil {
		return err
	}
	defer s.Close()

	if err := setup(s, id, src); err != nil {
		return err
	}

	nl, err := s.Listen("tcp", "localhost:0")
	if err != nil {
		return err
	}

	log.Println(nl.Addr())

	go serve(nl)

	time.Sleep(2 * time.Second)

	return nil
}
