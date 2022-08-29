package control

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/client/log"
	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/event"
	"github.com/NexClipper/sudory/pkg/server/event/managed_channel"
	"github.com/NexClipper/sudory/pkg/server/event/managed_event"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/pkg/errors"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmtex"

	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	authv2 "github.com/NexClipper/sudory/pkg/server/model/auth/v2"
	channelv3 "github.com/NexClipper/sudory/pkg/server/model/channel/v3"
	clusterv3 "github.com/NexClipper/sudory/pkg/server/model/cluster/v3"
	clusterinfov2 "github.com/NexClipper/sudory/pkg/server/model/cluster_infomation/v2"
	clustertokenv3 "github.com/NexClipper/sudory/pkg/server/model/cluster_token/v3"
	cryptov2 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v2"
	servicev3 "github.com/NexClipper/sudory/pkg/server/model/service/v3"

	sessionv3 "github.com/NexClipper/sudory/pkg/server/model/session/v3"
	"github.com/NexClipper/sudory/pkg/server/status/globvar"
	"github.com/golang-jwt/jwt/v4"

	"github.com/labstack/echo/v4"
)

// @Description get []Service
// @Security    ClientAuth
// @Accept      json
// @Produce     json
// @Tags        client/service
// @Router      /client/service [get]
// @Success     200 {array}  v3.HttpRsp_ClientServicePolling
// @Header      200 {string} x-sudory-client-token
func (ctl ControlVanilla) PollingService(ctx echo.Context) error {

	//get token claims
	claims, err := GetSudoryClientTokenClaims(ctx)
	err = errors.Wrapf(err, "failed to get client token")
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// get cluster
	cluster := new(clusterv3.Cluster)
	cluster.Uuid = claims.ClusterUuid
	eq_uuid := stmt.Equal("uuid", cluster.Uuid)

	err = stmtex.Select(cluster.TableName(), cluster.ColumnNames(), eq_uuid, nil, nil).
		QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(func(s stmtex.Scanner) (err error) {
		err = cluster.Scan(s)
		err = errors.Wrapf(err, "scan cluster")
		return
	})
	if err != nil {
		return errors.Wrapf(err, "failed to get cluster")
	}

	// gather service
	// get service offset
	var polling_offest vanilla.NullTime
	func() (err error) {
		cluster_info := clusterinfov2.ClusterInformation{}
		columnnames := []string{"polling_offset"}
		cond := stmt.Equal("cluster_uuid", claims.ClusterUuid)

		err = stmtex.Select(cluster_info.TableName(), columnnames, cond, nil, nil).
			QueryRowsContext(ctx.Request().Context(), ctl, ctl.Dialect())(func(scan stmtex.Scanner, _ int) error {
			return scan.Scan(&polling_offest)
		})
		if err != nil {
			return errors.Wrapf(err, "failed to get service offset")
		}
		return
	}()

	// polling limit filter
	polling_filter := newPollingFilter(cluster.PoliingLimit)
	services, steps, err := pollingService(ctx.Request().Context(), ctl, ctl.Dialect(), claims.ClusterUuid, polling_offest, polling_filter)

	// set polling_offest
	var polling_offest_ time.Time
	for _, service := range services {
		if polling_offest_.IsZero() {
			polling_offest_ = service.Created
		}

		if polling_offest_.After(service.Created) {
			polling_offest_ = service.Created
		}
	}

	if !polling_offest_.IsZero() {
		polling_offest = *vanilla.NewNullTime(polling_offest_)
	}

	// save polling_count to cluster_infomation
	cluster_info := clusterinfov2.ClusterInformation{}
	cluster_info.ClusterUuid = cluster.Uuid
	cluster_info.PollingCount = *vanilla.NewNullInt(len(services))
	cluster_info.PollingOffset = polling_offest
	cluster_info.Created = time.Now()
	cluster_info.Updated = *vanilla.NewNullTime(cluster_info.Created)

	cluster_info_update_columns := []string{
		"polling_count", "polling_offset", "updated",
	}
	_, _, err = stmtex.InsertOrUpdate(cluster_info.TableName(), cluster_info.ColumnNames(), cluster_info_update_columns, cluster_info.Values()).
		ExecContext(ctx.Request().Context(), ctl, ctl.Dialect())
	if err != nil {
		return errors.Wrapf(err, "exec insert|update statement")
	}

	UpdateServiceStatus := func(service servicev3.Service, assigned_client_uuid string, status servicev3.StepStatus, t time.Time) servicev3.Service {
		service.AssignedClientUuid = *vanilla.NewNullString(assigned_client_uuid)
		service.Status = status
		service.Timestamp = t
		return service
	}
	UpdateStepStatus := func(step servicev3.ServiceStep, status servicev3.StepStatus, t time.Time) servicev3.ServiceStep {
		step.Status = status
		step.Timestamp = t
		return step
	}

	type CALLBACK_Service func(a servicev3.Service)
	type CALLBACK_ServiceStep func(a servicev3.ServiceStep)

	MakeUpdate := func(
		t time.Time,
		services []servicev3.Service, CALLBACK_Service CALLBACK_Service,
		steps map[string][]servicev3.ServiceStep, CALLBACK_ServiceStep CALLBACK_ServiceStep) {

		for _, service := range services {
			if service.Status == servicev3.StepStatusRegist {
				v := UpdateServiceStatus(service, claims.Uuid, servicev3.StepStatusSend, t)
				CALLBACK_Service(v)

				for _, step := range steps[v.Uuid] {
					v := UpdateStepStatus(step, servicev3.StepStatusSend, t)
					CALLBACK_ServiceStep(v)
				}
			}
		}
	}

	var new_service_status = make([]vault.Table, 0, len(services))
	var new_step_status = make([]vault.Table, 0, len(steps))

	time_now := time.Now()
	MakeUpdate(time_now,
		services, func(a servicev3.Service) {
			new_service_status = append(new_service_status, a)
		},
		steps, func(a servicev3.ServiceStep) {
			new_step_status = append(new_step_status, a)
		})

	err = stmtex.ScopeTx(ctx.Request().Context(), ctl, func(tx *sql.Tx) error {

		if err = vault.SaveMultiTable(tx, ctl.Dialect(), new_service_status); err != nil {
			return errors.Wrapf(err, "faild to save service")
		}

		if err = vault.SaveMultiTable(tx, ctl.Dialect(), new_step_status); err != nil {
			return errors.Wrapf(err, "faild to save service step")
		}

		return nil
	})
	if err != nil {
		return err
	}

	// make response body
	rsp := make([]servicev3.HttpRsp_ClientServicePolling, len(services))
	for i, service := range services {
		rsp[i].Service = service
		rsp[i].Steps = steps[service.Uuid]
		i++
	}

	// invoke event (service-poll-out)
	for _, service := range rsp {
		const event_name = "service-poll-out"
		m := map[string]interface{}{}
		m["event_name"] = event_name
		m["service_uuid"] = service.Uuid
		m["service_name"] = service.Name
		m["template_uuid"] = service.TemplateUuid
		m["cluster_uuid"] = service.ClusterUuid
		m["assigned_client_uuid"] = service.AssignedClientUuid
		m["status"] = service.Status
		m["status_description"] = service.Status.String()
		m["result_type"] = servicev3.ResultSaveTypeNone.String()
		m["result"] = ""
		m["step_count"] = service.StepCount
		m["step_position"] = service.StepPosition

		event.Invoke(event_name, m)
		managed_event.Invoke(event_name, m)
		// invoke event by event category
		managed_channel.InvokeByEventCategory(channelv3.EventCategoryServicePollingOut, m)
	}

	return ctx.JSON(http.StatusOK, rsp)
}

