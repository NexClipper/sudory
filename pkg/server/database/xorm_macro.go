package database

import (
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
)

func XormGet(get func(i interface{}) (bool, error), i interface{}) error {
	has, err := get(i)
	if err != nil {
		return errors.WithMessagef(err, "xorm get%v",
			logs.KVL(
				"type_name", TypeName(i),
			))
	}
	if !has {
		return ErrorRecordWasNotFound()
	}
	return nil
}
