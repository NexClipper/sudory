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
		log.Fatalf("Client must have token('%s').\n", *token)
	}

	if len(*server) == 0 {
		log.Fatalf("Client must have server('%s').\n", *server)
	}

	// get k8s client
	// TODO: k8s client usage
	client, err := k8s.NewClient()
	if err != nil {
		log.Fatalf("Failed to create k8s.NewClient : %v.\n", err)
	}
	log.Printf("Created k8s clientset.\n")

	// check k8s's api-server status
	if err := client.RawRequest().CheckApiServerStatus(); err != nil {
		log.Fatalf("CheckApiServerStatus is failed : %v.\n", err)
	}
	log.Printf("Successed to check K8s's api-server status.\n")

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
