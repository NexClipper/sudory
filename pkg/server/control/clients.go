package control

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/client/log"
	"github.com/NexClipper/sudory/pkg/server/event"
	"github.com/NexClipper/sudory/pkg/server/event/managed_channel"
	"github.com/NexClipper/sudory/pkg/server/event/managed_event"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/pkg/errors"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/prepare"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	authv2 "github.com/NexClipper/sudory/pkg/server/model/auth/v2"
	channelv2 "github.com/NexClipper/sudory/pkg/server/model/channel/v2"
	clusterv2 "github.com/NexClipper/sudory/pkg/server/model/cluster/v2"
	clustertokenv2 "github.com/NexClipper/sudory/pkg/server/model/cluster_token/v2"
	cryptov2 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v2"
	servicev2 "github.com/NexClipper/sudory/pkg/server/model/service/v2"
	sessionv2 "github.com/NexClipper/sudory/pkg/server/model/session/v2"
	"github.com/NexClipper/sudory/pkg/server/status/globvar"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/golang-jwt/jwt/v4"

	"github.com/labstack/echo/v4"
)

// @Description get []Service
// @Accept      json
// @Produce     json
// @Tags        client/service
// @Router      /client/service [get]
// @Param       x-sudory-client-token header string true "client session token"
// @Success     200 {array}  v2.HttpRsp_ClientServicePolling
// @Header      200 {string} x-sudory-client-token
func (ctl ControlVanilla) PollingService(ctx echo.Context) error {

	//get token claims
	claims, err := GetSudoryClisentTokenClaims(ctx)
	err = errors.Wrapf(err, "failed to get client token")
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	cluster := clusterv2.Cluster{}
	Do(&err, func() (err error) {
		// uuid = ?
		eq_uuid := vanilla.Equal("uuid", claims.ClusterUuid).Parse()

		err = vanilla.Stmt.Select(cluster.TableName(), cluster.ColumnNames(), eq_uuid, nil, nil).
			QueryRowContext(ctx.Request().Context(), ctl)(func(s vanilla.Scanner) (err error) {
			err = cluster.Scan(s)
			err = errors.Wrapf(err, "scan cluster")
			return
		})
		err = errors.Wrapf(err, "failed to get cluster")
		return
	})

	// find services
	var services []servicev2.HttpRsp_ClientServicePolling = []servicev2.HttpRsp_ClientServicePolling{}
	Do(&err, func() (err error) {
		condition := pollingServiceCondition(cluster.Uuid)
		limit := pollingServiceLimit(cluster.PoliingLimit)
		service_statuses := make([]servicev2.Service_status, 0, state.ENV__INIT_SLICE_CAPACITY__())

		service_status := servicev2.Service_status{}
		err = vanilla.Stmt.Select(service_status.TableName(), service_status.ColumnNames(), condition, nil, limit).
			QueryRowsContext(ctx.Request().Context(), ctl)(func(scan vanilla.Scanner, _ int) (err error) {
			err = service_status.Scan(scan)
			err = errors.Wrapf(err, "scan service_status")
			if err != nil {
				return
			}

			service_statuses = append(service_statuses, service_status)
			return
		})

		err = errors.Wrapf(err, "failed to find services")
		Do(&err, func() (err error) {
			services = make([]servicev2.HttpRsp_ClientServicePolling, len(service_statuses))
			for i := range service_statuses {
				services[i].Service = service_statuses[i].Service
				services[i].ServiceStatus_essential = service_statuses[i].ServiceStatus_essential
				services[i].Steps = make([]servicev2.ServiceStep_tangled, 0, services[i].Service.StepCount)
				service_uuid := services[i].Service.Uuid

				step_query := `
SELECT A.uuid, A.sequence, A.created, A.name, A.summary, A.method, A.args, A.result_filter,
       IFNULL(B.status, 0) as status, B.started, B.ended, B.created AS updated
  FROM service_step A
  LEFT JOIN service_step_status B
      INNER JOIN (
                 SELECT uuid, MAX(created) AS Max_created, sequence
                   FROM service_step_status
                  WHERE UUID = ?
                  GROUP BY UUID, sequence
            ) C
        ON B.uuid = C.uuid AND B.created = C.Max_created AND B.sequence = C.sequence
    ON A.uuid = B.uuid AND A.sequence = B.sequence
 WHERE A.uuid = ?
`

				step_args := []interface{}{
					service_uuid,
					service_uuid,
				}

				err = vanilla.QueryRowsContext(ctx.Request().Context(), ctl, step_query, step_args)(func(scan vanilla.Scanner, _ int) error {
					step := servicev2.ServiceStep_tangled{}
					err = step.Scan(scan)
					err = errors.Wrapf(err, "scan service_step")
					if err != nil {
						return err
					}

					services[i].Steps = append(services[i].Steps, step)
					return err
				})

				// eq_uuid := vanilla.Equal("uuid", service_uuid).Parse()
				// step := servicev2.ServiceStep_tangled{}
				// err = vanilla.Stmt.Select(step.TableName(), step.ColumnNames(), eq_uuid, nil, nil).
				// 	QueryRows(ctl)(func(scan vanilla.Scanner, _ int) (err error) {
				// 	err = step.Scan(scan)
				// 	err = errors.Wrapf(err, "scan service_step")
				// 	if err != nil {
				// 		return
				// 	}

				// 	services[i].Steps = append(services[i].Steps, step)
				// 	return
				// })

				err = errors.Wrapf(err, "failed to get steps%v", logs.KVL(
					"service_uuid", service_uuid,
				))
				if err != nil {
					return
				}
			}
			return
		})
		return
	})

	time_now := time.Now()

	Do(&err, func() (err error) {

		var service_statuses [][]interface{}
		var service_step_statuses [][]interface{}

		service_statuses = make([][]interface{}, 0, len(services))
		service_step_statuses = make([][]interface{}, 0, len(services))
		for i := range services {
			service := &services[i]

			if service.ServiceStatus_essential.Status == servicev2.StepStatusRegist {
				service_status := servicev2.ServiceStatus{}
				service_status.ServiceStatus_essential = service.ServiceStatus_essential
				//Uuid
				service_status.Uuid = service.Uuid
				//Created
				service_status.Created = time_now
				//AssignedClientUuid
				//할당된 클라이언트 정보 추가
				service_status.AssignedClientUuid = claims.Uuid
				service_status.Status = servicev2.StepStatusSend
				service_status.Message = *vanilla.NewNullString(service_status.Status.String())

				service_statuses = append(service_statuses, service_status.Values())

				service.ServiceStatus_essential = service_status.ServiceStatus_essential
			}

			for i := range service.Steps {
				step := &service.Steps[i]

				//Status
				//StatusSend 보다 작으면 응답 전 업데이트
				if step.ServiceStepStatus_essential.Status == servicev2.StepStatusRegist {
					step_status := servicev2.ServiceStepStatus{}
					step_status.ServiceStepStatus_essential = step.ServiceStepStatus_essential
					//Uuid
					step_status.Uuid = step.Uuid
					//Sequence
					step_status.Sequence = step.Sequence
					//Created
					step_status.Created = time_now
					step_status.Status = servicev2.StepStatusSend

					service_step_statuses = append(service_step_statuses, step_status.Values())

					step.ServiceStepStatus_essential = step_status.ServiceStepStatus_essential
				}
			}
		}

		err = ctl.Scope(func(tx *sql.Tx) (err error) {
			Do(&err, func() (err error) {
				if len(service_statuses) == 0 {
					return
				}

				service_status := servicev2.ServiceStatus{}
				builder, err := vanilla.Stmt.Insert(service_status.TableName(), service_status.ColumnNames(), service_statuses...)
				err = errors.Wrapf(err, "cannot build a service_status insert statement")
				if err != nil {
					return err
				}

				affected, err := builder.ExecContext(ctx.Request().Context(), tx)
				if affected == 0 {
					err = errors.Wrapf(err, "no affected")
				}

				err = errors.Wrapf(err, "faild to save service status")
				return
			})
			Do(&err, func() (err error) {
				if len(service_step_statuses) == 0 {
					return
				}

				service_step_status := servicev2.ServiceStepStatus{}
				builder, err := vanilla.Stmt.Insert(service_step_status.TableName(), service_step_status.ColumnNames(), service_step_statuses...)
				err = errors.Wrapf(err, "cannot build a service_step_status insert statement")
				if err != nil {
					return err
				}

				affected, err := builder.ExecContext(ctx.Request().Context(), tx)
				if affected == 0 {
					err = errors.Wrapf(err, "no affected")
				}

				err = errors.Wrapf(err, "faild to save service_step status")
				return
			})
			return
		})

		return
	})

	Do(&err, func() (err error) {
		//invoke event (service-poll-out)
		for _, service := range services {
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
			m["result_type"] = servicev2.ResultSaveTypeNone.String()
			m["result"] = ""
			m["step_count"] = service.StepCount
			m["step_position"] = service.StepPosition

			event.Invoke(event_name, m)
			managed_event.Invoke(event_name, m)
			// invoke event by event category
			managed_channel.InvokeByEventCategory(channelv2.EventCategoryServicePollingOut, m)
		}
		return
	})

	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, services)
}

