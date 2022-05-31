package control

import (
	"fmt"
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/event"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	channelv1 "github.com/NexClipper/sudory/pkg/server/model/channel/v1"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

// @Description Create a channel notifier console
// @Accept      json
// @Produce     json
// @Tags        server/channel_notifier
// @Router      /server/channel_notifier/console [post]
// @Param       x_auth_token header string                          false "client session token"
// @Param       object       body   v1.NotifierConsole_create  true  "EventNotifierConsole_create"
// @Success     200 {object} v1.NotifierConsole
func (ctl Control) CreateChannelNotifierConsole(ctx echo.Context) error {
	body := new(channelv1.NotifierConsole_create)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if len(body.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.Name", TypeName(body)), body.Name)...,
				)))
	}

	//mime
	if err := body.MIME.Valid(); err != nil {
		return errors.Wrapf(err, "valid MIME")
	}

	notifier := channelv1.NotifierConsole{}
	notifier.UuidMeta = metav1.NewUuidMeta()
	notifier.LabelMeta = metav1.NewLabelMeta(body.Name, body.Summary)
	notifier.NotifierConsoleProperty = body.NotifierConsoleProperty
	notifier.MIME = body.MIME

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		notifier_, err := vault.NewNotifierConsole(tx).Create(notifier)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "create a console channel notifier"))
		}

		return notifier_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Create a channel notifier webhook
// @Accept      json
// @Produce     json
// @Tags        server/channel_notifier
// @Router      /server/channel_notifier/webhook [post]
// @Param       x_auth_token header string                          false "client session token"
// @Param       object       body   v1.NotifierWebhook_create  true  "EventNotifierWebhook_create"
// @Success     200 {object} v1.NotifierWebhook
func (ctl Control) CreateChannelNotifierWebhook(ctx echo.Context) error {
	body := new(channelv1.NotifierWebhook_create)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if len(body.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.Name", TypeName(body)), body.Name)...,
				)))
	}

	if len(body.Url) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.Url", TypeName(body)), body.Url)...,
				)))
	}

	//mime
	if err := body.MIME.Valid(); err != nil {
		return errors.Wrapf(err, "valid MIME")
	}

	notifier := channelv1.NotifierWebhook{}
	notifier.UuidMeta = metav1.NewUuidMeta()
	notifier.LabelMeta = metav1.NewLabelMeta(body.Name, body.Summary)
	notifier.NotifierWebhookProperty = body.NotifierWebhookProperty
	notifier.MIME = body.MIME

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		notifier_, err := vault.NewNotifierWebhook(tx).Create(notifier)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "create a webhook channel notifier"))
		}

		return notifier_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Create a channel notifier rabbitmq
// @Accept      json
// @Produce     json
// @Tags        server/channel_notifier
// @Router      /server/channel_notifier/rabbitmq [post]
// @Param       x_auth_token header string                           false "client session token"
// @Param       object       body   v1.NotifierRabbitMq_create  true  "NotifierRabbitMq_create"
// @Success     200 {object} v1.NotifierRabbitMq
func (ctl Control) CreateChannelNotifierRabbitMq(ctx echo.Context) error {
	body := new(channelv1.NotifierRabbitMq_create)
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

	notifier := channelv1.NotifierRabbitMq{}
	notifier.UuidMeta = metav1.NewUuidMeta()
	notifier.LabelMeta = metav1.NewLabelMeta(body.Name, body.Summary)
	notifier.NotifierRabbitMqProperty = body.NotifierRabbitMqProperty
	notifier.MIME = body.MIME

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		notifier_, err := vault.NewNotifierRabbitMq(tx).Create(notifier)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "create a rabbitmq channel notifier"))
		}

		return notifier_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Find channel notifier console
