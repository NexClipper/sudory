package control

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	channelv2 "github.com/NexClipper/sudory/pkg/server/model/channel/v2"
	"github.com/NexClipper/sudory/pkg/server/status/globvar"
	"github.com/itchyny/gojq"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// @Description Create a channel
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels [post]
// @Param       x_auth_token header string                           false "client session token"
// @Param       object       body   v2.HttpReq_ManagedChannel_create true  "HttpReq_ManagedChannel_create"
// @Success     200 {object} v2.HttpRsp_ManagedChannel
func (ctl ControlVanilla) CreateChannel(ctx echo.Context) (err error) {
	body := new(channelv2.HttpReq_ManagedChannel_create)
	err = echoutil.Bind(ctx, body)
	err = errors.Wrapf(err, "bind%s",
		logs.KVL(
			"type", TypeName(body),
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	if len(body.Name) == 0 {
		err = ErrorInvalidRequestParameter()
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog("Name", body.Name)...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}
	if false {
		if body.EventCategory == channelv2.EventCategoryNaV {
			err = ErrorInvalidRequestParameter()
		}
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog("EventCategory", body.EventCategory)...,
			))
		if err != nil {
			return HttpError(err, http.StatusBadRequest)
		}
	} else {
		if body.EventCategory == channelv2.EventCategoryNaV {
			body.EventCategory = channelv2.EventCategoryNonspecified
		}
	}
	if len(body.Uuid) == 0 {
		body.Uuid = NewUuidString() // len(body.Uuid) == 0; create uuid
	}

	created := time.Now()

	channel := channelv2.ManagedChannel{}
	channel.Created = created
	channel.Uuid = body.Uuid
	channel.Name = body.Name
	channel.Summary = body.Summary
	channel.EventCategory = body.EventCategory

	err = ctl.Scope(func(tx *sql.Tx) (err error) {
		insert_stmt, err := vanilla.Stmt.Insert(channel.TableName(), channel.ColumnNames(), channel.Values())
		if err != nil {
			return
		}
		affedted, err := insert_stmt.Exec(tx)
		if err != nil {
			return
		}
		if affedted == 0 {
			return errors.Errorf("no affedcted")
		}
		return
	})
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	rsp := channelv2.HttpRsp_ManagedChannel{}
	rsp.ManagedChannel = channel

	return ctx.JSON(http.StatusOK, rsp)
}

// @Description Find channel
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v2.HttpRsp_ManagedChannel
func (ctl ControlVanilla) FindChannel(ctx echo.Context) (err error) {
	q, o, p, err := ParseDecoration(echoutil.QueryParam(ctx))
	err = errors.Wrapf(err, "ParseDecoration%v", logs.KVL(
		"query", echoutil.QueryParamString(ctx),
	))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	rsp := make([]channelv2.HttpRsp_ManagedChannel, 0, __INIT_SLICE_CAPACITY__())
	channel_tangled := new(channelv2.ManagedChannel_tangled)
	err = vanilla.Stmt.Select(channel_tangled.TableName(), channel_tangled.ColumnNames(), q, o, p).
		QueryRows(ctl)(func(scan vanilla.Scanner, _ int) (err error) {
		err = channel_tangled.Scan(scan)
		if err == nil {
			rsp = append(rsp, channelv2.HttpRsp_ManagedChannel{ManagedChannel_tangled: *channel_tangled})
		}
		return
	})
	err = errors.Wrapf(err, "failed to query from channel tangled")
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, rsp)
}

