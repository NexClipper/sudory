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

	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/jwt"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
	authv1 "github.com/NexClipper/sudory/pkg/server/model/auth/v1"
	clientv1 "github.com/NexClipper/sudory/pkg/server/model/client/v1"
	clusterv1 "github.com/NexClipper/sudory/pkg/server/model/cluster/v1"
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
// @Param       service               body   []v1.HttpReqService_ClientSide true "HttpReqService_ClientSide"
// @Success     200 {array}  v1.HttpRspService_ClientSide
// @Header      200 {string} x-sudory-client-token "x-sudory-client-token"
func (ctl Control) PollService(ctx echo.Context) error {
	clientTokenPayload_once := func(ctx echo.Context) func() (*sessionv1.ClientSessionPayload, error) {
		var (
			err  error
			inst *sessionv1.ClientSessionPayload
			once sync.Once
		)
		return func() (*sessionv1.ClientSessionPayload, error) {
			once.Do(func() {
				inst, err = clientTokenPayload(ctx) //read client token
			})
			return inst, err
		}
	}(ctx)

	mapper_response_step_status := func(step stepv1.ServiceStep) stepv1.ServiceStep {
		//StatusSend 보다 작으면 응답 전 업데이트
		if nullable.Int32(step.Status).Value() < int32(servicev1.StatusSend) {
			step.Status = newist.Int32(int32(servicev1.StatusSend))
		}

		return step
	}
	mapper_response_assign_client_info := func(service servicev1.Service) servicev1.Service {
		//할당된 클라이언트 정보 추가
		payload, _ := clientTokenPayload_once()
		service.AssignedClientUuid = payload.ClientUuid

		return service
	}

	body := []servicev1.HttpReqService_ClientSide{}
	if err := echoutil.Bind(ctx, &body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	for i := range body {
		body[i].ServiceProperty = body[i].Service.ChaniningStep(body[i].Steps)
	}

	if _, err := ctl.Scope(func(ctx database.Context) (interface{}, error) {
		for _, iter := range body {
			service := iter.Service
			steps := iter.Steps

			if _, err := vault.NewService(ctx).Update(service); err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
					errors.Wrapf(err, "update request service"))
			}

			for _, step := range steps {
				if _, err := vault.NewServiceStep(ctx).Update(step); err != nil {
					return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
						errors.Wrapf(err, "update request service step"))
				}
			}
		}

		return nil, nil
	}); err != nil {
		return err
	}

	//invoke event (service-poll-in)
	for _, request := range body {
		m := map[string]interface{}{
			"event-name":           "service-poll-in",
			"service_uuid":         request.Uuid,
			"service-name":         request.Name,
			"cluster_uuid":         request.ClusterUuid,
			"assigned_client_uuid": request.AssignedClientUuid,
			"status":               nullable.Int32(request.Status).Value(),
			"result":               nullable.String(request.Result).Value(),
			"step_count":           nullable.Int32(request.StepCount).Value(),
			"step_position":        nullable.Int32(request.StepPosition).Value(),
		}
		event.Invoke("service-poll-in", m)
	}

	//make response
	cluster_service := make([]servicev1.HttpRspService, 0)
	response_callback := func(s servicev1.Service, ss []stepv1.ServiceStep) {
		//service; chaining step
		s.ServiceProperty = s.ChaniningStep(ss)
		//service; assign client info
		s = mapper_response_assign_client_info(s)
		//setp; update status
		for i := range ss {
			ss[i] = mapper_response_step_status(ss[i])
		}

		cluster_service = append(cluster_service, servicev1.HttpRspService{Service: s, Steps: ss})
	}
	//get token payload
	payload, err := clientTokenPayload_once()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "client token payload"))
	}
	if err := gatherClusterService(ctl.NewSession(), payload.ClusterUuid, response_callback); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "gather cluster service%s",
				logs.KVL(
					"cluster_uuid", payload.ClusterUuid,
				)))
	}

	//save response
	if _, err := ctl.Scope(func(ctx database.Context) (interface{}, error) {
		for i := range cluster_service {
			service := cluster_service[i].Service
			steps := cluster_service[i].Steps

			service_, err := vault.NewService(ctx).Update(service)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
					errors.Wrapf(err, "update response service"))
			}

			for i := range steps {
				step := steps[i]
				step_, err := vault.NewServiceStep(ctx).Update(step)
				if err != nil {
					return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
						errors.Wrapf(err, "update response service"))
				}
				steps[i] = *step_
			}

			cluster_service[i].Service = *service_
			cluster_service[i].Steps = steps
		}

		return nil, nil
	}); err != nil {
		return err
	}

	//invoke event (service-poll-out)
	for _, response := range cluster_service {
		m := map[string]interface{}{
			"event-name":           "service-poll-out",
			"service_uuid":         response.Uuid,
			"service-name":         response.Name,
			"cluster_uuid":         response.ClusterUuid,
			"assigned_client_uuid": response.AssignedClientUuid,
			"status":               nullable.Int32(response.Status).Value(),
			"result":               nullable.String(response.Result).Value(),
			"step_count":           nullable.Int32(response.StepCount).Value(),
			"step_position":        nullable.Int32(response.StepPosition).Value(),
		}
		event.Invoke("service-poll-out", m)
	}

	return ctx.JSON(http.StatusOK, cluster_service)
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
	auth := new(authv1.HttpReqAuth)
	if err := echoutil.Bind(ctx, auth); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(auth),
				)))
	}

	// auth := body
	cluster_uuid := auth.ClusterUuid

	//valid cluster
	if _, err := vault.NewCluster(ctl.NewSession()).Get(cluster_uuid); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "valid%s",
				logs.KVL(
					"cluster", cluster_uuid,
				)))
	}

	//valid token
	where := "user_kind = ? AND user_uuid = ? AND token = ?"
	tokens, err := vault.NewToken(ctl.NewSession()).
		Find(where, tokenv1.TokenUserKindCluster.String(), auth.ClusterUuid, auth.Assertion)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "NewToken Find"))
	}

	first := func() *tokenv1.Token {
		for _, it := range tokens {
			return &it
		}
		return nil
	}
	token := first()

	if token == nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			fmt.Errorf("record was not found: token"))
	}

	//만료 시간 검증
	if time.Until(token.ExpirationTime) < 0 {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			fmt.Errorf("token was expierd"))
	}

	//new session
	//make session payload
	token_uuid := macro.NewUuidString()
	iat := time.Now()
	exp := env.ClientSessionExpirationTime(iat)

	payload := &sessionv1.ClientSessionPayload{
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
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "jwt New payload=%+v", payload))
	}

	_, err = ctl.Scope(func(db database.Context) (interface{}, error) {
		//클라이언트를 조회 하여
		//레코드에 없으면 추가
		clients, err := vault.NewClient(db).
			Find("cluster_uuid = ? AND uuid = ?", auth.ClusterUuid, auth.ClientUuid)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "NewClient Find"))
		}

		//없으면 추가
		if len(clients) == 0 {
			name := fmt.Sprintf("client:%s", auth.ClientUuid)
			summary := fmt.Sprintf("client: %s, cluster: %s", auth.ClientUuid, auth.ClusterUuid)
			client := clientv1.Client{}
			client.Uuid = auth.ClientUuid
			client.LabelMeta = NewLabelMeta(name, newist.String(summary))
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
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
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

func newClientSession(payload sessionv1.ClientSessionPayload, token string) sessionv1.Session {
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
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "missing request header%s",
				logs.KVL(
					"header", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__,
				)))

	}

	//verify
	//jwt verify
	if err := jwt.Verify(token, []byte(env.ClientSessionSignatureSecret())); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(err, "jwt verify%s",
				logs.KVL(
					"header", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__,
					"token", token,
				)))
	}

	payload := new(sessionv1.ClientSessionPayload)
	if err := jwt.BindPayload(token, payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(err, "jwt bind payload%s",
				logs.KVL(
					"header", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__,
					"token", token,
				)))
	}

	//만료시간 비교
	if time.Until(payload.Exp) < 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			fmt.Errorf("token expierd%s",
				logs.KVL(
					"header", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__,
					"exp", payload.Exp.String(),
				)))
	}

	_, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		//smart polling
		where := "uuid = ?"
		clusters, err := vault.NewCluster(db).Find(where, payload.ClusterUuid)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "find cluster%s",
					logs.KVL(
						"uuid", payload.ClusterUuid,
					)))
		}

		if len(clusters) == 0 {
			return nil, echo.NewHTTPError(http.StatusBadRequest).SetInternal(
				errors.Wrapf(err, "not found cluster%s",
					logs.KVL(
						"uuid", payload.ClusterUuid,
					)))
		}

		service_count, err := countGatherClusterService(db, payload.ClusterUuid)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "count undone service%s",
					logs.KVL(
						"cluster_uuid", payload.ClusterUuid,
					)))
		}

		var cluster clusterv1.Cluster
		for _, iter := range clusters {
			cluster = iter
			break
		}

		//reflesh payload
		payload.Exp = env.ClientSessionExpirationTime(time.Now())
		// payload.PollInterval = env.ClientConfigPollInterval()
		payload.PollInterval = int(cluster.GetPollingOption().Interval(time.Duration(int64(env.ClientConfigPollInterval())*int64(time.Second)), int(service_count)) / time.Second)
		payload.Loglevel = env.ClientConfigLoglevel()

		//new jwt-new_token
		new_token, err := jwt.New(payload, []byte(env.ClientSessionSignatureSecret()))
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "new jwt"))
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