// @Accept      json
// @Produce     json
// @Tags        server/channel_notifier
// @Router      /server/channel_notifier/console [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q                   query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o                   query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p                   query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v1.NotifierConsole
func (ctl Control) FindChannelNotifierConsole(ctx echo.Context) error {
	r, err := vault.NewNotifierConsole(ctl.db.Engine().NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "query console channel notifier"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Find channel notifier webhook
// @Accept      json
// @Produce     json
// @Tags        server/channel_notifier
// @Router      /server/channel_notifier/webhook [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q                   query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o                   query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p                   query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v1.NotifierWebhook
func (ctl Control) FindChannelNotifierWebhook(ctx echo.Context) error {
	r, err := vault.NewNotifierWebhook(ctl.db.Engine().NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "query webhook channel notifier"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Find channel notifier
// @Accept      json
// @Produce     json
// @Tags        server/channel_notifier
// @Router      /server/channel_notifier/rabbitmq [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q                   query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o                   query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p                   query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v1.NotifierRabbitMq
func (ctl Control) FindChannelNotifierRabbitmq(ctx echo.Context) error {
	r, err := vault.NewNotifierRabbitMq(ctl.db.Engine().NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "query rabbitmq channel notifier"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Get a channel notifier
// @Accept      json
// @Produce     json
// @Tags        server/channel_notifier
// @Router      /server/channel_notifier/console/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid                path   string true  "channel notifier 의 Uuid"
// @Success     200 {object} v1.NotifierConsole
func (ctl Control) GetChannelNotifierConsole(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := vault.NewNotifierConsole(ctl.db.Engine().NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get console channel notifier"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Get a channel notifier webhook
// @Accept      json
// @Produce     json
// @Tags        server/channel_notifier
// @Router      /server/channel_notifier/webhook/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid                path   string true  "channel notifier 의 Uuid"
// @Success     200 {object} v1.NotifierWebhook
func (ctl Control) GetChannelNotifierWebhook(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := vault.NewNotifierWebhook(ctl.db.Engine().NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get webhook vent notifier"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Get a channel notifier rabbitmq
// @Accept      json
// @Produce     json
// @Tags        server/channel_notifier
// @Router      /server/channel_notifier/rabbitmq/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid                path   string true  "channel notifier 의 Uuid"
// @Success     200 {object} v1.NotifierRabbitMq
func (ctl Control) GetChannelNotifierRabbitmq(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := vault.NewNotifierRabbitMq(ctl.db.Engine().NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get rabbitmq channel notifier"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Update a console channel notifier
// @Accept      json
// @Produce     json
// @Tags        server/channel_notifier
// @Router      /server/channel_notifier/console/{uuid} [put]
// @Param       x_auth_token header string                         false "client session token"
// @Param       uuid         path   string                         true  "Channel 의 Uuid"
// @Param       object       body   v1.NotifierConsole_create true  "NotifierConsole_create"
// @Success     200 {object} v1.NotifierConsole
func (ctl Control) UpdateChannelNotifierConsole(ctx echo.Context) error {
	body := new(channelv1.NotifierConsole_create)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if err := body.MIME.Valid(); err != nil && 0 < len(body.MIME.ContentType) {
		return errors.Wrapf(err, "valid MIME")
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	notifier := channelv1.NotifierConsole{}
	notifier.Uuid = uuid
	notifier.LabelMeta = body.LabelMeta
	notifier.NotifierConsoleProperty = body.NotifierConsoleProperty
	notifier.MIME = body.MIME

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		notifier_, err := vault.NewNotifierConsole(tx).Update(notifier)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "upate console channel notifier"))
		}

		return notifier_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Update a webhook channel notifier
// @Accept      json
// @Produce     json
// @Tags        server/channel_notifier
// @Router      /server/channel_notifier/webhook/{uuid} [put]
// @Param       x_auth_token header string                         false "client session token"
// @Param       uuid         path   string                         true  "Channel 의 Uuid"
// @Param       object       body   v1.NotifierWebhook_create true  "NotifierWebhook_create"
// @Success     200 {object} v1.NotifierWebhook
func (ctl Control) UpdateChannelNotifierWebhook(ctx echo.Context) error {
	body := new(channelv1.NotifierWebhook_create)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if err := body.MIME.Valid(); err != nil && 0 < len(body.MIME.ContentType) {
		return errors.Wrapf(err, "valid MIME")
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	notifier := channelv1.NotifierWebhook{}
	notifier.Uuid = uuid
	notifier.LabelMeta = body.LabelMeta
	notifier.NotifierWebhookProperty = body.NotifierWebhookProperty
	notifier.MIME = body.MIME

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		notifier_, err := vault.NewNotifierWebhook(tx).Update(notifier)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "upate webhook channel notifier"))
		}

		return notifier_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Update a rabbitmq channel notifier
// @Accept      json
// @Produce     json
// @Tags        server/channel_notifier
// @Router      /server/channel_notifier/rabbitmq/{uuid} [put]
// @Param       x_auth_token header string                          false "client session token"
// @Param       uuid         path   string                          true  "Channel 의 Uuid"
// @Param       object       body   v1.NotifierRabbitMq_create true  "NotifierRabbitMq_create"
// @Success     200 {object} v1.NotifierRabbitMq
func (ctl Control) UpdateChannelNotifierRabbitMq(ctx echo.Context) error {
	body := new(channelv1.NotifierRabbitMq_create)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if err := body.MIME.Valid(); err != nil && 0 < len(body.MIME.ContentType) {
		return errors.Wrapf(err, "valid MIME")
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	notifier := channelv1.NotifierRabbitMq{}
	notifier.Uuid = uuid
	notifier.LabelMeta = body.LabelMeta
	notifier.NotifierRabbitMqProperty = body.NotifierRabbitMqProperty
	notifier.MIME = body.MIME

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		notifier_, err := vault.NewNotifierRabbitMq(tx).Update(notifier)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "upate rabbitMq channel notifier"))
		}

		return notifier_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Delete a channel notifier console
// @Accept json
// @Produce json
// @Tags server/channel_notifier
// @Router /server/channel_notifier/console/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid                path   string true  "Channel 의 Uuid"
// @Success 200
func (ctl Control) DeleteChannelNotifierConsole(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		if err := vault.NewNotifierConsole(tx).Delete(uuid); err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "delete console channel notifier"))
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, OK())
}

// @Description Delete a channel notifier webhook
// @Accept json
// @Produce json
// @Tags server/channel_notifier
// @Router /server/channel_notifier/webhook/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid                path   string true  "Channel 의 Uuid"
// @Success 200
func (ctl Control) DeleteChannelNotifierWebhook(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		if err := vault.NewNotifierWebhook(tx).Delete(uuid); err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "delete webhook channel notifier"))
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, OK())
}

// @Description Delete a channel notifier rabbitmq
// @Accept json
// @Produce json
// @Tags server/channel_notifier
// @Router /server/channel_notifier/rabbitmq/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid                path   string true  "Channel 의 Uuid"
// @Success 200
func (ctl Control) DeleteChannelNotifierRabbitmq(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		if err := vault.NewNotifierRabbitMq(tx).Delete(uuid); err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "delete rabbitmq channel notifier"))
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, OK())
}
