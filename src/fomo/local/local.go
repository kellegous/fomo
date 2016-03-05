package local

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"fomo/hosts"
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

type Conn struct {
	c  net.Conn
	l  net.Listener
	p  *os.Process
	d  string
	lk sync.Mutex
}

func (c *Conn) Submit(h *hosts.Host, n int64, r io.Reader) error {
	c.lk.Lock()
	defer c.lk.Unlock()

	_, err := io.CopyN(c.c, r, n)
	return err
}

func (c *Conn) Close() error {
	var ea, eb, ec error

	if c.p != nil {
		ea = c.p.Kill()
	}

	if c.c != nil {
		eb = c.c.Close()
	}

	if c.l != nil {
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

func Start(src string) (*Conn, error) {

	tmp, err := ioutil.TempDir(os.TempDir(), "fomo-")
	if err != nil {
		return nil, err
	}

	c := &Conn{d: tmp}

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
	c.c = nc

	return c, nil
}