// @Description update a service
// @Security    ClientAuth
// @Accept      json
// @Produce     json
// @Tags        client/service
// @Router      /client/service [put]
// @Param       body body v3.HttpReq_ClientServiceUpdate true "HttpReq_ClientServiceUpdate"
// @Success     200
// @Header      200 {string} x-sudory-client-token
func (ctl ControlVanilla) UpdateService(ctx echo.Context) (err error) {

	stepPosition := func(service_step servicev3.ServiceStep) int {
		// 스템 포지션 값은
		// ServiceStep.Sequence+1
		return service_step.Sequence + 1
	}
	stepStatus := func(body servicev3.HttpReq_ClientServiceUpdate, service servicev3.Service, service_step servicev3.ServiceStep) servicev3.StepStatus {
		// 스탭 포지션이 카운트와 같은 경우만
		if service.StepCount == stepPosition(service_step) {
			return body.Status
		}
		// 기본값; 처리중(Processing)
		return servicev3.StepStatusProcessing
	}
	serviceResult := func(body servicev3.HttpReq_ClientServiceUpdate) cryptov2.CryptoString {
		// 상태가 성공인 경우만
		if body.Status == servicev3.StepStatusSuccess {
			return cryptov2.CryptoString(body.Result)
		}
		//기본값; 공백 문자열
		return ""
	}
	stepMessage := func(body servicev3.HttpReq_ClientServiceUpdate) vanilla.NullString {
		// 상태가 실패인 경우만
		if body.Status == servicev3.StepStatusFail {
			return *vanilla.NewNullString(body.Result)
		}
		//기본값; 공백 문자열
		return vanilla.NullString{}
	}
	resultType := func(body servicev3.HttpReq_ClientServiceUpdate, service servicev3.Service, service_step servicev3.ServiceStep) (resultType servicev3.ResultSaveType) {
		// 마지막 스탭의 결과만 저장 한다
		if service.StepCount != stepPosition(service_step) {
			return
		}
		// 상태 값이 성공이 아닌 경우
		// 서비스 결과를 저장 하지 않는다
		if servicev3.StepStatusSuccess != stepStatus(body, service, service_step) {
			return
		}
		// 채널이 등록되어 있는 경우
		// 서비스 결과를 저장 하지 않는다
		if !service.SubscribedChannel.Valid || 0 < len(service.SubscribedChannel.String) {
			return
		}

		return servicev3.ResultSaveTypeDatabase
	}

	body := servicev3.HttpReq_ClientServiceUpdate{}
	if err := echoutil.Bind(ctx, &body); err != nil {
		err = errors.Wrapf(err, "bind%s",
			logs.KVL(
				"type", TypeName(body),
			))
		return HttpError(err, http.StatusBadRequest)
	}

	if err := body.Status.Valid(); err != nil {
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(fmt.Sprintf("%s.Status", TypeName(body)), body.Status)...,
			))
		return HttpError(err, http.StatusBadRequest)
	}

	// get client token claims
	claims, err := GetSudoryClientTokenClaims(ctx)
	if err != nil {
		err = errors.Wrapf(err, "failed to get client token")
		return HttpError(err, http.StatusBadRequest)
	}

	// get service
	service := new(servicev3.Service)
	service.ClusterUuid = claims.ClusterUuid
	service.Uuid = body.Uuid

	service, err = vault.Servicev3.GetService(context.Background(), ctl, ctl.Dialect(), service.ClusterUuid, service.Uuid)
	if err != nil {
		return errors.Wrapf(err, "failed to found service%v", logs.KVL(
			"cluster_uuid", claims.ClusterUuid,
			"uuid", body.Uuid,
		))
	}

	// get service service_step
	service_step := new(servicev3.ServiceStep)
	service_step.ClusterUuid = claims.ClusterUuid
	service_step.Uuid = body.Uuid
	service_step.Sequence = body.Sequence

	service_step, err = vault.Servicev3.GetServiceStep(context.Background(), ctl, ctl.Dialect(), service_step.ClusterUuid, service_step.Uuid, service_step.Sequence)
	if err != nil {
		return errors.Wrapf(err, "failed to found service step%v", logs.KVL(
			"cluster_uuid", claims.ClusterUuid,
			"uuid", body.Uuid,
			"sequence", uint8(body.Sequence),
		))
	}

	now_time := time.Now()
	// udpate service
	{
		// update key
		service.Timestamp = now_time
		// update value
		service.AssignedClientUuid = *vanilla.NewNullString(claims.Uuid)
		service.StepPosition = stepPosition(*service_step)
		service.Status = stepStatus(body, *service, *service_step)
		service.Message = stepMessage(body)
	}
	// new service result
	service_result := func() *servicev3.ServiceResult {
		// update key
		service_result := new(servicev3.ServiceResult)
		service_result.PartitionDate = service.PartitionDate
		service_result.ClusterUuid = service.ClusterUuid
		service_result.Uuid = service.Uuid
		service_result.Timestamp = now_time
		// update value
		service_result.ResultSaveType = resultType(body, *service, *service_step)
		service_result.Result = serviceResult(body)

		return service_result
	}()
	// udpate service step
	{
		// update key
		service_step.Timestamp = now_time
		// update value
		service_step.Status = body.Status                         // Status
		service_step.Started = *vanilla.NewNullTime(body.Started) // Started
		service_step.Ended = *vanilla.NewNullTime(body.Ended)     // Ended
	}

	// save to db
	err = stmtex.ScopeTx(context.Background(), ctl, func(tx *sql.Tx) (err error) {

		// save service
		if err = vault.SaveMultiTable(tx, ctl.Dialect(), []vault.Table{service}); err != nil {
			return errors.Wrapf(err, "failed to save service_status")
		}

		// save service step
		if err = vault.SaveMultiTable(tx, ctl.Dialect(), []vault.Table{service_step}); err != nil {
			return errors.Wrapf(err, "failed to save service_step_status")
		}

		// check servcie result save type
		if service_result.ResultSaveType != servicev3.ResultSaveTypeNone {
			// save service result
			if err = vault.SaveMultiTable(tx, ctl.Dialect(), []vault.Table{service_result}); err != nil {
				return errors.Wrapf(err, "failed to save service_result")
			}
		}

		return
	})
	err = errors.Wrapf(err, "failed to save")
	if err != nil {
		return err
	}

	// invoke event (service-poll-in)
	const event_name = "service-poll-in"
	m := map[string]interface{}{}
	m["event_name"] = event_name
	m["service_uuid"] = service.Uuid
	m["service_name"] = service.Name
	m["template_uuid"] = service.TemplateUuid
	m["cluster_uuid"] = service.ClusterUuid
	m["assigned_client_uuid"] = service.AssignedClientUuid
	m["status"] = service.Status
	m["status_description"] = service.Status.String()
	m["result_type"] = service_result.ResultSaveType.String()
	m["result"] = body.Result
	m["step_count"] = service.StepCount
	m["step_position"] = service.StepPosition

	if service.Status == servicev3.StepStatusSuccess && len(service_result.Result.String()) == 0 {
		log.Debugf("channel(poll-in-service): %+v", m)
	}

	event.Invoke(service.SubscribedChannel.String, m)         //Subscribe 등록된 구독 이벤트 이름으로 호출
	managed_event.Invoke(service.SubscribedChannel.String, m) //Subscribe 등록된 구독 이벤트 이름으로 호출
	// invoke event by channel uuid
	if 0 < len(service.SubscribedChannel.String) {
		// find channel
		mc := channelv3.ManagedChannel{}
		mc.Uuid = service.SubscribedChannel.String
		mc_cond := stmt.And(
			stmt.Equal("uuid", mc.Uuid),
			stmt.IsNull("deleted"),
		)
		found, err := stmtex.ExistContext(mc.TableName(), mc_cond)(context.Background(), ctl, ctl.Dialect())
		if err != nil {
			return err
		}
		if found {
			managed_channel.InvokeByChannelUuid(service.SubscribedChannel.String, m)
		}
	}
	// invoke event by event category
	managed_channel.InvokeByEventCategory(channelv3.EventCategoryServicePollingIn, m)

	return ctx.JSON(http.StatusOK, OK())
}

