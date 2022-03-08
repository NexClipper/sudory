package poll

import (
	"encoding/json"

	"github.com/NexClipper/sudory/pkg/client/log"
	authv1 "github.com/NexClipper/sudory/pkg/server/model/auth/v1"
)

func (p *Poller) HandShake() error {
	body := &authv1.HttpReqAuth{Auth: authv1.Auth{AuthProperty: authv1.AuthProperty{
		ClusterUuid: p.clusterId,
		ClientUuid:  p.machineID,
		Assertion:   p.bearerToken,
	}}}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = p.client.PostJson("/client/auth", nil, b)
	if err != nil {
		return err
	}
	log.Debugf("Successed to handshake: received token(%s) for polling.", p.client.GetToken())

	return nil
}
