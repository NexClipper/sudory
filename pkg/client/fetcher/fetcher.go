package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/panta/machineid"

	"github.com/NexClipper/sudory/pkg/client/executor"
	"github.com/NexClipper/sudory/pkg/client/httpclient"
	"github.com/NexClipper/sudory/pkg/client/log"
	"github.com/NexClipper/sudory/pkg/client/service"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
)

const (
	defaultPollingInterval = 5 // * time.Second
	minPollingInterval     = 5
)

type Fetcher struct {
	bearerToken     string
	machineID       string
	clusterId       string
	client          *httpclient.HttpClient
	pollingInterval int
	done            chan struct{}
}

// func NewFetcher(bearerToken, server, clusterId string, sch *scheduler.Scheduler) (*Fetcher, error) {
func NewFetcher(bearerToken, server, clusterId string) (*Fetcher, error) {
	id, err := machineid.ID()
	if err != nil {
		return nil, err
	}
	id = strings.ReplaceAll(id, "-", "")

	return &Fetcher{
		bearerToken:     bearerToken,
		machineID:       id,
		clusterId:       clusterId,
		client:          httpclient.NewHttpClient(server, "", 1, 5000), // 1 retry, 5s wait
		pollingInterval: defaultPollingInterval,
		done:            make(chan struct{})}, nil
}

func (f *Fetcher) ChangePollingInterval(interval int) (int, error) {
	if interval < minPollingInterval {
		return 0, fmt.Errorf("interval(%d) you want to change is less than the minimum interval(%d)", interval, minPollingInterval)
	}

	if f.pollingInterval == interval {
		return 0, nil
	}
	f.pollingInterval = interval

	return interval, nil
}

func (f *Fetcher) Done() <-chan struct{} {
	return f.done
}

func (f *Fetcher) Cancel() {
	close(f.done)
}

func (f *Fetcher) Polling(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			f.poll()
		}
	}()

	return nil
}

func (f *Fetcher) poll() {
	// get services from server
	body, err := f.client.Get("/client/service", nil)
	if err != nil {
		log.Errorf(err.Error())

		// if session token is expired, retry handshake
		if f.client.IsTokenExpired() {
			f.RetryHandshake()
		}

		return
	}

	f.ChangeClientConfigFromToken()

	respData := []servicev1.HttpRspService_ClientSide{}
	if body != nil {
		if err := json.Unmarshal(body, &respData); err != nil {
			log.Errorf(err.Error())
			return
		}
	}
	log.Debugf("Recived %d service from server.", len(respData))

	if len(respData) == 0 {
		<-time.After(time.Duration(f.pollingInterval) * time.Second)
		return
	}

	// respData -> services
	recvServices := ConvertServiceListServerToClient(respData)

	updateChan := make(chan service.Service)
	wg := &sync.WaitGroup{}
	for _, serv := range recvServices {
		se := executor.NewServiceExecutor(*serv, updateChan)

		wg.Add(1)
		go func(e *executor.ServiceExecutor) {
			defer wg.Done()
			e.Execute()
		}(se)
	}

	done := make(chan struct{})
	go func() {
		wgg := &sync.WaitGroup{}

		for serv := range updateChan {
			sendData := ConvertServiceClientToServer(&serv)

			jsonb, err := json.Marshal(sendData)
			if err != nil {
				log.Errorf(err.Error())
				continue
			}

			wgg.Add(1)

			go func(sendData *ReqUpdateService, jsonb []byte) {
				defer wgg.Done()
				if _, err := f.client.PutJson("/client/service", nil, jsonb); err != nil {
					log.Errorf(err.Error())

					if f.client.IsTokenExpired() {
						f.RetryHandshake()
					}
				}
			}(sendData, jsonb)
		}
		wgg.Wait()
		done <- struct{}{}
	}()

	wg.Wait()
	close(updateChan)
	<-done
}

func (f *Fetcher) ChangeClientConfigFromToken() {
	claims := new(sessionv1.ClientSessionPayload)
	jwt_token, _, err := jwt.NewParser().ParseUnverified(f.client.GetToken(), claims)
	if _, ok := jwt_token.Claims.(*sessionv1.ClientSessionPayload); !ok || err != nil {
		log.Warnf("Failed to bind payload : %v\n", err)
		return
	}

	if interval, err := f.ChangePollingInterval(claims.PollInterval); err != nil {
		log.Warnf("Failed to change polling interval : %v\n", err)
	} else {
		if interval != 0 {
			log.Debugf("Change polling interval to %v\n", claims.PollInterval)
		}
	}

	if log.GetLogger().GetLevel() != strings.ToLower(claims.Loglevel) {
		if err := log.GetLogger().SetLevel(claims.Loglevel); err == nil {
			log.Debugf("Changed logger level to %s\n", claims.Loglevel)
		}
	}
}
