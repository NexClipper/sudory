package control

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/event"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/status/env"
	"github.com/pkg/errors"

	"github.com/NexClipper/sudory/pkg/server/macro/jwt"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
	authv1 "github.com/NexClipper/sudory/pkg/server/model/auth/v1"
	clientv1 "github.com/NexClipper/sudory/pkg/server/model/client/v1"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
	tokenv1 "github.com/NexClipper/sudory/pkg/server/model/token/v1"

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
func (ctl Control) PollService(ctx echo.Context) error {
	clientTokenPayload := func(ctx echo.Context) func() (*sessionv1.ClientSessionPlayload, error) {
		var (
			err  error
			inst *sessionv1.ClientSessionPlayload
			once sync.Once
		)
		return func() (*sessionv1.ClientSessionPlayload, error) {
			once.Do(func() {
				inst, err = getClientTokenPayload(ctx) //read client token
			})
			return inst, err
		}
	}(ctx)

	unwrap := func(elems []servicev1.HttpReqClientSideService) []servicev1.ServiceAndSteps {
		out := make([]servicev1.ServiceAndSteps, len(elems))
		for n := range elems {
			out[n] = elems[n].ServiceAndSteps
		}
		return out
	}
	wrap := func(elems []servicev1.ServiceAndSteps) []servicev1.HttpRspClientSideService {
		out := make([]servicev1.HttpRspClientSideService, len(elems))
		for n := range elems {
			out[n] = servicev1.HttpRspClientSideService{ServiceAndSteps: elems[n]}
		}
		return out
	}
	unwrap_dbschema := func(elems []servicev1.DbSchemaServiceAndSteps) []servicev1.ServiceAndSteps {
		rst := make([]servicev1.ServiceAndSteps, len(elems))
		for n := range elems {
			rst[n] = servicev1.ServiceAndSteps{Service: elems[n].Service, Steps: stepv1.TransFormDbSchema(elems[n].Steps)}
		}
		return rst
	}

	map_step := func(elems []stepv1.ServiceStep, mapper func(stepv1.ServiceStep) stepv1.ServiceStep) []stepv1.ServiceStep {
		rst := make([]stepv1.ServiceStep, len(elems))
		for n := range elems {
			rst[n] = mapper(elems[n])
		}
		return rst
	}
	map_service_and_step := func(elems []servicev1.ServiceAndSteps, mapper func(servicev1.ServiceAndSteps) servicev1.ServiceAndSteps) []servicev1.ServiceAndSteps {
		rst := make([]servicev1.ServiceAndSteps, len(elems))
		for n := range elems {
			rst[n] = mapper(elems[n])
		}
		return rst
	}

	foreach_service_and_steps := func(elems []servicev1.ServiceAndSteps, fn func(servicev1.ServiceAndSteps) error) error {
		for _, it := range elems {
			if err := fn(it); err != nil {
				return err
			}
		}
		return nil
	}

	make_status_and_position := func(serviceAndSteps servicev1.ServiceAndSteps) servicev1.ServiceAndSteps {
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
				service.Status = newist.Int32(nullable.Int32(step.Status).Value()) //status
				*service.StepPosition = *service.StepPosition + 1                  //position
			}
			return step
		})
		return servicev1.ServiceAndSteps{Service: service, Steps: steps}
	}
	make_response_status := func(serviceAndSteps servicev1.ServiceAndSteps) servicev1.ServiceAndSteps {
		service := serviceAndSteps.Service
		steps := serviceAndSteps.Steps

		//StatusSend 보다 작으면 응답 전 업데이트
		steps = map_step(steps, func(step stepv1.ServiceStep) stepv1.ServiceStep {
			if nullable.Int32(step.Status).Value() < int32(servicev1.StatusSend) {
				step.Status = newist.Int32(int32(servicev1.StatusSend))
			}

			return step
		})

		return servicev1.ServiceAndSteps{Service: service, Steps: steps}
	}
	make_response_assign_info := func(serviceAndSteps servicev1.ServiceAndSteps) servicev1.ServiceAndSteps {
		service := serviceAndSteps.Service
		steps := serviceAndSteps.Steps

		//할당된 클라이언트 정보 추가
		payload, _ := clientTokenPayload()
		service.AssignedClientUuid = newist.String(payload.ClientUuid)

		return servicev1.ServiceAndSteps{Service: service, Steps: steps}
	}

	update_service := func(serviceAndSteps servicev1.ServiceAndSteps) error {
		_, err := ctl.Scope(func(db database.Context) (interface{}, error) {
			//save service
			if _, err := vault.NewService(db).Update(serviceAndSteps.Service); err != nil {
				return nil, errors.Wrapf(err, "update service")
			}
			//save steps
			for _, step := range serviceAndSteps.Steps {
				if _, err := vault.NewServiceStep(db).Update(step); err != nil {
					return nil, errors.Wrapf(err, "update service step")
				}
			}
			return nil, nil
		})

		return err
	}

	body := new([]servicev1.HttpReqClientSideService)
	if err := ctx.Bind(body); err != nil {
		return HttpError(http.StatusBadRequest,
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	// r, err := ctl.Scope(func(db database.Context) (interface{}, error) {

	//update request
	request := unwrap(*body)
	request = map_service_and_step(request, make_status_and_position)

	//save request
	if err := foreach_service_and_steps(request, update_service); err != nil {
		return HttpError(http.StatusInternalServerError,
			errors.Wrapf(err, "update request"))
	}

	//invoke event (service-poll-in)
	for _, request := range request {
		m := map[string]interface{}{
			"event-name":           "service-poll-in",
			"service_uuid":         request.Uuid,
			"service-name":         nullable.String(request.Name).Value(),
			"cluster_uuid":         nullable.String(request.ClusterUuid).Value(),
			"assigned_client_uuid": nullable.String(request.AssignedClientUuid).Value(),
			"status":               nullable.Int32(request.Status).Value(),
			"result":               nullable.String(request.Result).Value(),
			"step_count":           nullable.Int32(request.StepCount).Value(),
			"step_position":        nullable.Int32(request.StepPosition).Value(),
		}
		event.Invoke("service-poll-in", m)
	}

	//make response
	payload, _ := clientTokenPayload()
	//find service
	// where := "cluster_uuid = ? AND (status BETWEEN ? AND ?)"
	// args := []interface{}{
	// 	payload.ClusterUuid,
	// 	servicev1.StatusRegist,
	// 	servicev1.StatusProcessing,
	// }
	// response, err := vault.NewService(ctx.Database()).
	// 	Find(where, args...)
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "NewService Find")
	// }
	response, err := gatherClusterService(ctl.NewSession(), payload.ClusterUuid)
	if err != nil {
		// return nil, errors.Wrapf(err, "make cluster service")
		return HttpError(http.StatusInternalServerError,
			errors.Wrapf(err, "make service%s",
				logs.KVL(
					"cluster_uuid", payload.ClusterUuid,
				)))
	}

	//update response
	response_ := unwrap_dbschema(response)
	response_ = map_service_and_step(response_, make_response_status)
	response_ = map_service_and_step(response_, make_response_assign_info)
	response_ = map_service_and_step(response_, make_status_and_position)

	//save response
	if err := foreach_service_and_steps(response_, update_service); err != nil {
		// return nil, errors.Wrapf(err, "update response")
		return HttpError(http.StatusInternalServerError,
			errors.Wrapf(err, "update response"))
	}

	//invoke event (service-poll-out)
	for _, response := range response_ {
		m := map[string]interface{}{
			"event-name":           "service-poll-out",
			"service_uuid":         response.Uuid,
			"service-name":         nullable.String(response.Name).Value(),
			"cluster_uuid":         nullable.String(response.ClusterUuid).Value(),
			"assigned_client_uuid": nullable.String(response.AssignedClientUuid).Value(),
			"status":               nullable.Int32(response.Status).Value(),
			"result":               nullable.String(response.Result).Value(),
			"step_count":           nullable.Int32(response.StepCount).Value(),
			"step_position":        nullable.Int32(response.StepPosition).Value(),
		}
		event.Invoke("service-poll-out", m)
	}

	// })
	// if err != nil {
	// 	return HttpError(http.StatusInternalServerError, err)
	// }

	return ctx.JSON(http.StatusOK, wrap(response_))
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
func (ctl Control) AuthClient(ctx echo.Context) error {
	body := new(authv1.HttpReqAuth)
	if err := ctx.Bind(body); err != nil {
		return HttpError(http.StatusBadRequest,
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	auth := body.Auth
	cluster_uuid := auth.ClusterUuid

	//valid cluster
	if _, err := vault.NewCluster(ctl.NewSession()).Get(cluster_uuid); err != nil {
		return HttpError(http.StatusInternalServerError,
			errors.Wrapf(err, "valid%s",
				logs.KVL(
					"cluster", cluster_uuid,
				)))
	}

	//valid token
	where := "user_kind = ? AND user_uuid = ? AND token = ?"
	tokens, err := vault.NewToken(ctl.NewSession()).
		Find(where, TokenUserKindCluster.String(), auth.ClusterUuid, auth.Assertion)
	if err != nil {
		return HttpError(http.StatusInternalServerError,
			errors.Wrapf(err, "NewToken Find"))
	}

	first := func() *tokenv1.DbSchema {
		for _, it := range tokens {
			return &it
		}
		return nil
	}
	token := first()

	if token == nil {
		return HttpError(http.StatusInternalServerError,
			fmt.Errorf("record was not found: token"))
	}

	//만료 시간 검증
	if time.Until(nullable.Time(token.ExpirationTime).Value()) < 0 {
		return HttpError(http.StatusInternalServerError,
			fmt.Errorf("token was expierd"))
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
		return HttpError(http.StatusInternalServerError,
			errors.Wrapf(err, "jwt New payload=%+v", payload))
	}

	_, err = ctl.Scope(func(db database.Context) (interface{}, error) {
		//클라이언트를 조회 하여
		//레코드에 없으면 추가
		clients, err := vault.NewClient(db).
			Find("cluster_uuid = ? AND uuid = ?", auth.ClusterUuid, auth.ClientUuid)
		if err != nil {
			return nil, HttpError(http.StatusInternalServerError,
				errors.Wrapf(err, "NewClient Find"))
		}

		//없으면 추가
		if len(clients) == 0 {
			name := fmt.Sprintf("client:%s", auth.ClientUuid)
			summary := fmt.Sprintf("client: %s, cluster: %s", auth.ClientUuid, auth.ClusterUuid)
			client := clientv1.Client{}
			client.Uuid = auth.ClientUuid
			client.LabelMeta = NewLabelMeta(newist.String(name), newist.String(summary))
			client.ClusterUuid = auth.ClusterUuid

			if _, err := vault.NewClient(ctl.NewSession()).Create(client); err != nil {
				return nil, errors.Wrapf(err, "NewClient Create")
			}
		}

		//save session
		session := newClientSession(*payload, new_token)
		if _, err := vault.NewSession(db).Create(session); err != nil {
			return nil, errors.Wrapf(err, "NewSession Create")
		}

		return nil, nil
	})
	if err != nil {
		return HttpError(http.StatusInternalServerError, err)
	}

	//save token to header
	ctx.Response().Header().Add(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__, new_token)

	//invoke event (client-auth-accept)
	m := map[string]interface{}{
		"event-name":   "client-auth-accept",
		"cluster_uuid": auth.ClusterUuid,
		"client_uuid":  auth.ClientUuid,
	}
	event.Invoke("client-auth-accept", m)

	return ctx.JSON(http.StatusOK, OK())
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

func (ctl Control) VerifyClientSessionToken(ctx echo.Context) error {
	token := ctx.Request().Header.Get(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__)
	if len(token) == 0 {
		return ErrorInvalidRequestParameter()
	}

	//verify
	//jwt verify
	if err := jwt.Verify(token, []byte(env.ClientSessionSignatureSecret())); err != nil {
		return errors.Wrapf(err, "%s verify", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__)
	}

	payload := new(sessionv1.ClientSessionPlayload)
	if err := jwt.BindPayload(token, payload); err != nil {
		return errors.Wrapf(err, "%s bind payload", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__)
	}

	//만료시간 비교
	if time.Until(payload.Exp) < 0 {
		return fmt.Errorf("%s expierd", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__)
	}

	_, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		//smart polling
		cluster, err := vault.NewCluster(db).Get(payload.ClusterUuid)
		if err != nil {
			return nil, errors.Wrapf(err, "get cluster%s",
				logs.KVL(
					"uuid", payload.ClusterUuid,
				))
		}
		service_count, err := countGatherClusterService(db, payload.ClusterUuid)
		if err != nil {
			return nil, errors.Wrapf(err, "count undone service%s",
				logs.KVL(
					"cluster_uuid", payload.ClusterUuid,
				))
		}

		//reflesh payload
		payload.Exp = env.ClientSessionExpirationTime(time.Now())
		// payload.PollInterval = env.ClientConfigPollInterval()
		payload.PollInterval = int(cluster.GetPollingOption().Interval(time.Duration(int64(env.ClientConfigPollInterval())*int64(time.Second)), int(service_count)) / time.Second)
		payload.Loglevel = env.ClientConfigLoglevel()

		//new jwt-new_token
		new_token, err := jwt.New(payload, []byte(env.ClientSessionSignatureSecret()))
		if err != nil {
			return nil, errors.Wrapf(err, "new jwt")
		}

		//udpate session
		session := newClientSession(*payload, new_token)
		if _, err := vault.NewSession(db).Update(session); err != nil {
			return nil, errors.Wrapf(err, "update session-token")
		}

		//save client session-token to header
		ctx.Response().Header().Add(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__, new_token)

		return nil, nil
	})
	if err != nil {
		return err
	}

	return nil
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

func gatherClusterService(db database.Context, cluster_uuid string) ([]servicev1.DbSchemaServiceAndSteps, error) {
	where := "cluster_uuid = ? AND (status BETWEEN ? AND ?)"
	args := []interface{}{
		cluster_uuid,
		servicev1.StatusRegist,
		servicev1.StatusProcessing,
	}
	response, err := vault.NewService(db).
		Find(where, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "NewService Find%s",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return response, nil
}

func countGatherClusterService(db database.Context, cluster_uuid string) (int64, error) {
	where := "cluster_uuid = ? AND (status BETWEEN ? AND ?)"
	args := []interface{}{
		cluster_uuid,
		servicev1.StatusRegist,
		servicev1.StatusProcessing,
	}

	count, err := db.Where(where, args...).Count(new(servicev1.DbSchema))
	if err != nil {
		return 0, errors.Wrapf(err, "service count%s",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return count, nil
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
