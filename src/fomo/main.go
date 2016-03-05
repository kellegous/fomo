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

	h, err := hosts.New(*flagDir).Load(expr)
	if err != nil {
		log.Panic(err)
	}

	log.Println(h)

	if err := remote.Run(&hosts.Host{
		User: "knorton",
		Host: "localhost",
		Port: 22,
	}, os.Stdout, "task.py"); err != nil {
		log.Panic(err)
	}
}
