package main

import (
	"flag"

	"github.com/NexClipper/sudory-prototype-r1/pkg/config"
	"github.com/NexClipper/sudory-prototype-r1/pkg/database"
	"github.com/NexClipper/sudory-prototype-r1/pkg/route"
)

// @title SUDORY
// @version 0.0.1
// @description this is a sudory server.
// @contact.url https://nexclipper.io
// @contact.email jaehoon@nexclipper.io
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
