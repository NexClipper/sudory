package poll

import (
	"github.com/panta/machineid"
)

type Poll struct {
	token     string
	server    string
	machineID string
}

func New(token, server string) *Poll {
	id, err := machineid.ID()
	if err != nil {
		return nil
	}

	//log.Printf("machine id: %s", id)

	return &Poll{token: token, server: server, machineID: id}
}

func (p *Poll) Regist() error {
	return nil
}

func (p *Poll) Start() error {
	return nil
}