// @Description auth client
// @Accept      json
// @Produce     json
// @Tags        client/auth
// @Router      /client/auth [post]
// @Param       body body v2.HttpReqAuth true "HttpReqAuth"
// @Success     200 {string} ok
// @Header      200 {string} x-sudory-client-token
func (ctl ControlVanilla) AuthClient(ctx echo.Context) (err error) {
	auth := new(authv2.HttpReqAuth)
	err = func() (err error) {
		if err := echoutil.Bind(ctx, auth); err != nil {
			return errors.Wrapf(err, "bind%s",
				logs.KVL(
					"type", TypeName(auth),
				))
		}
		return
	}()
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	cluster := clusterv3.Cluster{}
	cluster.Uuid = auth.ClusterUuid
	cluster_eq_uuid := stmt.Equal("uuid", cluster.Uuid)
	cluster_found, err := stmtex.ExistContext(cluster.TableName(), cluster_eq_uuid)(ctx.Request().Context(), ctl, ctl.Dialect())
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}
	if !cluster_found {
		err = errors.Errorf("not found cluster%v",
			logs.KVL(
				"cluster_uuid", cluster.Uuid,
			))
		return HttpError(err, http.StatusBadRequest)
	}

	//valid token
	token := clustertokenv3.ClusterToken{}
	token.ClusterUuid = auth.ClusterUuid
	token.Token = cryptov2.CryptoString(auth.Assertion)

	token_cond := stmt.And(
		stmt.Equal("cluster_uuid", token.ClusterUuid),
		stmt.Equal("token", token.Token),
	)

	err = stmtex.Select(token.TableName(), token.ColumnNames(), token_cond, nil, nil).
		QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(func(scan stmtex.Scanner) (err error) {
		err = token.Scan(scan)
		return
	})
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	//만료 시간 검증
	if time.Until(token.ExpirationTime) < 0 {
		return HttpError(errors.Errorf("token was expierd"), http.StatusBadRequest)
	}

	//new session
	//make session payload
	session_token_uuid := macro.NewUuidString()
	created := time.Now()
	iat := time.Now()
	exp := globvar.ClientSession.ExpirationTime(iat)

	payload := &sessionv3.ClientSessionPayload{
		ExpiresAt:    exp.Unix(),
		IssuedAt:     iat.Unix(),
		Uuid:         session_token_uuid,
		ClusterUuid:  auth.ClusterUuid,
		PollInterval: globvar.ClientConfig.PollInterval(),
		Loglevel:     globvar.ClientConfig.Loglevel(),
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
		SignedString([]byte(globvar.ClientSession.SignatureSecret()))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "jwt New payload=%+v", payload))
	}

	session := sessionv3.Session{}
	session.Uuid = payload.Uuid
	session.ClusterUuid = payload.ClusterUuid
	session.Token = token_string
	session.IssuedAtTime = *vanilla.NewNullTime(time.Unix(payload.IssuedAt, 0))
	session.ExpirationTime = *vanilla.NewNullTime(time.Unix(payload.ExpiresAt, 0))
	session.Created = created

	err = func() (err error) {
		var affected int64
		affected, session.ID, err = stmtex.Insert(session.TableName(), session.ColumnNames(), session.Values()).ExecContext(ctx.Request().Context(), ctl, ctl.Dialect())
		if err != nil {
			return err
		}
		if affected == 0 {
			return errors.New("no affected")
		}
		return
	}()
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	//save token to header
	ctx.Response().Header().Set(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__, token_string)

	//invoke event (client-auth-accept)
	const event_name = "client-auth-accept"
	m := map[string]interface{}{
		"event_name":   event_name,
		"cluster_uuid": payload.ClusterUuid,
		"session_uuid": payload.Uuid,
	}
	event.Invoke(event_name, m)
	managed_event.Invoke(event_name, m)
	// invoke event by event category
	managed_channel.InvokeByEventCategory(channelv3.EventCategoryClientAuthAccept, m)

	return ctx.JSON(http.StatusOK, OK())
}