func clientTokenPayload(ctx echo.Context) (*sessionv1.ClientSessionPayload, error) {
	//get token
	token := ctx.Request().Header.Get(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__)
	if len(token) == 0 {
		return nil, fmt.Errorf("missing request header%s",
			logs.KVL(
				"key", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__,
			))
	}

	payload := new(sessionv1.ClientSessionPayload)
	if err := jwt.BindPayload(token, payload); err != nil {
		return nil, errors.Wrapf(err, "jwt bind payload%s",
			logs.KVL(
				"token", token,
			))
	}
	return payload, nil
}

func gatherClusterService(db database.Context, cluster_uuid string, fn func(servicev1.Service, []stepv1.ServiceStep)) error {
	where := "cluster_uuid = ? AND (status BETWEEN ? AND ?)"
	args := []interface{}{
		cluster_uuid,
		servicev1.StatusRegist,
		servicev1.StatusProcessing,
	}
	service, err := vault.NewService(db).Find(where, args...)
	if err != nil {
		return errors.Wrapf(err, "find service")
	}
	for _, service := range service {
		where := "service_uuid = ?"
		steps, err := vault.NewServiceStep(db).Find(where, service.Uuid)
		if err != nil {
			return errors.Wrapf(err, "find service step")
		}

		fn(service, steps)
	}

	return nil
}

func countGatherClusterService(db database.Context, cluster_uuid string) (int64, error) {
	where := "cluster_uuid = ? AND (status BETWEEN ? AND ?)"
	args := []interface{}{
		cluster_uuid,
		servicev1.StatusRegist,
		servicev1.StatusProcessing,
	}

	count, err := db.Where(where, args...).Count(new(servicev1.Service))
	if err != nil {
		return 0, errors.Wrapf(err, "service count")
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
