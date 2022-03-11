package env

import (
	"fmt"
	"os"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/status"
	"github.com/pkg/errors"

	envv1 "github.com/NexClipper/sudory/pkg/server/model/environment/v1"
)

type EnvironmentChron struct {
	ctx database.Context
}

func NewEnvironmentChron(ctx database.Context) *EnvironmentChron {
	return &EnvironmentChron{ctx: ctx}
}

var _ status.ChronUpdater = (*EnvironmentChron)(nil)

// Update
//  Update = read -> os.Setenv
func (chron *EnvironmentChron) Update() error {
	records := make([]envv1.DbSchema, 0)
	if err := chron.ctx.Find(&records); err != nil {
		return errors.Wrapf(err, "Database Find")
	}

	for i := range records {
		os.Setenv(*records[i].Name, *records[i].Value)
	}

	return nil
}

// WhiteListCheck
//  리스트 체크
func (chron *EnvironmentChron) WhiteListCheck() error {

	records := make([]envv1.DbSchema, 0)
	if err := chron.ctx.Find(&records); err != nil {
		return errors.Wrapf(err, "Database Find")
	}

	count := 0
	push, pos := macro.StringBuilder()
	for _, key := range EnvNames() {
		var found bool = false
	LOOP:
		for i := range records {
			if key == *records[i].Name {
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
		return fmt.Errorf("not exists environment keys=['%s']", pos("', '"))
	}

	return nil
}

func (chron *EnvironmentChron) Merge() error {
	records := make([]envv1.DbSchema, 0)
	if err := chron.ctx.Find(&records); err != nil {
		return errors.Wrapf(err, "Database Find")
	}

	for _, key := range EnvNames() {
		var found bool = false
	LOOP:
		for i := range records {
			if key == *records[i].Name {
				found = true
				break LOOP
			}
		}
		if !found {
			env, err := ParseEnv(key)
			if err != nil {
				return errors.Wrapf(err, "ParseEnv key=%s", key)
			}

			value, ok := DefaultEnvironmanets[env]
			if !ok {
				return fmt.Errorf("not found default environment key=%s", key)
			}

			value_ := Convert(env, value)
			if err = chron.ctx.Create(envv1.DbSchema{Environment: value_}); err != nil {
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
