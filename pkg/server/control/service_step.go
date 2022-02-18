package control

import (
	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/database"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
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
// @Success 200 {object} v1.HttpRspServiceStep
func (c *Control) CreateServiceStep() func(ctx echo.Context) error {
	binder := func(ctx Contexter) error {
		body := new(stepv1.HttpReqServiceStep)
		if err := ctx.Bind(body); err != nil {
			return ErrorBindRequestObject(err)
		}
		if len(body.Name) == 0 {
			return ErrorInvaliedRequestParameterName("Name")
		}
		if body.Method == nil {
			return ErrorInvaliedRequestParameterName("Method")
		}

		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__SERVICE_UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__SERVICE_UUID__)
		}
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		body, ok := ctx.Object().(*stepv1.HttpReqServiceStep)
		if !ok {
			return nil, ErrorFailedCast()
		}

		step := body.ServiceStep

		service_uuid := ctx.Params()[__SERVICE_UUID__]

		//property
		step.LabelMeta = NewLabelMeta(step.Name, step.Summary)
		step.ServiceUuid = service_uuid
		if step.Sequence == nil {
			//마직막 순서를 지정하기 위해서 스텝을 가져온다
			where := "service_uuid = ?"
			steps, err := operator.NewServiceStep(ctx.Database()).
				Find(where, service_uuid)
			if err != nil {
				return nil, err
			}
			//스탭 순서 지정
			step.Sequence = newist.Int32(int32(len(steps)))
		}
		//Status = Regist
		step.Status = newist.Int32(int32(servicev1.StatusRegist))

		//스탭 생성
		if err := operator.NewServiceStep(ctx.Database()).Create(step); err != nil {
			return nil, err
		}

		//Service Chaining
		if err := operator.NewService(ctx.Database()).Chaining(service_uuid); err != nil {
			return nil, err
		}
		//ServiceStep ChainingSequence
		if err := operator.NewServiceStep(ctx.Database()).ChainingSequence(service_uuid, step.Uuid); err != nil {
			return nil, err
		}

		return stepv1.HttpRspServiceStep{ServiceStep: step}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Lock(c.db.Engine()),
	})
}

// Get ServiceStep
// @Description Get a Service Step
// @Accept json
// @Produce json
// @Tags server/service_step
// @Router /server/service/{service_uuid}/step [get]
// @Param service_uuid path string true "ServiceStep 의 service_uuid"
// @Success 200 {array} v1.HttpRspServiceStep
func (c *Control) GetServiceSteps() func(ctx echo.Context) error {
	binder := func(ctx Contexter) error {
		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__SERVICE_UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__SERVICE_UUID__)
		}
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		where := "service_uuid = ?"
		service_uuid := ctx.Params()[__SERVICE_UUID__]

		record, err := operator.NewServiceStep(ctx.Database()).
			Find(where, service_uuid)
		if err != nil {
			return nil, err
		}
		return stepv1.TransToHttpRsp(record), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
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
// @Success 200 {object} v1.HttpRspServiceStep
func (c *Control) GetServiceStep() func(ctx echo.Context) error {
	binder := func(ctx Contexter) error {
		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__SERVICE_UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__SERVICE_UUID__)
		}
		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		_ = ctx.Params()[__SERVICE_UUID__]
		uuid := ctx.Params()[__UUID__]

		record, err := operator.NewServiceStep(ctx.Database()).
			Get(uuid)
		if err != nil {
			return nil, err
		}

		return &stepv1.HttpRspServiceStep{ServiceStep: *record}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
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
// @Success 200 {object} v1.HttpRspServiceStep
func (c *Control) UpdateServiceStep() func(ctx echo.Context) error {
	binder := func(ctx Contexter) error {
		body := new(stepv1.HttpReqServiceStep)
		if err := ctx.Bind(body); err != nil {
			return ErrorBindRequestObject(err)
		}

		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__SERVICE_UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__SERVICE_UUID__)
		}
		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}

		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		body, ok := ctx.Object().(*stepv1.HttpReqServiceStep)
		if !ok {
			return nil, ErrorFailedCast()
		}

		step := body.ServiceStep

		service_uuid := ctx.Params()[__SERVICE_UUID__]

		uuid := ctx.Params()[__UUID__]

		//set service uuid from path
		step.ServiceUuid = service_uuid
		//set uuid from path
		step.Uuid = uuid

		if err := operator.NewServiceStep(ctx.Database()).Update(step); err != nil {
			return nil, err
		}

		//ServiceStep ChainingSequence
		if err := operator.NewServiceStep(ctx.Database()).ChainingSequence(service_uuid, step.Uuid); err != nil {
			return nil, err
		}

		return stepv1.HttpRspServiceStep{ServiceStep: step}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Lock(c.db.Engine()),
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
	binder := func(ctx Contexter) error {
		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__SERVICE_UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__SERVICE_UUID__)
		}
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		_ = ctx.Params()[__SERVICE_UUID__]
		uuid := ctx.Params()[__UUID__]

		//조회 해서 레코드가 없으면 종료
		step, err := operator.NewServiceStep(ctx.Database()).
			Get(uuid)
		if Eqaul(err, database.ErrorRecordWasNotFound()) {
			return OK(), nil //idempotent
		} else if err != nil {
			return nil, err
		}

		//삭제
		err = operator.NewServiceStep(ctx.Database()).
			Delete(uuid)
		if err != nil {
			return nil, err
		}

		//Service Chaining
		if err := operator.NewService(ctx.Database()).Chaining(step.ServiceUuid); err != nil {
			return nil, err
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Lock(c.db.Engine()),
	})
}
