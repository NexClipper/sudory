package control

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/sqlex"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"

	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	channelv3 "github.com/NexClipper/sudory/pkg/server/model/channel/v3"
	"github.com/NexClipper/sudory/pkg/server/model/tenants/v3"
	"github.com/NexClipper/sudory/pkg/server/status/globvar"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/itchyny/gojq"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// @Description Create a channel
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels [post]
// @Param       object       body   v3.HttpReq_ManagedChannel_create true  "HttpReq_ManagedChannel_create"
// @Success     200 {object} v3.HttpRsp_ManagedChannel
func (ctl ControlVanilla) CreateChannel(ctx echo.Context) (err error) {
	body := new(channelv3.HttpReq_ManagedChannel_create)
	err = echoutil.WrapHttpError(http.StatusBadRequest,
		func() error {
			err = echoutil.Bind(ctx, body)
			return errors.Wrapf(err, "bind%s",
				logs.KVL(
					"type", TypeName(body),
				))
		},
		func() error {
			if len(body.Name) == 0 {
				err = ErrorInvalidRequestParameter
			}
			return errors.Wrapf(err, "valid param%s",
				logs.KVL(
					ParamLog("name", body.Name)...,
				))
		},
		func() error {
			if false {
				if body.EventCategory == channelv3.EventCategoryNaV {
					err = ErrorInvalidRequestParameter
				}
				return errors.Wrapf(err, "valid param%s",
					logs.KVL(
						ParamLog("event_category", body.EventCategory)...,
					))

			} else {
				if body.EventCategory == channelv3.EventCategoryNaV {
					body.EventCategory = channelv3.EventCategoryNonspecified
				}
				return nil
			}
		},
	)
	if err != nil {
		return errors.WithStack(err)
	}
	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	// gen uuid
	body.Uuid = genUuidString(body.Uuid)

	time_now := time.Now()

	new_channel := channelv3.ManagedChannel{}
	new_channel.Uuid = body.Uuid
	new_channel.Name = body.Name
	new_channel.Summary = body.Summary
	new_channel.EventCategory = body.EventCategory
	new_channel.Created = time_now

	err = sqlex.ScopeTx(ctx.Request().Context(), ctl, func(tx *sql.Tx) error {
		var affected int64
		affedted, _, err := ctl.dialect.Insert(new_channel.TableName(), new_channel.ColumnNames(), new_channel.Values())(
			ctx.Request().Context(), tx)
		if err != nil {
			return errors.Wrapf(err, "save a new channel")
		}
		if affedted == 0 {
			return errors.Wrapf(database.ErrorNoAffected, "save a new channel")
		}

		// save tenant_clusters
		tenant_channels := new(tenants.TenantChannels)
		tenant_channels.TenantId = claims.ID
		tenant_channels.ChannelUuid = new_channel.Uuid
		affected, _, err = ctl.dialect.Insert(tenant_channels.TableName(), tenant_channels.ColumnNames(), tenant_channels.Values())(
			ctx.Request().Context(), tx)
		if err != nil {
			return errors.Wrapf(err, "save a new tenant channel")
		}
		if affected == 0 {
			return errors.Wrapf(database.ErrorNoAffected, "save a new tenant channel")
		}

		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create a new channel")
	}

	rsp := channelv3.HttpRsp_ManagedChannel{}
	rsp.ManagedChannel = new_channel

	return ctx.JSON(http.StatusOK, channelv3.HttpRsp_ManagedChannel(rsp))
}

// @Description Find channel
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels [get]
// @Param       q            query  string false "query  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       o            query  string false "order  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       p            query  string false "paging github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Success     200 {array} v3.HttpRsp_ManagedChannel
func (ctl ControlVanilla) FindChannel(ctx echo.Context) (err error) {
	q, err := stmt.ConditionLexer.Parse(echoutil.QueryParam(ctx)["q"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return HttpError(err, http.StatusBadRequest)
	}
	o, err := stmt.OrderLexer.Parse(echoutil.QueryParam(ctx)["o"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return HttpError(err, http.StatusBadRequest)
	}
	p, err := stmt.PaginationLexer.Parse(echoutil.QueryParam(ctx)["p"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return HttpError(err, http.StatusBadRequest)
	}
	// // additional conditon
	// q = stmt.And(q,
	// 	stmt.IsNull("deleted"),
	// )
	// default pagination
	if p == nil {
		p = stmt.Limit(__DEFAULT_DECORATION_LIMIT__)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	var (
		tables = []string{
			channelv3.TableNameWithTenant_ManagedChannel(claims.Hash),
			channelv3.TableNameWithTenant_ChannelStatusOption(claims.Hash),
			channelv3.TableNameWithTenant_Format(claims.Hash),
			channelv3.TableNameWithTenant_NotifierEdge(claims.Hash),
			channelv3.TableNameWithTenant_NotifierConsole(claims.Hash),
			channelv3.TableNameWithTenant_NotifierWebhook(claims.Hash),
			channelv3.TableNameWithTenant_NotifierRabbitMq(claims.Hash),
			channelv3.TableNameWithTenant_NotifierSlackhook(claims.Hash),
		}
		columns = [][]string{
			new(channelv3.ManagedChannel).ColumnNames(),
			new(channelv3.ChannelStatusOption).ColumnNames(),
			new(channelv3.Format).ColumnNames(),
			new(channelv3.NotifierEdge).ColumnNames(),
			new(channelv3.NotifierConsole).ColumnNames(),
			new(channelv3.NotifierWebhook).ColumnNames(),
			new(channelv3.NotifierRabbitMq).ColumnNames(),
			new(channelv3.NotifierSlackhook).ColumnNames(),
		}
	)

	var cond_keys []string = make([]string, 0, state.ENV__INIT_SLICE_CAPACITY__())
	if q != nil {
		cond_keys = append(cond_keys, q.Keys()...)
	}
	if o != nil {
		cond_keys = append(cond_keys, o.Keys()...)
	}

	uuids := make([]string, 0, state.ENV__INIT_SLICE_CAPACITY__())
	uuidSet := map[string]int{}
	search := func(table string, columns_a, columns_b []string) error {

		var uuid string
		var column_uuid = []string{"uuid"}

		if !IsIncluded(columns_a, columns_b) {
			return nil
		}

		return ctl.dialect.QueryRows(table, column_uuid, q, o, p)(ctx.Request().Context(), ctl)(
			func(scan excute.Scanner, _ int) error {
				err := scan.Scan(&uuid)
				if err != nil {
					err = errors.WithStack(err)
					return err
				}

				uuids = append(uuids, uuid)
				uuidSet[uuid] = 0

				return err
			})
	}

	for i := range tables {
		t := tables[i]
		a := columns[i]
		b := cond_keys

		if err := search(t, a, b); err != nil {
			return errors.Wrapf(err, "failed to search channel keys")
		}
	}

	rsp := make([]channelv3.HttpRsp_ManagedChannel, 0, len(uuidSet))
	for _, uuid := range uuids {
		if 0 < uuidSet[uuid] {
			continue
		}

		uuidSet[uuid] += 1

		channel, err := vault.GetManagedChannel(ctx.Request().Context(), ctl.DB, ctl.dialect, uuid, claims.Hash)
		if err != nil {
			return errors.Wrapf(err, "failed to search channels")
		}
		if channel == nil {
			continue
		}

		rsp = append(rsp, *channel)
	}

	return ctx.JSON(http.StatusOK, []channelv3.HttpRsp_ManagedChannel(rsp))
}

// @Description Get a channel
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid} [get]
// @Param       uuid         path   string true  "Channel 의 Uuid"
// @Success     200 {object} v3.HttpRsp_ManagedChannel
func (ctl ControlVanilla) GetChannel(ctx echo.Context) (err error) {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	channel, err := vault.GetManagedChannel(ctx.Request().Context(), ctl.DB, ctl.dialect, uuid, claims.Hash)
	if err != nil {
		return errors.Wrapf(err, "failed to get a channel")
	}
	if channel == nil {
		return errors.Wrapf(database.ErrorRecordWasNotFound, "failed to get a channel")
	}

	return ctx.JSON(http.StatusOK, (*channelv3.HttpRsp_ManagedChannel)(channel))
}

// @Description Update a channel
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid} [put]
// @Param       uuid         path   string                           true  "Channel 의 Uuid"
// @Param       object       body   v3.HttpReq_ManagedChannel_update true  "HttpReq_ManagedChannel_update"
// @Success     200 {object} v3.ManagedChannel
func (ctl ControlVanilla) UpdateChannel(ctx echo.Context) (err error) {
	body := new(channelv3.HttpReq_ManagedChannel_update)
	err = echoutil.Bind(ctx, body)
	err = errors.Wrapf(err, "bind%s",
		logs.KVL(
			"type", TypeName(body),
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	time_now := time.Now()

	var channel channelv3.ManagedChannel
	channel.Uuid = uuid
	channel_cond := stmt.And(
		stmt.Equal("uuid", channel.Uuid),
		stmt.IsNull("deleted"),
	)

	channel_table := channelv3.TableNameWithTenant_ManagedChannel(claims.Hash)

	err = ctl.dialect.QueryRow(channel_table, channel.ColumnNames(), channel_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := channel.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "get channel")
	}

	updateSet := map[string]interface{}{}
	if 0 < len(body.Name) {
		channel.Name = body.Name
		updateSet["name"] = channel.Name
	}
	if body.Summary.Valid {
		channel.Summary = body.Summary
		updateSet["summary"] = channel.Summary
	}
	if body.EventCategory != channelv3.EventCategoryNaV {
		channel.EventCategory = body.EventCategory
		updateSet["event_category"] = channel.EventCategory
	}

	// valied update column counts
	if len(updateSet) == 0 {
		return HttpError(errors.New("noting to update"), http.StatusBadRequest)
	}

	channel.Updated = *vanilla.NewNullTime(time_now)
	updateSet["updated"] = channel.Updated

	err = func() error {
		// update channel
		_, err := ctl.dialect.Update(channel.TableName(), updateSet, channel_cond)(
			ctx.Request().Context(), ctl)
		if err != nil {
			return errors.Wrapf(err, "update channel")
		}
		// if affected == 0 {
		// 	return errors.Wrapf(database.ErrorNoAffected, "update channel")
		// }
		return nil
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to update a channel")
	}

	return ctx.JSON(http.StatusOK, channelv3.ManagedChannel(channel))
}

// @Description Delete a channel
// @Security    ServiceAuth
// @Accept json
// @Produce json
// @Tags server/channels
// @Router /server/channels/{uuid} [delete]
// @Param       uuid         path   string true  "Channel 의 Uuid"
// @Success 200
func (ctl ControlVanilla) DeleteChannel(ctx echo.Context) (err error) {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]
	time_now := time.Now()

	var channel channelv3.ManagedChannel
	channel.Uuid = uuid
	channel_cond := stmt.And(
		stmt.Equal("uuid", uuid),
	)
	channel_table := channelv3.TableNameWithTenant_ManagedChannel(claims.Hash)

	err = ctl.dialect.QueryRow(channel_table, channel.ColumnNames(), channel_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := channel.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "get channel")
	}

	err = func() error {

		//update channel
		channel.Deleted = *vanilla.NewNullTime(time_now)
		updateSet := map[string]interface{}{
			"deleted": channel.Deleted,
		}

		affected, err := ctl.dialect.Update(channel.TableName(), updateSet, channel_cond)(
			ctx.Request().Context(), ctl)
		if err != nil {
			return errors.Wrapf(err, "update channel")
		}
		if affected == 0 {
			return errors.Wrapf(database.ErrorNoAffected, "update channel")
		}

		// clear channel status
		status := channelv3.ChannelStatus{}
		_, err = ctl.dialect.Delete(status.TableName(), channel_cond)(
			ctx.Request().Context(), ctl)
		if err != nil {
			return errors.Wrapf(err, "clear channel status")
		}

		return nil
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to delete a channel")
	}

	return ctx.JSON(http.StatusOK, OK())
}

// @Description Get a channel notifier edge
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/notifiers/edge [get]
// @Param       uuid         path   string true  "Channel 의 Uuid"
// @Success     200 {object} v3.HttpRsp_ManagedChannel_NotifierEdge
func (ctl ControlVanilla) GetChannelNotifierEdge(ctx echo.Context) (err error) {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	// NotifierEdge
	var edge channelv3.NotifierEdge
	edge.Uuid = uuid
	channel_cond := stmt.And(
		stmt.Equal("uuid", edge.Uuid),
		// stmt.IsNull("deleted"),
	)
	edge_table := channelv3.TableNameWithTenant_NotifierEdge(claims.Hash)

	err = ctl.dialect.QueryRow(edge_table, edge.ColumnNames(), channel_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := edge.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "failed to get %v", edge.TableName())
	}

	// get edge
	notifier_edge, err := vault.GetChannelNotifierEdge(ctx.Request().Context(), ctl.DB, ctl.dialect, edge)
	if err != nil {
		return errors.Wrapf(err, "failed to get channel notifier edge")
	}

	return ctx.JSON(http.StatusOK, (*channelv3.HttpRsp_ManagedChannel_NotifierEdge)(notifier_edge))
}

// @Description Get a channel status option
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/status/option [get]
// @Param       uuid         path   string true  "Channel 의 Uuid"
// @Success     200 {object} v3.HttpRsp_ManagedChannel_ChannelStatusOption
func (ctl ControlVanilla) GetChannelStatusOption(ctx echo.Context) (err error) {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}
	uuid := echoutil.Param(ctx)[__UUID__]

	var status_option channelv3.ChannelStatusOption
	status_option.Uuid = uuid
	status_option_cond := stmt.Equal("uuid", status_option.Uuid)

	status_option_table := channelv3.TableNameWithTenant_ChannelStatusOption(claims.Hash)
	err = ctl.dialect.QueryRow(status_option_table, status_option.ColumnNames(), status_option_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := status_option.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "failed to get a channel status option")
	}

	return ctx.JSON(http.StatusOK, channelv3.HttpRsp_ManagedChannel_ChannelStatusOption(status_option))
}

