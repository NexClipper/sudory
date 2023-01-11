package control

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/NexClipper/sudory/pkg/client/log"
	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/event/managed_channel"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/pkg/errors"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/sqlex"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"

	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/model/auths/v2"
	channelv3 "github.com/NexClipper/sudory/pkg/server/model/channel/v3"
	clusterv3 "github.com/NexClipper/sudory/pkg/server/model/cluster/v3"
	clusterinfov2 "github.com/NexClipper/sudory/pkg/server/model/cluster_infomation/v2"
	"github.com/NexClipper/sudory/pkg/server/model/cluster_token/v3"
	cryptov2 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v2"
	servicev3 "github.com/NexClipper/sudory/pkg/server/model/service/v3"
	servicev4 "github.com/NexClipper/sudory/pkg/server/model/service/v4"
	sessions "github.com/NexClipper/sudory/pkg/server/model/session/v3"
	"github.com/NexClipper/sudory/pkg/server/model/tenants/v3"
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
// @Success     200 {array}  servicev4.HttpRsp_ClientServicePolling
// @Header      200 {string} x-sudory-client-token
func (ctl ControlVanilla) PollingService(ctx echo.Context) error {

	//get token claims
	// claims, err := GetSudoryClientTokenClaims(ctx)
	claims, err := GetClientSessionClaims(ctx, ctl, ctl.dialect)
	err = errors.Wrapf(err, "failed to get client token")
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// get cluster
	cluster := new(clusterv3.Cluster)
	cluster.Uuid = claims.ClusterUuid
	eq_uuid := stmt.Equal("uuid", cluster.Uuid)

	err = ctl.dialect.QueryRow(cluster.TableName(), cluster.ColumnNames(), eq_uuid, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := cluster.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "failed to get cluster")
	}

	// gather service
	// get service offset
	// var polling_offset vanilla.NullTime
	clusterInfo := func(fn func(clusterinfov2.ClusterInformation) (int, time.Time, []servicev4.HttpRsp_ClientServicePolling, error)) (cnt int, v []servicev4.HttpRsp_ClientServicePolling, err error) {
		var cluster_info_count int
		var cluster_info clusterinfov2.ClusterInformation
		cond := stmt.And(
			stmt.Equal("cluster_uuid", claims.ClusterUuid),
			stmt.IsNull("deleted"),
		)

		err = ctl.dialect.QueryRows(cluster_info.TableName(), cluster_info.ColumnNames(), cond, nil, nil)(ctx.Request().Context(), ctl)(
			func(scan excute.Scanner, i int) error {
				cluster_info_count = i + 1
				err = cluster_info.Scan(scan)
				err = errors.WithStack(err)

				return err
			})
		if err != nil {
			err = errors.Wrapf(err, "failed to get service offset")

			return
		}

		// build polling data
		cnt, offset, v, err := fn(cluster_info)
		if err != nil {
			return
		}

		err = func() (err error) {
			// save polling_count to cluster_infomation
			cluster_info := clusterinfov2.ClusterInformation{}
			cluster_info.ClusterUuid = cluster.Uuid
			cluster_info.PollingCount = *vanilla.NewNullInt(cnt)
			cluster_info.PollingOffset = *vanilla.NewNullTime(offset)
			cluster_info.Created = time.Now()
			cluster_info.Updated = *vanilla.NewNullTime(cluster_info.Created)

			switch cluster_info_count {
			case 0:
				affected, _, err := ctl.dialect.Insert(cluster_info.TableName(), cluster_info.ColumnNames(), cluster_info.Values())(
					ctx.Request().Context(), ctl)
				if err != nil {
					return errors.Wrapf(err, "insert")
				}
				if affected == 0 {
					err := errors.New("no affected")
					return err
				}

				return nil
			default:
				keys_values := map[string]interface{}{
					"cluster_uuid":   cluster_info.ClusterUuid,
					"polling_count":  cluster_info.PollingCount,
					"polling_offset": cluster_info.PollingOffset,
					"updated":        cluster_info.Updated,
				}
				if !cluster_info.PollingOffset.Valid {
					delete(keys_values, "polling_offset")
				}

				cond := stmt.And(
					stmt.Equal("cluster_uuid", cluster_info.ClusterUuid),
					stmt.IsNull("deleted"),
				)

				affected, err := ctl.dialect.Update(cluster_info.TableName(), keys_values, cond)(
					ctx.Request().Context(), ctl)
				if err != nil {
					err := errors.Wrapf(err, "update")
					return err
				}

				if 1 < affected {
					err := errors.Wrapf(err, "too many affected")

					return err
				}

				return nil
			}
		}()
		if err != nil {
			err = errors.Wrapf(err, "save cluster_information")

			return
		}

		return
	}

	makePollingServicev3 := func(cluster_info clusterinfov2.ClusterInformation) (polling_cnt int, polling_offset time.Time, rsp_ []servicev4.HttpRsp_ClientServicePolling, err error) {

		// 오프셋이 서비스 유효시간 보다 작은 경우
		// 혹은 오프셋이 없는 경우
		// 오프셋 시간을 서비스 유효시간으로 설정
		timelimit := globvar.ClientConfig.ServiceValidTimeLimit()
		ltime := time.Now().
			Truncate(time.Second).
			Add(time.Duration(timelimit) * time.Minute * -1)

		if !cluster_info.PollingOffset.Valid {
			cluster_info.PollingOffset = *vanilla.NewNullTime(ltime)
		}

		if ltime.After(cluster_info.PollingOffset.Time) {
			cluster_info.PollingOffset = *vanilla.NewNullTime(ltime)
		}

		// polling limit filter
		polling_filter := newPollingFilter(cluster.PoliingLimit, timelimit, ltime)
		servs, steps, err := pollingService(ctx.Request().Context(), ctl, ctl.dialect, claims.ClusterUuid, cluster_info.PollingOffset, polling_filter)

		// set polling_offest
		var polling_offest_ time.Time
		for _, serv := range servs {
			if polling_offest_.IsZero() {
				polling_offest_ = serv.Created
			}

			if polling_offest_.After(serv.Created) {
				polling_offest_ = serv.Created
			}
		}

		polling_cnt = len(servs)
		polling_offset = polling_offest_

		UpdateServiceStatus := func(serv servicev3.Service, assigned_client_uuid string, status servicev3.StepStatus, t time.Time) servicev3.Service {
			serv.AssignedClientUuid = *vanilla.NewNullString(assigned_client_uuid)
			serv.Status = status
			serv.Timestamp = t
			return serv
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
			servs []servicev3.Service, CALLBACK_Service CALLBACK_Service,
			steps map[string][]servicev3.ServiceStep, CALLBACK_ServiceStep CALLBACK_ServiceStep) {

			for _, serv := range servs {
				if serv.Status == servicev3.StepStatusRegist {
					v := UpdateServiceStatus(serv, claims.Uuid, servicev3.StepStatusSend, t)
					CALLBACK_Service(v)

					for _, step := range steps[v.Uuid] {
						v := UpdateStepStatus(step, servicev3.StepStatusSend, t)
						CALLBACK_ServiceStep(v)
					}
				}
			}
		}

		var new_service_status = make([]vault.Table, 0, len(servs))
		var new_step_status = make([]vault.Table, 0, len(steps))

		time_now := time.Now()
		MakeUpdate(time_now,
			servs, func(a servicev3.Service) {
				new_service_status = append(new_service_status, a)
			},
			steps, func(a servicev3.ServiceStep) {
				new_step_status = append(new_step_status, a)
			})

		err = sqlex.ScopeTx(ctx.Request().Context(), ctl, func(tx *sql.Tx) error {

			if err = vault.SaveMultiTable(tx, ctl.dialect, new_service_status); err != nil {
				return errors.Wrapf(err, "faild to save service")
			}

			if err = vault.SaveMultiTable(tx, ctl.dialect, new_step_status); err != nil {
				return errors.Wrapf(err, "faild to save service step")
			}

			return nil
		})
		if err != nil {
			return
		}

		// make response body
		rsp := make([]servicev3.HttpRsp_ClientServicePolling, len(servs))
		for i, serv := range servs {
			rsp[i].Service = serv
			rsp[i].Steps = steps[serv.Uuid]
		}

		// invoke event (service-poll-out)
		var mm = make([]map[string]interface{}, len(rsp))
		// invoke event (service-poll-out)
		const event_name = "service-poll-out"
		for i, service := range rsp {
			mm[i] = map[string]interface{}{}
			mm[i]["event_name"] = event_name
			mm[i]["service_uuid"] = service.Uuid
			mm[i]["service_name"] = service.Name
			mm[i]["template_uuid"] = service.TemplateUuid
			mm[i]["cluster_uuid"] = service.ClusterUuid
			mm[i]["assigned_client_uuid"] = service.AssignedClientUuid
			mm[i]["status"] = service.Status
			mm[i]["status_description"] = service.Status.String()
			mm[i]["result_type"] = servicev3.ResultSaveTypeNone.String()
			mm[i]["result"] = ""
			mm[i]["step_count"] = service.StepCount
			mm[i]["step_position"] = service.StepPosition
		}

		if 0 < len(mm) {
			// invoke event by event category
			managed_channel.InvokeByEventCategory(claims.TenantHash, channelv3.EventCategoryServicePollingOut, mm)
		}

		rsp_ = make([]servicev4.HttpRsp_ClientServicePolling, len(servs))
		for i := range rsp {
			rsp_[i].Version = "v3"
			rsp_[i].V3 = rsp[i]
		}

		return
	}

	makePollingServicev4 := func(cluster_info clusterinfov2.ClusterInformation) (polling_cnt int, polling_offset time.Time, rsp []servicev4.HttpRsp_ClientServicePolling, err error) {

		// 오프셋이 서비스 유효시간 보다 작은 경우
		// 혹은 오프셋이 없는 경우
		// 오프셋 시간을 서비스 유효시간으로 설정
		timelimit := globvar.ClientConfig.ServiceValidTimeLimit()
		ltime := time.Now().
			Truncate(time.Second).
			Add(time.Duration(timelimit) * time.Minute * -1)

		if !cluster_info.PollingOffset.Valid {
			cluster_info.PollingOffset = *vanilla.NewNullTime(ltime)
		}

		if ltime.After(cluster_info.PollingOffset.Time) {
			cluster_info.PollingOffset = *vanilla.NewNullTime(ltime)
		}

		// polling limit filter
		polling_filter := newPollingFilter_v2(cluster.PoliingLimit, timelimit, ltime)
		servs, stats, err := pollingService_v2(ctx.Request().Context(), ctl, ctl.dialect, claims.ClusterUuid, cluster_info.PollingOffset, polling_filter)

		// set polling_offest
		var polling_offest_ time.Time
		for _, serv := range servs {
			if polling_offest_.IsZero() {
				polling_offest_ = serv.Created
			}

			if polling_offest_.After(serv.Created) {
				polling_offest_ = serv.Created
			}
		}

		polling_cnt = len(servs)
		polling_offset = polling_offest_

		UpdateStepStatus := func(step servicev4.ServiceStatus, status servicev4.StepStatus, t time.Time) servicev4.ServiceStatus {
			step.Status = status
			step.Created = t
			return step
		}

		time_now := time.Now()
		var new_step_status = make([]vault.Table, 0, len(stats))
		for i := range stats {
			if stats[i].Status == servicev4.StepStatusRegist {
				v := UpdateStepStatus(stats[i], servicev4.StepStatusSent, time_now)

				new_step_status = append(new_step_status, v)
			}
		}

		err = sqlex.ScopeTx(ctx.Request().Context(), ctl, func(tx *sql.Tx) error {
			if err = vault.SaveMultiTable(tx, ctl.dialect, new_step_status); err != nil {
				err = errors.Wrapf(err, "faild to save a service status")
				return err
			}

			return nil
		})
		if err != nil {
			return
		}

		// make response body
		rsp = make([]servicev4.HttpRsp_ClientServicePolling, len(servs))
		for i := range servs {
			rsp[i].Version = "v4"
			rsp[i].V4 = servs[i]
		}

		// // get tenent by cluster_uuid
		// var tenant tenants.Tenant
		// tenant_table := clusterv3.TenantTableName(claims.ClusterUuid)
		// tenant_cond := stmt.And(
		// 	stmt.IsNull("deleted"),
		// )
		// err = ctl.dialect.QueryRows(tenant_table, tenant.ColumnNames(), tenant_cond, nil, nil)(
		// 	ctx.Request().Context(), ctl)(
		// 	func(scan excute.Scanner, _ int) error {
		// 		err := tenant.Scan(scan)
		// 		err = errors.WithStack(err)

		// 		return err
		// 	})
		// if err != nil {
		// 	err = errors.Wrapf(err, "failed to get a tenent by cluster_uuid")
		// 	return
		// }

		// invoke event (service-poll-out)
		const event_name = "service-poll-out"
		var mm = make([]map[string]interface{}, len(rsp))
		for i := range servs {
			serv := servs[i]
			stat := stats[i]

			mm[i] = map[string]interface{}{}
			mm[i]["event_name"] = event_name
			mm[i]["service_uuid"] = serv.Uuid
			mm[i]["service_name"] = serv.Name
			mm[i]["template_uuid"] = serv.TemplateUuid
			mm[i]["cluster_uuid"] = serv.ClusterUuid
			mm[i]["assigned_client_uuid"] = claims.Uuid
			mm[i]["status"] = stat.Status
			mm[i]["status_description"] = stat.Status.String()
			mm[i]["result_type"] = servicev3.ResultSaveTypeNone.String()
			mm[i]["result"] = ""
			mm[i]["step_count"] = stat.StepMax
			mm[i]["step_position"] = stat.StepSeq
		}

		if 0 < len(mm) {
			// invoke event by event category
			managed_channel.InvokeByEventCategory(claims.TenantHash, channelv3.EventCategoryServicePollingOut, mm)
		}

		return
	}

	// cnt, rsp, err := clusterInfo(makePollingServicev3)
	// if err != nil {
	// 	return err
	// }
	// if 0 < cnt {
	// 	return ctx.JSON(http.StatusOK, rsp)
	// }
	// _, rsp, err = clusterInfo(makePollingServicev4)
	// if err != nil {
	// 	return err
	// }

	cnt, rsp, err := clusterInfo(makePollingServicev4)
	if err != nil {
		return err
	}
	if 0 < cnt {
		return ctx.JSON(http.StatusOK, rsp)
	}
	_, rsp, err = clusterInfo(makePollingServicev3)
	if err != nil {
		return err
	}

	rspv3 := make([]servicev3.HttpRsp_ClientServicePolling, len(rsp))
	for i, s := range rsp {
		rspv3[i] = s.V3
	}

	return ctx.JSON(http.StatusOK, rspv3)
}

