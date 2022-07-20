package control

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/error_compose"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/prepare"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	clusterv2 "github.com/NexClipper/sudory/pkg/server/model/cluster/v2"
	servicev2 "github.com/NexClipper/sudory/pkg/server/model/service/v2"
	templatev2 "github.com/NexClipper/sudory/pkg/server/model/template/v2"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/qri-io/jsonschema"
)

// Create Service
// @Description Create a Service
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service [post]
// @Param       x_auth_token header string                   false "client session token"
// @Param       service      body   v2.HttpReq_Service_create true  "HttpReq_Service_create"
// @Success     200 {object} v2.HttpRsp_Service_create
func (ctl ControlVanilla) CreateService(ctx echo.Context) (err error) {
	body := new(servicev2.HttpReq_Service_create)
	Do(&err, func() (err error) {
		err = echoutil.Bind(ctx, body)
		err = errors.Wrapf(err, "bind%s",
			logs.KVL(
				"type", TypeName(body),
			))
		return
	})
	Do(&err, func() (err error) {
		if len(body.Name) == 0 {
			err = ErrorInvalidRequestParameter()
		}
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(fmt.Sprintf("%s.Name", TypeName(body)), body.Name)...,
			))
		return
	})
	Do(&err, func() (err error) {
		if len(body.TemplateUuid) == 0 {
			err = ErrorInvalidRequestParameter()
		}
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(fmt.Sprintf("%s.TemplateUuid", TypeName(body)), body.TemplateUuid)...,
			))
		return
	})
	Do(&err, func() (err error) {
		if len(body.ClusterUuid) == 0 {
			err = ErrorInvalidRequestParameter()
		}
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(fmt.Sprintf("%s.ClusterUuid", TypeName(body)), body.ClusterUuid)...,
			))
		return
	})

	//valid cluster
	Do(&err, func() (err error) {
		q := vanilla.And(
			vanilla.Equal("uuid", body.ClusterUuid),
			vanilla.IsNull("deleted"),
		).Parse()

		cluster := clusterv2.Cluster{}
		stmt := vanilla.Stmt.Select(cluster.TableName(), cluster.ColumnNames(), q, nil, nil)
		err = stmt.QueryRow(ctl)(func(s vanilla.Scanner) error {
			return cluster.Scan(s)
		})
		err = errors.Wrapf(err, "valid: cluster is not exists")
		return
	})

	//valid template
	template := templatev2.Template{}
	Do(&err, func() (err error) {
		q := vanilla.And(
			vanilla.Equal("uuid", body.TemplateUuid),
			vanilla.IsNull("deleted"),
		).Parse()

		stmt := vanilla.Stmt.Select(template.TableName(), template.ColumnNames(), q, nil, nil)
		err = stmt.QueryRow(ctl)(func(s vanilla.Scanner) error {
			return template.Scan(s)
		})

		err = errors.Wrapf(err, "valid: template is not exists")
		return
	})

	commands := make([]templatev2.TemplateCommand, 0, __INIT_SLICE_CAPACITY__())
	Do(&err, func() (err error) {
		q := vanilla.And(
			vanilla.Equal("uuid", body.TemplateUuid),
			vanilla.IsNull("deleted"),
		).Parse()
		o := vanilla.Asc("sequence").Parse()

		command := templatev2.TemplateCommand{}
		stmt := vanilla.Stmt.Select(command.TableName(), command.ColumnNames(), q, o, nil)
		err = stmt.QueryRow(ctl)(func(s vanilla.Scanner) (err error) {
			err = command.Scan(s)
			if err == nil {
				commands = append(commands, command)
			}
			return
		})

		err = errors.Wrapf(err, "failed to get template commands")
		return
	})

	Do(&err, func() (err error) {
		if len(body.Steps) != len(commands) {
			err = errors.Errorf("diff length of steps and commands%s",
				logs.KVL(
					"expected", len(commands),
					"actual", len(body.Steps),
				))
		}

		return
	})

	Do(&err, func() (err error) {
		for i := range body.Steps {
			step_args := body.Steps[i].Args
			command_args := commands[i].Args

			Do(&err, func() (err error) {
				if step_args == nil {
					err = errors.New("step.Args must have value")
					err = HttpError(err, http.StatusBadRequest) // bad request
				}
				return
			})
			Do(&err, func() (err error) {
				json_schema_validator := &jsonschema.Schema{}
				json_schema, err := json.Marshal(command_args)
				err = errors.Wrapf(err, "command.Args convert to json")
				Do(&err, func() (err error) {
					err = json.Unmarshal([]byte(json_schema), json_schema_validator)
					err = errors.Wrapf(err, "command.Args convert to json schema validator")
					return
				})
				Do(&err, func() (err error) {
					step_args, err := json.Marshal(step_args)
					err = errors.Wrapf(err, "step.Args convert to json")
					Do(&err, func() (err error) {
						timeout, cancel := context.WithTimeout(context.Background(), 333*time.Millisecond)
						defer cancel()

						verr, err := json_schema_validator.ValidateBytes(timeout, step_args)
						err = errors.Wrapf(err, "json schema validatebytes%s", logs.KVL(
							"step.args", string(step_args),
						))
						for _, verr := range verr {
							err = errors.Wrapf(verr, "valid step.args")
						}
						return
					})
					return
				})
				if err != nil {
					err = HttpError(err, http.StatusInternalServerError) // internal server error
				}
				return
			})
		}
		return
	})

	if err != nil {
		err = HttpError(err, http.StatusBadRequest)
	}

	rsp := servicev2.HttpRsp_Service_create{}
	rsp.Steps = make([]servicev2.ServiceStep, 0, __INIT_SLICE_CAPACITY__())
	Do(&err, func() (err error) {
		uuid := body.Uuid
		if len(uuid) == 0 {
			uuid = NewUuidString() // len(uuid) == 0; create uuid
		}

		//property service
		service := servicev2.Service{
			Uuid:    uuid,
			Created: time.Now(),
		}
		service.Name = body.Name
		service.Summary = *vanilla.NewNullString(body.Summary)
		service.ClusterUuid = body.ClusterUuid
		service.TemplateUuid = body.TemplateUuid
		service.StepCount = len(body.Steps)
		service.SubscribedChannel = *vanilla.NewNullString(body.SubscribedChannel)

		//create steps
		for i := range body.Steps {
			command := commands[i]
			body := body.Steps[i]

			// // optional; step.Name
			// name := body.Name
			// if len(name) == 0 {
			// 	name = command.Name
			// }
			// // optional; step.Summary
			// summary := body.Summary
			// if len(summary) == 0 {
			// 	summary = command.Summary.String()
			// }
			//property step
			step := servicev2.ServiceStep{
				Uuid:     uuid,
				Sequence: i,
				Created:  time.Now(),
			}
			step.Name = command.Name                 //
			step.Summary = command.Summary           //
			step.Method = command.Method.String      // command method
			step.Args = body.Args                    //
			step.ResultFilter = command.ResultFilter // command result filter

			// save step
			rsp.Steps = append(rsp.Steps, step)
		}

		//save service
		rsp.Service = service

		return
	})

	Do(&err, func() (err error) {
		err = ctl.Scope(func(tx *sql.Tx) (err error) {
			//save service
			Do(&err, func() (err error) {
				stmt, err := vanilla.Stmt.Insert(rsp.Service.TableName(), rsp.Service.ColumnNames(), rsp.Service.Values())
				err = errors.Wrapf(err, "can not build a service insert statement")
				if err != nil {
					return
				}

				affected, err := stmt.Exec(tx)
				err = errors.Wrapf(err, "failed to save service")
				if affected == 0 {
					err = error_compose.Compose(err, errors.New("no affected"))
				}
				return
			})
			//save steps
			Do(&err, func() (err error) {
				// flat values
				var step servicev2.ServiceStep
				values := make([][]interface{}, 0, len(step.ColumnNames())*len(rsp.Steps))
				for i := range rsp.Steps {
					step = rsp.Steps[i]
					values = append(values, step.Values())
				}

				if len(values) == 0 {
					return
				}

				stmt, err := vanilla.Stmt.Insert(step.TableName(), step.ColumnNames(), values...)
				err = errors.Wrapf(err, "can not build a service_step insert statement")
				if err != nil {
					return
				}

				affected, err := stmt.Exec(tx)
				err = errors.Wrapf(err, "failed to save service")
				if affected == 0 {
					err = error_compose.Compose(err, errors.New("no affected"))
				}

				err = errors.Wrapf(err, "failed to create service step")
				return
			})

			return
		})
		err = errors.Wrapf(err, "failed to create service && steps")
		return
	})

	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, rsp)
}

