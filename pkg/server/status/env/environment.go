package env

import (
	"os"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"

	envv1 "github.com/NexClipper/sudory/pkg/server/model/environment/v1"
)

type EnvironmentUpdate struct {
	ctx    database.Context
	offset time.Time //updated column
}

func NewEnvironmentUpdate(ctx database.Context) *EnvironmentUpdate {
	return &EnvironmentUpdate{ctx: ctx}
}

// Update
//  Update = read -> os.Setenv
func (worker *EnvironmentUpdate) Update() error {
	where := "updated > ?"
	args := []interface{}{
		worker.offset,
	}
	records := make([]envv1.Environment, 0)
	if err := worker.ctx.Where(where, args...).Find(&records); err != nil {
		return errors.Wrapf(err, "Database Find")
	}

	for i := range records {
		os.Setenv(records[i].Name, *records[i].Value)
	}

	//update offset
	worker.offset = time.Now()

	return nil
}

// WhiteListCheck
//  리스트 체크
func (worker *EnvironmentUpdate) WhiteListCheck() error {
	records := make([]envv1.Environment, 0)
	if err := worker.ctx.Find(&records); err != nil {
		return errors.Wrapf(err, "Database Find")
	}

	count := 0
	push, pop := macro.StringBuilder()
	for _, key := range EnvNames() {
		var found bool = false
	LOOP:
		for i := range records {
			if key == records[i].Name {
				found = true
				break LOOP
			}
		}
		if !found {
			count++
			push(key)
		}
	}
	if 0 < count {
		return errors.Errorf("not exists environment keys=['%s']", pop("', '"))
	}

	return nil
}

func (worker *EnvironmentUpdate) Merge() error {
	records := make([]envv1.Environment, 0)
	if err := worker.ctx.Find(&records); err != nil {
		return errors.Wrapf(err, "Database Find")
	}

	for _, key := range EnvNames() {
		var found bool = false
	LOOP:
		for i := range records {
			if key == records[i].Name {
				found = true
				break LOOP
			}
		}
		if !found {
			env, err := ParseEnv(key)
			if err != nil {
				return errors.Wrapf(err, "ParseEnv%s",
					logs.KVL(
						"key", key,
					))
			}

			value, ok := DefaultEnvironments[env]
			if !ok {
				return errors.Errorf("not found default environment%s",
					logs.KVL(
						"key", key,
					))
			}

			value_ := Convert(env, value)
			if err = worker.ctx.Create(value_); err != nil {
				return errors.Wrapf(err, "database create")
			}
		}
	}

	return nil
}

// func foreach_environment(elems []envv1.Environment, fn func(elem envv1.Environment) bool) {
// 	for n := range elems {
// 		ok := fn(elems[n])
// 		if !ok {
// 			return
// 		}
// 	}
// }
