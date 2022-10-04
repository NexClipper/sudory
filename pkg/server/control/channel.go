package control

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	channelv1 "github.com/NexClipper/sudory/pkg/server/model/channel/v1"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

// @deprecated
// @Description Create a channel
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/channel
// @Router      /server/channel [post]
// @Param       object       body   v1.Channel_create true  "Event_create"
// @Success     200 {object} v1.Channel
func (ctl Control) CreateChannel(ctx echo.Context) error {

	body := new(channelv1.Channel_create)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject, "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if len(body.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.Name", TypeName(body)), body.Name)...,
				)))
	}

	// if len(body.ClusterUuid) == 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
	// 		errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
	// 			logs.KVL(
	// 				ParamLog(fmt.Sprintf("%s.ClusterUuid", TypeName(body)), body.ClusterUuid)...,
	// 			)))
	// }

	// if len(body.Pattern) == 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
	// 		errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
	// 			logs.KVL(
	// 				ParamLog(fmt.Sprintf("%s.Pattern", TypeName(body)), body.Pattern)...,
	// 			)))
	// }

	// // exists cluster
	// if _, err := vault.NewCluster(ctl.db.Engine().NewSession()).Get(body.ClusterUuid); err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
	// 		errors.Wrapf(err, "exists cluster"))
	// }

	//pattern regex
	if _, err := regexp.Compile(body.Name); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(err, "channel name to regexp compile expr"))
	}

	channel := channelv1.Channel{}
	channel.UuidMeta = metav1.NewUuidMeta()
	channel.LabelMeta = metav1.NewLabelMeta(body.Name, body.Summary)
	channel.ChannelProperty = body.ChannelProperty

	r := channelv1.ChannelWithEdges{}
	_, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		//create channel
		channel_, err := vault.NewChannel(tx).Create(channel)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "create channel"))
		}
		r.Channel = *channel_
		r.NotifierEdges = body.NotifierEdges

		//create channel edges
		if err := AddChannelNotifierEdges(tx, channel_.Uuid, body.NotifierEdges); err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "create channel edge"))
		}

		return channel_, err
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// @deprecated
// @Description Find channel
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/channel
// @Router      /server/channel [get]
// @Param       q            query  string false "query  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       o            query  string false "order  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       p            query  string false "paging github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Success     200 {array} v1.Channel
func (ctl Control) FindChannel(ctx echo.Context) error {
	//find channel
	channels, err := vault.NewChannel(ctl.db.Engine().NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "query channel"))
	}

	return ctx.JSON(http.StatusOK, channels)

}

// @deprecated
// @Description Get a channel
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/channel
// @Router      /server/channel/{uuid} [get]
// @Param       uuid         path   string true  "Channel 의 Uuid"
// @Success     200 {object} v1.Channel
func (ctl Control) GetChannel(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	//get channel
	channel, err := vault.NewChannel(ctl.db.Engine().NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get channel"))
	}

	return ctx.JSON(http.StatusOK, channel)
}

// @deprecated
// @Description Get channel edges
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/channel
// @Router      /server/channel/{uuid}/notifier_edges [get]
// @Param       uuid         path   string true  "Channel 의 Uuid"
// @Success     200 {array} v1.ChannelNotifierEdge
func (ctl Control) ListChannelNotifierEdges(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	//get channel
	channel, err := vault.NewChannel(ctl.db.Engine().NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get channel"))
	}

	//find edge
	edges, err := vault.NewChannelNotifierEdge(ctl.db.Engine().NewSession()).Find("channel_uuid = ?", channel.Uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "find edge"))
	}

	return ctx.JSON(http.StatusOK, edges)
}

// @deprecated
// @Description Update a channel
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/channel
// @Router      /server/channel/{uuid} [put]
// @Param       uuid         path   string          true  "Channel 의 Uuid"
// @Param       object       body   v1.Channel_update true  "Event_update"
// @Success     200 {object} v1.Channel
func (ctl Control) UpdateChannel(ctx echo.Context) error {
	body := new(channelv1.Channel_update)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject, "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	channel := channelv1.Channel{}
	channel.Uuid = uuid
	channel.LabelMeta = body.LabelMeta
	channel.ChannelProperty = body.ChannelProperty

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		channel_, err := vault.NewChannel(tx).Update(channel)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "update channel"))
		}

		// event = *event_
		return channel_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// @deprecated
