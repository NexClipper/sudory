package control

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/event"
	"github.com/NexClipper/sudory/pkg/server/event/managed_event"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/pkg/errors"
	"xorm.io/xorm"

	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
	authv1 "github.com/NexClipper/sudory/pkg/server/model/auth/v1"
	clusterv1 "github.com/NexClipper/sudory/pkg/server/model/cluster/v1"
	clustertokenv1 "github.com/NexClipper/sudory/pkg/server/model/cluster_token/v1"
	cryptov1 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v1"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
	"github.com/NexClipper/sudory/pkg/server/status/globvar"
	"github.com/golang-jwt/jwt/v4"

	"github.com/labstack/echo/v4"
)

// @Description get []Service
// @Accept      json
// @Produce     json
// @Tags        client/service
// @Router      /client/service [get]
// @Param       x-sudory-client-token header string true "client session token"
// @Success     200 {array}  v1.HttpRspService_ClientSide
// @Header      200 {string} x-sudory-client-token
func (ctl Control) PollingService(ctx echo.Context) error {
	clientTokenPayload_once := func(ctx echo.Context) func() (*sessionv1.ClientSessionPayload, error) {
		var (
			once      sync.Once
			jwt_token *jwt.Token
			claims    *sessionv1.ClientSessionPayload
			err       error
		)
		return func() (*sessionv1.ClientSessionPayload, error) {
			once.Do(func() {
				token := ctx.Request().Header.Get(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__)
				if len(token) == 0 {
					err = errors.Errorf("missing request header%s",
						logs.KVL(
							"key", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__,
						))
				}

				claims = new(sessionv1.ClientSessionPayload)
				jwt_token, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(globvar.ClientSessionSignatureSecret()), nil
				})
				var ok bool
				if claims, ok = jwt_token.Claims.(*sessionv1.ClientSessionPayload); !ok || !jwt_token.Valid {
					err = errors.Wrapf(err, "jwt bind payload%s",
						logs.KVL(
							"token", token,
						))
				}
			})
			return claims, err
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
		service.AssignedClientUuid = payload.Uuid

		return service
	}

	//make response
	rsp_service := make([]servicev1.HttpRspService, 0)
	response_callback := func(s servicev1.Service, ss []stepv1.ServiceStep) {
		//service; chaining step
		s.ServiceProperty = s.ChaniningStep(ss)
		//service; assign client info
		s = mapper_response_assign_client_info(s)
		//setp; update status
		for i := range ss {
			ss[i] = mapper_response_step_status(ss[i])
		}

		rsp_service = append(rsp_service, servicev1.HttpRspService{Service: s, Steps: ss})
	}
	//get token payload
	payload, err := clientTokenPayload_once()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "client token payload"))
	}
	if err := gatherClusterService(ctl.db.Engine().NewSession(), payload.ClusterUuid, response_callback); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "gather cluster service%s",
				logs.KVL(
					"cluster_uuid", payload.ClusterUuid,
				)))
	}

	//save response
	if _, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		for i := range rsp_service {
			service := rsp_service[i].Service
			steps := rsp_service[i].Steps

			service_, err := vault.NewService(tx).Update(service)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
					errors.Wrapf(err, "update service"))
			}

			for i := range steps {
				step := steps[i]
				step_, err := vault.NewServiceStep(tx).Update(step)
				if err != nil {
					return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
						errors.Wrapf(err, "update step"))
				}
				steps[i] = *step_
			}

			rsp_service[i].Service = *service_
			rsp_service[i].Steps = steps
		}

		return nil, nil
	}); err != nil {
		return err
	}

	//invoke event (service-poll-out)
	for _, response := range rsp_service {
		const event_name = "service-poll-out"
		m := map[string]interface{}{}
		m["event_name"] = event_name
		m["service_uuid"] = response.Uuid
		m["service_name"] = response.Name
		m["template_uuid"] = response.TemplateUuid
		m["cluster_uuid"] = response.ClusterUuid
		m["assigned_client_uuid"] = response.AssignedClientUuid
		m["status"] = nullable.Int32(response.Status).Value()
		if response.Result != nil {
			m["result"] = string(*response.Result)
		}
		m["step_count"] = nullable.Int32(response.StepCount).Value()
		m["step_position"] = nullable.Int32(response.StepPosition).Value()

		event.Invoke(event_name, m)
		managed_event.Invoke(response.ClusterUuid, event_name, m)
	}

	return ctx.JSON(http.StatusOK, rsp_service)
}

