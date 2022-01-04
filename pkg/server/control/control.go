package control

import "github.com/NexClipper/sudory-prototype-r1/pkg/server/database"

type Control struct {
	db *database.DBManipulator
}

func New(d *database.DBManipulator) *Control {
	return &Control{db: d}
}