// @Description addtion channel notifier edge
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/channel
// @Router      /server/channel/{uuid}/notifier_edges/add [put]
// @Param       uuid         path   string            true  "Channel 의 Uuid"
// @Param       object       body   []v1.NotifierEdge true "NotifierEdge"
// @Success     200 {array} v1.ChannelNotifierEdge
func (ctl Control) AddChannelNotifierEdge(ctx echo.Context) error {
	body := []channelv1.NotifierEdge{}
	if err := echoutil.Bind(ctx, &body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject, "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	//get channel
	channel, err := vault.NewChannel(ctl.db.Engine().NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get channel"))
	}

	//addtion channel notifier edge
	_, err = ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		if err := AddChannelNotifierEdges(tx, channel.Uuid, body); err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "addtion channel notifier edge"))
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	edges, err := vault.NewChannelNotifierEdge(ctl.db.Engine().NewSession()).Find("channel_uuid = ?", uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "find channel edge"))
	}

	return ctx.JSON(http.StatusOK, edges)
}

// @deprecated
// @Description subtraction channel sub notifier
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/channel
// @Router      /server/channel/{uuid}/notifier_edges/sub [put]
// @Param       uuid         path   string            true  "Channel 의 Uuid"
// @Param       object       body   []v1.NotifierEdge true  "Channel 의 NotifierEdge"
// @Success     200 {array} v1.ChannelNotifierEdge
func (ctl Control) SubChannelNotifierEdge(ctx echo.Context) error {
	body := []channelv1.NotifierEdge{}
	if err := echoutil.Bind(ctx, &body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject, "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	//get channel
	channel, err := vault.NewChannel(ctl.db.Engine().NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get channel"))
	}

	//subtraction channel sub notifier
	_, err = ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		for _, edge := range body {
			if err := vault.NewChannelNotifierEdge(tx).Delete(channel.Uuid, edge.NotifierUuid); err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
					errors.Wrapf(err, "subtraction channel notifier edge"))
			}
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	edges, err := vault.NewChannelNotifierEdge(ctl.db.Engine().NewSession()).Find("channel_uuid = ?", uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "find event edge"))
	}

	return ctx.JSON(http.StatusOK, edges)
}

// @deprecated
// @Description Delete a channel
// @Security    XAuthToken
// @Accept json
// @Produce json
// @Tags server/channel
// @Router /server/channel/{uuid} [delete]
// @Param       uuid         path   string true  "Channel 의 Uuid"
// @Success 200
func (ctl Control) DeleteChannel(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		//delete channel
		if err := vault.NewChannel(tx).Delete(uuid); err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "delete channel"))
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, OK())
}

func AddChannelNotifierEdges(tx *xorm.Session, channel_uuid string, edges []channelv1.NotifierEdge) error {
	for _, edge := range edges {
		//check notifier
		_, err := vault.NewChannelNotifier(tx).Get(edge.NotifierType, edge.NotifierUuid)
		if err != nil {
			return errors.Wrapf(err, "get channel notifier")
		}

		bind_edges, err := vault.NewChannelNotifierEdge(tx).
			Find("channel_uuid = ? AND notifier_type = ? AND notifier_uuid = ?",
				channel_uuid, edge.NotifierType, edge.NotifierUuid)
		if err != nil {
			return errors.Wrapf(err, "find channel notifier edge")
		}

		if 0 < len(bind_edges) {
			continue //already has
		}

		//create edge
		edge_ := channelv1.ChannelNotifierEdge{}
		edge_.ChannelUuid = channel_uuid
		edge_.NotifierType = edge.NotifierType
		edge_.NotifierUuid = edge.NotifierUuid

		if _, err := vault.NewChannelNotifierEdge(tx).Create(edge_); err != nil {
			return errors.Wrapf(err, "create channel notifier edge")
		}
	}

	return nil
}
