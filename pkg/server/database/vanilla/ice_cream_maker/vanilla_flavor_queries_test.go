package ice_cream_maker_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/error_compose"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/ice_cream_maker"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/prepare"
	"github.com/pkg/errors"
)

func TestVanillaFlavorQueries(t *testing.T) {

	objs := []interface{}{
		ServiceStep_essential{},
		ServiceStep{},
	}

	s, err := ice_cream_maker.VanillaFlavorQueries(objs...)
	if err != nil {
		t.Fatal(err)
	}

	println(s)
}

const __DEFAULT_ARGS_CAPACITY__ = 10

type Condition interface {
	Query() string
	Args() []interface{}
}

type Order interface {
	Order() string
}

type Pagination interface {
	Offset() int
	Limit() int
}

func Demo_QueryRow(tx vanilla.Preparer, q *prepare.Condition, o *prepare.Order, p *prepare.Pagination) (r *ServiceStep, err error) {
	r = new(ServiceStep)
	column_names := strings.Join(ServiceStep_ColumnNames(), ", ")
	table_name := ServiceStep_TableName()
	args := make([]interface{}, 0, __DEFAULT_ARGS_CAPACITY__)
	s := fmt.Sprintf("SELECT %v FROM %v", column_names, table_name)
	if q != nil {
		s += "\nWHERE " + q.Query()
		args = append(args, q.Args()...)
	}
	if o != nil {
		s += "\nORDER BY " + o.Order()
	}
	if p != nil {
		s += fmt.Sprintf("\nLIMIT %v, %v", p.Offset(), p.Limit())
	}

	stmt, err := tx.Prepare(s)
	err = errors.Wrapf(err, "sql.DB.Prepare")
	if err != nil {
		return
	}
	defer func() {
		err = error_compose.Composef(err, stmt.Close(), "sql.Stmt.Close")
	}()

	row := stmt.QueryRow(args...)

	err = row.Scan(r.Dests()...)
	err = errors.Wrapf(err, "sql.Row.Scan")
	err = error_compose.Composef(err, row.Err(), "sql.Row; during scan")

	err = errors.Wrapf(err, "faild to query row table=\"%v\"",
		table_name,
	)

	return
}

func ServiceStep_TableName() string {
	return "service_step"
}

func ServiceStep_ColumnNames() []string {
	return []string{
		"foo", "bar",
	}
}