// @Description Get a channel
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Channel 의 Uuid"
// @Success     200 {object} v2.HttpRsp_ManagedChannel
func (ctl ControlVanilla) GetChannel(ctx echo.Context) (err error) {
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

	// // make HttpRsp_ManagedChannel
	// rsp, err := vault.MakeHttpRsp_ManagedChannel(ctl, uuid)
	// if err != nil {
	// 	return HttpError(err, http.StatusInternalServerError)
	// }

	channel_cond := vanilla.And(
		vanilla.Equal("uuid", uuid),
		vanilla.IsNull("deleted"),
	).Parse()

	rsp := channelv2.HttpRsp_ManagedChannel{}
	channel_tangled := new(channelv2.ManagedChannel_tangled)
	err = vanilla.Stmt.Select(channel_tangled.TableName(), channel_tangled.ColumnNames(), channel_cond, nil, nil).
		QueryRow(ctl)(func(scan vanilla.Scanner) (err error) {
		err = rsp.ManagedChannel_tangled.Scan(scan)
		return
	})
	err = errors.Wrapf(err, "failed to query from channel tangled")
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, rsp)
}

// @Description Update a channel
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid} [put]
// @Param       x_auth_token header string                           false "client session token"
// @Param       uuid         path   string                           true  "Channel 의 Uuid"
// @Param       object       body   v2.HttpReq_ManagedChannel_update true  "HttpReq_ManagedChannel_update"
// @Success     200
func (ctl ControlVanilla) UpdateChannel(ctx echo.Context) (err error) {
	body := new(channelv2.HttpReq_ManagedChannel_update)
	err = echoutil.Bind(ctx, body)
	err = errors.Wrapf(err, "bind%s",
		logs.KVL(
			"type", TypeName(body),
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

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

	updated := time.Now()

	channel := channelv2.ManagedChannel{}
	channel.Updated = *vanilla.NewNullTime(updated)
	channel.Name = body.Name
	channel.Summary = body.Summary
	channel.EventCategory = body.EventCategory

	channel_set := map[string]interface{}{}
	channel_set["updated"] = channel.Updated
	if 0 < len(channel.Name) {
		channel_set["name"] = channel.Name
	}
	if channel.Summary.Valid {
		channel_set["summary"] = channel.Summary
	}
	if channel.EventCategory != channelv2.EventCategoryNaV {
		channel_set["event_category"] = channel.EventCategory
	}

	channel_cond := vanilla.And(
		vanilla.Equal("uuid", uuid),
		vanilla.IsNull("deleted"),
	).Parse()

	err = ctl.Scope(func(tx *sql.Tx) (err error) {
		var affected int64
		// update channel
		affected, err = vanilla.Stmt.Update(channel.TableName(), channel_set, channel_cond).
			Exec(tx)
		if err != nil {
			return
		}
		if affected == 0 {
			err = errors.Errorf("no affected")
			return
		}
		return
	})
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, OK())
}

// @Description Delete a channel
// @Accept json
// @Produce json
// @Tags server/channels
// @Router /server/channels/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Channel 의 Uuid"
// @Success 200
func (ctl ControlVanilla) DeleteChannel(ctx echo.Context) (err error) {
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

	cond := vanilla.And(
		vanilla.Equal("uuid", uuid),
	).Parse()

	deleted := time.Now()

	stmts := []func(tx vanilla.Preparer) error{
		// managed_channel
		func(tx vanilla.Preparer) error {
			channel := channelv2.ManagedChannel{}
			channel.Deleted = *vanilla.NewNullTime(deleted)
			set := map[string]interface{}{
				"deleted": channel.Deleted,
			}
			affected, err := vanilla.Stmt.Update(channel.TableName(), set, cond).
				Exec(tx)
			if err != nil {
				return err
			}
			if affected == 0 {
				return errors.Errorf("no affected")
			}
			return nil
		},
		// managed_channel_status
		func(tx vanilla.Preparer) error {
			status := channelv2.ChannelStatus{}
			_, err := vanilla.Stmt.Delete(status.TableName(), cond).
				Exec(tx)
			if err != nil {
				return err
			}

			return nil
		},
	}

	err = ctl.Scope(func(tx *sql.Tx) (err error) {
		for _, item := range stmts {
			err = item(tx)
			if err != nil {
				return
			}
		}

		return
	})
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, OK())
}