// @Description update a service
// @Accept      json
// @Produce     json
// @Tags        client/service
// @Router      /client/service [put]
// @Param       x-sudory-client-token header string           true  "client session token"
// @Param       body body v1.HttpReq_ServiceUpdate_ClientSide true "HttpReq_ServiceUpdate_ClientSide"
// @Success     200
// @Header      200 {string} x-sudory-client-token
func (ctl Control) UpdateService(ctx echo.Context) error {
	body := servicev1.HttpReq_ServiceUpdate_ClientSide{}
	if err := echoutil.Bind(ctx, &body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}
	//get record; service, steps
	service, steps, err := vault.NewService(ctl.db.Engine().NewSession()).Get(body.Uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get service"))
	}

	//service; update body
	service.Result = func(s *string) *cryptov1.String {
		if s == nil {
			return nil
		}
		t := cryptov1.String(*body.Result)
		return &t
	}(body.Result)

	//steps; update body
	for i := range steps {
		steps[i].Uuid = body.Steps[i].Uuid
		steps[i].Status = body.Steps[i].Status
		steps[i].Started = body.Steps[i].Started
		steps[i].Ended = body.Steps[i].Ended
	}

	//service; ChaniningStep
	service.ServiceProperty = service.ChaniningStep(steps)

	if _, err := ctl.ScopeSession(func(tx *xorm.Session) (out interface{}, err error) {
		//save service
		service, err = vault.NewService(tx).Update(*service)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "update service"))
		}
		//save step
		for _, step := range steps {
			if _, err := vault.NewServiceStep(tx).Update(step); err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
					errors.Wrapf(err, "update service step"))
			}
		}
		//OnCompletion
		if servicev1.StatusSuccess <= servicev1.Status(nullable.Int32(service.Status).Value()) {
			switch servicev1.OnCompletion(nullable.Int8(service.OnCompletion).Value()) {
			case servicev1.OnCompletionRemove:
				//OnCompletionRemove
				if err := vault.NewService(tx).Delete(service.Uuid); err != nil {
					return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
						errors.Wrapf(err, "on completion a service%v",
							logs.KVL(
								"uuid", service.Uuid,
							)))
				}
			}
		}

		return nil, nil
	}); err != nil {
		return err
	}

	//invoke event (service-poll-in)
	const event_name = "service-poll-in"
	m := map[string]interface{}{}
	m["event_name"] = event_name
	m["service_uuid"] = service.Uuid
	m["service_name"] = service.Name
	m["template_uuid"] = service.TemplateUuid
	m["cluster_uuid"] = service.ClusterUuid
	m["assigned_client_uuid"] = service.AssignedClientUuid
	m["status"] = nullable.Int32(service.Status).Value()
	if service.Result != nil {
		m["result"] = string(*service.Result)
	}
	m["step_count"] = nullable.Int32(service.StepCount).Value()
	m["step_position"] = nullable.Int32(service.StepPosition).Value()

	event.Invoke(service.SubscribedChannel, m)                              //Subscribe 등록된 구독 이벤트 이름으로 호출
	managed_event.Invoke(service.ClusterUuid, service.SubscribedChannel, m) //Subscribe 등록된 구독 이벤트 이름으로 호출

	return ctx.JSON(http.StatusOK, OK())
}