// Find []Service
// @Description Find []Service
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v2.HttpRsp_Service_status
func (ctl ControlVanilla) FindService(ctx echo.Context) (err error) {
	q, o, p, err := ParseDecoration(echoutil.QueryParam(ctx))
	err = errors.Wrapf(err, "ParseDecoration%v", logs.KVL(
		"query", echoutil.QueryParamString(ctx),
	))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	rsps := make([]servicev2.HttpRsp_Service_status, 0, __INIT_SLICE_CAPACITY__())

	var servcie_status servicev2.Service_status
	err = vanilla.Stmt.Select(servcie_status.TableName(), servcie_status.ColumnNames(), q, o, p).
		QueryRows(ctl)(func(scan vanilla.Scanner, _ int) (err error) {
		err = servcie_status.Scan(scan)
		if err != nil {
			return
		}

		rst := servicev2.HttpRsp_Service_status{
			Service_status: servcie_status,
			Steps:          make([]servicev2.ServiceStep_tangled, 0, __INIT_SLICE_CAPACITY__()),
		}

		eq_uuid := vanilla.Equal("uuid", servcie_status.Uuid).Parse()

		step := servicev2.ServiceStep_tangled{}
		stmt := vanilla.Stmt.Select(step.TableName(), step.ColumnNames(), eq_uuid, nil, nil)
		err = stmt.QueryRows(ctl)(func(scan vanilla.Scanner, _ int) (err error) {
			err = step.Scan(scan)
			if err != nil {
				return
			}

			rst.Steps = append(rst.Steps, step)
			return
		})
		if err != nil {
			return
		}

		rsps = append(rsps, rst)
		return
	})

	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, rsps)
}

