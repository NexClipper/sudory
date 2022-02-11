package poll

import (
	"fmt"

	"github.com/NexClipper/sudory/pkg/client/log"
)

func (p *Poller) HandShake() error {
	b := []byte(fmt.Sprintf("assertion=%s&cluster_uuid=%s&client_uuid=%s", p.bearerToken, p.clusterId, p.machineID))

	_, err := p.client.PostForm("/client/auth", nil, b)
	if err != nil {
		return err
	}
	log.Debugf("Successed to handshake: received token(%s) for polling.", p.client.GetToken())

	return nil
}
