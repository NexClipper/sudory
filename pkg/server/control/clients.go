package control

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/status/env"
	"github.com/pkg/errors"

	"github.com/NexClipper/sudory/pkg/server/macro/jwt"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
	authv1 "github.com/NexClipper/sudory/pkg/server/model/auth/v1"
	clientv1 "github.com/NexClipper/sudory/pkg/server/model/client/v1"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
	tokenv1 "github.com/NexClipper/sudory/pkg/server/model/token/v1"
	"xorm.io/xorm"

	"github.com/labstack/echo/v4"
)

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
	clientTokenPayload := func() func(ctx echo.Context) (*sessionv1.ClientSessionPlayload, error) {
		var (
			err  error
			inst *sessionv1.ClientSessionPlayload
			once sync.Once
		)
		return func(ctx echo.Context) (*sessionv1.ClientSessionPlayload, error) {
			once.Do(func() {
				inst, err = getClientTokenPayload(ctx) //read client token
			})
			return inst, err
		}
	}()

	binder := func(ctx Contexter) error {
		body := new([]servicev1.HttpReqClientSideService)
		if err := ctx.Bind(body); err != nil {
			return ErrorBindRequestObject(err)
		}

		//read client token
		if _, err := clientTokenPayload(ctx.Echo()); err != nil {
			return errors.Wrapf(err, "clientTokenPayload")
		}
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
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

		make_request := func(serviceAndSteps servicev1.ServiceAndSteps) servicev1.ServiceAndSteps {
			map_step := func(elems []stepv1.ServiceStep, mapper func(stepv1.ServiceStep) stepv1.ServiceStep) []stepv1.ServiceStep {
				rst := make([]stepv1.ServiceStep, len(elems))
				for n := range elems {
					rst[n] = mapper(elems[n])
				}
				return rst
			}

			service := serviceAndSteps.Service
			steps := serviceAndSteps.Steps

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
			return servicev1.ServiceAndSteps{Service: service, Steps: steps}
		}
		update_service := func(serviceAndSteps servicev1.ServiceAndSteps) (*servicev1.DbSchemaServiceAndSteps, error) {
			//save service
			service, err := vault.NewService(ctx.Database()).Update(serviceAndSteps.Service)
			if err != nil {
				return nil, errors.Wrapf(err, "NewService Update")
			}
			//save steps
			steps := make([]stepv1.DbSchema, len(serviceAndSteps.Steps))
			for i := range serviceAndSteps.Steps {
				step, err := vault.NewServiceStep(ctx.Database()).Update(serviceAndSteps.Steps[i])
				if err != nil {
					return nil, errors.Wrapf(err, "NewServiceStep Update")
				}
				steps[i] = *step
			}
			return &servicev1.DbSchemaServiceAndSteps{DbSchema: *service, Steps: steps}, nil
		}
		map_service_req := func(elems []servicev1.ServiceAndSteps, mapper func(servicev1.ServiceAndSteps) servicev1.ServiceAndSteps) []servicev1.ServiceAndSteps {
			rst := make([]servicev1.ServiceAndSteps, len(elems))
			for n := range elems {
				rst[n] = mapper(elems[n])
			}
			return rst
		}
		map_service_rsp := func(elems []servicev1.DbSchemaServiceAndSteps, mapper func(servicev1.ServiceAndSteps) servicev1.ServiceAndSteps) []servicev1.ServiceAndSteps {
			rst := make([]servicev1.ServiceAndSteps, len(elems))
			for n := range elems {
				rst[n] = mapper(servicev1.ServiceAndSteps{Service: elems[n].Service, Steps: stepv1.TransFormDbSchema(elems[n].Steps)})
			}
			return rst
		}
		// foreach_client_service_and_steps := func(elems []servicev1.ServiceAndSteps, fn func(servicev1.ServiceAndSteps) error) error {
		// 	for _, it := range elems {
		// 		if err := fn(it); err != nil {
		// 			return err
		// 		}
		// 	}
		// 	return nil
		// }

		make_response := func(serviceAndSteps servicev1.ServiceAndSteps) servicev1.ServiceAndSteps {
			map_step := func(elems []stepv1.ServiceStep, mapper func(stepv1.ServiceStep) stepv1.ServiceStep) []stepv1.ServiceStep {
				rst := make([]stepv1.ServiceStep, len(elems))
				for n := range elems {
					rst[n] = mapper(elems[n])
				}
				return rst
			}

			service := serviceAndSteps.Service
			steps := serviceAndSteps.Steps

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
			payload, _ := clientTokenPayload(ctx.Echo())
			service.AssignedClientUuid = newist.String(payload.ClientUuid)

			return servicev1.ServiceAndSteps{Service: service, Steps: steps}
		}

		body, ok := ctx.Object().(*[]servicev1.HttpReqClientSideService)
		if !ok {
			return nil, ErrorFailedCast()
		}

		//update request
		request := map_service_req(unwarp(*body), make_request)
		for i := range request {
			_, err := update_service(request[i])
			if err != nil {
				return nil, errors.Wrapf(err, "update request")
			}
		}

		//make response
		payload, _ := clientTokenPayload(ctx.Echo())
		//find service
		response, err := vault.NewService(ctx.Database()).
			Find("cluster_uuid = ? AND status BETWEEN ? AND ?", payload.ClusterUuid, servicev1.StatusRegist, servicev1.StatusProcessing)
		if err != nil {
			return nil, errors.Wrapf(err, "NewService Find")
		}
		//update response
		response_ := map_service_rsp(response, make_response)
		for i := range response_ {
			service := response_[i].Service
			steps := response_[i].Steps
			record, err := update_service(servicev1.ServiceAndSteps{Service: service, Steps: steps})
			if err != nil {
				return nil, errors.Wrapf(err, "update response")
			}
			response_[i].Service = record.Service
		}

		return warp(response_), nil
	}

	return MakeMiddlewareFunc(Option{
		TokenVerifier: func(ctx Contexter) error {
			if err := verifyClientSessionToken(c.db.Engine(), Nolock)(ctx); err != nil {
				return errors.Wrapf(err, "PollService TokenVerifier")
			}
			return nil
		},
		Binder: func(ctx Contexter) error {
			if err := binder(ctx); err != nil {
				return errors.Wrapf(err, "PollService binder")
			}
			return nil
		},
		Operator: func(ctx Contexter) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "PollService operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Lock(c.db.Engine()),
	})
}