// @Description Update a console channel notifier
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/notifiers/console [put]
// @Param       uuid         path   string                         true  "Channel 의 Uuid"
// @Param       object       body   v3.HttpReq_ManagedChannel_NotifierConsole_update true  "HttpReq_ManagedChannel_NotifierConsole_update"
// @Success     200 {object} v3.HttpRsp_ManagedChannel_NotifierEdge
func (ctl ControlVanilla) UpdateChannelNotifierConsole(ctx echo.Context) (err error) {
	body := new(channelv3.HttpReq_ManagedChannel_NotifierConsole_update)
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
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	// time_now := time.Now()

	var notifier channelv3.NotifierConsole
	notifier.Uuid = uuid
	// notifier_cond := stmt.Equal("uuid", notifier.Uuid)
	// notifier_table := channelv3.TableNameWithTenant_NotifierConsole(claims.Hash)

	// err = ctl.dialect.Select(notifier_table, notifier.ColumnNames(), notifier_cond, nil, nil).
	// 	QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(
	// 	func(scan excute.Scanner) error {
	// 		return notifier.Scan(scan)
	// 	})
	// if err != nil {
	// 	return errors.Wrapf(err, "get a notifier console")
	// }
	// notifier.Created = *vanilla.NewNullTime(time_now)
	// notifier.Updated = *vanilla.NewNullTime(time_now)
	// notifier.Deleted = vanilla.NullTime{} // set null

	update_columns := []string{"uuid"}

	var edge *channelv3.NotifierEdge
	err = sqlex.ScopeTx(ctx.Request().Context(), ctl, func(tx *sql.Tx) error {
		var err error
		edge, err = updateChannelNotifier(ctx.Request().Context(), tx, ctl.dialect, claims.Hash, uuid, notifier, update_columns)
		return err
	})
	if err != nil {
		return errors.Wrapf(err, "failed to set a Console notifier")
	}
	if edge == nil {
		return errors.Wrapf(database.ErrorRecordWasNotFound, "failed to set a Console notifier")
	}

	rsp := channelv3.NotifierEdge_option{
		NotifierEdge: *edge,
		Console:      &notifier.NotifierConsole_property,
	}

	return ctx.JSON(http.StatusOK, channelv3.HttpRsp_ManagedChannel_NotifierEdge(rsp))
}

