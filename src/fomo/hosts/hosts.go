package hosts

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

const rootDir = "~/.fomo/hosts"

type Host struct {
	User string
	Host string
	Port int
}

type Loader struct {
	dir string
}

func (h *Host) Addr() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

func toHost(desc string, username string) (*Host, error) {
	h := &Host{
		User: username,
		Port: 22,
	}

	ix := strings.Index(desc, "@")
	if ix >= 0 {
		h.User = desc[:ix]
		desc = desc[ix+1:]
	}

	ix = strings.LastIndex(desc, ":")
	if ix >= 0 {
		p, err := strconv.ParseInt(desc[ix+1:], 10, 64)
		if err != nil {
			return nil, err
		}
		h.Port = int(p)
		desc = desc[:ix]
	}

	h.Host = desc

	return h, nil
}

func toHosts(r io.Reader, username string) ([]*Host, error) {
	var hosts []*Host

	s := bufio.NewScanner(r)
	for s.Scan() {
		h, err := toHost(s.Text(), username)
		if err != nil {
			return nil, err
		}

		hosts = append(hosts, h)
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}

func (l *Loader) Load(exp string) ([]*Host, error) {
	// TODO(knorton): Handle set expressions
	r, err := os.Open(filepath.Join(l.dir, exp))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	return toHosts(r, u.Username)
}

func New(path string) *Loader {
	return &Loader{dir: path}
}
