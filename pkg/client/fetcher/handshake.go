package fetcher

import (
	"encoding/json"

	"github.com/NexClipper/sudory/pkg/client/log"
	authv1 "github.com/NexClipper/sudory/pkg/server/model/auth/v1"
	"github.com/NexClipper/sudory/pkg/version"
)

func (f *Fetcher) HandShake() error {
	body := &authv1.HttpReqAuth{Auth: authv1.Auth{AuthProperty: authv1.AuthProperty{
		ClusterUuid:   f.clusterId,
		ClientUuid:    f.machineID,
		Assertion:     f.bearerToken,
		ClientVersion: version.Version,
	}}}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = f.client.PostJson("/client/auth", nil, b)
	if err != nil {
		return err
	}
	log.Debugf("Successed to handshake: received token(%s) for polling.", f.client.GetToken())

	return nil
}
