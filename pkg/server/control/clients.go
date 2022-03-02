package control

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/database/query_parser"

	//lint:ignore ST1001 auto-generated
	. "github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/jwt"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
	clientv1 "github.com/NexClipper/sudory/pkg/server/model/client/v1"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
	tokenv1 "github.com/NexClipper/sudory/pkg/server/model/token/v1"
	"github.com/NexClipper/sudory/pkg/server/status/env"
	"xorm.io/xorm"

	"github.com/labstack/echo/v4"
)

type ClientPlayload struct {
	Exp          time.Time `json:"exp,omitempty"`           //expiration_time
	Iat          time.Time `json:"iat,omitempty"`           //issued_at_time
	Uuid         string    `json:"uuid,omitempty"`          //token_uuid
	ClusterUuid  string    `json:"cluster-uuid,omitempty"`  //cluster_uuid
	ClientUuid   string    `json:"client-uuid,omitempty"`   //client_uuid
	PollInterval int       `json:"poll-interval,omitempty"` //config_poll_interval
	Loglevel     string    `json:"Loglevel,omitempty"`      //config_loglevel
}

// Poll []Service (client)
// @Description Poll a Service
// @Accept      json
// @Produce     json
// @Tags        client/service
// @Router      /client/service [put]
// @Param       x-sudory-client-token header string                        true "client session token"
// @Param       service               body   []v1.HttpReqClientSideService true "HttpReqClientSideService"
// @Success     200 {array}  v1.HttpRspClientSideService
// @Header      200 {string} x-sudory-client-token "x-sudory-client-token"
func (c *Control) PollService() func(ctx echo.Context) error {

	unwarp := func(elems []servicev1.HttpReqClientSideService) []servicev1.ServiceAndSteps {
		out := make([]servicev1.ServiceAndSteps, len(elems))
		for n := range elems {
			out[n] = elems[n].ServiceAndSteps
		}
		return out
	}

	warp := func(elems []servicev1.ServiceAndSteps) []servicev1.HttpRspClientSideService {
		out := make([]servicev1.HttpRspClientSideService, len(elems))
		for n := range elems {
			out[n] = servicev1.HttpRspClientSideService{ServiceAndSteps: elems[n]}
		}
		return out
	}

	binder := func(ctx Contexter) error {
		body := new([]servicev1.HttpReqClientSideService)
		if err := ctx.Bind(body); err != nil {
			return ErrorBindRequestObject(err)
		}

		payload := getClientTokenPayload(ctx.Echo()) //read client token
		if payload == nil {
			return ErrorInvaliedRequestParameter()
		}
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		body, ok := ctx.Object().(*[]servicev1.HttpReqClientSideService)
		if !ok {
			return nil, ErrorFailedCast()
		}

		//update service
		err := foreach_client_service_and_steps(unwarp(*body), func(service servicev1.Service, steps []stepv1.ServiceStep) error {

			//service.Status 초기화
			//service.Status; 상태가 가장큰 step의 Status
			//service.StepPosition; 상태가 가장큰 step의 Sequence
			service.Status = newist.Int32(int32(servicev1.StatusRegist))
			service.StepPosition = newist.Int32(0)
			steps = map_step(steps, func(step stepv1.ServiceStep) stepv1.ServiceStep {
				//step.Status 상태가 service.Status 보다 크다면
				//서비스의 상태 정보를 해당 값으로 덮어쓰기
				if nullable.Int32(service.Status).Value() < nullable.Int32(step.Status).Value() {
					service.Status = newist.Int32(nullable.Int32(step.Status).Value())         //status
					service.StepPosition = newist.Int32(nullable.Int32(step.Sequence).Value()) //position
				}
				return step
			})

			//save step
			if err := foreach_step(steps, operator.NewServiceStep(ctx.Database()).Update); err != nil {
				return err
			}

			//save service
			if err := operator.NewService(ctx.Database()).Update(service); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		payload := getClientTokenPayload(ctx.Echo()) //read client token

		//find service
		where := "cluster_uuid = ? AND status BETWEEN ? AND ?"
		args := []interface{}{
			payload.ClusterUuid,
			servicev1.StatusRegist,
			servicev1.StatusProcessing,
		}

		services, err := operator.NewService(ctx.Database()).
			Find(where, args...)
		if err != nil {
			return nil, err
		}

		//make response
		push, pop := servicev1.HttpRspBuilder(len(services))
		err = foreach_service(services, func(service servicev1.Service) error {

			//find steps
			where := "service_uuid = ?"
			steps, err := operator.NewServiceStep(ctx.Database()).
				Find(where, service.Uuid)
			if err != nil {
				return err
			}

			push(service, steps) //push
			return nil
		})
		if err != nil {
			return nil, err
		}

		service_rsp := pop()

		service_rsp = map_client_service_and_steps(service_rsp, func(service servicev1.Service, steps []stepv1.ServiceStep) servicev1.ServiceAndSteps {
			//service.Status 초기화
			service.Status = newist.Int32(int32(servicev1.StatusRegist))
			service.StepPosition = newist.Int32(0)
			steps = map_step(steps, func(step stepv1.ServiceStep) stepv1.ServiceStep {
				//StatusSend 보다 작으면 응답 전 업데이트
				if nullable.Int32(step.Status).Value() < int32(servicev1.StatusSend) {
					step.Status = newist.Int32(int32(servicev1.StatusSend))
				}

				//step.Status 상태가 service.Status 보다 크다면
				//서비스의 상태 정보를 해당 값으로 덮어쓰기
				if nullable.Int32(service.Status).Value() < nullable.Int32(step.Status).Value() {
					service.Status = nullable.Int32(step.Status).Ptr()         //status
					service.StepPosition = nullable.Int32(step.Sequence).Ptr() //position
				}
				return step
			})

			//할당된 클라이언트 정보 추가
			service.AssignedClientUuid = newist.String(payload.ClientUuid)

			return servicev1.ServiceAndSteps{Service: service, Steps: steps}
		})

		err = foreach_client_service_and_steps(service_rsp, func(service servicev1.Service, steps []stepv1.ServiceStep) error {
			//save step
			if err := foreach_step(steps, operator.NewServiceStep(ctx.Database()).Update); err != nil {
				return err
			}
			//save service
			if err := operator.NewService(ctx.Database()).Update(service); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		return warp(service_rsp), nil //pop
	}

	return MakeMiddlewareFunc(Option{
		TokenVerifier: verifyClientSessionToken(c.db.Engine(), Nolock),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Lock(c.db.Engine()),
	})
}

// Auth Client
// @Description Auth a client
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Tags        client/auth
// @Router      /client/auth [post]
// @Param       assertion    formData string true "assertion=<bearer-token>"
// @Param       cluster_uuid formData string true "Cluster 의 Uuid"
// @Param       client_uuid  formData string true "Client 의 Uuid"
// @Success     200 {string} ok
// @Header      200 {string} x-sudory-client-token "x-sudory-client-token"
func (c *Control) AuthClient() func(ctx echo.Context) error {

	binder := func(ctx Contexter) error {

		if len(ctx.Forms()) == 0 {
			return ErrorInvaliedRequestParameter()
		}

		if len(ctx.Forms()[__ASSERTION__]) == 0 {
			return ErrorInvaliedRequestParameterName(__ASSERTION__)
		}
		if len(ctx.Forms()[__CLUSTER_UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__CLUSTER_UUID__)
		}
		if len(ctx.Forms()[__CLIENT_UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__CLIENT_UUID__)
		}
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {

		assertion := ctx.Forms()[__ASSERTION__]
		cluster_uuid := ctx.Forms()[__CLUSTER_UUID__]
		client_uuid := ctx.Forms()[__CLIENT_UUID__]

		//valid cluster
		_, err := operator.NewCluster(ctx.Database()).
			Get(cluster_uuid)
		if err != nil {
			return nil, err
		}

		//클라이언트를 조회 하여
		//레코드에 없으면 추가
		where := "cluster_uuid = ? AND uuid = ?"
		clients, err := operator.NewClient(ctx.Database()).
			Find(where, cluster_uuid, client_uuid)
		if err != nil {
			return nil, err
		}
		//없으면 추가
		if len(clients) == 0 {
			name := fmt.Sprintf("client:%s", client_uuid)
			summary := fmt.Sprintf("client: %s, cluster: %s", client_uuid, cluster_uuid)
			client := clientv1.Client{}
			client.Uuid = client_uuid
			client.LabelMeta = NewLabelMeta(newist.String(name), newist.String(summary))
			client.ClusterUuid = cluster_uuid

			if err = operator.NewClient(ctx.Database()).Create(client); err != nil {
				return nil, err
			}
		}

		//valid token
		m := map[string]string{
			"user_kind": token_user_kind_cluster,
			"user_uuid": cluster_uuid,
			"token":     assertion,
		}

		cond := query_parser.NewCondition(m, func(key string) (string, string, bool) {
			return "=", "%s", true
		})

		tokens, err := operator.NewToken(ctx.Database()).
			Find(cond.Where(), cond.Args()...)
		if err != nil {
			return nil, err
		}
		first := func() *tokenv1.Token {
			for _, it := range tokens {
				return &it
			}
			return nil
		}
		token := first()

		if token == nil {
			return nil, fmt.Errorf("record was not found: token")
		}

		//만료 시간 검증
		if time.Until(token.ExpirationTime) < 0 {
			return nil, fmt.Errorf("token was expierd")
		}

		//new session
		//make session payload
		token_uuid := NewUuidString()
		iat := time.Now()
		exp := env.ClientSessionExpirationTime(iat)

		payload := &ClientPlayload{
			Exp:          exp,
			Iat:          iat,
			Uuid:         token_uuid,
			ClusterUuid:  cluster_uuid,
			ClientUuid:   client_uuid,
			PollInterval: env.ClientConfigPollInterval(),
			Loglevel:     env.ClientConfigLoglevel(),
		}

		if false {
			json_mashal := func(v interface{}) []byte {
				json_mashal := func(v interface{}) ([]byte, error) { return json.MarshalIndent(v, "", " ") }
				// json_mashal := json.Marshal
				right := func(b []byte, err error) []byte {
					if err != nil {
						panic(err)
					}
					return b
				}
				return right(json_mashal(v))
			}
			println("payload=", string(json_mashal(payload)))
		}

		//new jwt
		session_token_value, err := jwt.New(payload, []byte(env.ClientSessionSignatureSecret()))
		if err != nil {
			return nil, err
		}

		session := newClientSession(*payload, session_token_value)

		//save session
		err = operator.NewSession(ctx.Database()).
			Create(session)
		if err != nil {
			return nil, err
		}

		//save token to header
		ctx.Echo().Response().Header().Add(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__, session_token_value)

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Lock(c.db.Engine()),
	})
}

func newClientSession(payload ClientPlayload, token string) sessionv1.Session {
	session := sessionv1.Session{}
	session.Uuid = payload.Uuid
	session.UserUuid = payload.ClientUuid
	session.UserKind = "client"
	session.Token = token
	session.IssuedAtTime = payload.Iat
	session.ExpirationTime = payload.Exp
	return session
}

func verifyClientSessionToken(engine *xorm.Engine, behave func(*xorm.Engine) func(Contexter, func(Contexter) (interface{}, error)) (interface{}, error)) func(ctx Contexter) error {

	operate := func(ctx Contexter) error {
		var err error
		token := ctx.Echo().Request().Header.Get(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__)

		if len(token) == 0 {
			return ErrorInvaliedRequestParameter()
		}

		//verify
		//jwt verify
		err = jwt.Verify(token, []byte(env.ClientSessionSignatureSecret()))
		if err != nil {
			return WithCode(err, "jwt verify", http.StatusForbidden)
		}

		payload := new(ClientPlayload)
		err = jwt.BindPayload(token, payload)
		if err != nil {
			return WithCode(err, "jwt bind payload", http.StatusForbidden)
		}

		//만료시간 비교
		if time.Until(payload.Exp) < 0 {
			return WithCode(err, "token expierd", http.StatusForbidden)
		}

		//reflesh payload
		payload.Exp = env.ClientSessionExpirationTime(time.Now())
		payload.PollInterval = env.ClientConfigPollInterval()
		payload.Loglevel = env.ClientConfigLoglevel()

		//new jwt-new_token
		new_token, err := jwt.New(payload, []byte(env.ClientSessionSignatureSecret()))
		if err != nil {
			return WithCode(err, "new jwt", http.StatusInternalServerError)
		}

		session := newClientSession(*payload, new_token)
		//udpate session
		err = operator.NewSession(ctx.Database()).
			Update(session)
		if err != nil {
			return WithCode(err, "update session-token", http.StatusInternalServerError)
		}

		//save client session-token to header
		ctx.Echo().Response().Header().Add(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__, new_token)

		return nil
	}

	return func(ctx Contexter) error {

		fackRight := func(fn func(ctx Contexter) error) func(ctx Contexter) (interface{}, error) {
			return func(ctx Contexter) (interface{}, error) {
				return nil, fn(ctx) //return fack and error
			}
		}

		left := func(v interface{}, err error) error {
			return err
		}

		err := left(behave(engine)(ctx, fackRight(operate)))
		return err
	}
}

func getClientTokenPayload(ctx echo.Context) *ClientPlayload {
	//get token
	token := ctx.Request().Header.Get(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__)
	if len(token) == 0 {
		return nil
	}

	payload := new(ClientPlayload)

	err := jwt.BindPayload(token, payload)
	if err != nil {
		return nil
	}
	return payload
}

// setCookie
//lint:ignore U1000 auto-generated
func setCookie(ctx echo.Context, key, value string, exp time.Duration) {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.Value = value
	cookie.Expires = time.Now().Add(exp)
	ctx.SetCookie(cookie)
}

// setCookie
//lint:ignore U1000 auto-generated
func getCookie(ctx echo.Context, key string) (string, error) {
	cookie, err := ctx.Cookie(key)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