// func newClientSession(payload sessionv1.ClientSessionPayload, token string) sessionv1.Session {
// 	session := sessionv1.Session{}
// 	session.Uuid = payload.Uuid
// 	session.ClusterUuid = payload.ClusterUuid
// 	session.Token = token
// 	session.IssuedAtTime = time.Unix(payload.IssuedAt, 0)
// 	session.ExpirationTime = time.Unix(payload.ExpiresAt, 0)
// 	return session
// }

func (ctl ControlVanilla) VerifyClientSessionToken(ctx echo.Context) (err error) {
	token := ctx.Request().Header.Get(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__)
	err = func() (err error) {
		if len(token) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter, "missing request header%v", logs.KVL(
				"header", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__,
			))
		}

		claims := new(sessionv3.ClientSessionPayload)
		jwt_token, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(globvar.ClientSession.SignatureSecret()), nil
		})
		if err != nil {
			return errors.Wrapf(err, "jwt parse claims")
		}

		if _, ok := jwt_token.Claims.(*sessionv3.ClientSessionPayload); !ok || !jwt_token.Valid {
			return errors.Errorf("jwt verify%v", logs.KVL(
				"header", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__,
				"token", token,
			))
		}
		return
	}()
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	return nil
}

func (ctl ControlVanilla) RefreshClientSessionToken(ctx echo.Context) (err error) {
	token := ctx.Request().Header.Get(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__)
	var jwt_token *jwt.Token
	claims := new(sessionv3.ClientSessionPayload)
	err = func() (err error) {
		if len(token) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter, "missing request header%s",
				logs.KVL(
					"key", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__,
				))
		}

		jwt_token, _, err = jwt.NewParser().ParseUnverified(token, claims)
		if err != nil {
			return errors.Wrapf(err, "failed to jwt.ParseUnverified%v", logs.KVL(
				"token", token,
			))
		}

		if _, ok := jwt_token.Claims.(*sessionv3.ClientSessionPayload); !ok {
			return errors.Errorf("is not valid type%v",
				logs.KVL(
					"jwt.Token.Claims", TypeName(jwt_token.Claims),
					"jwt.Token.Method.Alg", jwt_token.Method.Alg(),
					"token", token,
				))
		}
		return
	}()
	if err != nil {
		return HttpError(err, http.StatusBadRequest) // StatusBadRequest
	}

	time_now := time.Now()

	// polling interval
	cluster := clusterv3.Cluster{}
	err = func() (err error) {
		eq_uuid := stmt.Equal("uuid", claims.ClusterUuid)
		err = stmtex.Select(cluster.TableName(), cluster.ColumnNames(), eq_uuid, nil, nil).
			QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(
			func(s stmtex.Scanner) (err error) {
				return cluster.Scan(s)
			})

		err = errors.Wrapf(err, "failed to get cluster")
		return
	}()
	if err != nil {
		return HttpError(err, http.StatusInternalServerError) // StatusInternalServerError
	}

	var service_count vanilla.NullInt
	err = func() (err error) {
		cluster_info := clusterinfov2.ClusterInformation{}
		columnnames := []string{"polling_count"}
		cond := stmt.Equal("cluster_uuid", claims.ClusterUuid)

		err = stmtex.Select(cluster_info.TableName(), columnnames, cond, nil, nil).
			QueryRowsContext(ctx.Request().Context(), ctl, ctl.Dialect())(func(scan stmtex.Scanner, _ int) error {
			return scan.Scan(&service_count)
		})

		// cond := pollingServiceCondition(claims.ClusterUuid)
		// service := servicev3.Service{}
		// service_count, err = stmtex.Count(service.TableName(), cond.Parse(), nil)(ctx.Request().Context(), ctl)
		if err != nil {
			return errors.Wrapf(err, "failed to get cluster service count")
		}

		return
	}()
	if err != nil {
		return HttpError(err, http.StatusInternalServerError) // StatusInternalServerError
	}

	// polling interval 해더 저장
	polling_interval := int(clusterv3.ConvPollingOption(cluster.PollingOption).Interval(time.Duration(int64(globvar.ClientConfig.PollInterval())*int64(time.Second)), service_count.Int()) / time.Second)

	//reflesh payload
	claims.PollInterval = polling_interval
	claims.ExpiresAt = globvar.ClientSession.ExpirationTime(time_now).Unix()
	claims.Loglevel = globvar.ClientConfig.Loglevel()

	var new_token_string string
	err = func() (err error) {
		//client auth 에서 사용된 알고리즘 그대로 사용
		new_token_string, err = jwt.NewWithClaims(usedJwtSigningMethod(*jwt_token, jwt.SigningMethodHS256), claims).
			SignedString([]byte(globvar.ClientSession.SignatureSecret()))
		if err != nil {
			return errors.Wrapf(err, "failed to make session token to formed jwt")
		}
		return
	}()
	if err != nil {
		return HttpError(err, http.StatusInternalServerError) // StatusInternalServerError
	}

	//save client session-token to header
	ctx.Response().Header().Set(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__, new_token_string)

	session := sessionv3.Session{}
	session.Uuid = claims.Uuid
	session.Token = new_token_string
	session.ExpirationTime = *vanilla.NewNullTime(time.Unix(claims.ExpiresAt, 0))
	session.Updated = *vanilla.NewNullTime(time_now)

	// uuid = ? AND deleted IS NULL
	cond_session := stmt.And(
		stmt.Equal("uuid", session.Uuid),
		stmt.IsNull("deleted"),
	)
	err = func() (err error) {
		session_found, err := stmtex.ExistContext(session.TableName(), cond_session)(ctx.Request().Context(), ctl, ctl.Dialect())
		if err != nil {
			return err
		}
		if !session_found {
			return errors.Errorf("could not found session%v", logs.KVL(
				"uuid", session.Uuid,
			))
		}
		return
	}()
	if err != nil {
		return HttpError(err, http.StatusInternalServerError) // StatusInternalServerError
	}

	err = func() (err error) {
		keys_values := map[string]interface{}{
			"token":           session.Token,
			"expiration_time": session.ExpirationTime,
			"updated":         session.Updated,
		}

		_, err = stmtex.Update(session.TableName(), keys_values, cond_session).
			Exec(ctl, ctl.Dialect())
		if err != nil {
			return errors.Wrapf(err, "failed to refresh client session%v", logs.KVL(
				"uuid", claims.Uuid,
				"data", keys_values,
			))
		}

		return
	}()
	if err != nil {
		return HttpError(err, http.StatusInternalServerError) // StatusInternalServerError
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

func GetSudoryClientTokenClaims(ctx echo.Context) (claims *sessionv3.ClientSessionPayload, err error) {
	var token string
	if token = ctx.Request().Header.Get(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__); len(token) == 0 {
		err = errors.Errorf("missing request header%s",
			logs.KVL(
				"key", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__,
			))
		return
	}

	claims = new(sessionv3.ClientSessionPayload)
	var fn func() error
	fn = func() error {
		var jwt_token *jwt.Token
		// JWT unverify
		jwt_token, _, err = jwt.NewParser().ParseUnverified(token, claims)
		if err != nil {
			return errors.Wrapf(err, "jwt.Parser.ParseUnverified")
		}

		var ok bool
		claims, ok = jwt_token.Claims.(*sessionv3.ClientSessionPayload)
		if !ok {
			return errors.New("jwt.Token.Claims not matched")
		}

		return nil
	}
	if false {
		fn = func() error {
			var jwt_token *jwt.Token
			// JWT verify
			jwt_token, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(globvar.ClientSession.SignatureSecret()), nil
			})
			if err != nil {
				return errors.Wrapf(err, "jwt.Parser.ParseWithClaims")
			}

			var ok bool
			claims, ok = jwt_token.Claims.(*sessionv3.ClientSessionPayload)
			if !ok {
				return errors.New("jwt.Token.Claims not matched")
			}
			if !jwt_token.Valid {
				return errors.New("jwt.Token.Valid false")
			}

			return nil
		}
	}

	err = fn()
	err = errors.Wrapf(err, "failed to parse header token%v", logs.KVL(
		"header_token", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__,
		"token", token,
	))
	return
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

