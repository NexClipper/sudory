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
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
	tokenv1 "github.com/NexClipper/sudory/pkg/server/model/token/v1"
	"github.com/NexClipper/sudory/pkg/server/status/env"
	"xorm.io/xorm"

	"github.com/labstack/echo/v4"
)

type ClienPlayload struct {
	Exp          time.Time `json:"exp,omitempty"`           //expiration_time
	Iat          time.Time `json:"iat,omitempty"`           //issued_at_time
	Uuid         string    `json:"uuid,omitempty"`          //token_uuid
	ClusterUuid  string    `json:"cluster-uuid,omitempty"`  //cluster_uuid
	ClientUuid   string    `json:"client-uuid,omitempty"`   //client_uuid
	PollInterval int       `json:"poll-interval,omitempty"` //config_poll_interval
	Loglevel     string    `json:"Loglevel,omitempty"`      //config_loglevel
}

type SessionTokenError struct {
	HttpStatus int
	Err        error
}

func NewClientSessionTokenError(httpStatus int, err error) *SessionTokenError {
	return &SessionTokenError{HttpStatus: httpStatus, Err: err}
}

func (e SessionTokenError) Error() string {
	return fmt.Errorf("client session-token error: %w", e.Err).Error()
}

// Poll []Service (client)
// @Description Poll a Service
// @Accept json
// @Produce json
// @Tags client/service
// @Router /client/service [put]
// @Param x-sudory-client-token header string                        true "client session token"
// @Param service               body   []v1.HttpReqClientSideService true "HttpReqClientSideService"
// @Success 200 {array}  v1.HttpRspClientSideService
// @Header  200 {string} x-sudory-client-token "x-sudory-client-token"
func (c *Control) PollService() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for key := range ctx.QueryParams() {
			req[key] = ctx.QueryParam(key)
		}
		for _, key := range ctx.ParamNames() {
			req[key] = ctx.Param(key)
		}

		body := make([]servicev1.HttpReqClientSideService, 0) //bind body
		err := ctx.Bind(&body)
		if err != nil {
			return nil, ErrorBindRequestObject(err)
		}
		req[__BODY__] = body //save body

		payload := getClientTokenPayload(ctx) //read client token
		if payload == nil {
			return nil, ErrorInvaliedRequestParameter()
		}
		req[__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__] = payload //save client token

		return req, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}
		body, ok := req[__BODY__].([]servicev1.HttpReqClientSideService)
		if !ok {
			return nil, ErrorFailedCast()
		}

		//update service
		err := foreach_client_service(body, func(service servicev1.Service, steps []stepv1.ServiceStep) error {

			//update step
			err := foreach_step(steps, func(step stepv1.ServiceStep) error {

				//스탭의 상태가 StatusSend 보다 크다면
				//서비스의 상태 정보를 해당 값으로 덮어쓰기
				if int32(servicev1.StatusSend) < *step.Status {
					service.StepPosition = newist.Int32(*step.Sequence) //position
					service.Result = newist.String(*step.Result)        //result
				}
				//update step
				err := operator.NewServiceStep(ctx.Database).
					Update(step)
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}

			//update service
			err = operator.NewService(ctx.Database).
				Update(service)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		payload, ok := req[__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__].(*ClienPlayload)
		if !ok {
			return nil, ErrorFailedCast()
		}

		//find service
		where := "cluster_uuid = ? AND ? < status AND status < ?"
		cluster_uuid := payload.ClusterUuid
		if !ok {
			return nil, ErrorFailedCast()
		}
		status_regist := servicev1.StatusRegist
		status_success := servicev1.StatusSuccess

		services, err := operator.NewService(ctx.Database).
			Find(where, cluster_uuid, status_regist, status_success)
		if err != nil {
			return nil, err
		}
		//make response
		push, pop := servicev1.HttpRspBuilder(len(services))
		err = foreach_service(services, func(service servicev1.Service) error {
			service_uuid := service.Uuid
			where := "service_uuid = ?"
			//find steps
			steps, err := operator.NewServiceStep(ctx.Database).
				Find(where, service_uuid)
			if err != nil {
				return err
			}
			push(service, steps) //push
			return nil
		})
		if err != nil {
			return nil, err
		}

		return pop(), nil //pop
	}

	return MakeMiddlewareFunc_experimental(Option_experimental{
		TokenVerifier: verifyClientSessionToken(c.db.Engine()),
		Binder:        binder,
		Operator:      MakeBlockWithLock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// Auth Client
// @Description Auth a client
// @Accept x-www-form-urlencoded
// @Produce json
// @Tags client/auth
// @Router /client/auth [post]
// @Param assertion    formData string true "assertion=<bearer-token>"
// @Param cluster_uuid formData string true "Cluster 의 Uuid"
// @Param client_uuid  formData string true "Client 의 Uuid"
// @Success 200 {string} ok
// @Header  200 {string} x-sudory-client-token "x-sudory-client-token"
func (c *Control) AuthClient() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {

		req := make(map[string]string)
		formdatas, err := ctx.FormParams()
		if err != nil {
			return nil, err
		}
		for key := range formdatas {
			req[key] = ctx.FormValue(key)
		}
		// if len(req[__GRANT_TYPE__]) == 0 {
		// 	return nil, ErrorInvaliedRequestParameter()
		// }
		if len(req[__ASSERTION__]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		if len(req[__CLUSTER_UUID__]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		if len(req[__CLIENT_UUID__]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		return req, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		assertion := req[__ASSERTION__]
		cluster_uuid := req[__CLUSTER_UUID__]
		client_uuid := req[__CLIENT_UUID__]

		//valid cluster
		_, err := operator.NewCluster(ctx.Database).
			Get(cluster_uuid)
		if err != nil {
			return nil, err
		}

		//valid client
		_, err = operator.NewClient(ctx.Database).
			Get(client_uuid)
		if err != nil {
			return nil, err
		}

		//valid token
		m := map[string]interface{}{
			"user_kind": token_user_kind_cluster,
			"user_uuid": cluster_uuid,
			"token":     assertion,
		}
		query := query_parser.NewQueryParser(m, func(key string) (string, string) {
			return "=", "%s"
		})

		tokens, err := operator.NewToken(ctx.Database).
			Query(query)
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

		payload := &ClienPlayload{
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

		session := newSession(*payload, session_token_value)

		//save session
		err = operator.NewSession(ctx.Database).
			Create(session)
		if err != nil {
			return nil, err
		}

		//save token to header
		ctx.Http.Response().Header().Add(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__, session_token_value)

		return OK(), nil
	}

	return MakeMiddlewareFunc_experimental(Option_experimental{
		Binder:        binder,
		Operator:      MakeBlockWithLock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

func newSession(payload ClienPlayload, token string) sessionv1.Session {
	session := sessionv1.Session{}
	session.Uuid = payload.Uuid
	session.UserUuid = payload.ClientUuid
	session.UserKind = token_user_kind_cluster
	session.Token = token
	session.IssuedAtTime = payload.Iat
	session.ExpirationTime = payload.Exp
	return session
}

func verifyClientSessionToken(engine *xorm.Engine) func(ctx echo.Context) error {

	operator := func(ctx OperateContext) error {
		var err error
		token := ctx.Http.Request().Header.Get(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__)

		if len(token) == 0 {
			return ErrorInvaliedRequestParameter()
		}

		//verify
		//jwt verify
		err = jwt.Verify(token, []byte(env.ClientSessionSignatureSecret()))
		if err != nil {
			return NewClientSessionTokenError(http.StatusForbidden, fmt.Errorf("jwt verify: %w", err))
		}

		payload := new(ClienPlayload)
		err = jwt.BindPayload(token, payload)
		if err != nil {
			return NewClientSessionTokenError(http.StatusForbidden, fmt.Errorf("jwt bind payload: %w", err))
		}

		//만료시간 비교
		if time.Until(payload.Exp) < 0 {
			return NewClientSessionTokenError(http.StatusForbidden, fmt.Errorf("expierd"))
		}

		//reflesh payload
		payload.Exp = env.ClientSessionExpirationTime(time.Now())
		payload.PollInterval = env.ClientConfigPollInterval()
		payload.Loglevel = env.ClientConfigLoglevel()

		//new jwt-new_token
		new_token, err := jwt.New(payload, []byte(env.ClientSessionSignatureSecret()))
		if err != nil {
			return NewClientSessionTokenError(http.StatusInternalServerError, fmt.Errorf("create new jwt: %w", err))
		}

		session := newSession(*payload, new_token)
		//udpate session
		err = operator.NewSession(ctx.Database).
			Update(session)
		if err != nil {
			return NewClientSessionTokenError(http.StatusInternalServerError, fmt.Errorf("update: %w", err))
		}

		//save client session-token to header
		ctx.Http.Response().Header().Add(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__, new_token)

		return nil
	}

	return func(ctx echo.Context) error {

		block := MakeBlockNoLock

		err := left(block(engine, func(ctx OperateContext) (interface{}, error) {
			err := operator(ctx)
			return nil, err
		})(ctx, nil))
		return err
	}
}

func left(v interface{}, err error) error {
	return err
}

func getClientTokenPayload(ctx echo.Context) *ClienPlayload {
	//get token
	token := ctx.Request().Header.Get(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__)
	if len(token) == 0 {
		return nil
	}

	payload := new(ClienPlayload)

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
