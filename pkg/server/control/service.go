package control

import (
	"fmt"
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Create Service
// @Description Create a Service
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service [post]
// @Param       service body v1.HttpReqServiceCreate true "HttpReqServiceCreate"
// @Success     200 {object} v1.HttpRspServiceWithSteps
func (ctl Control) CreateService(ctx echo.Context) error {
	foreach_step := func(elems []stepv1.StepCreate, fn func(stepv1.StepCreate) error) error {
		for _, it := range elems {
			if err := fn(it); err != nil {
				return err
			}
		}
		return nil
	}

	// valid_step := func(step stepv1.StepCreate) error {
	// 	//Method
	// 	if step.Args == nil {
	// 		return errors.Wrapf(ErrorBindRequestObject(), "valid%s",
	// 			logs.KVL(
	// 				"pram", TypeName(step.Args),
	// 			))
	// 	}
	// 	return nil
	// }

	map_step_create := func(elems []stepv1.StepCreate, mapper func(stepv1.StepCreate, int) stepv1.ServiceStep) []stepv1.ServiceStep {
		rst := make([]stepv1.ServiceStep, len(elems))
		for n := range elems {
			rst[n] = mapper(elems[n], n)
		}
		return rst
	}

	body := new(servicev1.HttpReqServiceCreate)
	if err := ctx.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if len(nullable.String(body.Name).Value()) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", TypeName(body.Name),
				)))
	}
	if len(body.TemplateUuid) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", TypeName(body.TemplateUuid),
				)))
	}
	if len(body.ClusterUuid) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", TypeName(body.ClusterUuid),
				)))
	}
	//valid steps
	if err := foreach_step(body.Steps, func(step stepv1.StepCreate) error {
		//Method
		if step.Args == nil {
			return errors.Wrapf(ErrorBindRequestObject(), "valid%s",
				logs.KVL(
					"pram", TypeName(step.Args),
				))
		}
		return nil
	}); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", TypeName(body.Steps),
				)))
	}

	service_create := body.ServiceCreate
	steps_create := body.Steps
	template_uuid := service_create.TemplateUuid

	//valid template
	if _, err := vault.NewTemplate(ctl.NewSession()).Get(template_uuid); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrapf(err, "found template%s",
			logs.KVL(
				"OriginUuid", service_create.TemplateUuid,
			)))
	}
	//valid cluster
	if _, err := vault.NewCluster(ctl.NewSession()).Get(service_create.ClusterUuid); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrapf(err, "found template%s",
			logs.KVL(
				"ClusterUuid", service_create.ClusterUuid,
			)))
	}

	//property
	//property service
	service := servicev1.Service{}
	service.Name = nullable.String(service_create.Name).Ptr()
	service.Summary = nullable.String(service_create.Summary).Ptr()
	service.ClusterUuid = newist.String(service_create.ClusterUuid)
	service.TemplateUuid = newist.String(template_uuid)
	service.UuidMeta = NewUuidMeta()                                //meta uuid
	service.LabelMeta = NewLabelMeta(service.Name, service.Summary) //meta label
	service.Status = newist.Int32(int32(servicev1.StatusRegist))    //Status
	service.StepCount = newist.Int32(int32(len(steps_create)))      //step count

	//get commands
	where := "template_uuid = ?"
	commands, err := vault.NewTemplateCommand(ctl.NewSession()).Find(where, template_uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			errors.Wrapf(err, "NewTemplateCommand Find"))
	}
	if len(steps_create) != len(commands) {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Errorf("diff length of steps and commands expect=%d actual=%d", len(commands), len(steps_create)))
	}

	//property step
	steps := map_step_create(steps_create, func(sse stepv1.StepCreate, i int) stepv1.ServiceStep {
		command := commands[i]

		step := stepv1.ServiceStep{}
		step.UuidMeta = NewUuidMeta()                                //meta uuid
		step.LabelMeta = NewLabelMeta(command.Name, command.Summary) //meta label
		step.ServiceUuid = newist.String(service.Uuid)               //ServiceUuid
		step.Sequence = newist.Int32(int32(i))                       //Sequence
		step.Status = newist.Int32(int32(servicev1.StatusRegist))    //Status(Regist)
		step.Method = command.Method                                 //Method
		step.Args = sse.Args                                         //Args
		step.ResultFilter = command.ResultFilter                     //ResultFilter
		// step.Result = newist.String("")
		return step
	})

	serviceAndSteps := servicev1.ServiceAndSteps{Service: service, Steps: steps}
	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		serviceAndSteps_, err := vault.NewService(db).Create(serviceAndSteps)
		if err != nil {
			return nil, errors.Wrapf(err, "database create")
		}
		return &servicev1.HttpRspServiceWithSteps{DbSchemaServiceAndSteps: *serviceAndSteps_}, nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, r)
}

// Find []Service
// @Description Find []Service
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service [get]
// @Param       q query string false "query  pkg/server/database/prepared/README.md"
// @Param       o query string false "order  pkg/server/database/prepared/README.md"
// @Param       p query string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v1.HttpRspServiceWithSteps
func (ctl Control) FindService(ctx echo.Context) error {
	map_service := func(elems []servicev1.DbSchemaServiceAndSteps, mapper func(servicev1.DbSchemaServiceAndSteps) servicev1.DbSchemaServiceAndSteps) []servicev1.DbSchemaServiceAndSteps {
		rst := make([]servicev1.DbSchemaServiceAndSteps, len(elems))
		for n := range elems {
			rst[n] = mapper(elems[n])
		}
		return rst
	}

	services, err := vault.NewService(ctl.NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrapf(err, "find service%s",
			logs.KVL(
				"query", echoutil.QueryParamString(ctx),
			)))
	}

	//필터링 적용
	// 서비스 Result 필드 값 제거
	services = map_service(services, service_exclude_result)

	return ctx.JSON(http.StatusOK, servicev1.TransToHttpRspServiceAndSteps(services))
}