// @Description update a service
// @Security    ClientAuth
// @Accept      json
// @Produce     json
// @Tags        client/service
// @Router      /client/service [put]
// @Param       body body servicev4.HttpReq_ClientServiceUpdate true "HttpReq_ClientServiceUpdate"
// @Success     200
// @Header      200 {string} x-sudory-client-token
func (ctl ControlVanilla) UpdateService(ctx echo.Context) (err error) {
	type temp struct {
		Version string `json:"version,omitempty"`
	}
	bodyTemp := temp{}
	if err := echoutil.Bind(ctx, &bodyTemp); err != nil {
		err = errors.Wrapf(err, "bind%s",
			logs.KVL(
				"type", TypeName(bodyTemp),
			))
		return HttpError(err, http.StatusBadRequest)
	}

	body := servicev4.HttpReq_ClientServiceUpdate{}

	if len(bodyTemp.Version) > 0 {
		if err := echoutil.Bind(ctx, &body); err != nil {
			err = errors.Wrapf(err, "bind%s",
				logs.KVL(
					"type", TypeName(body),
				))
			return HttpError(err, http.StatusBadRequest)
		}
	} else {
		body.Version = "v3"
		if err := echoutil.Bind(ctx, &body.V3); err != nil {
			err = errors.Wrapf(err, "bind%s",
				logs.KVL(
					"type", TypeName(body),
				))
			return HttpError(err, http.StatusBadRequest)
		}
	}

	updateServiceStatus_v3 := func(body servicev3.HttpReq_ClientServiceUpdate) error {

		stepPosition := func(service_step servicev3.ServiceStep) int {
			// 스탭 포지션 값은
			// ServiceStep.Sequence+1
			return service_step.Sequence + 1
		}
		stepStatus := func(body servicev3.HttpReq_ClientServiceUpdate, serv servicev3.Service, service_step servicev3.ServiceStep) servicev3.StepStatus {
			// 스탭 포지션이 카운트와 같은 경우만
			if serv.StepCount == stepPosition(service_step) {
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
		resultType := func(body servicev3.HttpReq_ClientServiceUpdate, serv servicev3.Service, service_step servicev3.ServiceStep) (resultType servicev3.ResultSaveType) {
			// 마지막 스탭의 결과만 저장 한다
			if serv.StepCount != stepPosition(service_step) {
				return
			}
			// 상태 값이 성공이 아닌 경우
			// 서비스 결과를 저장 하지 않는다
			if servicev3.StepStatusSuccess != stepStatus(body, serv, service_step) {
				return
			}
			// 채널이 등록되어 있는 경우
			// 서비스 결과를 저장 하지 않는다
			if !serv.SubscribedChannel.Valid || 0 < len(serv.SubscribedChannel.String) {
				return
			}

			return servicev3.ResultSaveTypeDatabase
		}

		if err := body.Status.Valid(); err != nil {
			err = errors.Wrapf(err, "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.Status", TypeName(body)), body.Status)...,
				))
			return HttpError(err, http.StatusBadRequest)
		}

		// get client token claims
		// claims, err := GetSudoryClientTokenClaims(ctx)
		claims, err := GetClientSessionClaims(ctx, ctl.DB, ctl.dialect)
		if err != nil {
			err = errors.Wrapf(err, "failed to get client token")
			return HttpError(err, http.StatusBadRequest)
		}

		// get serv
		serv := new(servicev3.Service)
		serv.ClusterUuid = claims.ClusterUuid
		serv.Uuid = body.Uuid

		serv, err = vault.GetService(context.Background(), ctl, ctl.dialect, serv.ClusterUuid, serv.Uuid)
		if err != nil {
			return errors.Wrapf(err, "failed to found service%v", logs.KVL(
				"cluster_uuid", claims.ClusterUuid,
				"uuid", body.Uuid,
			))
		}

		// get service step
		step := new(servicev3.ServiceStep)
		step.ClusterUuid = claims.ClusterUuid
		step.Uuid = body.Uuid
		step.Sequence = body.Sequence

		step, err = vault.GetServiceStep(context.Background(), ctl, ctl.dialect, step.ClusterUuid, step.Uuid, step.Sequence)
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
			serv.Timestamp = now_time
			// update value
			serv.AssignedClientUuid = *vanilla.NewNullString(claims.Uuid)
			serv.StepPosition = stepPosition(*step)
			serv.Status = stepStatus(body, *serv, *step)
			serv.Message = stepMessage(body)
		}
		// new service result
		result := func() *servicev3.ServiceResult {
			// update key
			result := new(servicev3.ServiceResult)
			result.PartitionDate = serv.PartitionDate
			result.ClusterUuid = serv.ClusterUuid
			result.Uuid = serv.Uuid
			result.Timestamp = now_time
			// update value
			result.ResultSaveType = resultType(body, *serv, *step)
			result.Result = serviceResult(body)

			return result
		}()
		// udpate service step
		{
			// update key
			step.Timestamp = now_time
			// update value
			step.Status = body.Status                         // Status
			step.Started = *vanilla.NewNullTime(body.Started) // Started
			step.Ended = *vanilla.NewNullTime(body.Ended)     // Ended
		}

		// save to db
		err = sqlex.ScopeTx(context.Background(), ctl, func(tx *sql.Tx) (err error) {

			// save service
			if err = vault.SaveMultiTable(tx, ctl.dialect, []vault.Table{serv}); err != nil {
				return errors.Wrapf(err, "failed to save service_status")
			}

			// save service step
			if err = vault.SaveMultiTable(tx, ctl.dialect, []vault.Table{step}); err != nil {
				return errors.Wrapf(err, "failed to save service_step_status")
			}

			// check servcie result save type
			if result.ResultSaveType != servicev3.ResultSaveTypeNone {
				// save service result
				if err = vault.SaveMultiTable(tx, ctl.dialect, []vault.Table{result}); err != nil {
					return errors.Wrapf(err, "failed to save service_result")
				}
			}

			return
		})
		err = errors.Wrapf(err, "failed to save")
		if err != nil {
			return err
		}

		// // get tenent by cluster_uuid
		// var tenant tenants.Tenant
		// tenant_table := clusterv3.TenantTableName(claims.ClusterUuid)
		// tenant_cond := stmt.And(
		// 	stmt.IsNull("deleted"),
		// )
		// err = ctl.dialect.QueryRows(tenant_table, tenant.ColumnNames(), tenant_cond, nil, nil)(
		// 	context.Background(), ctl.DB)(
		// 	func(scan excute.Scanner, _ int) error {
		// 		err := tenant.Scan(scan)
		// 		err = errors.WithStack(err)

		// 		return err
		// 	})
		// if err != nil {
		// 	return errors.Wrapf(err, "failed to get a tenent by cluster_uuid")
		// }

		// invoke event (service-poll-in)
		const event_name = "service-poll-in"
		mm := make([]map[string]interface{}, 1)
		mm[0] = map[string]interface{}{}
		mm[0]["event_name"] = event_name
		mm[0]["service_uuid"] = serv.Uuid
		mm[0]["service_name"] = serv.Name
		mm[0]["template_uuid"] = serv.TemplateUuid
		mm[0]["cluster_uuid"] = serv.ClusterUuid
		mm[0]["assigned_client_uuid"] = serv.AssignedClientUuid
		mm[0]["status"] = serv.Status
		mm[0]["status_description"] = serv.Status.String()
		mm[0]["result_type"] = result.ResultSaveType.String()
		mm[0]["result"] = func() interface{} {
			var b = []byte(strings.TrimSpace(body.Result))
			if len(b) == 0 {
				// empty string
				return map[string]interface{}{
					"message": body.Result,
				}
			}

			head, tail := b[0], b[len(b)-1]
			if ok := (head == '{' && tail == '}') || (head == '[' && tail == ']'); !ok {
				// not json
				return map[string]interface{}{
					"message": body.Result,
				}
			}

			// json
			return json.RawMessage(body.Result)
		}()
		mm[0]["step_count"] = serv.StepCount
		mm[0]["step_position"] = serv.StepPosition

		if serv.Status == servicev3.StepStatusSuccess && len(result.Result.String()) == 0 {
			log.Debugf("channel(poll-in-service): %+v", mm)
		}

		// invoke event by channel uuid
		if 0 < len(serv.SubscribedChannel.String) {
			// find channel
			channel := channelv3.ManagedChannel{}
			channel.Uuid = serv.SubscribedChannel.String
			channel_cond := stmt.And(
				stmt.Equal("uuid", channel.Uuid),
				stmt.IsNull("deleted"),
			)
			channel_table := channelv3.TableNameWithTenant_ManagedChannel(claims.TenantHash)
			found, err := ctl.dialect.Exist(channel_table, channel_cond)(context.Background(), ctl)
			if err != nil {
				return err
			}
			if found {
				managed_channel.InvokeByChannelUuid(claims.TenantHash, serv.SubscribedChannel.String, mm)
			}
		}

		// invoke event by event category
		managed_channel.InvokeByEventCategory(claims.TenantHash, channelv3.EventCategoryServicePollingIn, mm)

		return nil
	}

	updateServiceStatus_v4 := func(body servicev4.HttpReq_ClientServiceUpdate_multistep) error {
		if err := body.Status.Valid(); err != nil {
			err = errors.Wrapf(err, "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.Status", TypeName(body)), body.Status)...,
				))
			return HttpError(err, http.StatusBadRequest)
		}

		// get client token claims
		// claims, err := GetSudoryClientTokenClaims(ctx)
		claims, err := GetClientSessionClaims(ctx, ctl.DB, ctl.dialect)
		if err != nil {
			err = errors.Wrapf(err, "failed to get client token")
			return HttpError(err, http.StatusBadRequest)
		}

		// get serv
		var serv servicev4.Service
		serv.ClusterUuid = claims.ClusterUuid
		serv.Uuid = body.Uuid

		serv_cond := stmt.And(
			stmt.Equal("cluster_uuid", serv.ClusterUuid),
			stmt.Equal("uuid", serv.Uuid))

		err = ctl.dialect.QueryRow(serv.TableName(), serv.ColumnNames(), serv_cond, nil, nil)(
			ctx.Request().Context(), ctl)(
			func(scan excute.Scanner) error {
				err = serv.Scan(scan)
				err = errors.WithStack(err)

				return err
			})
		if err != nil {
			return errors.Wrapf(err, "failed to found a service%v", logs.KVL(
				"cluster_uuid", claims.ClusterUuid,
				"uuid", body.Uuid,
			))
		}

		now_time := time.Now()

		newStatus := servicev4.ServiceStatus{}
		newStatus.PartitionDate = serv.PartitionDate
		newStatus.ClusterUuid = serv.ClusterUuid
		newStatus.Uuid = serv.Uuid
		newStatus.Created = now_time
		newStatus.StepMax = serv.StepMax
		newStatus.StepSeq = body.Sequence + 1
		newStatus.Status = body.Status
		newStatus.Started = *vanilla.NewNullTime(body.Started)
		newStatus.Ended = *vanilla.NewNullTime(body.Ended)
		if body.Status == servicev4.StepStatusFailed {
			newStatus.Message = *vanilla.NewNullString(body.Result)
		}

		newResult := servicev4.ServiceResult{}
		newResult.PartitionDate = serv.PartitionDate
		newResult.ClusterUuid = serv.ClusterUuid
		newResult.Uuid = serv.Uuid
		// 채널이 등록되어 있는 경우
		// 서비스 결과를 저장 하지 않는다
		if !serv.SubscribedChannel.Valid || 0 < len(serv.SubscribedChannel.String) {
			newResult.ResultSaveType = servicev4.ResultSaveTypeDatabase
		}
		if body.Status == servicev4.StepStatusSucceeded {
			newResult.Result = cryptov2.CryptoString(body.Result)
		}
		newResult.Created = now_time

		// save to db
		err = sqlex.ScopeTx(context.Background(), ctl, func(tx *sql.Tx) (err error) {
			// save service step
			_, _, err = ctl.dialect.Insert(newStatus.TableName(), newStatus.ColumnNames(), newStatus.Values())(
				context.Background(), ctl)
			if err != nil {
				err = errors.Wrapf(err, "failed to save a status")
				return err
			}

			if body.Status == servicev4.StepStatusSucceeded {
				updateColumn := []string{"result_type", "result", "created"}
				_, _, err = ctl.dialect.InsertOrUpdate(newResult.TableName(), newResult.ColumnNames(), updateColumn, newResult.Values())(
					context.Background(), ctl)
				if err != nil {
					err = errors.Wrapf(err, "failed to save a result")
					return
				}
			}

			return
		})
		if err != nil {
			err = errors.Wrapf(err, "failed to save")
			return err
		}

		// // get tenent by cluster_uuid
		// var tenant tenants.Tenant
		// tenant_table := clusterv3.TenantTableName(claims.ClusterUuid)
		// tenant_cond := stmt.And(
		// 	stmt.IsNull("deleted"),
		// )
		// err = ctl.dialect.QueryRows(tenant_table, tenant.ColumnNames(), tenant_cond, nil, nil)(
		// 	context.Background(), ctl.DB)(
		// 	func(scan excute.Scanner, _ int) error {
		// 		err := tenant.Scan(scan)
		// 		err = errors.WithStack(err)

		// 		return err
		// 	})
		// if err != nil {
		// 	return errors.Wrapf(err, "failed to get a tenent by cluster_uuid")
		// }

		// invoke event (service-poll-in)
		const event_name = "service-poll-in"
		mm := make([]map[string]interface{}, 1)
		mm[0] = map[string]interface{}{}
		mm[0]["event_name"] = event_name
		mm[0]["service_uuid"] = serv.Uuid
		mm[0]["service_name"] = serv.Name
		mm[0]["template_uuid"] = serv.TemplateUuid
		mm[0]["cluster_uuid"] = serv.ClusterUuid
		mm[0]["assigned_client_uuid"] = claims.Uuid
		mm[0]["status"] = newStatus.Status
		mm[0]["status_description"] = newStatus.Status.String()
		mm[0]["result_type"] = newResult.ResultSaveType.String()
		mm[0]["result"] = func() interface{} {

			var b = []byte(strings.TrimSpace(body.Result))
			if len(b) == 0 {
				// empty string
				return map[string]interface{}{
					"message": body.Result,
				}
			}

			head, tail := b[0], b[len(b)-1]
			if ok := (head == '{' && tail == '}') || (head == '[' && tail == ']'); !ok {
				// not json
				return map[string]interface{}{
					"message": body.Result,
				}
			}

			// json
			return json.RawMessage(body.Result)
		}()
		mm[0]["step_count"] = newStatus.StepMax
		mm[0]["step_position"] = newStatus.StepSeq

		if newStatus.Status == servicev4.StepStatusSucceeded && len(newResult.Result.String()) == 0 {
			log.Debugf("channel(poll-in-service): %+v", mm)
		}

		// invoke event by channel uuid
		if 0 < len(serv.SubscribedChannel.String) {
			// find channel
			channel := channelv3.ManagedChannel{}
			channel.Uuid = serv.SubscribedChannel.String
			channel_cond := stmt.And(
				stmt.Equal("uuid", channel.Uuid),
				stmt.IsNull("deleted"),
			)
			channel_table := channelv3.TableNameWithTenant_ManagedChannel(claims.TenantHash)
			found, err := ctl.dialect.Exist(channel_table, channel_cond)(context.Background(), ctl)
			if err != nil {
				return err
			}
			if found {
				managed_channel.InvokeByChannelUuid(claims.TenantHash, serv.SubscribedChannel.String, mm)
			}
		}
		// invoke event by event category
		managed_channel.InvokeByEventCategory(claims.TenantHash, channelv3.EventCategoryServicePollingIn, mm)

		return nil
	}

	switch body.Version {
	case "v3":
		updateServiceStatus_v3(body.V3)
	case "v4":
		updateServiceStatus_v4(body.V4)
	default:
		err := errors.New("invalid version")
		return err
	}

	return ctx.JSON(http.StatusOK, OK())
}