// @Description auth client
// @Accept      json
// @Produce     json
// @Tags        client/auth
// @Router      /client/auth [post]
// @Param       body body v1.HttpReqAuth true "HttpReqAuth"
// @Success     200 {string} ok
// @Header      200 {string} x-sudory-client-token
func (ctl Control) AuthClient(ctx echo.Context) error {
	auth := new(authv1.HttpReqAuth)
	if err := echoutil.Bind(ctx, auth); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(auth),
				)))
	}

	//valid cluster
	cluster := clusterv1.Cluster{}
	if err := database.XormGet(
		ctl.db.Engine().NewSession().
			Where("uuid = ?", auth.ClusterUuid),
		&cluster); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "valid cluster%s",
				logs.KVL(
					"cluster_uuid", auth.ClusterUuid,
				)))
	}

	//valid token
	cluster_token := clustertokenv1.ClusterToken{}
	if err := database.XormGet(
		ctl.db.Engine().NewSession().
			Where("token = ? AND cluster_uuid = ?", auth.Assertion, auth.ClusterUuid),
		&cluster_token); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "valid cluster token%s",
				logs.KVL(
					"cluster_uuid", auth.ClusterUuid,
					"assertion", auth.Assertion,
				)))
	}

	//만료 시간 검증
	if time.Until(cluster_token.ExpirationTime) < 0 {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Errorf("token was expierd"))
	}

	//new session
	//make session payload
	session_token_uuid := macro.NewUuidString()
	iat := time.Now()
	exp := globvar.ClientSessionExpirationTime(iat)

	payload := &sessionv1.ClientSessionPayload{
		ExpiresAt:    exp.Unix(),
		IssuedAt:     iat.Unix(),
		Uuid:         session_token_uuid,
		ClusterUuid:  auth.ClusterUuid,
		PollInterval: globvar.ClientConfigPollInterval(),
		Loglevel:     globvar.ClientConfigLoglevel(),
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

	token_string, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).
		SignedString([]byte(globvar.ClientSessionSignatureSecret()))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "jwt New payload=%+v", payload))
	}

	session := newClientSession(*payload, token_string)

	_, err = ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		//save session
		if _, err := tx.Insert(&session); err != nil {
			return nil, errors.Wrapf(err, "session insert")
		}

		return nil, nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	//save token to header
	ctx.Response().Header().Add(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__, token_string)

	//invoke event (client-auth-accept)
	const event_name = "client-auth-accept"
	m := map[string]interface{}{
		"event_name":   event_name,
		"cluster_uuid": payload.ClusterUuid,
		"session_uuid": payload.Uuid,
	}
	event.Invoke(event_name, m)
	managed_event.Invoke(payload.ClusterUuid, event_name, m)

	return ctx.JSON(http.StatusOK, OK())
}