// Auth Client
// @Description Auth a client
// @Accept      json
// @Produce     json
// @Tags        client/auth
// @Router      /client/auth [post]
// @Param       auth body v1.HttpReqAuth true "HttpReqAuth"
// @Success     200 {string} ok
// @Header      200 {string} x-sudory-client-token "x-sudory-client-token"
func (c *Control) AuthClient() func(ctx echo.Context) error {
	binder := func(ctx Contexter) error {
		body := new(authv1.HttpReqAuth)
		if err := ctx.Bind(body); err != nil {
			return ErrorBindRequestObject(err)
		}
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		body, ok := ctx.Object().(*authv1.HttpReqAuth)
		if !ok {
			return nil, ErrorFailedCast()
		}

		auth := body.Auth

		//valid cluster
		if _, err := vault.NewCluster(ctx.Database()).Get(auth.ClusterUuid); err != nil {
			return nil, errors.Wrapf(err, "NewCluster Get")
		}

		//클라이언트를 조회 하여
		//레코드에 없으면 추가
		clients, err := vault.NewClient(ctx.Database()).
			Find("cluster_uuid = ? AND uuid = ?", auth.ClusterUuid, auth.ClientUuid)
		if err != nil {
			return nil, errors.Wrapf(err, "NewClient Find")
		}
		//없으면 추가
		if len(clients) == 0 {
			name := fmt.Sprintf("client:%s", auth.ClientUuid)
			summary := fmt.Sprintf("client: %s, cluster: %s", auth.ClientUuid, auth.ClusterUuid)
			client := clientv1.Client{}
			client.Uuid = auth.ClientUuid
			client.LabelMeta = NewLabelMeta(newist.String(name), newist.String(summary))
			client.ClusterUuid = auth.ClusterUuid

			if _, err := vault.NewClient(ctx.Database()).Create(client); err != nil {
				return nil, errors.Wrapf(err, "NewClient Create")
			}
		}

		//valid token
		tokens, err := vault.NewToken(ctx.Database()).
			Find("user_kind = ? AND user_uuid = ? AND token = ?", token_user_kind_cluster, auth.ClusterUuid, auth.Assertion)
		if err != nil {
			return nil, errors.Wrapf(err, "NewToken Find")
		}

		first := func() *tokenv1.DbSchema {
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
		token_uuid := macro.NewUuidString()
		iat := time.Now()
		exp := env.ClientSessionExpirationTime(iat)

		payload := &sessionv1.ClientSessionPlayload{
			Exp:          exp,
			Iat:          iat,
			Uuid:         token_uuid,
			ClusterUuid:  auth.ClusterUuid,
			ClientUuid:   auth.ClientUuid,
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
		new_token, err := jwt.New(payload, []byte(env.ClientSessionSignatureSecret()))
		if err != nil {
			return nil, errors.Wrapf(err, "jwt New payload=%+v", payload)
		}

		//save session
		session := newClientSession(*payload, new_token)
		if _, err := vault.NewSession(ctx.Database()).Create(session); err != nil {
			return nil, errors.Wrapf(err, "NewSession Create")
		}
		//save token to header
		ctx.Echo().Response().Header().Add(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__, new_token)

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Contexter) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "AuthClient binder")
			}
			return nil
		},
		Operator: func(ctx Contexter) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "AuthClient operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Lock(c.db.Engine()),
	})
}

