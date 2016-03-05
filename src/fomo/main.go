package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"fomo/hosts"
	"fomo/local"
	"fomo/remote"
)

func defaultDir() string {
	return fmt.Sprintf("%s/.fomo", os.Getenv("HOME"))
}

func runAll(hs []*hosts.Host, loc *local.Conn, src string) error {
	errs := make(chan error, len(hs))
	for _, h := range hs {
		go func(h *hosts.Host) {
			errs <- remote.Run(h, loc, src)
		}(h)
	}

	for i, n := 0, len(hs); i < n; i++ {
		if err := <-errs; err != nil {
			return err
		}
	}

	return nil
}

// fomo script.py all - dbs
func main() {
	flagDir := flag.String("prefs", defaultDir(), "")
	flag.Parse()

	src, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		log.Panic(err)
	}

	loc, err := local.Start(src)
	if err != nil {
		log.Panic(err)
	}
	defer loc.Close()

	expr := strings.Join(flag.Args()[1:], " ")

	hs, err := hosts.New(*flagDir).Load(expr)
	if err != nil {
		log.Panic(err)
	}

	if err := runAll(hs, loc, src); err != nil {
		log.Panic(err)
	}
}