func newClientSession(payload sessionv1.ClientSessionPayload, token string) sessionv1.Session {
	session := sessionv1.Session{}
	session.Uuid = payload.Uuid
	session.ClusterUuid = payload.ClusterUuid
	session.Token = token
	session.IssuedAtTime = time.Unix(payload.IssuedAt, 0)
	session.ExpirationTime = time.Unix(payload.ExpiresAt, 0)
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

	claims := new(sessionv1.ClientSessionPayload)
	jwt_token, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(globvar.ClientSessionSignatureSecret()), nil
	})

	if _, ok := jwt_token.Claims.(*sessionv1.ClientSessionPayload); !ok || !jwt_token.Valid {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(err, "jwt verify%s",
				logs.KVL(
					"header", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__,
					"token", token,
				)))
	}

	if _, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		//smart polling
		where := "uuid = ?"
		clusters, err := vault.NewCluster(tx).Find(where, claims.ClusterUuid)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "find cluster%s",
					logs.KVL(
						"uuid", claims.ClusterUuid,
					)))
		}

		if len(clusters) == 0 {
			return nil, echo.NewHTTPError(http.StatusBadRequest).SetInternal(
				errors.Wrapf(err, "not found cluster%s",
					logs.KVL(
						"uuid", claims.ClusterUuid,
					)))
		}

		service_count, err := countGatherClusterService(tx, claims.ClusterUuid)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "count undone service%s",
					logs.KVL(
						"cluster_uuid", claims.ClusterUuid,
					)))
		}

		var cluster clusterv1.Cluster
		for _, iter := range clusters {
			cluster = iter
			break
		}

		//reflesh payload
		claims.ExpiresAt = globvar.ClientSessionExpirationTime(time.Now()).Unix()
		// payload.PollInterval = globvar.ClientConfigPollInterval()
		claims.PollInterval = int(cluster.GetPollingOption().Interval(time.Duration(int64(globvar.ClientConfigPollInterval())*int64(time.Second)), int(service_count)) / time.Second)
		claims.Loglevel = globvar.ClientConfigLoglevel()

		//new jwt-new_token
		// new_token, err := jwt.New(claims, []byte(globvar.ClientSessionSignatureSecret()))
		// if err != nil {
		// 	return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
		// 		errors.Wrapf(err, "new jwt"))
		// }

		//client auth 에서 사용된 알고리즘 그대로 사용
		token_string, err := jwt.NewWithClaims(usedJwtSigningMethod(*jwt_token, jwt.SigningMethodHS256), claims).
			SignedString([]byte(globvar.ClientSessionSignatureSecret()))
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "new jwt"))
		}
		//udpate session
		session := newClientSession(*claims, token_string)
		if _, err := vault.NewSession(tx).Update(session); err != nil {
			return nil, errors.Wrapf(err, "update session-token")
		}

		//save client session-token to header
		ctx.Response().Header().Add(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__, token_string)

		return nil, nil
	}); err != nil {
		return err
	}

	return nil
}
func usedJwtSigningMethod(token jwt.Token, init jwt.SigningMethod) jwt.SigningMethod {
	alg, _ := token.Header["alg"].(string)

	if jwt.GetSigningMethod(alg) != nil {
		init = jwt.GetSigningMethod(alg)
	}

	return init
}

func gatherClusterService(tx *xorm.Session, cluster_uuid string, fn func(servicev1.Service, []stepv1.ServiceStep)) error {
	cluster, err := vault.NewCluster(tx).Get(cluster_uuid)
	if err != nil {
		return errors.Wrapf(err, "gather service for cluster")
	}
	if nullable.Int16(cluster.PoliingLimit).Has() && 0 < nullable.Int16(cluster.PoliingLimit).Value() {
		tx = tx.Limit(int(nullable.Int16(cluster.PoliingLimit).Value()))
	}

	where := "cluster_uuid = ? AND (status BETWEEN ? AND ?)"
	args := []interface{}{
		cluster_uuid,
		servicev1.StatusRegist,
		servicev1.StatusProcessing,
	}
	service, steps, err := vault.NewService(tx).Find(where, args...)
	if err != nil {
		return errors.Wrapf(err, "gather service for cluster")
	}
	for _, service := range service {
		fn(service, steps[service.Uuid])
	}

	return nil
}

func countGatherClusterService(tx *xorm.Session, cluster_uuid string) (int64, error) {
	cluster, err := vault.NewCluster(tx).Get(cluster_uuid)
	if err != nil {
		return 0, errors.Wrapf(err, "count service for cluster")
	}
	if nullable.Int16(cluster.PoliingLimit).Has() && 0 < nullable.Int16(cluster.PoliingLimit).Value() {
		tx = tx.Limit(int(nullable.Int16(cluster.PoliingLimit).Value()))
	}

	where := "cluster_uuid = ? AND (status BETWEEN ? AND ?)"
	args := []interface{}{
		cluster_uuid,
		servicev1.StatusRegist,
		servicev1.StatusProcessing,
	}
	count, err := tx.Where(where, args...).Count(new(servicev1.Service))
	if err != nil {
		return 0, errors.Wrapf(err, "count service for cluster")
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
