package local

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
)

func sockName() (string, error) {
	var r [8]byte
	if _, err := rand.Read(r[:]); err != nil {
		return "", err
	}

	return filepath.Join(
		os.TempDir(),
		fmt.Sprintf("fomo-%s", hex.EncodeToString(r[:]))), nil
}

type conn struct {
	net.Conn
	l net.Listener
	p *os.Process
	d string
}

func (c *conn) Close() error {
	var ea, eb, ec error

	if c.p != nil {
		log.Println("kill process")
		ea = c.p.Kill()
	}

	if c.Conn != nil {
		log.Println("close conn")
		eb = c.Conn.Close()
	}

	if c.l != nil {
		log.Println("close listener")
		ec = c.l.Close()
	}

	if err := os.RemoveAll(c.d); err != nil {
		return err
	}

	if ea != nil {
		return ea
	}

	if eb != nil {
		return eb
	}

	return ec
}

func unpack(dst string, srcs ...string) error {
	for _, src := range srcs {
		b, err := Asset(src)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(filepath.Join(dst, src), b, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func startProc(tmp, src, sck string) (*os.Process, error) {
	c := exec.Command("python", filepath.Join(tmp, "local.py"), sck, src)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	if err := c.Start(); err != nil {
		return nil, err
	}
	return c.Process, nil
}

func Start(src string) (io.WriteCloser, error) {

	tmp, err := ioutil.TempDir(os.TempDir(), "fomo-")
	if err != nil {
		return nil, err
	}

	c := &conn{d: tmp}

	if err := unpack(tmp, "local.py"); err != nil {
		c.Close()
		return nil, err
	}

	sck := filepath.Join(tmp, "sock")
	nl, err := net.Listen("unix", sck)
	if err != nil {
		c.Close()
		return nil, err
	}
	c.l = nl

	p, err := startProc(tmp, src, sck)
	if err != nil {
		c.Close()
		return nil, err
	}
	c.p = p

	nc, err := nl.Accept()
	if err != nil {
		c.Close()
		return nil, err
	}
	c.Conn = nc

	return c, nil
}