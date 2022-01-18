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
	server := flag.String("server", "http://localhost:8099", "sudory server url")
	clusterid := flag.String("clusterid", "", "sudory client's cluster id")

	flag.Parse()

	if len(*server) == 0 {
		log.Fatalf("Client must have server('%s').\n", *server)
	}

	if len(*clusterid) == 0 {
		log.Fatalf("Client must have server('%s').\n", *clusterid)
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

	poller := poll.NewPoller("", *server, serviceScheduler)

	// polling
	poller.Start()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	select {
	case <-quit:
		os.Exit(1)
	}
}
