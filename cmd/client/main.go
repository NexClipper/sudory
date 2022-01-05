package main

import (
	"flag"
	"os"

	"github.com/NexClipper/sudory-prototype-r1/pkg/client/poll"
)

func main() {
	token := flag.String("token", "token", "sudory token")
	server := flag.String("server", "http://localhost:8099", "sudory server url")

	flag.Parse()

	if len(*token) == 0 {
		os.Exit(1)
	}

	if len(*server) == 0 {
		os.Exit(1)
	}

	poller := poll.New(*token, *server)

	// regist

	if err := poller.Regist(); err != nil {
		os.Exit(1)
	}

	// polling
	if err := poller.Start(); err != nil {
		os.Exit(1)
	}
}
