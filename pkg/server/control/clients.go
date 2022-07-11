package control

import (
	"database/sql"
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/event"
	"github.com/NexClipper/sudory/pkg/server/event/managed_event"
	"github.com/pkg/errors"
	"xorm.io/xorm"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/prepare"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	authv1 "github.com/NexClipper/sudory/pkg/server/model/auth/v1"
	clusterv1 "github.com/NexClipper/sudory/pkg/server/model/cluster/v1"
	clusterv2 "github.com/NexClipper/sudory/pkg/server/model/cluster/v2"
	clustertokenv1 "github.com/NexClipper/sudory/pkg/server/model/cluster_token/v1"
	crypto "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v2"
	servicev2 "github.com/NexClipper/sudory/pkg/server/model/service/v2"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
	sessionv2 "github.com/NexClipper/sudory/pkg/server/model/session/v2"
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
			QueryRow(ctl.DB())(func(s vanilla.Scanner) (err error) {
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
		service_statuses := make([]servicev2.Service_status, 0, __INIT_RECORD_CAPACITY__)

		service_status := servicev2.Service_status{}
		err = vanilla.Stmt.Select(service_status.TableName(), service_status.ColumnNames(), condition, nil, limit).
			QueryRows(ctl.DB())(func(scan vanilla.Scanner, _ int) (err error) {
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

				eq_uuid := vanilla.Equal("uuid", service_uuid).Parse()
				step := servicev2.ServiceStep_tangled{}
				err = vanilla.Stmt.Select(step.TableName(), step.ColumnNames(), eq_uuid, nil, nil).
					QueryRows(ctl.DB())(func(scan vanilla.Scanner, _ int) (err error) {
					err = step.Scan(scan)
					err = errors.Wrapf(err, "scan service_step")
					if err != nil {
						return
					}

					services[i].Steps = append(services[i].Steps, step)
					return
				})

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
			service := services[i]

			if service.ServiceStatus_essential.Status == servicev2.StepStatusRegist {
				service_status := servicev2.ServiceStatus{}
				service_status.ServiceStatus_essential = service.ServiceStatus_essential
				//Uuid
				service_status.Uuid = service.Service.Uuid
				//Created
				service_status.Created = time_now
				//AssignedClientUuid
				//할당된 클라이언트 정보 추가
				service_status.AssignedClientUuid = claims.Uuid
				service_status.Status = servicev2.StepStatusSend

				service_statuses = append(service_statuses, service_status.Values())
			}

			for i := range service.Steps {
				step := service.Steps[i]

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
				err = errors.Wrapf(err, "can not build a service_status insert statement")
				if err != nil {
					return err
				}

				affected, err := builder.Exec(tx)
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
				err = errors.Wrapf(err, "can not build a service_step_status insert statement")
				if err != nil {
					return err
				}

				affected, err := builder.Exec(tx)
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
			// if 0 < len(service.Result) {
			// 	m["result_type"] = service.ResultType.String()
			// 	m["result"] = service.Result
			// }
			m["result"] = ""
			m["step_count"] = service.StepCount
			m["step_position"] = service.StepPosition

			logs.KVL(
				"channel(poll-out-service_uuid)", m["service_uuid"],
				"channel(poll-out-status)", m["status"],
				"channel(poll-out-result-length)", 0,
			)

			event.Invoke(event_name, m)
			managed_event.Invoke(event_name, m)
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
	Do(&err, func() (err error) {
		err = echoutil.Bind(ctx, &body)
		err = errors.Wrapf(err, "bind%s",
			logs.KVL(
				"type", TypeName(body),
			))
		return
	})

	claims, err := GetSudoryClisentTokenClaims(ctx)
	err = errors.Wrapf(err, "failed to get client token")

	//request check point
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	//get service
	service := servicev2.Service_tangled{}
	Do(&err, func() (err error) {
		// uuid = ?
		eq_uuid := vanilla.Equal("uuid", body.Uuid).Parse()

		err = vanilla.Stmt.Select(service.TableName(), service.ColumnNames(), eq_uuid, nil, nil).
			QueryRow(ctl.DB())(func(s vanilla.Scanner) (err error) {
			err = service.Scan(s)
			err = errors.Wrapf(err, "scan service")
			return

		})
		return
	})

	step := servicev2.ServiceStep_tangled{}
	Do(&err, func() (err error) {

		// uuid = ? AND sequence = ?
		unique_step := vanilla.And(
			vanilla.Equal("uuid", body.Uuid),
			vanilla.Equal("sequence", body.Sequence),
		).Parse()

		err = vanilla.Stmt.Select(step.TableName(), step.ColumnNames(), unique_step, nil, nil).
			QueryRow(ctl.DB())(func(s vanilla.Scanner) (err error) {
			err = step.Scan(s)
			err = errors.Wrapf(err, "scan service_step")
			return
		})

		return
	})

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
	serviceResult := func() crypto.CryptoString {
		// 상태가 실패인 경우만
		if body.Status == servicev2.StepStatusSuccess {
			return (crypto.CryptoString)(body.Result)
		}
		//기본값; 공백 문자열
		return ""
	}
	stepMessage := func() vanilla.NullString {
		// 상태가 실패인 경우만
		if body.Status == servicev2.StepStatusFail {
			return (vanilla.NullString)(body.Result)
		}
		//기본값; 공백 문자열
		return ""
	}

	time_now := time.Now()

	// service status
	service_status := func() servicev2.ServiceStatus {
		service.AssignedClientUuid = claims.Uuid
		service.StepPosition = stepPosition()
		service.Status = stepStatus()
		service.Message = stepMessage()
		return servicev2.ServiceStatus{
			Uuid:                    service.Uuid,
			Created:                 time_now,
			ServiceStatus_essential: service.ServiceStatus_essential,
		}
	}()
	// service result
	service_result := func() servicev2.ServiceResult {
		service.ServiceResults_essential.ResultType = servicev2.ResultTypeDatabase //default
		service.ServiceResults_essential.Result = serviceResult()
		return servicev2.ServiceResult{
			Uuid:                     service.Uuid,
			Created:                  time_now,
			ServiceResults_essential: service.ServiceResults_essential,
		}
	}()
	// step status
	step_status := func() servicev2.ServiceStepStatus {
		step.ServiceStepStatus_essential.Status = body.Status                       // Status
		step.ServiceStepStatus_essential.Started = (vanilla.NullTime)(body.Started) // Started
		step.ServiceStepStatus_essential.Ended = (vanilla.NullTime)(body.Ended)     // Ended
		return servicev2.ServiceStepStatus{
			Uuid:                        step.Uuid,
			Sequence:                    step.Sequence, // missing
			Created:                     time_now,
			ServiceStepStatus_essential: step.ServiceStepStatus_essential,
		}
	}()

	//save status
	Do(&err, func() (err error) {
		err = ctl.Scope(func(tx *sql.Tx) (err error) {
			Do(&err, func() (err error) {
				// 서비스 상태 저장
				builder, err := vanilla.Stmt.Insert(service_status.TableName(), service_status.ColumnNames(), service_status.Values())
				err = errors.Wrapf(err, "can not build a service_status insert statement")
				if err != nil {
					return
				}

				affected, err := builder.Exec(tx)
				if affected == 0 {
					err = errors.Wrapf(err, "no affected")
				}

				err = errors.Wrapf(err, "faild to save service_status")
				return
			})
			Do(&err, func() (err error) {
				// 마지막 스탭의 결과만 저장 한다
				if service.StepCount != stepPosition() {
					return
				}
				// 상태 값이 성공이 아닌 경우
				// 서비스 결과를 저장 하지 않는다
				if service_status.Status != servicev2.StepStatusSuccess {
					return
				}
				// 채널이 등록되어 있는 경우
				// 서비스 결과를 저장 하지 않는다
				if 0 < len(service.SubscribedChannel) {
					return
				}

				// 서비스 결과 저장
				builder, err := vanilla.Stmt.Insert(service_result.TableName(), service_result.ColumnNames(), service_result.Values())
				err = errors.Wrapf(err, "can not build a service_result insert statement")
				if err != nil {
					return
				}

				affected, err := builder.Exec(tx)
				if affected == 0 {
					err = errors.Wrapf(err, "no affected")
				}

				err = errors.Wrapf(err, "faild to save service_result")
				return
			})
			Do(&err, func() (err error) {
				// 서비스 스탭 저장
				builder, err := vanilla.Stmt.Insert(step_status.TableName(), step_status.ColumnNames(), step_status.Values())
				err = errors.Wrapf(err, "can not build a service_step_status insert statement")
				if err != nil {
					return
				}

				affected, err := builder.Exec(tx)
				if affected == 0 {
					err = errors.Wrapf(err, "no affected")
				}

				err = errors.Wrapf(err, "faild to save service_step_status")
				return
			})

			return
		})

		err = errors.Wrapf(err, "failed to save updated service and service_step status")
		return
	})

	//invoke event (service-poll-in)
	Do(&err, func() (err error) {
		const event_name = "service-poll-in"
		m := map[string]interface{}{}
		m["event_name"] = event_name
		m["service_uuid"] = service.Uuid
		m["service_name"] = service.Name
		m["template_uuid"] = service.TemplateUuid
		m["cluster_uuid"] = service.ClusterUuid
		m["assigned_client_uuid"] = service_status.AssignedClientUuid
		m["status"] = service_status.Status
		// if 0 < len(service_result.Result) {
		m["result_type"] = service_result.ResultType.String()
		m["result"] = service_result.Result.String()
		// }
		// if 0 < len(service_status.Message) {
		m["message"] = service_status.Message.String()
		// }
		m["step_count"] = service.StepCount
		m["step_position"] = service_status.StepPosition

		logs.KVL(
			"channel(poll-in-service_uuid)", m["service_uuid"],
			"channel(poll-in-status)", m["status"],
			"channel(poll-in-result-length)", len(service_result.Result.String()),
		)

		event.Invoke(service.SubscribedChannel.String(), m)         //Subscribe 등록된 구독 이벤트 이름으로 호출
		managed_event.Invoke(service.SubscribedChannel.String(), m) //Subscribe 등록된 구독 이벤트 이름으로 호출

		return
	})
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

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
	session_token_uuid := NewUuidString()
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

	return nil
}

func (ctl ControlVanilla) RefreshClientSessionToken(ctx echo.Context) (err error) {
	token := ctx.Request().Header.Get(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__)
	Do(&err, func() (err error) {
		if len(token) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter(), "missing request header%s",
				logs.KVL(
					"key", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__,
				))
		}
		return
	})

	var jwt_token *jwt.Token
	claims := new(sessionv1.ClientSessionPayload)
	Do(&err, func() (err error) {
		jwt_token, _, err = jwt.NewParser().ParseUnverified(token, claims)
		err = errors.Wrapf(err, "failed to jwt.ParseUnverified%v", logs.KVL(
			"token", token,
		))
		if _, ok := jwt_token.Claims.(*sessionv1.ClientSessionPayload); !ok {
			return ErrorCompose(err, errors.Wrapf(err, "is not valid type%v",
				logs.KVL(
					"jwt.Token.Claims", TypeName(jwt_token.Claims),
					"jwt.Token.Method.Alg", jwt_token.Method.Alg(),
					"token", token,
				)))
		}
		return
	})

	if err != nil {
		return HttpError(err, http.StatusBadRequest) // StatusBadRequest
	}

	time_now := time.Now()

	// polling interval
	cluster := clusterv2.Cluster{}
	Do(&err, func() (err error) {
		eq_uuid := vanilla.Equal("uuid", claims.ClusterUuid).Parse()
		err = vanilla.Stmt.Select(cluster.TableName(), cluster.ColumnNames(), eq_uuid, nil, nil).
			QueryRow(ctl.DB())(func(s vanilla.Scanner) (err error) {
			err = cluster.Scan(s)
			err = errors.Wrapf(err, "cluster Scan")
			return
		})

		err = errors.Wrapf(err, "failed to get cluster")
		return
	})

	var service_count int
	Do(&err, func() (err error) {
		eq_uuid := pollingServiceCondition(claims.ClusterUuid)
		service := servicev2.Service_tangled{}
		service_count, err = vanilla.Count(ctl.DB(), service.TableName(), eq_uuid)
		err = errors.Wrapf(err, "failed to get cluster service count")
		return
	})

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
	Do(&err, func() (err error) {
		//client auth 에서 사용된 알고리즘 그대로 사용
		new_token_string, err = jwt.NewWithClaims(usedJwtSigningMethod(*jwt_token, jwt.SigningMethodHS256), claims).
			SignedString([]byte(globvar.ClientSessionSignatureSecret()))
		if err != nil {
			return errors.Wrapf(err, "failed to make session token to formed jwt")
		}
		return
	})

	Do(&err, func() (err error) {
		//save client session-token to header
		ctx.Response().Header().Set(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__, new_token_string)
		return
	})

	session := sessionv2.Session{}
	// uuid = ? AND deleted IS NULL
	condition_update_session := vanilla.And(
		vanilla.Equal("uuid", claims.Uuid),
		vanilla.IsNull("deleted"),
	).Parse()
	var affected int64
	Do(&err, func() (err error) {
		err = ctl.Scope(func(tx *sql.Tx) (err error) {
			keys_values := map[string]interface{}{
				"token":           new_token_string,
				"expiration_time": time.Unix(claims.ExpiresAt, 0),
				"updated":         time_now,
			}

			affected, err = vanilla.Stmt.Update(session.TableName(), keys_values, condition_update_session).
				Exec(tx)
			err = errors.Wrapf(err, "failed to update client session for refresh client session%v", logs.KVL(
				"uuid", claims.Uuid,
				"data", keys_values,
			))
			return
		})

		return
	})

	Do(&err, func() (err error) {
		// check client session record
		if 0 < affected {
			// exists record
			return
		}

		columns := []string{
			"COUNT(1)",
		}
		var count int
		err = vanilla.Stmt.Select(session.TableName(), columns, condition_update_session, nil, nil).
			QueryRow(ctl.DB())(func(s vanilla.Scanner) (err error) {
			err = s.Scan(&count)
			err = errors.Wrapf(err, "session Scan")
			return
		})
		err = errors.Wrapf(err, "not found session record%v", logs.KVL(
			"claims_uuid", claims.Uuid,
		))
		if count == 0 {
			err = ErrorCompose(err, errors.New("no affected"))
		}
		return
	})

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

func GetSudoryClisentTokenClaims(ctx echo.Context) (claims *sessionv1.ClientSessionPayload, err error) {
	token := ctx.Request().Header.Get(__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__)
	if len(token) == 0 {
		err = errors.Errorf("missing request header%s",
			logs.KVL(
				"key", __HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__,
			))
	}

	claims = new(sessionv1.ClientSessionPayload)
	if true {
		Do(&err, func() (err error) {
			var jwt_token *jwt.Token
			// JWT unverify
			jwt_token, _, err = jwt.NewParser().ParseUnverified(token, claims)
			err = errors.Wrapf(err, "jwt.Parser.ParseUnverified")
			Do(&err, func() (err error) {
				var ok bool
				claims, ok = jwt_token.Claims.(*sessionv1.ClientSessionPayload)
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
				claims, ok = jwt_token.Claims.(*sessionv1.ClientSessionPayload)
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

	return vanilla.And(
		vanilla.Equal("cluster_uuid", cluster_uuid),
		vanilla.Or(
			vanilla.Equal("status", servicev2.StepStatusRegist),
			vanilla.Equal("status", servicev2.StepStatusSend),
			vanilla.Equal("status", servicev2.StepStatusProcessing),
		)).Parse()

}

// pollingServiceLimit
//  LIMIT ?
func pollingServiceLimit(poliing_limit int) *prepare.Pagination {
	if poliing_limit == 0 {
		poliing_limit = math.MaxInt8 // 127
	}

	return vanilla.Limit(poliing_limit).Parse()
}