// @Description Get a channel notifier edge
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/notifiers/edge [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Channel 의 Uuid"
// @Success     200 {object} v2.HttpRsp_ManagedChannel_NotifierEdge
func (ctl ControlVanilla) GetChannelNotifierEdge(ctx echo.Context) (err error) {
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

	eq_uuid := vanilla.Equal("uuid", uuid)

	notifier_edge := channelv2.NotifierEdge_option{}
	err = vanilla.Stmt.Select(notifier_edge.TableName(), notifier_edge.ColumnNames(), eq_uuid.Parse(), nil, nil).
		QueryRow(ctl)(func(scan vanilla.Scanner) (err error) {
		err = notifier_edge.Scan(scan)
		return
	})
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, channelv2.HttpRsp_ManagedChannel_NotifierEdge{
		NotifierEdge_option: notifier_edge,
	})
}

// @Description Get a channel status option
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/status/option [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Channel 의 Uuid"
// @Success     200 {object} v2.HttpRsp_ManagedChannel_ChannelStatusOption
func (ctl ControlVanilla) GetChannelStatusOption(ctx echo.Context) (err error) {
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

	eq_uuid := vanilla.Equal("uuid", uuid)

	status_option := channelv2.ChannelStatusOption{}
	err = vanilla.Stmt.Select(status_option.TableName(), status_option.ColumnNames(), eq_uuid.Parse(), nil, nil).
		QueryRow(ctl)(func(scan vanilla.Scanner) (err error) {
		err = status_option.Scan(scan)
		return
	})
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, channelv2.HttpRsp_ManagedChannel_ChannelStatusOption{
		ChannelStatusOption: status_option,
	})
}

// @Description Update a console channel notifier
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/notifiers/console [put]
// @Param       x_auth_token header string                         false "client session token"
// @Param       uuid         path   string                         true  "Channel 의 Uuid"
// @Param       object       body   v2.HttpReq_ManagedChannel_NotifierConsole_update true  "HttpReq_ManagedChannel_NotifierConsole_update"
// @Success     200
func (ctl ControlVanilla) UpdateChannelNotifierConsole(ctx echo.Context) (err error) {
	body := new(channelv2.HttpReq_ManagedChannel_NotifierConsole_update)
	if false {
		err = echoutil.Bind(ctx, body)
		err = errors.Wrapf(err, "bind%s",
			logs.KVL(
				"type", TypeName(body),
			))
		if err != nil {
			return HttpError(err, http.StatusBadRequest)
		}
	}

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

	updated := time.Now()

	notifier := channelv2.NotifierConsole{}
	notifier.Uuid = uuid
	notifier.Created = *vanilla.NewNullTime(updated)
	notifier.Updated = *vanilla.NewNullTime(updated)
	// notifier.Deleted = vanilla.NullTime{} // set null

	update_columns := []string{"updated"}

	err = ctl.Scope(func(tx *sql.Tx) (err error) {
		return updateChannelNotifier(tx, uuid, notifier, update_columns)
	})

	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	// // make HttpRsp_ManagedChannel
	// rsp, err := makeHttpRsp_ManagedChannel(
	// 	func() vanilla.Preparer { return ctl.DB() },
	// 	uuid,
	// 	channelv2.NotifierTypeConsole,
	// )
	// if err != nil {
	// 	return HttpError(err, http.StatusInternalServerError)
	// }

	return ctx.JSON(http.StatusOK, OK())
}

