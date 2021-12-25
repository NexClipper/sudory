package main

import "github.com/NexClipper/sudory-prototype-r1/pkg/route"

// @title SUDORY
// @version 0.0.1
// @description this is a sudory server.
// @contact.url https://nexclipper.io
// @contact.email jaehoon@nexclipper.io
func main() {
	r := route.New()

	r.Start()
}