// @Description Update a rabbitmq channel notifier
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/notifiers/rabbitmq [put]
// @Param       uuid         path   string                                            true  "Channel 의 Uuid"
// @Param       object       body   v3.HttpReq_ManagedChannel_NotifierRabbitMq_update true  "HttpReq_ManagedChannel_NotifierRabbitMq_update"
// @Success     200 {object} v3.HttpRsp_ManagedChannel_NotifierEdge
func (ctl ControlVanilla) UpdateChannelNotifierRabbitMq(ctx echo.Context) (err error) {
	body := new(channelv3.HttpReq_ManagedChannel_NotifierRabbitMq_update)
	err = echoutil.Bind(ctx, body)
	err = errors.Wrapf(err, "bind%s",
		logs.KVL(
			"type", TypeName(body),
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	body.Url = strings.TrimSpace(body.Url)
	if len(body.Url) == 0 {
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog("Url", body.Url)...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	if len(body.ChannelPublish.Exchange.String) == 0 && len(body.ChannelPublish.RoutingKey.String) == 0 {
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog("Exchange", body.ChannelPublish.Exchange)...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// vaild RabbitMq config
	if err := body.Valid(); err != nil {
		return HttpError(
			errors.Wrapf(err, "failed to vaild RabbitMq config"),
			http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}
	uuid := echoutil.Param(ctx)[__UUID__]

	// updated := time.Now()

	var notifier channelv3.NotifierRabbitMq
	notifier.Uuid = uuid

	// notifier_cond := stmt.Equal("uuid", notifier.Uuid)
	// notifier_table := channelv3.TableNameWithTenant_NotifierRabbitMq(claims.Hash)

	// err = ctl.dialect.Select(notifier_table, notifier.ColumnNames(), notifier_cond, nil, nil).
	// 	QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(
	// 	func(scan excute.Scanner) error {
	// 		return notifier.Scan(scan)
	// 	})
	// if err != nil {
	// 	return errors.Wrapf(err, "get a notifier rabbitmq")
	// }

	// notifier.Created = *vanilla.NewNullTime(updated)
	// notifier.Updated = *vanilla.NewNullTime(updated)
	// notifier.Deleted = vanilla.NullTime{} // set null
	notifier.NotifierRabbitMq_property = *body

	update_columns := notifier.NotifierRabbitMq_property.ColumnNames()

	var edge *channelv3.NotifierEdge
	err = sqlex.ScopeTx(ctx.Request().Context(), ctl, func(tx *sql.Tx) error {
		var err error
		edge, err = updateChannelNotifier(ctx.Request().Context(), tx, ctl.dialect, claims.Hash, uuid, notifier, update_columns)
		return err
	})
	if err != nil {
		return errors.Wrapf(err, "failed to set a RabbitMq notifier")
	}
	if edge == nil {
		return errors.Wrapf(database.ErrorRecordWasNotFound, "failed to set a RabbitMq notifier")
	}

	rsp := channelv3.NotifierEdge_option{
		NotifierEdge: *edge,
		RabbitMq:     &notifier.NotifierRabbitMq_property,
	}

	return ctx.JSON(http.StatusOK, channelv3.HttpRsp_ManagedChannel_NotifierEdge(rsp))
}

// @Description Update a webhook channel notifier
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/notifiers/webhook [put]
// @Param       uuid         path   string                                           true  "Channel 의 Uuid"
// @Param       object       body   v3.HttpReq_ManagedChannel_NotifierWebhook_update true  "HttpReq_ManagedChannel_NotifierWebhook_update"
// @Success     200 {object} v3.HttpRsp_ManagedChannel_NotifierEdge
func (ctl ControlVanilla) UpdateChannelNotifierWebhook(ctx echo.Context) (err error) {
	body := new(channelv3.HttpReq_ManagedChannel_NotifierWebhook_update)
	err = echoutil.Bind(ctx, body)
	err = errors.Wrapf(err, "bind%s",
		logs.KVL(
			"type", TypeName(body),
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	body.Method = strings.TrimSpace(body.Method)
	body.Url = strings.TrimSpace(body.Url)

	if len(body.Method) == 0 {
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog("Method", body.Method)...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	if len(body.Url) == 0 {
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog("Url", body.Url)...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// vaild Webhook config
	if err := body.Valid(); err != nil {
		return HttpError(
			errors.Wrapf(err, "failed to vaild Webhook config"),
			http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}
	uuid := echoutil.Param(ctx)[__UUID__]

	// updated := time.Now()

	notifier := channelv3.NotifierWebhook{}
	notifier.Uuid = uuid

	// notifier_cond := stmt.Equal("uuid", notifier.Uuid)
	// notifier_table := channelv3.TableNameWithTenant_NotifierWebhook(claims.Hash)

	// err = ctl.dialect.Select(notifier_table, notifier.ColumnNames(), notifier_cond, nil, nil).
	// 	QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(
	// 	func(scan excute.Scanner) error {
	// 		return notifier.Scan(scan)
	// 	})
	// if err != nil {
	// 	return errors.Wrapf(err, "get a notifier webhook")
	// }

	// notifier.Created = *vanilla.NewNullTime(updated)
	// notifier.Updated = *vanilla.NewNullTime(updated)
	// notifier.Deleted = vanilla.NullTime{} // set null

	notifier.NotifierWebhook_property = *body

	notifier.RequestTimeout = func() uint {
		if notifier.RequestTimeout == 0 {
			return 10
		}
		return notifier.RequestTimeout
	}()

	update_columns := notifier.NotifierWebhook_property.ColumnNames()

	var edge *channelv3.NotifierEdge
	err = sqlex.ScopeTx(ctx.Request().Context(), ctl, func(tx *sql.Tx) error {
		var err error
		edge, err = updateChannelNotifier(ctx.Request().Context(), tx, ctl.dialect, claims.Hash, uuid, notifier, update_columns)
		return err
	})
	if err != nil {
		return errors.Wrapf(err, "failed to set a Webhook notifier")
	}
	if edge == nil {
		return errors.Wrapf(database.ErrorRecordWasNotFound, "failed to set a Webhook notifier")
	}

	rsp := channelv3.NotifierEdge_option{
		NotifierEdge: *edge,
		Webhook:      &notifier.NotifierWebhook_property,
	}

	return ctx.JSON(http.StatusOK, channelv3.HttpRsp_ManagedChannel_NotifierEdge(rsp))
}

// @Description Update a slackhook channel notifier
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/notifiers/slackhook [put]
// @Param       uuid         path   string                                             true  "Channel 의 Uuid"
// @Param       object       body   v3.HttpReq_ManagedChannel_NotifierSlackhook_update true  "HttpReq_ManagedChannel_NotifierSlackhook_update"
// @Success     200 {object} v3.HttpRsp_ManagedChannel_NotifierEdge
func (ctl ControlVanilla) UpdateChannelNotifierSlackhook(ctx echo.Context) (err error) {
	body := new(channelv3.HttpReq_ManagedChannel_NotifierSlackhook_update)
	err = echoutil.Bind(ctx, body)
	err = errors.Wrapf(err, "bind%s",
		logs.KVL(
			"type", TypeName(body),
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	body.Url = strings.TrimSpace(body.Url)
	if len(body.Url) == 0 {
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog("Url", body.Url)...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// vaild SlackWebhook config
	if err := body.Valid(); err != nil {
		return HttpError(
			errors.Wrapf(err, "failed to vaild SlackWebhook config"),
			http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}
	uuid := echoutil.Param(ctx)[__UUID__]

	// updated := time.Now()

	var notifier channelv3.NotifierSlackhook
	notifier.Uuid = uuid

	// notifier_cond := stmt.Equal("uuid", notifier.Uuid)
	// notifier_table := channelv3.TableNameWithTenant_NotifierRabbitMq(claims.Hash)

	// err = ctl.dialect.Select(notifier_table, notifier.ColumnNames(), notifier_cond, nil, nil).
	// 	QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(
	// 	func(scan excute.Scanner) error {
	// 		return notifier.Scan(scan)
	// 	})
	// if err != nil {
	// 	return errors.Wrapf(err, "get a notifier rabbitmq")
	// }

	// notifier.Created = *vanilla.NewNullTime(updated)
	// notifier.Updated = *vanilla.NewNullTime(updated)
	notifier.NotifierSlackhook_property = *body

	notifier.RequestTimeout = func() uint {
		if notifier.RequestTimeout == 0 {
			return 3
		}
		return notifier.RequestTimeout
	}()

	update_columns := notifier.NotifierSlackhook_property.ColumnNames()

	var edge *channelv3.NotifierEdge
	err = sqlex.ScopeTx(ctx.Request().Context(), ctl, func(tx *sql.Tx) error {
		var err error
		edge, err = updateChannelNotifier(ctx.Request().Context(), tx, ctl.dialect, claims.Hash, uuid, notifier, update_columns)
		return err
	})
	if err != nil {
		return errors.Wrapf(err, "failed to set a SlackHook notifier")
	}
	if edge == nil {
		return errors.Wrapf(database.ErrorRecordWasNotFound, "failed to set a SlackHook notifier")
	}

	rsp := channelv3.NotifierEdge_option{
		NotifierEdge: *edge,
		Slackhook:    &notifier.NotifierSlackhook_property,
	}

	return ctx.JSON(http.StatusOK, channelv3.HttpRsp_ManagedChannel_NotifierEdge(rsp))
}

// @Description Get a channel format
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/format [get]
// @Param       uuid         path   string true  "Channel 의 Uuid"
// @Success     200 {object} v3.HttpRsq_ManagedChannel_Format
func (ctl ControlVanilla) GetChannelFormat(ctx echo.Context) (err error) {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	var channel_format channelv3.Format
	channel_format.Uuid = uuid

	channel_format_cond := stmt.Equal("uuid", uuid)
	channel_format_table := channelv3.TableNameWithTenant_Format(claims.Hash)
	err = ctl.dialect.QueryRow(channel_format_table, channel_format.ColumnNames(), channel_format_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := channel_format.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "get a channel format")
	}
	return ctx.JSON(http.StatusOK, channelv3.HttpRsq_ManagedChannel_Format(channel_format))
}

// @Description Update a channel format
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/format [put]
// @Param       uuid         path   string                                  true  "Channel 의 Uuid"
// @Param       object       body   v3.HttpReq_ManagedChannel_Format_update true  "HttpReq_ManagedChannel_Format_update"
// @Success     200 {object} v3.HttpRsq_ManagedChannel_Format
func (ctl ControlVanilla) UpdateChannelFormat(ctx echo.Context) (err error) {
	body := new(channelv3.HttpReq_ManagedChannel_Format_update)
	err = echoutil.Bind(ctx, body)
	err = errors.Wrapf(err, "bind%s",
		logs.KVL(
			"type", TypeName(body),
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	switch body.FormatType {
	case channelv3.FormatTypeDisable:
		body.FormatData = "" // remove data
	case channelv3.FormatTypeFields:
		var ss []string
		err = json.Unmarshal([]byte(body.FormatData), &ss)
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog("FormatData", body.FormatData)...,
			))
		if err != nil {
			return HttpError(err, http.StatusBadRequest)
		}
	case channelv3.FormatTypeJq:
		_, err := gojq.Parse(body.FormatData)
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog("FormatData", body.FormatData)...,
			))
		if err != nil {
			return HttpError(err, http.StatusBadRequest)
		}
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	// updated := time.Now()

	// var channel channelv3.ManagedChannel
	// channel.Uuid = uuid
	// channel_cond := stmt.And(
	// 	stmt.Equal("uuid", channel.Uuid),
	// 	stmt.IsNull("deleted"),
	// )
	// channel_table := channelv3.TableNameWithTenant_ManagedChannel(claims.Hash)
	// exist, err := ctl.dialect.ExistContext(channel_table, channel_cond)(ctx.Request().Context(), ctl, ctl.Dialect())
	// if err != nil {
	// 	return errors.Wrapf(err, "failed to check a channel")
	// }
	// if !exist {
	// 	return errors.Wrapf(database.ErrorRecordWasNotFound, "failed to check a channel")
	// }
	err = vault.CheckManagedChannel(ctx.Request().Context(), ctl.DB, ctl.dialect, uuid, claims.Hash)
	if err != nil {
		return errors.Wrapf(err, "failed to check a channel")
	}

	var channel_format channelv3.Format
	channel_format.Uuid = uuid

	// channel_format_cond := stmt.Equal("uuid", uuid)
	// channel_format_table := channelv3.TableNameWithTenant_Format(claims.Hash)
	// err = ctl.dialect.Select(channel_format_table, channel_format.ColumnNames(), channel_format_cond, nil, nil).
	// 	QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(
	// 	func(scan excute.Scanner) (err error) {
	// 		return channel_format.Scan(scan)
	// 	})
	// if err != nil {
	// 	return errors.Wrapf(err, "get a channel format")
	// }

	// channel_format.Created = *vanilla.NewNullTime(updated)
	// channel_format.Updated = *vanilla.NewNullTime(updated)
	channel_format.Format_property = *body
	update_columns := channel_format.Format_property.ColumnNames()

	err = sqlex.ScopeTx(ctx.Request().Context(), ctl, func(tx *sql.Tx) error {
		// channel
		time_now := time.Now()
		var channel channelv3.ManagedChannel
		channel.Uuid = uuid
		channel_cond := stmt.And(
			stmt.Equal("uuid", channel.Uuid),
		)

		channel.Updated = *vanilla.NewNullTime(time_now)

		updateSet := map[string]interface{}{}
		updateSet["updated"] = channel.Updated

		_, err = ctl.dialect.Update(channel.TableName(), updateSet, channel_cond)(
			ctx.Request().Context(), tx)
		if err != nil {
			return errors.Wrapf(err, "update channel")
		}

		// channel format
		_, _, err = ctl.dialect.InsertOrUpdate(channel_format.TableName(), channel_format.ColumnNames(), update_columns, channel_format.Values())(
			ctx.Request().Context(), tx)
		if err != nil {
			return errors.Wrapf(err, "update channel format")
		}

		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "failed to update a channel format")
	}

	return ctx.JSON(http.StatusOK, channelv3.HttpRsq_ManagedChannel_Format(channel_format))
}

// @Description Update a channel status option
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/status/option [put]
// @Param       uuid         path   string                                               true  "Channel 의 Uuid"
// @Param       object       body   v3.HttpReq_ManagedChannel_ChannelStatusOption_update true  "HttpReq_ManagedChannel_ChannelStatusOption_update"
// @Success     200 {object} v3.HttpRsp_ManagedChannel_ChannelStatusOption
func (ctl ControlVanilla) UpdateChannelStatusOption(ctx echo.Context) (err error) {
	body := new(channelv3.HttpReq_ManagedChannel_ChannelStatusOption_update)
	err = echoutil.Bind(ctx, body)
	err = errors.Wrapf(err, "bind%s",
		logs.KVL(
			"type", TypeName(body),
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	// var channel channelv3.ManagedChannel
	// channel.Uuid = uuid
	// channel_cond := stmt.And(
	// 	stmt.Equal("uuid", channel.Uuid),
	// 	stmt.IsNull("deleted"),
	// )
	// channel_table := channelv3.TableNameWithTenant_ManagedChannel(claims.Hash)
	// exist, err := ctl.dialect.ExistContext(channel_table, channel_cond)(ctx.Request().Context(), ctl, ctl.Dialect())
	// if err != nil {
	// 	return errors.Wrapf(err, "failed to check a channel")
	// }
	// if !exist {
	// 	return errors.Wrapf(database.ErrorRecordWasNotFound, "failed to check a channel")
	// }
	err = vault.CheckManagedChannel(ctx.Request().Context(), ctl.DB, ctl.dialect, uuid, claims.Hash)
	if err != nil {
		return errors.Wrapf(err, "failed to check a channel")
	}

	status_option := channelv3.ChannelStatusOption{}
	status_option.Uuid = uuid
	// status_option.Created = *vanilla.NewNullTime(updated)
	// status_option.Updated = *vanilla.NewNullTime(updated)
	status_option.ChannelStatusOption_property = *body
	update_columns := status_option.ChannelStatusOption_property.ColumnNames()

	err = sqlex.ScopeTx(ctx.Request().Context(), ctl, func(tx *sql.Tx) error {
		// channel
		time_now := time.Now()
		var channel channelv3.ManagedChannel
		channel.Uuid = uuid
		channel_cond := stmt.And(
			stmt.Equal("uuid", channel.Uuid),
		)

		channel.Updated = *vanilla.NewNullTime(time_now)

		updateSet := map[string]interface{}{}
		updateSet["updated"] = channel.Updated

		_, err = ctl.dialect.Update(channel.TableName(), updateSet, channel_cond)(
			ctx.Request().Context(), tx)
		if err != nil {
			return errors.Wrapf(err, "update channel")
		}

		// channel status option
		_, _, err := ctl.dialect.InsertOrUpdate(status_option.TableName(), status_option.ColumnNames(), update_columns, status_option.Values())(
			ctx.Request().Context(), tx)
		if err != nil {
			return errors.Wrapf(err, "update channel status option")
		}

		return nil
	})

	if err != nil {
		return errors.Wrapf(err, "failed to update a channel status option")
	}

	return ctx.JSON(http.StatusOK, channelv3.HttpRsp_ManagedChannel_ChannelStatusOption(status_option))
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
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	if len(echoutil.QueryParam(ctx)["message"]) == 0 {
		err = ErrorInvalidRequestParameter
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

	err = vault.CreateChannelStatus(ctx.Request().Context(), ctl.DB, ctl.dialect, uuid, message, created, globvar.Event.NofitierStatusRotateLimit())
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, OK())
}

// @Description List channel status
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/status [get]
// @Param       uuid  path  string                                   true  "channel status 의 Uuid"
// @Success     200 {array} v3.HttpRsp_ManagedChannel_ChannelStatus
func (ctl ControlVanilla) ListChannelStatus(ctx echo.Context) (err error) {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
		))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	order := stmt.Asc("created")

	rsp := make([]channelv3.HttpRsp_ManagedChannel_ChannelStatus, 0, state.ENV__INIT_SLICE_CAPACITY__())
	var status channelv3.ChannelStatus
	status.Uuid = uuid
	status_cond := stmt.And(
		stmt.Equal("uuid", status.Uuid),
	)
	status_table := channelv3.TableNameWithTenant_ChannelStatus(claims.Hash)

	err = ctl.dialect.QueryRows(status_table, status.ColumnNames(), status_cond, order, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner, _ int) error {
			err := status.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			rsp = append(rsp, status)

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "failed to get channel status")
	}

	return ctx.JSON(http.StatusOK, []channelv3.HttpRsp_ManagedChannel_ChannelStatus(rsp))
}

// @Description Purge channel status
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/{uuid}/status/purge [delete]
// @Param       uuid         path   string true  "channel status 의 Uuid"
// @Success 200
func (ctl ControlVanilla) PurgeChannelStatus(ctx echo.Context) (err error) {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		err = ErrorInvalidRequestParameter
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
		))

	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}
	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	err = vault.CheckManagedChannel(ctx.Request().Context(), ctl.DB, ctl.dialect, uuid, claims.Hash)
	if err != nil {
		return errors.Wrapf(err, "failed to check a channel")
	}

	var status channelv3.ChannelStatus
	status.Uuid = uuid
	status_cond := stmt.Equal("uuid", uuid)

	err = func() error {
		_, err := ctl.dialect.Delete(status.TableName(), status_cond)(
			ctx.Request().Context(), ctl)
		return err
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to purge channel status")
	}

	return ctx.JSON(http.StatusOK, OK())
}

// @Description Find channel status
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/channels
// @Router      /server/channels/status [get]
// @Param       q            query  string false "query  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       o            query  string false "order  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       p            query  string false "paging github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Success     200 {array} v3.HttpRsp_ManagedChannel_ChannelStatus
func (ctl ControlVanilla) FindChannelStatus(ctx echo.Context) error {
	q, err := stmt.ConditionLexer.Parse(echoutil.QueryParam(ctx)["q"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return HttpError(err, http.StatusBadRequest)
	}
	o, err := stmt.OrderLexer.Parse(echoutil.QueryParam(ctx)["o"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return HttpError(err, http.StatusBadRequest)
	}
	p, err := stmt.PaginationLexer.Parse(echoutil.QueryParam(ctx)["p"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return HttpError(err, http.StatusBadRequest)
	}
	// additional conditon
	q = stmt.And(q,
		stmt.IsNull("deleted"),
	)
	// default pagination
	if p == nil {
		p = stmt.Limit(__DEFAULT_DECORATION_LIMIT__)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	rsp := make([]channelv3.HttpRsp_ManagedChannel_ChannelStatus, 0, state.ENV__INIT_SLICE_CAPACITY__())

	var status channelv3.ChannelStatus
	status_table := channelv3.TableNameWithTenant_ChannelStatus(claims.Hash)
	err = ctl.dialect.QueryRows(status_table, status.ColumnNames(), q, o, p)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner, _ int) error {
			err = status.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			rsp = append(rsp, status)

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "failed to search channel status")
	}

	return ctx.JSON(http.StatusOK, []channelv3.HttpRsp_ManagedChannel_ChannelStatus(rsp))
}

type notifier_table interface {
	TableName() string
	ColumnNames() []string
	Values() []interface{}
	Type() channelv3.NotifierType
}

func updateChannelNotifier(ctx context.Context, tx *sql.Tx, dialect excute.SqlExcutor, cluster_hash, channel_uuid string, notifier notifier_table, notifierUpdateColumns []string) (*channelv3.NotifierEdge, error) {

	time_now := time.Now()

	var channel channelv3.ManagedChannel
	channel.Uuid = channel_uuid
	channel.Updated = *vanilla.NewNullTime(time_now)

	channel_cond := stmt.And(
		stmt.Equal("uuid", channel.Uuid),
		stmt.IsNull("deleted"),
	)

	channel_table := channelv3.TableNameWithTenant_ManagedChannel(cluster_hash)
	err := dialect.QueryRow(channel_table, channel.ColumnNames(), channel_cond, nil, nil)(ctx, tx)(
		func(scan excute.Scanner) error {
			err := channel.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return nil, err
	}

	updateSet := map[string]interface{}{}
	updateSet["updated"] = channel.Updated

	_, err = dialect.Update(channel.TableName(), updateSet, channel_cond)(ctx, tx)
	if err != nil {
		return nil, err
	}

	edge := channelv3.NotifierEdge{}
	edge.Uuid = channel_uuid
	edge.NotifierType = notifier.Type()
	// edge.Created = *vanilla.NewNullTime(time_now)
	// edge.Updated = *vanilla.NewNullTime(time_now)

	edgeUpdateColumns := []string{
		"notifier_type",
	}

	// insert or update; notifier edge
	_, _, err = dialect.InsertOrUpdate(edge.TableName(), edge.ColumnNames(), edgeUpdateColumns, edge.Values())(ctx, tx)
	if err != nil {
		return nil, err
	}
	// if affected == 0 {
	// 	return errors.WithStack(database.ErrorNoAffected)
	// }

	// insert or update; notifier X
	_, _, err = dialect.InsertOrUpdate(notifier.TableName(), notifier.ColumnNames(), notifierUpdateColumns, notifier.Values())(ctx, tx)
	if err != nil {
		return nil, err
	}
	// if affected == 0 {
	// 	return errors.WithStack(database.ErrorNoAffected)
	// }

	return (*channelv3.NotifierEdge)(&edge), nil
}

// Is the slice B included in the A?
func IsIncluded(a, b []string) bool {
	for _, b := range b {
		ok := func(s string) bool {
			for _, a := range a {
				if s == a {
					return true
				}
			}
			return false
		}(b)
		if !ok {
			return false
		}
	}

	return true
}