// Get Service
// @Description Get a Service
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid} [get]
// @Param       uuid path string true "Service 의 Uuid"
// @Success     200 {object} v1.HttpRspServiceWithSteps
func (ctl Control) GetService(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	service, err := vault.NewService(ctl.NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrapf(err, "get service%s",
			logs.KVL(
				"uuid", uuid,
			)))
	}

	//서비스 조회에 결과 필드는 제거
	*service = service_exclude_result(*service)

	return ctx.JSON(http.StatusOK, servicev1.HttpRspServiceWithSteps{DbSchemaServiceAndSteps: *service})
}

// Get Service Result
// @Description Get a Service with Result
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid}/result [get]
// @Param       uuid path string true "Service 의 Uuid"
// @Success     200 {object} v1.HttpRspServiceWithSteps
func (ctl Control) GetServiceResult(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	service, err := vault.NewService(ctl.NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrapf(err, "get service%s",
			logs.KVL(
				"uuid", uuid,
			)))
	}

	return ctx.JSON(http.StatusOK, servicev1.HttpRspServiceWithSteps{DbSchemaServiceAndSteps: *service})
}

// // Update Service
// // @Description Update a Service
// // @Accept      json
// // @Produce     json
// // @Tags        server/service
// // @Router      /server/service/{uuid} [put]
// // @Param       uuid    path string true "Service 의 Uuid"
// // @Param       service body v1.HttpReqService true "HttpReqService"
// // @Success     200 {object} v1.HttpRspService
// func (ctl Control) UpdateService(ctx echo.Context) error {
// 	binder := func(ctx Contexter) error {
// 		body := new(servicev1.HttpReqService)
// 		if err := ctx.Bind(body); err != nil {
// 			return ErrorBindRequestObject(err)
// 		}

// 		if len(ctx.Params()) == 0 {
// 			return ErrorInvaliedRequestParameter()
// 		}
// 		if len(ctx.Params()[__UUID__]) == 0 {
// 			return ErrorInvaliedRequestParameterName(__UUID__)
// 		}

// 		return nil
// 	}
// 	operator := func(ctx Contexter) (interface{}, error) {
// 		body, ok := ctx.Object().(*servicev1.HttpReqService)
// 		if !ok {
// 			return nil, ErrorFailedCast()
// 		}

// 		service := body.Service

// 		//set uuid from path
// 		service.Uuid = ctx.Params()[__UUID__]

// 		//update service
// 		service_, err := vault.NewService(ctx.Database()).Update(service)
// 		if err != nil {
// 			return nil, errors.Wrapf(err, "NewService Update")
// 		}

// 		return servicev1.HttpRspService{DbSchema: *service_}, nil
// 	}

// 	return MakeMiddlewareFunc(Option{
// 		Binder: func(ctx Contexter) error {
// 			err := binder(ctx)
// 			if err != nil {
// 				return errors.Wrapf(err, "UpdateService binder")
// 			}
// 			return nil
// 		},
// 		Operator: func(ctx Contexter) (interface{}, error) {
// 			v, err := operator(ctx)
// 			if err != nil {
// 				return nil, errors.Wrapf(err, "UpdateService operator")
// 			}
// 			return v, nil
// 		},
// 		HttpResponsor: HttpJsonResponsor,
// 		Behavior:      Lock(c.db.Engine()),
// 	})
// }

// Delete Service
// @Description Delete a Service
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid} [delete]
// @Param       uuid path string true "Service 의 Uuid"
// @Success     200
func (ctl Control) DeleteService(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		err := vault.NewService(db).Delete(uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "delete service%s",
				logs.KVL(
					"uuid", uuid,
				))
		}

		return nil, nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, OK())
}

// func TraceServiceOrigin(ctx database.Context, originkind, originuuid string) ([]string, error) {
// 	trace := make([]string, 0)
// 	trace_append := func(kind, uuid string) {
// 		trace = append(trace, fmt.Sprintf("%s:%s", kind, uuid))
// 	}

// LOOP:
// 	for {
// 		switch originkind {
// 		case "template":
// 			template, err := vault.NewTemplate(ctx).Get(originuuid)
// 			if err != nil {
// 				return nil, errors.Wrapf(err, "NewTemplate Get")
// 			}
// 			trace_append("template", template.Uuid)
// 			break LOOP
// 		case "service":
// 			service, err := vault.NewService(ctx).Get(originuuid)
// 			if err != nil {
// 				return nil, errors.Wrapf(err, "NewService Get")
// 			}
// 			trace_append("service", service.Uuid)
// 			originkind, originuuid = *service.OriginKind, *service.OriginUuid
// 		default:
// 			return nil, fmt.Errorf("unknown origin_kind=%s", originkind)
// 		}
// 	}

// 	if len(trace) == 0 {
// 		return nil, fmt.Errorf("not found service origin")
// 	}

// 	return trace, nil
// }

// service_exclude_result

//  서비스 조회에 결과 필드는 제거
func service_exclude_result(service servicev1.DbSchemaServiceAndSteps) servicev1.DbSchemaServiceAndSteps {
	service.Result = nil //서비스 조회에 결과 필드는 제거
	return service
}