// @Description Update a rabbitmq channel notifier
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/notifiers/rabbitmq [put]
// @Param       x_auth_token header string                                            false "client session token"
// @Param       uuid         path   string                                            true  "Channel 의 Uuid"
// @Param       object       body   v2.HttpReq_ManagedChannel_NotifierRabbitMq_update true  "HttpReq_ManagedChannel_NotifierRabbitMq_update"
// @Success     200
func (ctl ControlVanilla) UpdateChannelNotifierRabbitMq(ctx echo.Context) (err error) {
	body := new(channelv2.HttpReq_ManagedChannel_NotifierRabbitMq_update)
	err = echoutil.Bind(ctx, body)
	err = errors.Wrapf(err, "bind%s",
		logs.KVL(
			"type", TypeName(body),
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	if len(body.Url) == 0 {
		err = ErrorInvalidRequestParameter()
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog("Url", body.Url)...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	if len(body.ChannelPublish.Exchange.String) == 0 && len(body.ChannelPublish.RoutingKey.String) == 0 {
		err = ErrorInvalidRequestParameter()
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog("Exchange", body.ChannelPublish.Exchange)...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

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

	updated := time.Now()

	notifier := channelv2.NotifierRabbitMq{}
	notifier.Uuid = uuid
	notifier.Created = *vanilla.NewNullTime(updated)
	notifier.Updated = *vanilla.NewNullTime(updated)
	// notifier.Deleted = vanilla.NullTime{} // set null
	notifier.NotifierRabbitMq_essential = body.NotifierRabbitMq_essential

	update_columns := []string{
		"updated",
		"url",
		"exchange",
		"routing_key",
		"mandatory",
		"immediate",
		"message_headers",
		"message_content_type",
		"message_content_encoding",
		"message_delivery_mode",
		"message_priority",
		"message_correlation_id",
		"message_reply_to",
		"message_expiration",
		"message_message_id",
		"message_timestamp",
		"message_type",
		"message_user_id",
		"message_app_id",
	}

	err = ctl.Scope(func(tx *sql.Tx) (err error) {
		return updateChannelNotifier(tx, uuid, notifier, update_columns)
	})
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	// // make HttpRsp_ManagedChannel
	// rsp, err := makeHttpRsp_ManagedChannel(
	// 	func() vanilla.Preparer { return ctl.DB() },
	// 	uuid,
	// 	channelv2.NotifierTypeConsole,
	// )
	// if err != nil {
	// 	return HttpError(err, http.StatusInternalServerError)
	// }

	return ctx.JSON(http.StatusOK, OK())
}

// @Description Update a webhook channel notifier
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/notifiers/webhook [put]
// @Param       x_auth_token header string                         false "client session token"
// @Param       uuid         path   string                         true  "Channel 의 Uuid"
// @Param       object       body   v2.HttpReq_ManagedChannel_NotifierWebhook_update true  "HttpReq_ManagedChannel_NotifierWebhook_update"
// @Success     200
func (ctl ControlVanilla) UpdateChannelNotifierWebhook(ctx echo.Context) (err error) {
	body := new(channelv2.HttpReq_ManagedChannel_NotifierWebhook_update)
	err = echoutil.Bind(ctx, body)
	err = errors.Wrapf(err, "bind%s",
		logs.KVL(
			"type", TypeName(body),
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	if len(body.Method) == 0 {
		err = ErrorInvalidRequestParameter()
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog("Method", body.Method)...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	if len(body.Url) == 0 {
		err = ErrorInvalidRequestParameter()
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog("Url", body.Url)...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

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

	updated := time.Now()

	notifier := channelv2.NotifierWebhook{}
	notifier.Uuid = uuid
	notifier.Created = *vanilla.NewNullTime(updated)
	notifier.Updated = *vanilla.NewNullTime(updated)
	// notifier.Deleted = vanilla.NullTime{} // set null
	notifier.Method = body.Method
	notifier.Url = body.Url
	notifier.RequestHeaders = body.RequestHeaders
	notifier.RequestTimeout = func() uint {
		if body.RequestTimeout == 0 {
			return 10
		}
		return body.RequestTimeout
	}()

	update_columns := []string{
		"updated",
		"method",
		"url",
		"request_headers",
		"request_timeout",
	}

	err = ctl.Scope(func(tx *sql.Tx) (err error) {
		return updateChannelNotifier(tx, uuid, notifier, update_columns)
	})

	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	// // make HttpRsp_ManagedChannel
	// rsp, err := makeHttpRsp_ManagedChannel(
	// 	func() vanilla.Preparer { return ctl.DB() },
	// 	uuid,
	// 	channelv2.NotifierTypeConsole,
	// )
	// if err != nil {
	// 	return HttpError(err, http.StatusInternalServerError)
	// }

	return ctx.JSON(http.StatusOK, OK())
}

// @Description Get a channel format
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/format [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Channel 의 Uuid"
// @Success     200 {object} v2.HttpRsq_ManagedChannel_Format
func (ctl ControlVanilla) GetChannelFormat(ctx echo.Context) (err error) {
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

	eq_uuid := vanilla.Equal("uuid", uuid)

	channel_format := channelv2.Format{}
	err = vanilla.Stmt.Select(channel_format.TableName(), channel_format.ColumnNames(), eq_uuid.Parse(), nil, nil).
		QueryRow(ctl)(func(scan vanilla.Scanner) (err error) {
		err = channel_format.Scan(scan)
		return
	})
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, channelv2.HttpRsq_ManagedChannel_Format{
		Format: channel_format,
	})
}

// @Description Update a channel format
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/format [put]
// @Param       x_auth_token header string                                  false "client session token"
// @Param       uuid         path   string                                  true  "Channel 의 Uuid"
// @Param       object       body   v2.HttpReq_ManagedChannel_Format_update true  "HttpReq_ManagedChannel_Format_update"
// @Success     200
func (ctl ControlVanilla) UpdateChannelFormat(ctx echo.Context) (err error) {
	body := new(channelv2.HttpReq_ManagedChannel_Format_update)
	err = echoutil.Bind(ctx, body)
	err = errors.Wrapf(err, "bind%s",
		logs.KVL(
			"type", TypeName(body),
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

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

	switch body.FormatType {
	case channelv2.FormatTypeDisable:
		body.FormatData = "" // remove data
	case channelv2.FormatTypeFields:
		var ss []string
		err = json.Unmarshal([]byte(body.FormatData), &ss)
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog("FormatData", body.FormatData)...,
			))
		if err != nil {
			return HttpError(err, http.StatusBadRequest)
		}
	case channelv2.FormatTypeJq:
		_, err := gojq.Parse(body.FormatData)
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog("FormatData", body.FormatData)...,
			))
		if err != nil {
			return HttpError(err, http.StatusBadRequest)
		}
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	updated := time.Now()

	channel_format := channelv2.Format{}
	channel_format.Uuid = uuid
	channel_format.Created = *vanilla.NewNullTime(updated)
	channel_format.Updated = *vanilla.NewNullTime(updated)
	channel_format.FormatType = body.FormatType
	channel_format.FormatData = body.FormatData

	update_columns := []string{
		"format_type", "format_data", "updated",
	}

	err = ctl.Scope(func(tx *sql.Tx) (err error) {
		var insert_stmt *vanilla.StmtBuild
		insert_stmt, err = vanilla.Stmt.InsertOrUpdate(channel_format.TableName(), channel_format.ColumnNames(), update_columns, channel_format.Values())
		err = errors.Wrapf(err, "failed to build sql statement")
		if err != nil {
			return
		}

		_, err = insert_stmt.Exec(tx)
		if err != nil {
			return
		}

		return
	})

	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, OK())
}

// @Description Update a channel status option
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/status/option [put]
// @Param       x_auth_token header string                                               false "client session token"
// @Param       uuid         path   string                                               true  "Channel 의 Uuid"
// @Param       object       body   v2.HttpReq_ManagedChannel_ChannelStatusOption_update true  "HttpReq_ManagedChannel_ChannelStatusOption_update"
// @Success     200
func (ctl ControlVanilla) UpdateChannelStatusOption(ctx echo.Context) (err error) {
	body := new(channelv2.HttpReq_ManagedChannel_ChannelStatusOption_update)
	err = echoutil.Bind(ctx, body)
	err = errors.Wrapf(err, "bind%s",
		logs.KVL(
			"type", TypeName(body),
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

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

	updated := time.Now()

	status_option := channelv2.ChannelStatusOption{}
	status_option.Uuid = uuid
	status_option.Created = *vanilla.NewNullTime(updated)
	status_option.Updated = *vanilla.NewNullTime(updated)
	status_option.StatusMaxCount = body.StatusMaxCount

	update_columns := []string{"updated", "status_max_count"}

	err = ctl.Scope(func(tx *sql.Tx) (err error) {
		var insert_stmt *vanilla.StmtBuild
		insert_stmt, err = vanilla.Stmt.InsertOrUpdate(status_option.TableName(), status_option.ColumnNames(), update_columns, status_option.Values())
		err = errors.Wrapf(err, "failed to build sql statement")
		if err != nil {
			return
		}

		_, err = insert_stmt.Exec(tx)
		if err != nil {
			return
		}

		return
	})

	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, OK())
}

// @@Description Create a channel status
// @@Accept      json
// @@Produce     json
// @@Tags        server/channels
// @@Router      /server/channels/{uuid}/status [post]
// @@Param       x_auth_token header                           string false "client session token"
// @@Param       uuid    path                                  string true  "channel status 의 Uuid"
// @@Param       message query                                 string true  "message"
// @@Success     200
func (ctl ControlVanilla) CreateChannelStatus(ctx echo.Context) (err error) {
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

	if len(echoutil.QueryParam(ctx)["message"]) == 0 {
		err = ErrorInvalidRequestParameter()
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog("message", echoutil.Param(ctx)["message"])...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	uuid := echoutil.Param(ctx)[__UUID__]
	message := echoutil.QueryParam(ctx)["message"]
	created := time.Now()

	err = vault.CreateChannelStatus(ctl.DB, uuid, message, created, globvar.EventNofitierStatusRotateLimit())
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, OK())
}

// @Description List channel status
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/status [get]
// @Param       x_auth_token header string                           false "client session token"
// @Param       uuid  path  string                                   true  "channel status 의 Uuid"
// @Success     200 {array} v2.HttpRsp_ManagedChannel_ChannelStatus
func (ctl ControlVanilla) ListChannelStatus(ctx echo.Context) (err error) {
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

	cond := vanilla.And(
		vanilla.Equal("uuid", uuid),
	).Parse()
	order := vanilla.Asc("created").Parse()

	rsp := make([]channelv2.HttpRsp_ManagedChannel_ChannelStatus, 0, __INIT_SLICE_CAPACITY__())
	var status channelv2.ChannelStatus
	err = vanilla.Stmt.Select(status.TableName(), status.ColumnNames(), cond, order, nil).
		QueryRows(ctl)(func(scan vanilla.Scanner, _ int) (err error) {
		err = status.Scan(scan)
		err = errors.Wrapf(err, "list channel status%v", logs.KVL(
			"uuid", uuid,
		))
		if err != nil {
			return
		}

		rsp = append(rsp, channelv2.HttpRsp_ManagedChannel_ChannelStatus{
			ChannelStatus: status,
		})
		return
	})
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, rsp)
}

// // @Description Delete a channel status
// // @Accept json
// // @Produce json
// // @Tags server/channels
// // @Router /server/channels/{uuid}/status/{count} [delete]
// // @Param       x_auth_token header string false "client session token"
// // @Param       uuid         path   string true  "channel status 의 Uuid"
// // @Param       count         path   int true  "channel status 의 Uuid"
// // @Success 200
// func (ctl ControlVanilla) DeleteChannelStatus(ctx echo.Context) (err error) {
// 	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
// 		err = ErrorInvalidRequestParameter()
// 	}
// 	err = errors.Wrapf(err, "valid param%s",
// 		logs.KVL(
// 			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
// 		))

// 	if err != nil {
// 		return HttpError(err, http.StatusBadRequest)
// 	}

// 	if len(echoutil.Param(ctx)[__SEQUENCE__]) == 0 {
// 		err = ErrorInvalidRequestParameter()
// 	}
// 	err = errors.Wrapf(err, "valid param%s",
// 		logs.KVL(
// 			ParamLog(__SEQUENCE__, echoutil.Param(ctx)[__SEQUENCE__])...,
// 		))

// 	if err != nil {
// 		return HttpError(err, http.StatusBadRequest)
// 	}

// 	uuid := echoutil.Param(ctx)[__UUID__]
// 	sequence := echoutil.Param(ctx)[__SEQUENCE__]

// 	eq_uuid := vanilla.Equal("uuid", uuid)
// 	eq_sequence := vanilla.Equal("sequence", sequence)
// 	and := vanilla.And(eq_uuid, eq_sequence).Parse()

// 	var status channelv2.ChannelStatus

// 	err = ctl.Scope(func(tx *sql.Tx) (err error) {
// 		_, err = vanilla.Stmt.Delete(status.TableName(), and).
// 			Exec(tx)

// 		err = errors.Wrapf(err, "delete channel status%v", logs.KVL(
// 			"uuid", uuid,
// 		))

// 		return
// 	})
// 	if err != nil {
// 		return HttpError(err, http.StatusInternalServerError)
// 	}

// 	return ctx.JSON(http.StatusOK, OK())

// }

// @Description Purge channel status
// @Accept json
// @Produce json
// @Tags server/channels
// @Router /server/channels/{uuid}/status/purge [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "channel status 의 Uuid"
// @Success 200
func (ctl ControlVanilla) PurgeChannelStatus(ctx echo.Context) (err error) {
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

	var status channelv2.ChannelStatus

	err = ctl.Scope(func(tx *sql.Tx) (err error) {
		_, err = vanilla.Stmt.Delete(status.TableName(), eq_uuid).
			Exec(tx)

		err = errors.Wrapf(err, "purge channel status%v", logs.KVL(
			"uuid", uuid,
		))

		return
	})
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, OK())
}

// @Description Find channel status
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/status [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v2.HttpRsp_ManagedChannel_ChannelStatus
func (ctl ControlVanilla) FindChannelStatus(ctx echo.Context) error {
	q, o, p, err := ParseDecoration(echoutil.QueryParam(ctx))
	err = errors.Wrapf(err, "ParseDecoration%v", logs.KVL(
		"query", echoutil.QueryParamString(ctx),
	))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	rsp := make([]channelv2.HttpRsp_ManagedChannel_ChannelStatus, 0, __INIT_SLICE_CAPACITY__())

	var status channelv2.ChannelStatus
	err = vanilla.Stmt.Select(status.TableName(), status.ColumnNames(), q, o, p).
		QueryRows(ctl)(func(scan vanilla.Scanner, _ int) (err error) {
		err = status.Scan(scan)
		err = errors.Wrapf(err, "find channel status")
		if err != nil {
			return
		}

		rsp = append(rsp, channelv2.HttpRsp_ManagedChannel_ChannelStatus{
			ChannelStatus: status,
		})
		return
	})
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, rsp)
}

type notifier_table interface {
	TableName() string
	ColumnNames() []string
	Values() []interface{}
	Type() channelv2.NotifierType
}

func updateChannelNotifier(tx *sql.Tx, channel_uuid string, notifier notifier_table, notifier_update_columns []string) (err error) {

	updated := time.Now()

	edge := channelv2.NotifierEdge{}
	edge.Uuid = channel_uuid
	edge.NotifierType = notifier.Type()
	edge.Created = *vanilla.NewNullTime(updated)
	edge.Updated = *vanilla.NewNullTime(updated)

	edge_update_columns := []string{
		"notifier_type",
		"updated",
	}

	// insert or update; notifier edge
	var insert_stmt *vanilla.StmtBuild
	insert_stmt, err = vanilla.Stmt.InsertOrUpdate(edge.TableName(), edge.ColumnNames(), edge_update_columns, edge.Values())
	err = errors.Wrapf(err, "failed to build sql statement")
	if err != nil {
		return
	}

	_, err = insert_stmt.Exec(tx)
	if err != nil {
		return
	}

	// insert or update; notifier X
	insert_stmt, err = vanilla.Stmt.InsertOrUpdate(notifier.TableName(), notifier.ColumnNames(), notifier_update_columns, notifier.Values())
	err = errors.Wrapf(err, "failed to build sql statement")
	if err != nil {
		return
	}

	_, err = insert_stmt.Exec(tx)
	if err != nil {
		return
	}

	return
}