// // pollingServiceCondition
// //  cluster_uuid = ? AND ((status = ? OR status = ? OR status = ?))
// func pollingServiceCondition(cluster_uuid string) vanilla.PrepareCondition {
// 	// 서비스의 polling 조건

// 	// return vanilla.And(
// 	// 	vanilla.Equal("cluster_uuid", cluster_uuid),
// 	// 	vanilla.Or(
// 	// 		vanilla.Equal("status", servicev3.StepStatusRegist),
// 	// 		vanilla.Equal("status", servicev3.StepStatusSend),
// 	// 		vanilla.Equal("status", servicev3.StepStatusProcessing),
// 	// 	))

// 	return vanilla.And(
// 		vanilla.Equal("cluster_uuid", cluster_uuid),
// 		vanilla.LessThan("status", servicev3.StepStatusSuccess),
// 	)
// }

// // pollingServiceLimit
// //  LIMIT ?
// func pollingServiceLimit(poliing_limit int) *vanilla.PreparePagination {
// 	if poliing_limit == 0 {
// 		poliing_limit = math.MaxInt8 // 127
// 	}

// 	return vanilla.Limit(poliing_limit)
// }

type PollingFilter = func(status servicev3.StepStatus) bool

func newPollingFilter(limit int) PollingFilter {
	limit = func(limit int) int {
		if limit == 0 {
			limit = math.MaxInt8 // 127
		}
		return limit + 1
	}(limit)

	return func(status servicev3.StepStatus) bool {
		if 0 < limit && status < servicev3.StepStatusSuccess {
			limit--
			return true
		}
		return false
	}
}

