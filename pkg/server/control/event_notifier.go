package control

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/event"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	eventv1 "github.com/NexClipper/sudory/pkg/server/model/event/v1"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

// @Description Create a event notifier console
// @Accept      json
// @Produce     json
// @Tags        server/event_notifier
// @Router      /server/event_notifier/console [post]
// @Param       x_auth_token header string                          false "client session token"
// @Param       object       body   v1.EventNotifierConsole_create  true  "EventNotifierConsole_create"
// @Success     200 {object} v1.EventNotifierConsole
func (ctl Control) CreateEventNotifierConsole(ctx echo.Context) error {
	body := new(eventv1.EventNotifierConsole_create)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if err := body.MIME.Valid(); err != nil {
		return errors.Wrapf(err, "valid MIME")
	}

	notifier := eventv1.EventNotifierConsole{}
	notifier.UuidMeta = metav1.NewUuidMeta()
	notifier.LabelMeta = metav1.NewLabelMeta(body.Name, body.Summary)
	notifier.EventNotifierConsoleProperty = body.EventNotifierConsoleProperty
	notifier.MIME = body.MIME

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		notifier_, err := vault.NewEventNotifierConsole(tx).Create(notifier)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "create a console event notifier"))
		}

		return notifier_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Create a event notifier webhook
// @Accept      json
// @Produce     json
// @Tags        server/event_notifier
// @Router      /server/event_notifier/webhook [post]
// @Param       x_auth_token header string                          false "client session token"
// @Param       object       body   v1.EventNotifierWebhook_create  true  "EventNotifierWebhook_create"
// @Success     200 {object} v1.EventNotifierWebhook
func (ctl Control) CreateEventNotifierWebhook(ctx echo.Context) error {
	body := new(eventv1.EventNotifierWebhook_create)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if err := body.MIME.Valid(); err != nil {
		return errors.Wrapf(err, "valid MIME")
	}

	notifier := eventv1.EventNotifierWebhook{}
	notifier.UuidMeta = metav1.NewUuidMeta()
	notifier.LabelMeta = metav1.NewLabelMeta(body.Name, body.Summary)
	notifier.EventNotifierWebhookProperty = body.EventNotifierWebhookProperty
	notifier.MIME = body.MIME

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		notifier_, err := vault.NewEventNotifierWebhook(tx).Create(notifier)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "create a webhook event notifier"))
		}

		return notifier_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Create a event notifier rabbitmq
// @Accept      json
// @Produce     json
// @Tags        server/event_notifier
// @Router      /server/event_notifier/rabbitmq [post]
// @Param       x_auth_token header string                           false "client session token"
// @Param       object       body   v1.EventNotifierRabbitMq_create  true  "EventNotifierRabbitMq_create"
// @Success     200 {object} v1.EventNotifierRabbitMq
func (ctl Control) CreateEventNotifierRabbitMq(ctx echo.Context) error {
	body := new(eventv1.EventNotifierRabbitMq_create)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if err := body.MIME.Valid(); err != nil {
		return errors.Wrapf(err, "valid MIME")
	}

	//rabbitMQ 연결 테스트
	if false {
		//valid rabbitmq connection
		conn, _, err := new(event.RabbitMQNotifier).Dial(body.Url)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "valid rabbitmq connection%s",
					logs.KVL(
						"url", body.Url,
					)))
		}
		if conn.IsClosed() {
			return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "rabbitmq connection is closed%s",
					logs.KVL(
						"url", body.Url,
					)))
		}
		conn.Close() //no more use
	}

	notifier := eventv1.EventNotifierRabbitMq{}
	notifier.UuidMeta = metav1.NewUuidMeta()
	notifier.LabelMeta = metav1.NewLabelMeta(body.Name, body.Summary)
	notifier.EventNotifierRabbitMqProperty = body.EventNotifierRabbitMqProperty
	notifier.MIME = body.MIME

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		notifier_, err := vault.NewEventNotifierRabbitMq(tx).Create(notifier)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "create a rabbitmq event notifier"))
		}

		return notifier_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Find event notifier console
