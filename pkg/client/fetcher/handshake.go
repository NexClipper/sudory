package fetcher

import (
	"encoding/json"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/NexClipper/sudory/pkg/client/log"
	authv1 "github.com/NexClipper/sudory/pkg/server/model/auth/v1"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
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

	f.ChangeClientConfigFromToken()

	// save session_uuid from token
	claims := new(sessionv1.ClientSessionPayload)
	jwt_token, _, err := jwt.NewParser().ParseUnverified(f.client.GetToken(), claims)
	if _, ok := jwt_token.Claims.(*sessionv1.ClientSessionPayload); !ok || err != nil {
		log.Warnf("Failed to bind payload : %v\n", err)
		return err
	}
	if err := writeFile(".sudory", []byte(claims.Uuid)); err != nil {
		return err
	}

	return nil
}

func (f *Fetcher) RetryHandshake() {
	maxRetryCnt := 5

	ticker := time.NewTicker(time.Second * time.Duration(f.pollingInterval))
	defer ticker.Stop()

	for retry := 0; ; <-ticker.C {
		log.Debugf("retry handshake : count(%d)\n", retry+1)
		if err := f.HandShake(); err != nil {
			log.Warnf("Failed to Handshake Retry : count(%d), error(%v)\n", retry, err)
		} else {
			return
		}
		retry++

		if maxRetryCnt <= retry {
			f.Cancel()
			return
		}
	}
}

func writeFile(filename string, data []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}
