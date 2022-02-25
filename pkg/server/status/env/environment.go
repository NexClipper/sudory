package env

import (
	"fmt"
	"os"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/status"

	envv1 "github.com/NexClipper/sudory/pkg/server/model/environment/v1"
)

type EnvironmentChron struct {
	ctx database.Context
}

func NewEnvironmentChron(ctx database.Context) status.ChronUpdater {
	return &EnvironmentChron{ctx: ctx}
}

var _ status.ChronUpdater = (*EnvironmentChron)(nil)

// Update
//  Update = read -> os.Setenv
func (chron *EnvironmentChron) Update() error {

	const where = ""

	records, err := chron.ctx.FindEnvironment(where)
	if err != nil {
		return fmt.Errorf("enviroment chron error: %w", err)
	}

	foreach_environment(envv1.TransFormDbSchema(records), func(elem envv1.Environment) bool {
		os.Setenv(*elem.Name, *elem.Value)
		return true
	})

	return nil
}

// WhiteListCheck
//  리스트 체크
func (chron *EnvironmentChron) WhiteListCheck() error {

	const where = ""

	records, err := chron.ctx.FindEnvironment(where)
	if err != nil {
		return fmt.Errorf("enviroment chron error: %w", err)
	}

	count := 0
	appnd, build := macro.StringBuilder()

	for _, key := range EnvNames() {
		var found bool = false
		foreach_environment(envv1.TransFormDbSchema(records), func(elem envv1.Environment) bool {

			if key == *elem.Name {
				found = true
				return false
			}
			return true
		})
		if !found {
			count++
			appnd(key)
		}
	}
	if 0 < count {
		return fmt.Errorf("not found environment keys=['%s']", build("', '"))
	}

	return nil
}

func (chron *EnvironmentChron) Merge() error {
	const where = ""

	records, err := chron.ctx.FindEnvironment(where)
	if err != nil {
		return fmt.Errorf("enviroment chron error: %w", err)
	}

	for _, key := range EnvNames() {
		var found bool = false
		foreach_environment(envv1.TransFormDbSchema(records), func(elem envv1.Environment) bool {

			if key == *elem.Name {
				found = true
				return false
			}
			return true
		})
		if !found {
			key, err := ParseEnv(key)
			if err != nil {
				return err
			}
			value, ok := DefaultEnvironmanets[key]
			if !ok {
				return fmt.Errorf("not found environment key")
			}
			value_ := Convert(key, value)

			err = chron.ctx.CreateEnvironment(envv1.DbSchemaEnvironment{Environment: value_})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func foreach_environment(elems []envv1.Environment, fn func(elem envv1.Environment) bool) {
	for n := range elems {
		ok := fn(elems[n])
		if !ok {
			return
		}
	}
}
