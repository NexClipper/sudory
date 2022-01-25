package control

import (
	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/database"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	"github.com/labstack/echo/v4"
)

// Create ServiceStep
// @Description Create a Service Step
// @Accept json
// @Produce json
// @Tags server/service_step
// @Router /server/service/{service_uuid}/step [post]
// @Param service_uuid path string true "ServiceStep 의 service_uuid"
// @Param step         body stepv1.HttpReqServiceStep true "HttpReqServiceStep"
// @Success 200
func (c *Control) CreateServiceStep() func(ctx echo.Context) error {
	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req["service_uuid"].(string)) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}

		body := new(stepv1.HttpReqServiceStep)
		err := ctx.Bind(body)
		if err != nil {
			return nil, ErrorBindRequestObject(err)
		}
		req["_"] = body
		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}
		service_uuid, ok := req["service_uuid"].(string)
		if !ok {
			return nil, ErrorFailedCast()
		}
		body, ok := req["_"].(*stepv1.HttpReqServiceStep)
		if !ok {
			return nil, ErrorFailedCast()
		}

		body.ServiceStep.ServiceUuid = service_uuid

		//스탭 생성
		err := operator.NewServiceStep(ctx).
			Create(body.ServiceStep)
		if err != nil {
			return nil, err
		}

		//Service Chaining
		operator.NewService(ctx).
			Chaining(service_uuid)

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    Lock,
	})
}

// Get ServiceStep
// @Description Get a Service Step
// @Accept json
// @Produce json
// @Tags server/service_step
// @Router /server/service/{service_uuid}/step [get]
// @Param service_uuid path string true "ServiceStep 의 service_uuid"
// @Success 200 {array} stepv1.HttpRspServiceStep
func (c *Control) GetServiceSteps() func(ctx echo.Context) error {
	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req["service_uuid"]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		where := "service_uuid = ?"
		service_uuid := req["service_uuid"]

		record, err := operator.NewServiceStep(ctx).
			Find(where, service_uuid)
		if err != nil {
			return nil, err
		}
		return stepv1.TransToHttpRsp(record), nil
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    NoLock,
	})
}

// Get ServiceStep
// @Description Get a Service Step
// @Accept json
// @Produce json
// @Tags server/service_step
// @Router /server/service/{service_uuid}/step/{uuid} [get]
// @Param service_uuid path string true "ServiceStep 의 service_uuid"
// @Param uuid         path string true "ServiceStep 의 Uuid"
// @Success 200 {object} stepv1.HttpRspServiceStep
func (c *Control) GetServiceStep() func(ctx echo.Context) error {
	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req["service_uuid"]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		if len(req["uuid"]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		_ = req["service_uuid"]
		uuid := req["uuid"]

		record, err := operator.NewServiceStep(ctx).
			Get(uuid)
		if err != nil {
			return nil, err
		}

		return &stepv1.HttpRspServiceStep{ServiceStep: *record}, nil
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    NoLock,
	})
}

// Update ServiceStep
// @Description Update a Service Step
// @Accept json
// @Produce json
// @Tags server/service_step
// @Router /server/service/{service_uuid}/step/{uuid} [put]
// @Param service_uuid path string true "ServiceStep 의 service_uuid"
// @Param uuid         path string true "ServiceStep 의 Uuid"
// @Param step         body stepv1.HttpReqServiceStep true "HttpReqServiceStep"
// @Success 200
func (c *Control) UpdateServiceStep() func(ctx echo.Context) error {
	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req["service_uuid"].(string)) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		if len(req["uuid"].(string)) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}

		body := new(stepv1.HttpReqServiceStep)
		err := ctx.Bind(body)
		if err != nil {
			return nil, ErrorBindRequestObject(err)
		}
		req["_"] = body

		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}
		service_uuid, ok := req["service_uuid"].(string)
		if !ok {
			return nil, ErrorFailedCast()
		}
		uuid, ok := req["uuid"].(string)
		if !ok {
			return nil, ErrorFailedCast()
		}
		body, ok := req["_"].(*stepv1.HttpReqServiceStep)
		if !ok {
			return nil, ErrorFailedCast()
		}

		body.ServiceStep.ServiceUuid = service_uuid
		body.ServiceStep.Uuid = uuid
		err := operator.NewServiceStep(ctx).
			Update(body.ServiceStep)
		if err != nil {
			return nil, err
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    NoLock,
	})
}

// Delete Service
// @Description Delete a Service
// @Accept json
// @Produce json
// @Tags server/service_step
// @Router /server/service/{service_uuid}/step/{uuid} [delete]
// @Param service_uuid path string true "ServiceStep 의 service_uuid"
// @Param uuid         path string true "ServiceStep 의 Uuid"
// @Success 200
func (c *Control) DeleteServiceStep() func(ctx echo.Context) error {
	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req["service_uuid"]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		if len(req["uuid"]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		_ = req["service_uuid"]
		uuid := req["uuid"]

		//조회 해서 레코드가 없으면 종료
		step, err := operator.NewServiceStep(ctx).
			Get(uuid)
		if Eqaul(err, database.ErrorRecordWasNotFound()) {
			return OK(), nil //idempotent
		} else if err != nil {
			return nil, err
		}

		//삭제
		err = operator.NewServiceStep(ctx).
			Delete(uuid)
		if err != nil {
			return nil, err
		}

		//Service Chaining
		operator.NewService(ctx).
			Chaining(step.ServiceUuid)

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    Lock,
	})
}

// // Update ServiceStep (client)
// // @Description Update a Service Step
// // @Accept json
// // @Produce json
// // @Tags server/service_step
// // @Router /client/service/{service_uuid}/step/{uuid} [put]
// // @Param service_uuid path string true "ServiceStep 의 Uuid"
// // @Param uuid         path string true "ServiceStep 의 Uuid"
// // @Param step         body stepv1.HttpReqServiceStep true "HttpReqServiceStep"
// // @Success 200
// func (c *Control) UpdateClientServiceStep() func(ctx echo.Context) error {
// 	binder := func(ctx echo.Context) (interface{}, error) {
// 		req := make(map[string]interface{})
// 		for _, it := range ctx.ParamNames() {
// 			req[it] = ctx.Param(it)
// 		}
// 		if len(req["service_uuid"].(string)) == 0 {
// 			return nil, ErrorInvaliedRequestParameter()
// 		}
// 		if len(req["uuid"].(string)) == 0 {
// 			return nil, ErrorInvaliedRequestParameter()
// 		}
// 		body := new(stepv1.HttpReqServiceStep)
// 		err := ctx.Bind(body)
// 		if err != nil {
// 			return nil, ErrorBindRequestObject(err)
// 		}
// 		req["_"] = body
// 		return req, nil
// 	}
// 	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
// 		req, ok := v.(map[string]interface{})
// 		if !ok {
// 			return nil, ErrorFailedCast()
// 		}
// 		service_uuid, ok := req["service_uuid"].(string)
// 		if !ok {
// 			return nil, ErrorFailedCast()
// 		}
// 		uuid, ok := req["uuid"].(string)
// 		if !ok {
// 			return nil, ErrorFailedCast()
// 		}
// 		body, ok := req["_"].(*stepv1.HttpReqServiceStep)
// 		if !ok {
// 			return nil, ErrorFailedCast()
// 		}

// 		body.ServiceStep.ServiceUuid = service_uuid
// 		body.ServiceStep.Uuid = uuid

// 		err := operator.NewServiceStep(ctx).
// 			Update(body.ServiceStep)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return nil, nil
// 	}

// 	return MakeMiddlewareFunc(binder, operator, HttpResponse)
// }
