package control

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"

	"github.com/NexClipper/sudory/pkg/server/control/vanilla"
	clusterv2 "github.com/NexClipper/sudory/pkg/server/model/cluster/v2"
	noxorm "github.com/NexClipper/sudory/pkg/server/model/noxorm/v2"
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
// @Param       service      body   v2.HttpReq_ServiceCreate true  "HttpReq_ServiceCreate"
// @Success     200 {object} v2.HttpRsp_Service
func (ctl ControlVanilla) CreateService(ctx echo.Context) (err error) {
	body := new(servicev2.HttpReq_ServiceCreate)
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
		cond := vanilla.NewCond(
			"WHERE uuid = ? AND deleted IS NULL",
			body.ClusterUuid,
		)

		cluster := clusterv2.Cluster{}
		err = vanilla.QueryRow(ctl.DB(), cluster.TableName(), cluster.ColumnNames(), *cond)(func(s vanilla.Scanner) error {
			return cluster.Scan(s)
		})

		err = errors.Wrapf(err, "valid: cluster is not exists")
		return
	})

	//valid template
	template := templatev2.Template{}
	Do(&err, func() (err error) {
		cond := vanilla.NewCond(
			"WHERE uuid = ? AND deleted IS NULL",
			body.TemplateUuid,
		)

		err = vanilla.QueryRow(ctl.DB(), template.TableName(), template.ColumnNames(), *cond)(func(s vanilla.Scanner) error {
			return template.Scan(s)
		})

		err = errors.Wrapf(err, "valid: template is not exists")
		return
	})

	commands := make([]templatev2.TemplateCommand, 0, __INIT_RECORD_CAPACITY__)
	Do(&err, func() (err error) {
		command := templatev2.TemplateCommand{}

		cond := vanilla.NewCond(
			"WHERE template_uuid = ? AND deleted IS NULL",
			body.TemplateUuid,
		)
		order := vanilla.NewCond(
			"ORDER BY sequence",
		)

		err = vanilla.QueryRows(ctl.DB(), command.TableName(), command.ColumnNames(), *cond, *order)(func(s vanilla.Scanner) (err error) {
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

	rsp := servicev2.HttpRsp_Service{}
	rsp.Steps = make([]servicev2.ServiceStep_tangled, 0, __INIT_RECORD_CAPACITY__)
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
		service.Summary = noxorm.NullString(body.Summary)
		service.ClusterUuid = body.ClusterUuid
		service.TemplateUuid = body.TemplateUuid
		service.StepCount = len(body.Steps)
		service.SubscribedChannel = noxorm.NullString(body.SubscribedChannel)
		service.OnCompletion = body.OnCompletion

		//create steps
		for i := range body.Steps {
			command := commands[i]
			body := body.Steps[i]

			// optional; step.Name
			name := body.Name
			if len(name) == 0 {
				name = command.Name
			}
			// optional; step.Summary
			summary := body.Summary
			if len(summary) == 0 {
				summary = command.Summary.String()
			}
			//property step
			step := servicev2.ServiceStep{
				Uuid:     uuid,
				Sequence: i,
				Created:  time.Now(),
			}
			step.Name = name                          //
			step.Summary = noxorm.NullString(summary) //
			step.Method = string(command.Method)      // command method
			step.Args = body.Args                     //
			step.ResultFilter = command.ResultFilter  // command result filter

			// save step
			rsp.Steps = append(rsp.Steps, servicev2.ServiceStep_tangled{ServiceStep: step})
		}

		//save service
		rsp.Service_tangled = servicev2.Service_tangled{Service: service}

		return
	})

	Do(&err, func() (err error) {
		err = ctl.Scope(func(tx *sql.Tx) (err error) {
			//save service
			Do(&err, func() (err error) {
				err = vanilla.InsertRow(tx, rsp.Service.TableName(), rsp.Service.ColumnNames())(func(e vanilla.Executor) (sql.Result, error) {
					return e.Exec(rsp.Service.Values()...)
				})
				err = errors.Wrapf(err, "failed to save service")
				return
			})
			//save steps
			// Do(&err, func() (err error) {
			// 	for i := range rsp.Steps {
			// 		step := rsp.Steps[i].ServiceStep
			// 		Do(&err, func() (err error) {
			// 			err = create_service_step(tx, step)
			// 			err = errors.Wrapf(err, "failed to save service step")
			// 			return
			// 		})
			// 	}
			// 	err = errors.Wrapf(err, "failed to create service step")
			// 	return
			// })
			Do(&err, func() (err error) {
				// for i := range rsp.Steps {
				// 	step := rsp.Steps[i].ServiceStep
				// 	Do(&err, func() (err error) {
				// 		err = create_service_step(tx, rsp.Steps)
				// 		err = errors.Wrapf(err, "failed to save service step")
				// 		return
				// 	})
				// }

				step := servicev2.ServiceStep{}
				err = vanilla.InsertRows(tx, step.TableName(), step.ColumnNames())(func(i int) ([]interface{}, bool) {
					if i == len(rsp.Steps) {
						return nil, false
					}
					return rsp.Steps[i].ServiceStep.Values(), true
				})
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

	conditions, args, err := ParseDecoration(echoutil.QueryParam(ctx))
	err = errors.Wrapf(err, "ParseDecoration%v", logs.KVL(
		"query", echoutil.QueryParamString(ctx),
	))

	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	rsps := make([]servicev2.HttpRsp_Service_status, 0, __INIT_RECORD_CAPACITY__)

	Do(&err, func() (err error) {
		cond := vanilla.NewCond(
			strings.Join(conditions, "\n"),
			args...,
		)
		var servcie_status []servicev2.Service_status
		servcie_status, err = find_services_status(ctl.DB(), *cond)
		err = errors.Wrapf(err, "find service")

		Do(&err, func() (err error) {
			for i := range servcie_status {
				service_uuid := servcie_status[i].Service.Uuid
				var service_steps []servicev2.ServiceStep_tangled
				service_steps, err = get_service_steps(ctl.DB(), service_uuid)
				err = errors.Wrapf(err, "get steps")
				if err != nil {
					break
				}

				rsps = append(rsps, servicev2.HttpRsp_Service_status{
					Service_status: servcie_status[i],
					Steps:          service_steps,
				})
			}
			return
		})

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
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	var service *servicev2.Service_tangled
	Do(&err, func() (err error) {
		service, err = get_service(ctl.DB(), uuid)
		return
	})

	var steps []servicev2.ServiceStep_tangled
	Do(&err, func() (err error) {
		steps, err = get_service_steps(ctl.DB(), uuid)
		return
	})

	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, servicev2.HttpRsp_Service{Service_tangled: *service, Steps: steps})
}

const __DEFAULT_DECORATION_LIMIT__ = 20

func ParseDecoration(m map[string]string) (conditions []string, args []interface{}, err error) {
	conditions = make([]string, 0, 3)
	args = make([]interface{}, 0, 3)

	Do(&err, func() (err error) {
		deco, err := prepare.NewDecoration(m)
		Do(&err, func() (err error) {
			if deco.Condition != nil {
				cond := "WHERE " + deco.Condition.Query()
				conditions = append(conditions, cond)
				args = deco.Condition.Args()
			}
			if deco.Orders != nil {
				order := ""
				for i := range ([]prepare.Order)(*deco.Orders) {
					if len(order) == 0 {
						order += "ORDER BY "
					}
					order += ([]prepare.Order)(*deco.Orders)[i].Order()
				}
				conditions = append(conditions, order)
			}
			if deco.Pagination != nil {
				page := fmt.Sprintf("LIMIT %v, %v", deco.Pagination.Offset(), deco.Pagination.Limit())
				conditions = append(conditions, page)
			}
			if deco.Pagination == nil {
				page := fmt.Sprintf("LIMIT %v", __DEFAULT_DECORATION_LIMIT__)
				conditions = append(conditions, page)
			}
			return
		})
		return
	})

	return
}

func get_service_step(tx vanilla.Preparer, service_uuid string, sequence int) (step servicev2.ServiceStep_tangled, err error) {
	cond := vanilla.Condition{
		Condition: "WHERE uuid = ? AND sequence = ?",
		Args: []interface{}{
			service_uuid,
			sequence,
		},
	}

	err = vanilla.QueryRow(tx, step.TableName(), step.ColumnNames(), cond)(func(s vanilla.Scanner) (err error) {
		err = step.Scan(s)
		err = errors.Wrapf(err, "step Scan")
		return
	})

	err = errors.Wrapf(err, "failed to get a step")
	return
}

func get_service_steps(tx vanilla.Preparer, service_uuid string) (steps []servicev2.ServiceStep_tangled, err error) {
	steps = make([]servicev2.ServiceStep_tangled, 0, __INIT_RECORD_CAPACITY__)

	cond := vanilla.Condition{
		Condition: "WHERE uuid = ?",
		Args: []interface{}{
			service_uuid,
		},
	}
	step := servicev2.ServiceStep_tangled{}

	err = vanilla.QueryRows(tx, step.TableName(), step.ColumnNames(), cond)(func(s vanilla.Scanner) (err error) {
		err = step.Scan(s)
		err = errors.Wrapf(err, "step Scan")
		Do(&err, func() (err error) {
			steps = append(steps, step)
			return
		})
		return
	})

	err = errors.Wrapf(err, "failed to get step lists")
	return
}

func find_services_status(tx vanilla.Preparer, condition ...vanilla.Condition) (rsps []servicev2.Service_status, err error) {
	rsps = make([]servicev2.Service_status, 0, __INIT_RECORD_CAPACITY__)

	// cond := vanilla.Condition{
	// 	Condition: condition,
	// 	Args:      args,
	// }
	service_status := servicev2.Service_status{}

	err = vanilla.QueryRows(tx, service_status.TableName(), service_status.ColumnNames(), condition...)(func(s vanilla.Scanner) (err error) {
		err = service_status.Scan(s)
		err = errors.Wrapf(err, "service Scan")
		Do(&err, func() (err error) {
			rsps = append(rsps, service_status)
			return
		})
		return
	})

	err = errors.Wrapf(err, "failed to find services")
	return
}

func find_services_tangled(tx vanilla.Preparer, condition ...vanilla.Condition) (rsps []servicev2.Service_tangled, err error) {
	rsps = make([]servicev2.Service_tangled, 0, __INIT_RECORD_CAPACITY__)

	// cond := vanilla.Condition{
	// 	Condition: condition,
	// 	Args:      args,
	// }
	service := servicev2.Service_tangled{}

	err = vanilla.QueryRows(tx, service.TableName(), service.ColumnNames(), condition...)(func(s vanilla.Scanner) (err error) {
		err = service.Scan(s)
		err = errors.Wrapf(err, "service Scan")
		Do(&err, func() (err error) {
			rsps = append(rsps, service)
			return
		})
		return
	})

	err = errors.Wrapf(err, "failed to find services")
	return
}

func find_service_steps(tx vanilla.Preparer, condition ...vanilla.Condition) (rsps []servicev2.HttpRsp_ServiceStep, err error) {
	rsps = make([]servicev2.HttpRsp_ServiceStep, 0, __INIT_RECORD_CAPACITY__)

	// cond := vanilla.Condition{
	// 	Condition: condition,
	// 	Args:      args,
	// }
	step := servicev2.ServiceStep_tangled{}

	err = vanilla.QueryRows(tx, step.TableName(), step.ColumnNames(), condition...)(func(s vanilla.Scanner) (err error) {
		err = step.Scan(s)
		err = errors.Wrapf(err, "service step Scan")
		Do(&err, func() (err error) {
			rsps = append(rsps, servicev2.HttpRsp_ServiceStep{ServiceStep_tangled: step})
			return
		})
		return
	})

	err = errors.Wrapf(err, "failed to find service steps")
	return
}

func get_service(tx vanilla.Preparer, service_uuid string) (service *servicev2.Service_tangled, err error) {
	cond := vanilla.NewCond(
		"WHERE uuid = ?",
		service_uuid,
	)
	service = &servicev2.Service_tangled{}
	err = vanilla.QueryRow(tx, service.TableName(), service.ColumnNames(), *cond)(func(s vanilla.Scanner) (err error) {
		err = service.Scan(s)
		err = errors.Wrapf(err, "service Scan")
		return
	})

	err = errors.Wrapf(err, "failed to get a service")
	return
}

func get_cluster(tx vanilla.Preparer, cluster_uuid string) (cluster *clusterv2.Cluster, err error) {
	cond := vanilla.NewCond(
		"WHERE uuid = ?",
		cluster_uuid,
	)
	cluster = &clusterv2.Cluster{}
	err = vanilla.QueryRow(tx, cluster.TableName(), cluster.ColumnNames(), *cond)(func(s vanilla.Scanner) (err error) {
		err = cluster.Scan(s)
		err = errors.Wrapf(err, "cluster Scan")
		return
	})

	err = errors.Wrapf(err, "faild to get a cluster")
	return
}
