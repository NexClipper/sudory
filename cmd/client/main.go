package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/NexClipper/sudory/pkg/client/poll"
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

	poller := poll.NewPoller(*token, *server)

	// polling
	poller.Start()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	select {
	case <-quit:
		os.Exit(1)
	}
}
