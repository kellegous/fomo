package main

import (
	"flag"
	"log"

	"fomo/auth"

	"golang.org/x/crypto/ssh"
)

func main() {
	flagUser := flag.String("user", "pi", "")
	flag.Parse()

	agt, err := auth.Agent()
	if err != nil {
		log.Panic(err)
	}

	cfg := ssh.ClientConfig{
		User: *flagUser,
		Auth: []ssh.AuthMethod{agt},
	}

	c, err := ssh.Dial("tcp", "pz:22", &cfg)
	if err != nil {
		log.Panic(err)
	}
	defer c.Close()

	log.Println(c.ClientVersion())
}