// @Description auth client
// @Accept      json
// @Produce     json
// @Tags        client/auth
// @Router      /client/auth [post]
// @Param       body body auths.HttpReqAuth true "HttpReqAuth"
// @Success     200 {string} ok
// @Header      200 {string} x-sudory-client-token
func (ctl ControlVanilla) AuthClient(ctx echo.Context) (err error) {
	auth := new(auths.HttpReqAuth)
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
	cluster_found, err := ctl.dialect.Exist(cluster.TableName(), cluster_eq_uuid)(
		ctx.Request().Context(), ctl)
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

	// get tenent by cluster_uuid
	var tenant tenants.Tenant
	tenant_table := clusterv3.TenantTableName(cluster.Uuid)
	tenant_cond := stmt.And(
		stmt.IsNull("deleted"),
	)

	err = ctl.dialect.QueryRows(tenant_table, tenant.ColumnNames(), tenant_cond, nil, nil)(
		ctx.Request().Context(), ctl.DB)(
		func(scan excute.Scanner, _ int) error {
			err := tenant.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "failed to get a tenent by cluster_uuid")
	}

	//valid token
	token := cluster_token.ClusterToken{}
	token.ClusterUuid = auth.ClusterUuid
	token.Token = cryptov2.CryptoString(auth.Assertion)

	token_cond := stmt.And(
		stmt.Equal("cluster_uuid", token.ClusterUuid),
		stmt.Equal("token", token.Token),
	)

	err = ctl.dialect.QueryRow(token.TableName(), token.ColumnNames(), token_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := token.Scan(scan)
			err = errors.WithStack(err)

			return err
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

	payload := &sessions.ClientSessionPayload{
		ExpiresAt:    exp.Unix(),
		IssuedAt:     iat.Unix(),
		Uuid:         session_token_uuid,
		ClusterUuid:  auth.ClusterUuid,
		TenantHash:   tenant.Hash,
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

	session := sessions.Session{}
	session.Uuid = payload.Uuid
	session.ClusterUuid = payload.ClusterUuid
	session.Token = token_string
	session.IssuedAtTime = *vanilla.NewNullTime(time.Unix(payload.IssuedAt, 0))
	session.ExpirationTime = *vanilla.NewNullTime(time.Unix(payload.ExpiresAt, 0))
	session.Created = created

	err = func() (err error) {
		var affected int64
		affected, session.ID, err = ctl.dialect.Insert(session.TableName(), session.ColumnNames(), session.Values())(
			ctx.Request().Context(), ctl)
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
	mm := make([]map[string]interface{}, 1)
	mm[0] = map[string]interface{}{}
	mm[0]["event_name"] = event_name
	mm[0]["cluster_uuid"] = payload.ClusterUuid
	mm[0]["session_uuid"] = payload.Uuid

	// invoke event by event category
	managed_channel.InvokeByEventCategory(tenant.Hash, channelv3.EventCategoryClientAuthAccept, mm)

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

		claims := new(sessions.ClientSessionPayload)
		jwt_token, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(globvar.ClientSession.SignatureSecret()), nil
		})
		if err != nil {
			return errors.Wrapf(err, "jwt parse claims")
		}

		if _, ok := jwt_token.Claims.(*sessions.ClientSessionPayload); !ok || !jwt_token.Valid {
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
	claims := new(sessions.ClientSessionPayload)
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

		if _, ok := jwt_token.Claims.(*sessions.ClientSessionPayload); !ok {
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

		err = ctl.dialect.QueryRow(cluster.TableName(), cluster.ColumnNames(), eq_uuid, nil, nil)(
			ctx.Request().Context(), ctl)(
			func(s excute.Scanner) error {
				err := cluster.Scan(s)
				err = errors.WithStack(err)

				return err
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

		err = ctl.dialect.QueryRows(cluster_info.TableName(), columnnames, cond, nil, nil)(
			ctx.Request().Context(), ctl)(
			func(scan excute.Scanner, _ int) error {
				err := scan.Scan(&service_count)
				err = errors.WithStack(err)

				return err
			})

		// cond := pollingServiceCondition(claims.ClusterUuid)
		// service := service.Service{}
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

	session := sessions.Session{}
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
		session_found, err := ctl.dialect.Exist(session.TableName(), cond_session)(
			ctx.Request().Context(), ctl)
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

		_, err = ctl.dialect.Update(session.TableName(), keys_values, cond_session)(ctx.Request().Context(), ctl)
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

// func GetSudoryClientTokenClaims(ctx echo.Context) (claims *sessionv3.ClientSessionPayload, err error) {
// 	var token string
// 	if token = ctx.Request().Header.Get(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__); len(token) == 0 {
// 		err = errors.Errorf("missing request header%s",
// 			logs.KVL(
// 				"key", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__,
// 			))
// 		return
// 	}

// 	claims = new(sessionv3.ClientSessionPayload)
// 	var fn func() error
// 	fn = func() error {
// 		var jwt_token *jwt.Token
// 		// JWT unverify
// 		jwt_token, _, err = jwt.NewParser().ParseUnverified(token, claims)
// 		if err != nil {
// 			return errors.Wrapf(err, "jwt.Parser.ParseUnverified")
// 		}

// 		var ok bool
// 		claims, ok = jwt_token.Claims.(*sessionv3.ClientSessionPayload)
// 		if !ok {
// 			return errors.New("jwt.Token.Claims not matched")
// 		}

// 		return nil
// 	}
// 	if false {
// 		fn = func() error {
// 			var jwt_token *jwt.Token
// 			// JWT verify
// 			jwt_token, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
// 				return []byte(globvar.ClientSession.SignatureSecret()), nil
// 			})
// 			if err != nil {
// 				return errors.Wrapf(err, "jwt.Parser.ParseWithClaims")
// 			}

// 			var ok bool
// 			claims, ok = jwt_token.Claims.(*sessionv3.ClientSessionPayload)
// 			if !ok {
// 				return errors.New("jwt.Token.Claims not matched")
// 			}
// 			if !jwt_token.Valid {
// 				return errors.New("jwt.Token.Valid false")
// 			}

// 			return nil
// 		}
// 	}

// 	err = fn()
// 	err = errors.Wrapf(err, "failed to parse header token%v", logs.KVL(
// 		"header_token", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__,
// 		"token", token,
// 	))
// 	return
// }

// setCookie
//
//lint:ignore U1000 auto-generated
func setCookie(ctx echo.Context, key, value string, exp time.Duration) {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.Value = value
	cookie.Expires = time.Now().Add(exp)
	ctx.SetCookie(cookie)
}

// setCookie
//
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
// 	// 		vanilla.Equal("status", service.StepStatusRegist),
// 	// 		vanilla.Equal("status", service.StepStatusSend),
// 	// 		vanilla.Equal("status", service.StepStatusProcessing),
// 	// 	))

// 	return vanilla.And(
// 		vanilla.Equal("cluster_uuid", cluster_uuid),
// 		vanilla.LessThan("status", service.StepStatusSuccess),
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

type PollingFilter = func(service servicev3.Service_polling) bool

func newPollingFilter(limit int, timelimit int, ltime time.Time) PollingFilter {
	limit = func(limit int) int {
		if limit == 0 {
			limit = math.MaxInt8 // 127
		}
		return limit + 1
	}(limit)

	return func(serv servicev3.Service_polling) bool {
		status := serv.Status
		created := serv.Created

		if !(0 < limit) {
			return false
		}

		if !(status < servicev3.StepStatusSuccess) {
			return false
		}

		if !created.After(ltime) && 0 < timelimit {
			return false
		}

		limit--
		return true
	}
}

type PollingFilter_v2 = func(service vault.ServicePolling) bool

func newPollingFilter_v2(limit int, timelimit int, ltime time.Time) PollingFilter_v2 {
	limit = func(limit int) int {
		if limit == 0 {
			limit = math.MaxInt8 // 127
		}
		return limit + 1
	}(limit)

	return func(service vault.ServicePolling) bool {
		status := service.LastStatus.Status
		created := service.Service.Created

		if !(0 < limit) {
			return false
		}

		if !(status < servicev4.StepStatusSucceeded) {
			return false
		}

		// println(ltime.String())
		// println(created.String())

		if !created.After(ltime) && 0 < timelimit {
			return false
		}

		limit--
		return true
	}
}

func pollingService(ctx context.Context, tx excute.Preparer, dialect excute.SqlExcutor,
	cluster_uuid string, polling_offest vanilla.NullTime,
	polling_filter PollingFilter,
) (servs []servicev3.Service, stepSet map[string][]servicev3.ServiceStep, err error) {

	// check polling
	polling_keys, err := vault.GetServicesPolling(ctx, tx, dialect, cluster_uuid, polling_offest)
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
		if polling_filter(service_key) {
			filtered_keys = append(filtered_keys, service_key)
		}
	}

	// get polling service detail
	servs = make([]servicev3.Service, 0, len(filtered_keys))
	for _, service_key := range filtered_keys {
		var serv *servicev3.Service
		serv, err = vault.GetService(ctx, tx, dialect, cluster_uuid, service_key.Uuid)
		if err != nil {
			err = errors.Wrapf(err, "failed to found service %v", logs.KVL(
				"cluster_uuid", cluster_uuid,
				"uuid", service_key.Uuid,
			))
			return
		}
		// append to
		servs = append(servs, *serv)
	}

	// gather service step
	stepSet = make(map[string][]servicev3.ServiceStep)
	for _, serv := range filtered_keys {
		var steps []servicev3.ServiceStep
		steps, err = vault.GetServiceSteps(ctx, tx, dialect, cluster_uuid, serv.Uuid)
		if err != nil {
			err = errors.Wrapf(err, "failed to found service steps%v", logs.KVL(
				"cluster_uuid", cluster_uuid,
				"uuid", serv.Uuid,
			))
			return
		}

		stepSet[serv.Uuid] = steps

	}

	return
}

func pollingService_v2(ctx context.Context, tx excute.Preparer, dialect excute.SqlExcutor,
	cluster_uuid string, polling_offest vanilla.NullTime,
	polling_filter PollingFilter_v2,
) (servs []servicev4.Service, stats []servicev4.ServiceStatus, err error) {

	// check pollingServices
	pollingServices, err := vault.GetServicesPolling_v4(ctx, tx, dialect, cluster_uuid, polling_offest)
	if err != nil {
		err = errors.Wrapf(err, "failed to found services%v", logs.KVL(
			"cluster_uuid", cluster_uuid,
			"polling_offest", polling_offest,
		))
		return
	}

	// filtering
	servs = make([]servicev4.Service, 0, len(pollingServices))
	stats = make([]servicev4.ServiceStatus, 0, len(pollingServices))
	for i := range pollingServices {
		if polling_filter(pollingServices[i]) {
			servs = append(servs, pollingServices[i].Service)
			stats = append(stats, pollingServices[i].LastStatus)
		}
	}

	return
}
