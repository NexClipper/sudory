package fetcher

import (
	"encoding/json"
	"time"

	"github.com/NexClipper/sudory/pkg/client/log"
	authv1 "github.com/NexClipper/sudory/pkg/server/model/auth/v1"
	"github.com/NexClipper/sudory/pkg/version"
)

func (f *Fetcher) HandShake() error {
	body := &authv1.HttpReqAuth{AuthProperty: authv1.AuthProperty{
		ClusterUuid:   f.clusterId,
		Assertion:     f.bearerToken,
		ClientVersion: version.Version,
	}}

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

func (f *Fetcher) RetryHandshake() {
	maxRetryCnt := 5
	// retry := 0

	ticker := time.NewTicker(time.Second * time.Duration(f.pollingInterval))
	defer ticker.Stop()

	for retry := 0; ; <-ticker.C {
		log.Debugf("retry handshake : count(%d)\n", retry+1)
		if err := f.HandShake(); err != nil {
			log.Warnf("Failed to Handshake Retry : count(%d), error(%v)\n", retry, err)
		} else {
			f.ticker.Reset(time.Second * time.Duration(f.pollingInterval))
			return
		}
		retry++

		if maxRetryCnt <= retry {
			f.Cancel()
			return
		}
	}
}
