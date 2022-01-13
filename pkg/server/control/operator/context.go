package operator

import "github.com/NexClipper/sudory/pkg/server/database"

type OperateContext struct {
	Db       *database.DBManipulator //Database
	Response ResponseFn              //Respose closuer
}

type KeyValueParam struct {
	OperateContext

	Params map[string]string //key value pair
}

func NewKeyValueParam(ctx OperateContext) *KeyValueParam {
	return &KeyValueParam{
		OperateContext: ctx,
		Params:         make(map[string]string),
	}
}