// @Accept      json
// @Produce     json
// @Tags        server/event_notifier
// @Router      /server/event_notifier/console [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q                   query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o                   query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p                   query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v1.EventNotifierConsole
func (ctl Control) FindEventNotifierConsole(ctx echo.Context) error {
	r, err := vault.NewEventNotifierConsole(ctl.db.Engine().NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "query console event notifier"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Find event notifier webhook
// @Accept      json
// @Produce     json
// @Tags        server/event_notifier
// @Router      /server/event_notifier/webhook [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q                   query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o                   query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p                   query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v1.EventNotifierWebhook
func (ctl Control) FindEventNotifierWebhook(ctx echo.Context) error {
	r, err := vault.NewEventNotifierWebhook(ctl.db.Engine().NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "query webhook event notifier"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Find event notifier
// @Accept      json
// @Produce     json
// @Tags        server/event_notifier
// @Router      /server/event_notifier/rabbitmq [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q                   query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o                   query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p                   query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v1.EventNotifierRabbitMq
func (ctl Control) FindEventNotifierRabbitmq(ctx echo.Context) error {
	r, err := vault.NewEventNotifierRabbitMq(ctl.db.Engine().NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "query rabbitmq event notifier"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Get a event notifier
// @Accept      json
// @Produce     json
// @Tags        server/event_notifier
// @Router      /server/event_notifier/console/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid                path   string true  "event notifier 의 Uuid"
// @Success     200 {object} v1.EventNotifierConsole
func (ctl Control) GetEventNotifierConsole(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := vault.NewEventNotifierConsole(ctl.db.Engine().NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get console event notifier"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Get a event notifier webhook
// @Accept      json
// @Produce     json
// @Tags        server/event_notifier
// @Router      /server/event_notifier/webhook/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid                path   string true  "event notifier 의 Uuid"
// @Success     200 {object} v1.EventNotifierWebhook
func (ctl Control) GetEventNotifierWebhook(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := vault.NewEventNotifierWebhook(ctl.db.Engine().NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get webhook vent notifier"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Get a event notifier rabbitmq
// @Accept      json
// @Produce     json
// @Tags        server/event_notifier
// @Router      /server/event_notifier/rabbitmq/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid                path   string true  "event notifier 의 Uuid"
// @Success     200 {object} v1.EventNotifierRabbitMq
func (ctl Control) GetEventNotifierRabbitmq(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := vault.NewEventNotifierRabbitMq(ctl.db.Engine().NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get rabbitmq event notifier"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Update a console event notifier
// @Accept      json
// @Produce     json
// @Tags        server/event_notifier
// @Router      /server/event_notifier/console/{uuid} [put]
// @Param       x_auth_token header string                         false "client session token"
// @Param       uuid         path   string                         true  "Event 의 Uuid"
// @Param       object       body   v1.EventNotifierConsole_create true  "EventNotifierConsole_create"
// @Success     200 {object} v1.EventNotifierConsole
func (ctl Control) UpdateEventNotifierConsole(ctx echo.Context) error {
	body := new(eventv1.EventNotifierConsole_create)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if err := body.MIME.Valid(); err != nil {
		return errors.Wrapf(err, "valid MIME")
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	notifier := eventv1.EventNotifierConsole{}
	notifier.Uuid = uuid
	notifier.LabelMeta = body.LabelMeta
	notifier.EventNotifierConsoleProperty = body.EventNotifierConsoleProperty
	notifier.MIME = body.MIME

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		notifier_, err := vault.NewEventNotifierConsole(tx).Update(notifier)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "upate console event notifier"))
		}

		return notifier_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Update a webhook event notifier
// @Accept      json
// @Produce     json
// @Tags        server/event_notifier
// @Router      /server/event_notifier/webhook/{uuid} [put]
// @Param       x_auth_token header string                         false "client session token"
// @Param       uuid         path   string                         true  "Event 의 Uuid"
// @Param       object       body   v1.EventNotifierWebhook_create true  "EventNotifierWebhook_create"
// @Success     200 {object} v1.EventNotifierWebhook
func (ctl Control) UpdateEventNotifierWebhook(ctx echo.Context) error {
	body := new(eventv1.EventNotifierWebhook_create)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if err := body.MIME.Valid(); err != nil {
		return errors.Wrapf(err, "valid MIME")
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	notifier := eventv1.EventNotifierWebhook{}
	notifier.Uuid = uuid
	notifier.LabelMeta = body.LabelMeta
	notifier.EventNotifierWebhookProperty = body.EventNotifierWebhookProperty
	notifier.MIME = body.MIME

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		notifier_, err := vault.NewEventNotifierWebhook(tx).Update(notifier)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "upate webhook event notifier"))
		}

		return notifier_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Update a rabbitmq event notifier
// @Accept      json
// @Produce     json
// @Tags        server/event_notifier
// @Router      /server/event_notifier/rabbitmq/{uuid} [put]
// @Param       x_auth_token header string                          false "client session token"
// @Param       uuid         path   string                          true  "Event 의 Uuid"
// @Param       object       body   v1.EventNotifierRabbitMq_create true  "EventNotifierRabbitMq_create"
// @Success     200 {object} v1.EventNotifierRabbitMq
func (ctl Control) UpdateEventNotifierRabbitMq(ctx echo.Context) error {
	body := new(eventv1.EventNotifierRabbitMq_create)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if err := body.MIME.Valid(); err != nil {
		return errors.Wrapf(err, "valid MIME")
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	notifier := eventv1.EventNotifierRabbitMq{}
	notifier.Uuid = uuid
	notifier.LabelMeta = body.LabelMeta
	notifier.EventNotifierRabbitMqProperty = body.EventNotifierRabbitMqProperty
	notifier.MIME = body.MIME

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		notifier_, err := vault.NewEventNotifierRabbitMq(tx).Update(notifier)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "upate rabbitMq event notifier"))
		}

		return notifier_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Delete a event notifier console
// @Accept json
// @Produce json
// @Tags server/event_notifier
// @Router /server/event_notifier/console/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid                path   string true  "Event 의 Uuid"
// @Success 200
func (ctl Control) DeleteEventNotifierConsole(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		if err := vault.NewEventNotifierConsole(tx).Delete(uuid); err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "delete console event notifier"))
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, OK())
}

// @Description Delete a event notifier webhook
// @Accept json
// @Produce json
// @Tags server/event_notifier
// @Router /server/event_notifier/webhook/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid                path   string true  "Event 의 Uuid"
// @Success 200
func (ctl Control) DeleteEventNotifierWebhook(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		if err := vault.NewEventNotifierWebhook(tx).Delete(uuid); err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "delete webhook event notifier"))
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, OK())
}

// @Description Delete a event notifier rabbitmq
// @Accept json
// @Produce json
// @Tags server/event_notifier
// @Router /server/event_notifier/rabbitmq/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid                path   string true  "Event 의 Uuid"
// @Success 200
func (ctl Control) DeleteEventNotifierRabbitmq(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		if err := vault.NewEventNotifierRabbitMq(tx).Delete(uuid); err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "delete rabbitmq event notifier"))
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, OK())
}