func pollingService(ctx context.Context, tx stmtex.Preparer, dialect string,
	cluster_uuid string, polling_offest vanilla.NullTime,
	polling_filter PollingFilter,
) (services []servicev3.Service, stepSet map[string][]servicev3.ServiceStep, err error) {

	// check polling
	polling_keys, err := vault.Servicev3.GetServicesPolling(ctx, tx, dialect, cluster_uuid, polling_offest)
	if err != nil {
		err = errors.Wrapf(err, "failed to found services%v", logs.KVL(
			"cluster_uuid", cluster_uuid,
			"polling_offest", polling_offest,
		))
		return
	}

	// filtering
	filtered_keys := make([]servicev3.Service_polling, 0, len(polling_keys))
	for _, service_key := range polling_keys {
		if polling_filter(service_key.Status) {
			filtered_keys = append(filtered_keys, service_key)
		}
	}

	// get polling service detail
	services = make([]servicev3.Service, 0, len(filtered_keys))
	for _, service_key := range filtered_keys {
		var service *servicev3.Service
		service, err = vault.Servicev3.GetService(ctx, tx, dialect, cluster_uuid, service_key.Uuid)
		if err != nil {
			err = errors.Wrapf(err, "failed to found service %v", logs.KVL(
				"cluster_uuid", cluster_uuid,
				"uuid", service_key.Uuid,
			))
			return
		}
		// append to
		services = append(services, *service)
	}

	// gather service step
	stepSet = make(map[string][]servicev3.ServiceStep)
	for _, service := range filtered_keys {
		var steps []servicev3.ServiceStep
		steps, err = vault.Servicev3.GetServiceSteps(ctx, tx, dialect, cluster_uuid, service.Uuid)
		if err != nil {
			err = errors.Wrapf(err, "failed to found service steps%v", logs.KVL(
				"cluster_uuid", cluster_uuid,
				"uuid", service.Uuid,
			))
			return
		}

		stepSet[service.Uuid] = steps

	}

	return
}
