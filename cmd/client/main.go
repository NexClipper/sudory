package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/NexClipper/sudory/pkg/client/k8s"
	"github.com/NexClipper/sudory/pkg/client/poll"
	"github.com/NexClipper/sudory/pkg/client/service"
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

	// get k8s client
	// TODO: k8s client usage
	_, err := k8s.NewClient()
	if err != nil {
		log.Panicln(err)
	}

	serviceScheduler := service.NewScheduler()
	serviceScheduler.Start()

	poller := poll.NewPoller(*token, *server, serviceScheduler)

	// polling
	poller.Start()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	select {
	case <-quit:
		os.Exit(1)
	}
}