func newClientSession(payload sessionv1.ClientSessionPlayload, token string) sessionv1.Session {
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
		token := ctx.Echo().Request().Header.Get(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__)
		if len(token) == 0 {
			return ErrorInvaliedRequestParameter()
		}

		//verify
		//jwt verify
		if err := jwt.Verify(token, []byte(env.ClientSessionSignatureSecret())); err != nil {
			return errors.Wrapf(err, "jwt verify")
		}

		payload := new(sessionv1.ClientSessionPlayload)
		if err := jwt.BindPayload(token, payload); err != nil {
			return errors.Wrapf(err, "jwt bind payload")
		}

		//만료시간 비교
		if time.Until(payload.Exp) < 0 {
			return fmt.Errorf("token expierd")
		}

		//reflesh payload
		payload.Exp = env.ClientSessionExpirationTime(time.Now())
		payload.PollInterval = env.ClientConfigPollInterval()
		payload.Loglevel = env.ClientConfigLoglevel()

		//new jwt-new_token
		new_token, err := jwt.New(payload, []byte(env.ClientSessionSignatureSecret()))
		if err != nil {
			return errors.Wrapf(err, "new jwt")
		}

		//udpate session
		session := newClientSession(*payload, new_token)
		if _, err := vault.NewSession(ctx.Database()).Update(session); err != nil {
			return errors.Wrapf(err, "update session-token")
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

		return left(behave(engine)(ctx, fackRight(operate)))
	}
}

func getClientTokenPayload(ctx echo.Context) (*sessionv1.ClientSessionPlayload, error) {
	//get token
	token := ctx.Request().Header.Get(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__)
	if len(token) == 0 {
		return nil, fmt.Errorf("Echo Request Header Get key=%s", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__)
	}

	payload := new(sessionv1.ClientSessionPlayload)
	if err := jwt.BindPayload(token, payload); err != nil {
		return nil, errors.Wrapf(err, "jwt BindPayload token=%s", token)
	}
	return payload, nil
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
