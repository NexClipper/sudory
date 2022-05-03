package globvar

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"

	globvarv1 "github.com/NexClipper/sudory/pkg/server/model/global_variant/v1"
)

type GlobalVariantUpdate struct {
	ctx    database.Context
	offset time.Time //updated column
}

func NewGlobalVariantUpdate(ctx database.Context) *GlobalVariantUpdate {
	return &GlobalVariantUpdate{ctx: ctx}
}

// Update
//  Update = read -> os.Setenv
func (worker *GlobalVariantUpdate) Update() error {
	where := "updated > ?"
	args := []interface{}{
		worker.offset,
	}
	records := make([]globvarv1.GlobalVariant, 0)
	if err := worker.ctx.Where(where, args...).Find(&records); err != nil {
		return errors.Wrapf(err, "database Find")
	}

	for i := range records {
		gv, err := ParseKey(records[i].Name)
		if err != nil {
			return errors.Wrapf(err, "parse record name to key%v",
				logs.KVL(
					"key", records[i].Name,
				))
		}

		if err := storeManager.Call(gv, *records[i].Value); err != nil {
			return errors.Wrapf(err, "store globalVariant%v",
				logs.KVL(
					"key", records[i].Name,
					"value", *records[i].Value,
				))
		}

	}

	//update offset
	worker.offset = time.Now()

	return nil
}

// WhiteListCheck
//  리스트 체크
func (worker *GlobalVariantUpdate) WhiteListCheck() error {
	records := make([]globvarv1.GlobalVariant, 0)
	if err := worker.ctx.Find(&records); err != nil {
		return errors.Wrapf(err, "database Find")
	}

	count := 0
	push, pop := macro.StringBuilder()
	for _, key := range KeyNames() {
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
		return errors.Errorf("not exists global_variant keys=['%s']", pop("', '"))
	}

	return nil
}

func (worker *GlobalVariantUpdate) Merge() error {
	records := make([]globvarv1.GlobalVariant, 0)
	if err := worker.ctx.Find(&records); err != nil {
		return errors.Wrapf(err, "database Find")
	}

	for _, key := range KeyNames() {
		var found bool = false
	LOOP:
		for i := range records {
			if key == records[i].Name {
				found = true
				break LOOP
			}
		}
		if !found {
			env, err := ParseKey(key)
			if err != nil {
				return errors.Wrapf(err, "ParseGlobVar%s",
					logs.KVL(
						"key", key,
					))
			}

			value, ok := defaultValueSet[env]
			if !ok {
				return errors.Errorf("not found global_variant variant%s",
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
