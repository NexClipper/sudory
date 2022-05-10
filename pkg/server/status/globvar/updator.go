package globvar

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type GlobalVariantUpdate struct {
	tx     *xorm.Session
	offset time.Time //updated column
}

func NewGlobalVariantUpdate(tx *xorm.Session) *GlobalVariantUpdate {
	return &GlobalVariantUpdate{tx: tx}
}

// Update
//  Update = read -> os.Setenv
func (worker *GlobalVariantUpdate) Update() error {
	where := "updated > ?"
	args := []interface{}{
		worker.offset,
	}

	records, err := vault.NewGlobalVariant(worker.tx).Find(where, args...)
	if err != nil {
		return errors.Wrapf(err, "find global variant")
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
	records, err := vault.NewGlobalVariant(worker.tx).Query(map[string]string{})
	if err != nil {
		return errors.Wrapf(err, "find global variant")
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
	records, err := vault.NewGlobalVariant(worker.tx).Query(map[string]string{})
	if err != nil {
		return errors.Wrapf(err, "find global variant")
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
			if _, err := vault.NewGlobalVariant(worker.tx).Create(value_); err != nil {
				return errors.Wrapf(err, "create global variant")
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
