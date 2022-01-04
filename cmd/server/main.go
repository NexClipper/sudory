package main

import (
	"flag"

	"github.com/NexClipper/sudory-prototype-r1/pkg/server/config"
	"github.com/NexClipper/sudory-prototype-r1/pkg/server/database"
	"github.com/NexClipper/sudory-prototype-r1/pkg/server/route"
)

func main() {
	configPath := flag.String("config", "../../conf/sudory-server.yml", "Path to sudory-server's config file")
	flag.Parse()

	cfg, err := config.New(*configPath)
	if err != nil {
		panic(err)
	}

	db, err := database.New(cfg)
	if err != nil {
		panic(err)
	}

	r := route.New(cfg, db)

	r.Start(cfg.Host.Port)
}
