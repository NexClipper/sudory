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
	authv1 "github.com/NexClipper/sudory/pkg/server/model/auth/v1"
	servicev2 "github.com/NexClipper/sudory/pkg/server/model/service/v2"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
)

const sudoryAuthTokenHeaderName = "x-sudory-client-token"

type SudoryAPI struct {
	client    *httpclient.HttpClient
	authToken atomic.Value
}

func NewSudoryAPI(address string) (*SudoryAPI, error) {
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
	// jwt_token, _, err := jwt.NewParser().ParseUnverified(s.authToken, claims)
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
	// return s.authToken
}

func (s *SudoryAPI) Auth(ctx context.Context, auth *authv1.HttpReqAuth) error {
	if auth == nil {
		return fmt.Errorf("auth is nil")
	}

	b, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	result := s.client.Post("/client/auth").SetBody("application/json", b).Do(ctx)

	// get session token
	if headers := result.Headers(); headers != nil {
		if token := headers.Get(sudoryAuthTokenHeaderName); token != "" {
			s.authToken.Store(token)
			// s.authToken = token
		}
	}

	if err := result.Error(); err != nil {
		return err
	}

	return nil
}

func (s *SudoryAPI) GetServices(ctx context.Context) ([]servicev2.HttpRsp_ClientServicePolling, error) {
	var services []servicev2.HttpRsp_ClientServicePolling

	token := s.GetToken()
	if token == "" {
		return nil, fmt.Errorf("session token is empty")
	}

	result := s.client.Get("/client/service").SetHeader(sudoryAuthTokenHeaderName, token).Do(ctx)

	// get session token
	if headers := result.Headers(); headers != nil {
		if token := headers.Get(sudoryAuthTokenHeaderName); token != "" {
			s.authToken.Store(token)
			// s.authToken = token
		}
	}

	if err := result.IntoJson(&services); err != nil {
		return nil, err
	}

	return services, nil
}

func (s *SudoryAPI) UpdateServices(ctx context.Context, service *servicev2.HttpReq_ClientServiceUpdate) error {
	if service == nil {
		return fmt.Errorf("service is nil")
	}

	b, err := json.Marshal(service)
	if err != nil {
		return err
	}

	token := s.GetToken()
	if token == "" {
		return fmt.Errorf("session token is empty")
	}

	result := s.client.Put("/client/service").SetHeader(sudoryAuthTokenHeaderName, token).SetBody("application/json", b).Do(ctx)

	// get session token
	if headers := result.Headers(); headers != nil {
		if token := headers.Get(sudoryAuthTokenHeaderName); token != "" {
			s.authToken.Store(token)
			// s.authToken = token
		}
	}

	if err := result.Error(); err != nil {
		return err
	}

	return nil
}
