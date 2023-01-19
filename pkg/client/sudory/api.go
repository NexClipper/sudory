package sudory

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/NexClipper/sudory/pkg/client/httpclient"
	"github.com/NexClipper/sudory/pkg/client/log"
	"github.com/NexClipper/sudory/pkg/server/model/auths/v2"
	servicev3 "github.com/NexClipper/sudory/pkg/server/model/service/v3"
	servicev4 "github.com/NexClipper/sudory/pkg/server/model/service/v4"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
)

const sudoryAuthTokenHeaderName = "x-sudory-client-token"

type SudoryAPI struct {
	client    *httpclient.HttpClient
	authToken atomic.Value
}

func NewSudoryAPI(address string) (*SudoryAPI, error) {
	log.Debugf("address in NewSudoryAPI : %s\n", address)
	client, err := httpclient.NewHttpClient(address, false, 0, 0)
	if err != nil {
		return nil, err
	}

	return &SudoryAPI{client: client}, nil
}

func NewSudoryAPIWithClient(client *httpclient.HttpClient) *SudoryAPI {
	return &SudoryAPI{client: client}
}

func (s *SudoryAPI) IsTokenExpired() bool {
	claims := new(sessionv1.ClientSessionPayload)
	jwt_token, _, err := jwt.NewParser().ParseUnverified(s.GetToken(), claims)
	if _, ok := jwt_token.Claims.(*sessionv1.ClientSessionPayload); !ok || err != nil {
		log.Warnf("jwt.ParseUnverified error : %v\n", err)
		return true
	}

	return !claims.VerifyExpiresAt(time.Now().Unix(), true)
}

func (s *SudoryAPI) GetToken() string {
	x := s.authToken.Load()

	return x.(string)
}

func (s *SudoryAPI) Auth(ctx context.Context, auth *auths.HttpReqAuth) error {
	if auth == nil {
		return fmt.Errorf("auth is nil")
	}

	b, err := json.Marshal(auth)
	if err != nil {
		return err
	}
	log.Debugf("sudory API in Auth : %s\n", s.client)
	result := s.client.Post("/client/auth").SetBody("application/json", b).Do(ctx)

	// get session token
	if headers := result.Headers(); headers != nil {
		if token := headers.Get(sudoryAuthTokenHeaderName); token != "" {
			s.authToken.Store(token)
		}
	}

	if err := result.Error(); err != nil {
		return err
	}

	return nil
}

func (s *SudoryAPI) GetServices(ctx context.Context) ([]servicev4.HttpRsp_ClientServicePolling, error) {
	var services HttpPollingDataset

	token := s.GetToken()
	if token == "" {
		return nil, fmt.Errorf("session token is empty")
	}

	result := s.client.Get("/client/service").
		SetHeader(sudoryAuthTokenHeaderName, token).
		Do(ctx)

	// get session token
	if headers := result.Headers(); headers != nil {
		if token := headers.Get(sudoryAuthTokenHeaderName); token != "" {
			s.authToken.Store(token)
		}
	}

	if err := result.IntoJson(&services); err != nil {
		return nil, err
	}

	return services, nil
}

func (s *SudoryAPI) UpdateServices(ctx context.Context, service *servicev4.HttpReq_ClientServiceUpdate) error {
	if service == nil {
		return fmt.Errorf("service is nil")
	}

	switch service.Version {
	case "v3":
		log.Debugf("request update_service: service{version:%s, uuid:%s, status:%d, result_len:%d}\n", service.Version, service.V3.Uuid, service.V3.Status, len(service.V3.Result))
	case "v4":
		log.Debugf("request update_service: service{version:%s, uuid:%s, status:%d, result_len:%d}\n", service.Version, service.V4.Uuid, service.V4.Status, len(service.V4.Result))
	default:
		return fmt.Errorf("unknown service version: %s", service.Version)
	}

	b, err := json.Marshal(service)
	if err != nil {
		return err
	}

	token := s.GetToken()
	if token == "" {
		return fmt.Errorf("session token is empty")
	}

	result := s.client.Put("/client/service").
		SetHeader(sudoryAuthTokenHeaderName, token).
		SetGzip(true).
		SetBody("application/json", b).
		Do(ctx)

	// get session token
	if headers := result.Headers(); headers != nil {
		if token := headers.Get(sudoryAuthTokenHeaderName); token != "" {
			s.authToken.Store(token)
		}
	}

	if err := result.Error(); err != nil {
		return err
	}

	return nil
}

type HttpPollingDataset []servicev4.HttpRsp_ClientServicePolling

func (s *HttpPollingDataset) UnmarshalJSON(b []byte) error {
	var l []json.RawMessage

	if err := json.Unmarshal(b, &l); err != nil {
		return err
	}

	for _, e := range l {
		data := servicev4.HttpRsp_ClientServicePolling{}
		if err := json.Unmarshal(e, &data); err != nil {
			// older version
			datav3 := servicev3.HttpRsp_ClientServicePolling{}
			if err := json.Unmarshal(e, &datav3); err != nil {
				return err
			}
			data.Version = "v3"
			data.V3 = datav3
		}

		*s = append(*s, data)
	}

	return nil
}