// Get Service
// @Description Get a Service
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Service 의 Uuid"
// @Success     200 {object} v2.HttpRsp_Service
func (ctl ControlVanilla) GetService(ctx echo.Context) (err error) {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		err = ErrorInvalidRequestParameter()
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
		))

	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	eq_uuid := vanilla.Equal("uuid", uuid).Parse()

	rst := servicev2.HttpRsp_Service{}

	var servcie servicev2.Service_tangled
	err = vanilla.Stmt.Select(servcie.TableName(), servcie.ColumnNames(), eq_uuid, nil, nil).
		QueryRow(ctl)(func(s vanilla.Scanner) (err error) {
		err = servcie.Scan(s)
		if err == nil {
			rst.Service_tangled = servcie
			rst.Steps = make([]servicev2.ServiceStep_tangled, 0, __INIT_SLICE_CAPACITY__())

			step := servicev2.ServiceStep_tangled{}
			stmt := vanilla.Stmt.Select(step.TableName(), step.ColumnNames(), eq_uuid, nil, nil)
			err = stmt.QueryRows(ctl)(func(scan vanilla.Scanner, _ int) (err error) {

				err = step.Scan(scan)
				if err == nil {
					rst.Steps = append(rst.Steps, step)
				}
				return
			})
		}
		return
	})

	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, rst)
}

const __DEFAULT_DECORATION_LIMIT__ = 20

func ParseDecoration(m map[string]string) (q *prepare.Condition, o *prepare.Orders, p *prepare.Pagination, err error) {
	q, o, p, err = prepare.NewParser(m)
	if p == nil {
		p = vanilla.Limit(__DEFAULT_DECORATION_LIMIT__).Parse()
	}

	return
}
