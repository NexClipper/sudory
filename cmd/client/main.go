package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/NexClipper/sudory/pkg/client/fetcher"
	"github.com/NexClipper/sudory/pkg/client/httpclient"
	"github.com/NexClipper/sudory/pkg/client/k8s"
	"github.com/NexClipper/sudory/pkg/client/log"
	"github.com/NexClipper/sudory/pkg/client/scheduler"
)

func main() {
	server := flag.String("server", "http://localhost:8099", "sudory server url")
	clusterid := flag.String("clusterid", "", "sudory client's cluster id")
	token := flag.String("token", "", "sudory client's token for server connection")
	loglevel := flag.String("loglevel", "debug", "sudory client's log level. One of: debug(defualt)|info|warn|error")

	flag.Parse()

	log.New(*loglevel)

	if len(*server) == 0 {
		log.Fatalf("Client must have server('%s').\n", *server)
	}

	if len(*clusterid) == 0 {
		log.Fatalf("Client must have clusterid('%s').\n", *clusterid)
	}

	if len(*token) == 0 {
		log.Fatalf("Client must have token('%s').\n", *token)
	}

	if err := httpclient.ValidateURL(*server); err != nil {
		log.Fatalf(err.Error())
	}

	// get k8s client
	// TODO: k8s client usage
	client, err := k8s.NewClient()
	if err != nil {
		log.Fatalf("Failed to create k8s.NewClient : %v.\n", err)
	}
	log.Debugf("Created k8s clientset.\n")

	// check k8s's api-server status
	if err := client.RawRequest().CheckApiServerStatus(); err != nil {
		log.Fatalf("CheckApiServerStatus is failed : %v.\n", err)
	}
	log.Debugf("Successed to check K8s's api-server status.\n")

	scheduler := scheduler.NewScheduler()
	scheduler.Start()

	poller, err := fetcher.NewFetcher(*token, *server, *clusterid, scheduler)
	if err != nil {
		log.Fatalf("Failed to create poller : %v.\n", err)
	}

	if err := poller.HandShake(); err != nil {
		log.Fatalf("Failed to handshake : %v.\n", err)
	}

	// polling
	poller.Polling()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	select {
	case <-quit:
		os.Exit(1)
	}
}
