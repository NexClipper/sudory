package operator

import "github.com/NexClipper/sudory/pkg/server/database"

type OperateContext struct {
	Db       *database.DBManipulator //Database
	Response ResponseFn              //Respose closuer
}
