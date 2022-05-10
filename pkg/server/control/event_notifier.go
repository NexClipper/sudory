package control

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/event"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	eventv1 "github.com/NexClipper/sudory/pkg/server/model/event/v1"
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

	notifier := eventv1.EventNotifierConsole{}
	notifier.UuidMeta = NewUuidMeta()
	notifier.LabelMeta = NewLabelMeta(body.Name, body.Summary)
	notifier.EventNotifierConsoleProperty = body.EventNotifierConsoleProperty

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		notifier_, err := vault.NewEventNotifierConsole(tx).Create(notifier)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "create event notifier to console"))
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

	notifier := eventv1.EventNotifierWebhook{}
	notifier.UuidMeta = NewUuidMeta()
	notifier.LabelMeta = NewLabelMeta(body.Name, body.Summary)
	notifier.EventNotifierWebhookProperty = body.EventNotifierWebhookProperty

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		notifier_, err := vault.NewEventNotifierWebhook(tx).Create(notifier)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "create event notifier to webhook"))
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
	notifier.UuidMeta = NewUuidMeta()
	notifier.LabelMeta = NewLabelMeta(body.Name, body.Summary)
	notifier.EventNotifierRabbitMqProperty = body.EventNotifierRabbitMqProperty

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		notifier_, err := vault.NewEventNotifierRabbitMq(tx).Create(notifier)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "create event notifier to rabbitmq"))
		}

		return notifier_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Find event notifier
// @Accept      json
// @Produce     json
// @Tags        server/event_notifier
// @Router      /server/event_notifier/{event_notifier_type} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       event_notifier_type path   string true  "v1.EventNotifierType"
// @Param       q                   query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o                   query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p                   query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v1.EventWithNotifier
func (ctl Control) FindEventNotifier(ctx echo.Context) error {
	type_, err := eventv1.ParseEventNotifierType(echoutil.Param(ctx)[__EVENT_NOTIFIER_TYPE__])
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "invalid event notifier type"))
	}

	finder := func(tx *xorm.Session) (chan interface{}, error) {
		switch type_ {
		case eventv1.EventNotifierTypeConsole:
			notifier, err := vault.NewEventNotifierConsole(tx).Query(echoutil.QueryParam(ctx))
			if err != nil {
				return nil, errors.Wrapf(err, "query console event notifier")
			}

			c := make(chan interface{}, len(notifier))
			defer close(c)
			for _, notifier := range notifier {
				c <- notifier
			}

			return c, nil
		case eventv1.EventNotifierTypeWebhook:
			notifier, err := vault.NewEventNotifierWebhook(tx).Query(echoutil.QueryParam(ctx))
			if err != nil {
				return nil, errors.Wrapf(err, "query webhook event notifier")
			}

			c := make(chan interface{}, len(notifier))
			defer close(c)
			for _, notifier := range notifier {
				c <- notifier
			}

			return c, nil
		case eventv1.EventNotifierTypeRabbitmq:
			notifier, err := vault.NewEventNotifierRabbitMq(tx).Query(echoutil.QueryParam(ctx))
			if err != nil {
				return nil, errors.Wrapf(err, "query rabbitmq event notifier")
			}

			c := make(chan interface{}, len(notifier))
			defer close(c)
			for _, notifier := range notifier {
				c <- notifier
			}

			return c, nil
		}
		return nil, errors.Wrapf(err, "invalid event notifier type")
	}

	c, err := finder(ctl.db.Engine().NewSession())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "query event notifier"))
	}

	r := make([]interface{}, 0)
	for i := range c {
		r = append(r, i)
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Get a event notifier
// @Accept      json
// @Produce     json
// @Tags        server/event_notifier
// @Router      /server/event_notifier/{event_notifier_type}/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       event_notifier_type path   string true  "v1.EventNotifierType"
// @Param       uuid                path   string true  "event notifier 의 Uuid"
// @Success     200 {object} v1.EventWithNotifier
func (ctl Control) GetEventNotifier(ctx echo.Context) error {
	//valid notifier type
	type_, err := eventv1.ParseEventNotifierType(echoutil.Param(ctx)[__EVENT_NOTIFIER_TYPE__])
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "invalid event notifier type"))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	getter := func(tx *xorm.Session) (interface{}, error) {
		switch type_ {
		case eventv1.EventNotifierTypeConsole:
			notifier, err := vault.NewEventNotifierConsole(ctl.db.Engine().NewSession()).Get(uuid)
			if err != nil {
				return nil, errors.Wrapf(err, "get console event notifier")
			}

			return notifier, nil
		case eventv1.EventNotifierTypeWebhook:
			notifier, err := vault.NewEventNotifierWebhook(ctl.db.Engine().NewSession()).Get(uuid)
			if err != nil {
				return nil, errors.Wrapf(err, "get webhook event notifier")
			}

			return notifier, nil
		case eventv1.EventNotifierTypeRabbitmq:
			notifier, err := vault.NewEventNotifierRabbitMq(ctl.db.Engine().NewSession()).Get(uuid)
			if err != nil {
				return nil, errors.Wrapf(err, "get rabbitmq event notifier")
			}

			return notifier, nil
		}
		return nil, errors.Errorf("invalid event notifier type")
	}

	r, err := getter(ctl.db.Engine().NewSession())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get event notifier"))
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

	uuid := echoutil.Param(ctx)[__UUID__]

	notifier := eventv1.EventNotifierConsole{}
	notifier.Uuid = uuid
	notifier.LabelMeta = body.LabelMeta
	notifier.EventNotifierConsoleProperty = body.EventNotifierConsoleProperty

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

	uuid := echoutil.Param(ctx)[__UUID__]

	notifier := eventv1.EventNotifierWebhook{}
	notifier.Uuid = uuid
	notifier.LabelMeta = body.LabelMeta
	notifier.EventNotifierWebhookProperty = body.EventNotifierWebhookProperty

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

	uuid := echoutil.Param(ctx)[__UUID__]

	notifier := eventv1.EventNotifierRabbitMq{}
	notifier.Uuid = uuid
	notifier.LabelMeta = body.LabelMeta
	notifier.EventNotifierRabbitMqProperty = body.EventNotifierRabbitMqProperty

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

// @Description Delete a event notifier
// @Accept json
// @Produce json
// @Tags server/event_notifier
// @Router /server/event_notifier/{event_notifier_type}/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       event_notifier_type path   string true  "v1.EventNotifierType"
// @Param       uuid                path   string true  "Event 의 Uuid"
// @Success 200
func (ctl Control) DeleteEventNotifier(ctx echo.Context) error {
	//valid notifier type
	type_, err := eventv1.ParseEventNotifierType(echoutil.Param(ctx)[__EVENT_NOTIFIER_TYPE__])
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "invalid event notifier type"))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	delete := func(tx *xorm.Session) error {
		switch type_ {
		case eventv1.EventNotifierTypeConsole:
			err := vault.NewEventNotifierConsole(tx).Delete(uuid)
			if err != nil {
				return errors.Wrapf(err, "delete console event notifier")
			}
		case eventv1.EventNotifierTypeWebhook:
			err := vault.NewEventNotifierWebhook(tx).Delete(uuid)
			if err != nil {
				return errors.Wrapf(err, "delete webhook event notifier")
			}
		case eventv1.EventNotifierTypeRabbitmq:
			err := vault.NewEventNotifierRabbitMq(tx).Delete(uuid)
			if err != nil {
				return errors.Wrapf(err, "delete rabbitmq event notifier")
			}
		}

		return nil
	}

	_, err = ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		if err := delete(tx); err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "delete event notifier"))
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, OK())
}
