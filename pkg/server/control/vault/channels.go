package vault

import (
	"database/sql"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	channelv2 "github.com/NexClipper/sudory/pkg/server/model/channel/v2"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/pkg/errors"
)

// CreateChannelStatus
func CreateChannelStatus(db *sql.DB, uuid, message string, created time.Time, max_count uint) (err error) {
	// created := time.Now()
	channel_opt := channelv2.ManagedChannel_option{}
	channel_status := channelv2.ChannelStatus{
		Uuid:    uuid,
		Created: created,
		Message: message,
	}

	channel_status_opt_cond := vanilla.And(
		vanilla.Equal("uuid", uuid),
		vanilla.IsNull("deleted"),
	).Parse()

	// valid; exist channel
	err = vanilla.Stmt.Select(channel_opt.TableName(), channel_opt.ColumnNames(), channel_status_opt_cond, nil, nil).
		QueryRow(db)(func(scan vanilla.Scanner) (err error) {
		err = channel_opt.Scan(scan)
		return
	})
	err = errors.Wrapf(err, "failed to query to channel opt")
	if err != nil {
		return
	}

	if channel_opt.StatusOption.StatusMaxCount == 0 {
		channel_opt.StatusOption.StatusMaxCount = max_count // default channel status max_count
	}

	// rotation; find
	rotation_cond := vanilla.And(
		vanilla.Equal("uuid", uuid),
	).Parse()
	rotation_order := vanilla.Desc("created").Parse()
	rotation_limit := vanilla.Limit(int(channel_opt.StatusOption.StatusMaxCount)-1, 2).Parse()
	rotation_columns := []string{"uuid", "created"}

	uuids, createds := make([]string, 0, state.ENV__INIT_SLICE_CAPACITY__()), make([]vanilla.NullTime, 0, state.ENV__INIT_SLICE_CAPACITY__())
	err = vanilla.Stmt.Select(channel_status.TableName(), rotation_columns, rotation_cond, rotation_order, rotation_limit).
		QueryRows(db)(func(scan vanilla.Scanner, _ int) (err error) {
		var uuid string
		var created vanilla.NullTime
		err = scan.Scan(&uuid, &created)
		if err == nil {
			uuids = append(uuids, uuid)
			createds = append(createds, created)
		}
		return
	})
	if err != nil {
		return
	}

	// insert
	// err = vanilla.Scope(db, func(tx *sql.Tx) (err error) {
	err = func() (err error) {
		stmt, err := vanilla.Stmt.Insert(channel_status.TableName(), channel_status.ColumnNames(), channel_status.Values())
		if err != nil {
			return
		}

		affected, err := stmt.Exec(db)
		if err != nil {
			return
		}
		if affected == 0 {
			err = errors.Errorf("no affected")
			return
		}

		return
	}()
	if err != nil {
		return
	}

	err = func() (err error) {
		// rotation; remove
		for i := range uuids {
			uuid, created := uuids[i], createds[i]
			rm_cond := vanilla.And(
				vanilla.Equal("uuid", uuid),
				vanilla.Equal("created", created),
			).Parse()
			_, err = vanilla.Stmt.Delete(channel_status.TableName(), rm_cond).Exec(db)
			if err != nil {
				return
			}
		}

		return
	}()
	// })
	if err != nil {
		return
	}

	return
}

// func MakeHttpRsp_ManagedChannel(
// 	tx vanilla.Preparer,
// 	channel_uuid string,
// ) (rsp *channelv2.HttpRsp_ManagedChannel, err error) {

// 	// NS_EG := vanilla.NS("EG")
// 	// NS_CN := vanilla.NS("CN")

// 	// notifier_cond := vanilla.And(
// 	// 	vanilla.Equal(NS_EG("uuid"), channel_uuid),
// 	// ).Parse()
// 	channel_cond := vanilla.And(
// 		vanilla.Equal("uuid", channel_uuid),
// 		vanilla.IsNull("deleted"),
// 	).Parse()

// 	channel_tangled := new(channelv2.ManagedChannel_tangled)

// 	err = vanilla.Stmt.Select(channel_tangled.TableName(), channel_tangled.ColumnNames(), channel_cond, nil, nil).
// 		QueryRow(tx)(func(scan vanilla.Scanner) (err error) {
// 		err = channel_tangled.Scan(scan)
// 		return
// 	})
// 	err = errors.Wrapf(err, "failed to query from channel")
// 	if err != nil {
// 		return
// 	}

// 	rsp = &channelv2.HttpRsp_ManagedChannel{
// 		ManagedChannel_tangled: *channel_tangled,
// 	}

// 	return
// }