// @Description update a service
// @Accept      json
// @Produce     json
// @Tags        client/service
// @Router      /client/service [put]
// @Param       x-sudory-client-token header string           true  "client session token"
// @Param       body body v2.HttpReq_ClientServiceUpdate true "HttpReq_ClientServiceUpdate"
// @Success     200
// @Header      200 {string} x-sudory-client-token
func (ctl ControlVanilla) UpdateService(ctx echo.Context) (err error) {
	body := servicev2.HttpReq_ClientServiceUpdate{}
	err = echoutil.Bind(ctx, &body)
	err = errors.Wrapf(err, "bind%s",
		logs.KVL(
			"type", TypeName(body),
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	claims, err := GetSudoryClisentTokenClaims(ctx)
	err = errors.Wrapf(err, "failed to get client token")

	//request check point
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	//get service
	service_query := `
SELECT A.uuid, A.created, A.name, A.summary, A.cluster_uuid, A.template_uuid, A.step_count, A.subscribed_channel,
       IFNULL(B.assigned_client_uuid, '') AS assigned_client_uuid, IFNULL(B.step_position, 0) AS step_position, IFNULL(B.status, 0) AS status, B.message, B.created AS updated
  FROM service A
  LEFT JOIN (
            SELECT C.uuid, C.created, C.assigned_client_uuid, C.step_position, C.status, C.message
              FROM service_status C
             WHERE C.uuid = ?
       ) B 
    ON A.uuid = B.uuid 
 WHERE A.uuid = ?
`

	service_args := []interface{}{
		body.Uuid,
		body.Uuid,
	}

	var service *servicev2.Service_status

	err = vanilla.QueryRowsContext(ctx.Request().Context(), ctl, service_query, service_args)(func(scan vanilla.Scanner, _ int) error {
		t := servicev2.Service_status{}
		err := t.Scan(scan)
		err = errors.Wrapf(err, "scan service")
		if err != nil {
			return err
		}

		if service == nil {
			service = &t
		}

		// replace latest service
		if service.Updated.Time.Before(t.Updated.Time) {
			service = &t
		}

		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "failed to found service")
	}

	if service == nil {
		return errors.New("cannot found service")
	}

	// service := servicev2.Service_tangled{}
	// // uuid = ?
	// eq_uuid := vanilla.Equal("uuid", body.Uuid).Parse()

	// err = vanilla.Stmt.Select(service.TableName(), service.ColumnNames(), eq_uuid, nil, nil).
	// 	QueryRow(ctl)(func(s vanilla.Scanner) (err error) {
	// 	err = service.Scan(s)
	// 	err = errors.Wrapf(err, "scan service")
	// 	return

	// })
	// if err != nil {
	// 	return err
	// }

	step_query := `
SELECT A.uuid, A.sequence, A.created, A.name, A.summary, A.method, A.args, A.result_filter,
       IFNULL(B.status, 0) AS status, B.started, B.ended, B.created AS updated
  FROM service_step A
  LEFT JOIN (
            SELECT C.uuid, C.sequence, C.status,
                   C.started, C.ended, C.created
              FROM service_step_status C
             WHERE C.uuid = ? AND C.sequence = ?
       ) B
    ON A.uuid = B.uuid
   AND A.sequence = B.sequence
 WHERE A.uuid = ? AND A.sequence = ?
`

	step_args := []interface{}{
		body.Uuid, body.Sequence,
		body.Uuid, body.Sequence,
	}

	var step *servicev2.ServiceStep_tangled

	err = vanilla.QueryRowsContext(ctx.Request().Context(), ctl, step_query, step_args)(func(scan vanilla.Scanner, _ int) error {
		t := servicev2.ServiceStep_tangled{}
		err := t.Scan(scan)
		err = errors.Wrapf(err, "scan service step")
		if err != nil {
			return err
		}

		if step == nil {
			step = &t
		}

		fmt.Println(step.Updated.Time)
		fmt.Println(t.Updated.Time)

		// replace latest service
		if step.Updated.Time.Before(t.Updated.Time) {
			step = &t
		}

		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "failed to found service step")
	}

	if step == nil {
		return errors.New("cannot found service step")
	}

	// step := servicev2.ServiceStep_tangled{}

	// // uuid = ? AND sequence = ?
	// unique_step := vanilla.And(
	// 	vanilla.Equal("uuid", body.Uuid),
	// 	vanilla.Equal("sequence", body.Sequence),
	// ).Parse()

	// err = vanilla.Stmt.Select(step.TableName(), step.ColumnNames(), unique_step, nil, nil).
	// 	QueryRow(ctl.DB)(func(s vanilla.Scanner) (err error) {
	// 	err = step.Scan(s)
	// 	err = errors.Wrapf(err, "scan service_step")
	// 	return
	// })
	// if err != nil {
	// 	return err
	// }

	stepPosition := func() int {
		// 스템 포지션 값은
		// ServiceStep.Sequence+1
		return body.Sequence + 1
	}
	stepStatus := func() servicev2.StepStatus {
		// 스탭 포지션이 카운트와 같은 경우만
		if service.StepCount == stepPosition() {
			return body.Status
		}
		// 기본값; 처리중(Processing)
		return servicev2.StepStatusProcessing
	}
	serviceResult := func() cryptov2.CryptoString {
		// 상태가 성공인 경우만
		if body.Status == servicev2.StepStatusSuccess {
			return cryptov2.CryptoString(body.Result)
		}
		//기본값; 공백 문자열
		return ""
	}
	stepMessage := func() vanilla.NullString {
		// 상태가 실패인 경우만
		if body.Status == servicev2.StepStatusFail {
			return *vanilla.NewNullString(body.Result)
		}
		//기본값; 공백 문자열
		return vanilla.NullString{}
	}
	// eventMessage := func() vanilla.NullString {
	// 	// 상태가 실패인 경우만
	// 	if body.Status == servicev2.StepStatusFail {
	// 		return *vanilla.NewNullString(body.Result)
	// 	}
	// 	//기본값; 공백 문자열
	// 	return *vanilla.NewNullString(body.Status.String())
	// }
	resultType := func() (resultType servicev2.ResultSaveType) {
		// 마지막 스탭의 결과만 저장 한다
		if service.StepCount != stepPosition() {
			return
		}
		// 상태 값이 성공이 아닌 경우
		// 서비스 결과를 저장 하지 않는다
		if servicev2.StepStatusSuccess != stepStatus() {
			return
		}
		// 채널이 등록되어 있는 경우
		// 서비스 결과를 저장 하지 않는다
		if !service.SubscribedChannel.Valid || 0 < len(service.SubscribedChannel.String) {
			return
		}

		return servicev2.ResultSaveTypeDatabase
	}

	time_now := time.Now()

	// service status
	service_status := func() *servicev2.ServiceStatus {

		service_status := new(servicev2.ServiceStatus)
		service_status.Uuid = service.Uuid
		service_status.Created = time_now
		service_status.AssignedClientUuid = claims.Uuid
		service_status.StepPosition = stepPosition()
		service_status.Status = stepStatus()
		service_status.Message = stepMessage()
		return service_status
	}()
	// service result
	service_result := func() *servicev2.ServiceResult {

		service_result := new(servicev2.ServiceResult)
		service_result.Uuid = service.Uuid
		service_result.Created = time_now
		service_result.Result = serviceResult()
		service_result.ResultSaveType = resultType()
		return service_result
	}()
	// step status
	step_status := func() *servicev2.ServiceStepStatus {

		step_status := new(servicev2.ServiceStepStatus)
		step_status.Uuid = step.Uuid
		step_status.Sequence = step.Sequence // missing
		step_status.Created = time_now
		step_status.Status = body.Status                         // Status
		step_status.Started = *vanilla.NewNullTime(body.Started) // Started
		step_status.Ended = *vanilla.NewNullTime(body.Ended)     // Ended
		return step_status
	}()

	save_service_status := func(tx *sql.Tx) (err error) {
		// 서비스 상태 저장
		builder, err := vanilla.Stmt.Insert(service_status.TableName(), service_status.ColumnNames(), service_status.Values())
		err = errors.Wrapf(err, "cannot build a service_status insert statement")
		if err != nil {
			return err
		}

		affected, err := builder.ExecContext(ctx.Request().Context(), tx)
		err = errors.Wrapf(err, "exec insert statement")
		if err != nil {
			return err
		}
		if affected == 0 {
			return errors.New("no affected")
		}

		return
	}
	save_service_result := func(tx *sql.Tx) (err error) {
		if service_result.ResultSaveType != servicev2.ResultSaveTypeDatabase {
			return
		}

		// 서비스 결과 저장
		builder, err := vanilla.Stmt.Insert(service_result.TableName(), service_result.ColumnNames(), service_result.Values())
		err = errors.Wrapf(err, "cannot build a service_result insert statement")
		if err != nil {
			return err
		}

		affected, err := builder.ExecContext(ctx.Request().Context(), tx)
		err = errors.Wrapf(err, "exec insert statement")
		if err != nil {
			return err
		}
		if affected == 0 {
			return errors.New("no affected")
		}

		return
	}
	save_service_step_status := func(tx *sql.Tx) (err error) {
		// 서비스 스탭 저장
		builder, err := vanilla.Stmt.Insert(step_status.TableName(), step_status.ColumnNames(), step_status.Values())
		err = errors.Wrapf(err, "cannot build a service_step_status insert statement")
		if err != nil {
			return err
		}

		affected, err := builder.ExecContext(ctx.Request().Context(), tx)
		err = errors.Wrapf(err, "exec insert statement")
		if err != nil {
			return err
		}
		if affected == 0 {
			return errors.New("no affected")
		}

		return
	}

	//save status
	err = ctl.Scope(func(tx *sql.Tx) (err error) {
		if err = save_service_status(tx); err != nil {
			return errors.Wrapf(err, "failed to save service_status")
		}
		if err = save_service_result(tx); err != nil {
			return errors.Wrapf(err, "failed to save service_result")
		}
		if err = save_service_step_status(tx); err != nil {
			return errors.Wrapf(err, "failed to save service_step_status")
		}
		return
	})

	err = errors.Wrapf(err, "failed to save updated service and service_step status")
	if err != nil {
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
	m["assigned_client_uuid"] = service_status.AssignedClientUuid
	m["status"] = service_status.Status
	m["status_description"] = service_status.Status.String()
	m["result_type"] = service_result.ResultSaveType.String()
	m["result"] = body.Result
	m["step_count"] = service.StepCount
	m["step_position"] = service_status.StepPosition

	if service_status.Status == servicev2.StepStatusSuccess && len(service_result.Result.String()) == 0 {
		log.Debugf("channel(poll-in-service): %+v", m)
	}

	event.Invoke(service.SubscribedChannel.String, m)         //Subscribe 등록된 구독 이벤트 이름으로 호출
	managed_event.Invoke(service.SubscribedChannel.String, m) //Subscribe 등록된 구독 이벤트 이름으로 호출
	// invoke event by channel uuid
	if 0 < len(service.SubscribedChannel.String) {
		// find channel
		mc := channelv2.ManagedChannel{}
		mc.Uuid = service.SubscribedChannel.String
		mc_cond := vanilla.And(
			vanilla.Equal("uuid", mc.Uuid),
			vanilla.IsNull("deleted"),
		)
		found, err := vanilla.Stmt.Exist(mc.TableName(), mc_cond.Parse())(ctx.Request().Context(), ctl)
		if err != nil {
			return err
		}
		if found {
			managed_channel.InvokeByChannelUuid(service.SubscribedChannel.String, m)
		}
	}
	// invoke event by event category
	managed_channel.InvokeByEventCategory(channelv2.EventCategoryServicePollingIn, m)

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

	cluster := clusterv2.Cluster{}
	cluster.Uuid = auth.ClusterUuid
	cluster_eq_uuid := vanilla.Equal("uuid", cluster.Uuid)
	cluster_found, err := vanilla.Stmt.Exist(cluster.TableName(), cluster_eq_uuid.Parse())(ctx.Request().Context(), ctl)
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
	token := clustertokenv2.ClusterToken{}
	token.ClusterUuid = auth.ClusterUuid
	token.Token = cryptov2.CryptoString(auth.Assertion)

	token_cond := vanilla.And(
		vanilla.Equal("cluster_uuid", token.ClusterUuid),
		vanilla.Equal("token", token.Token),
	)

	err = vanilla.Stmt.Select(token.TableName(), token.ColumnNames(), token_cond.Parse(), nil, nil).
		QueryRow(ctl)(func(scan vanilla.Scanner) (err error) {
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
	exp := globvar.ClientSessionExpirationTime(iat)

	payload := &sessionv2.ClientSessionPayload{
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

	session := sessionv2.Session{}
	session.Uuid = payload.Uuid
	session.ClusterUuid = payload.ClusterUuid
	session.Token = token_string
	session.IssuedAtTime = *vanilla.NewNullTime(time.Unix(payload.IssuedAt, 0))
	session.ExpirationTime = *vanilla.NewNullTime(time.Unix(payload.ExpiresAt, 0))
	session.Created = created

	err = func() (err error) {
		stmt, err := vanilla.Stmt.Insert(session.TableName(), session.ColumnNames(), session.Values())
		if err != nil {
			return err
		}
		affected, err := stmt.Exec(ctl)
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
	managed_channel.InvokeByEventCategory(channelv2.EventCategoryClientAuthAccept, m)

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
			return errors.Wrapf(ErrorInvalidRequestParameter(), "missing request header%v", logs.KVL(
				"header", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__,
			))
		}

		claims := new(sessionv2.ClientSessionPayload)
		jwt_token, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(globvar.ClientSessionSignatureSecret()), nil
		})
		if err != nil {
			return errors.Wrapf(err, "jwt parse claims")
		}

		if _, ok := jwt_token.Claims.(*sessionv2.ClientSessionPayload); !ok || !jwt_token.Valid {
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
	claims := new(sessionv2.ClientSessionPayload)
	err = func() (err error) {
		if len(token) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter(), "missing request header%s",
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

		if _, ok := jwt_token.Claims.(*sessionv2.ClientSessionPayload); !ok {
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
	cluster := clusterv2.Cluster{}
	err = func() (err error) {
		eq_uuid := vanilla.Equal("uuid", claims.ClusterUuid).Parse()
		err = vanilla.Stmt.Select(cluster.TableName(), cluster.ColumnNames(), eq_uuid, nil, nil).
			QueryRowContext(ctx.Request().Context(), ctl)(func(s vanilla.Scanner) (err error) {
			err = cluster.Scan(s)
			err = errors.Wrapf(err, "cluster Scan")
			return
		})

		err = errors.Wrapf(err, "failed to get cluster")
		return
	}()
	if err != nil {
		return HttpError(err, http.StatusInternalServerError) // StatusInternalServerError
	}

	var service_count int
	err = func() (err error) {
		eq_uuid := pollingServiceCondition(claims.ClusterUuid)
		service := servicev2.Service_status{}
		service_count, err = vanilla.Stmt.Count(service.TableName(), eq_uuid, nil)(ctx.Request().Context(), ctl)
		if err != nil {
			return errors.Wrapf(err, "failed to get cluster service count")
		}

		return
	}()
	if err != nil {
		return HttpError(err, http.StatusInternalServerError) // StatusInternalServerError
	}

	// polling interval 해더 저장
	polling_interval := int(cluster.GetPollingOption().Interval(time.Duration(int64(globvar.ClientConfigPollInterval())*int64(time.Second)), service_count) / time.Second)

	//reflesh payload
	claims.PollInterval = polling_interval
	claims.ExpiresAt = globvar.ClientSessionExpirationTime(time_now).Unix()
	claims.Loglevel = globvar.ClientConfigLoglevel()

	var new_token_string string
	err = func() (err error) {
		//client auth 에서 사용된 알고리즘 그대로 사용
		new_token_string, err = jwt.NewWithClaims(usedJwtSigningMethod(*jwt_token, jwt.SigningMethodHS256), claims).
			SignedString([]byte(globvar.ClientSessionSignatureSecret()))
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

	session := sessionv2.Session{}
	session.Uuid = claims.Uuid
	session.Token = new_token_string
	session.ExpirationTime = *vanilla.NewNullTime(time.Unix(claims.ExpiresAt, 0))
	session.Updated = *vanilla.NewNullTime(time_now)

	// uuid = ? AND deleted IS NULL
	cond_session := vanilla.And(
		vanilla.Equal("uuid", session.Uuid),
		vanilla.IsNull("deleted"),
	)
	err = func() (err error) {
		session_found, err := vanilla.Stmt.Exist(session.TableName(), cond_session.Parse())(ctx.Request().Context(), ctl)
		if err != nil {
			return err
		}
		if !session_found {
			return errors.Errorf("not found session%v", logs.KVL(
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

		_, err = vanilla.Stmt.Update(session.TableName(), keys_values, cond_session.Parse()).
			Exec(ctl)
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

func GetSudoryClisentTokenClaims(ctx echo.Context) (claims *sessionv2.ClientSessionPayload, err error) {
	token := ctx.Request().Header.Get(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__)
	if len(token) == 0 {
		err = errors.Errorf("missing request header%s",
			logs.KVL(
				"key", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__,
			))
	}

	claims = new(sessionv2.ClientSessionPayload)
	if true {
		Do(&err, func() (err error) {
			var jwt_token *jwt.Token
			// JWT unverify
			jwt_token, _, err = jwt.NewParser().ParseUnverified(token, claims)
			err = errors.Wrapf(err, "jwt.Parser.ParseUnverified")
			Do(&err, func() (err error) {
				var ok bool
				claims, ok = jwt_token.Claims.(*sessionv2.ClientSessionPayload)
				if !ok {
					err = errors.New("jwt.Token.Claims not matched")
				}
				return
			})
			return
		})
	} else {
		Do(&err, func() (err error) {
			var jwt_token *jwt.Token
			// JWT verify
			jwt_token, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(globvar.ClientSessionSignatureSecret()), nil
			})
			err = errors.Wrapf(err, "jwt.Parser.ParseWithClaims")
			Do(&err, func() (err error) {
				var ok bool
				claims, ok = jwt_token.Claims.(*sessionv2.ClientSessionPayload)
				if !ok {
					err = errors.New("jwt.Token.Claims not matched")
				}
				if !jwt_token.Valid {
					err = errors.New("jwt.Token.Valid false")
				}
				return
			})
			return
		})
	}

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

// pollingServiceCondition
//  cluster_uuid = ? AND ((status = ? OR status = ? OR status = ?))
func pollingServiceCondition(cluster_uuid string) *prepare.Condition {
	// condition := vanilla.NewCond(
	// 	// 서비스의 polling 조건
	// 	"WHERE cluster_uuid = ? AND ((status = ? OR status = ? OR status = ?))",
	// 	cluster_uuid,
	// 	servicev2.StepStatusRegist,
	// 	servicev2.StepStatusSend,
	// 	servicev2.StepStatusProcessing,
	// )

	if false {
		return vanilla.And(
			vanilla.Equal("cluster_uuid", cluster_uuid),
			vanilla.Or(
				vanilla.Equal("status", servicev2.StepStatusRegist),
				vanilla.Equal("status", servicev2.StepStatusSend),
				vanilla.Equal("status", servicev2.StepStatusProcessing),
			)).Parse()
	}

	return vanilla.And(
		vanilla.Equal("cluster_uuid", cluster_uuid),
		vanilla.LessThan("status", servicev2.StepStatusSuccess),
	).Parse()
}

// pollingServiceLimit
//  LIMIT ?
func pollingServiceLimit(poliing_limit int) *prepare.Pagination {
	if poliing_limit == 0 {
		poliing_limit = math.MaxInt8 // 127
	}

	return vanilla.Limit(poliing_limit).Parse()
}
