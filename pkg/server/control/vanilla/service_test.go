package vanilla_test

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/control/vanilla"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	clusterv2 "github.com/NexClipper/sudory/pkg/server/model/cluster/v2"
	servicev2 "github.com/NexClipper/sudory/pkg/server/model/service/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const (
	__INIT_RECORD_CAPACITY__ = 5
)

var (
	db, err = sql.Open("mysql", "root:verysecret@tcp(localhost:3306)/sudory?parseTime=true")
)

func TestQueryRowCluster(t *testing.T) {

	uuid := "7dcf4f5eb62641c083614dccab3a044b"

	cond := vanilla.NewCond(
		"WHERE uuid = ?",
		uuid,
	)

	cluster := clusterv2.Cluster{}
	sacn := vanilla.QueryRow(db, cluster.TableName(), cluster.ColumnNames(), *cond)
	err = sacn(func(s vanilla.Scanner) (err error) {
		err = cluster.Scan(s)
		err = errors.Wrapf(err, "cluster Scan")
		return
	})

	Do(&err, func() (err error) {
		b, err := json.Marshal(cluster)
		if err != nil {
			println(err.Error())
		}
		println(string(b))

		err = json.Unmarshal(b, &cluster)
		if err != nil {
			println(err.Error())
		}
		return
	})

	err = errors.Wrapf(err, "query multi rows")
	t.Error(err)

}

func TestQueryRow(t *testing.T) {

	uuid := "s:1"

	cond := vanilla.NewCond(
		"WHERE uuid = ?",
		uuid,
	)

	service := servicev2.Service_tangled{}
	sacn := vanilla.QueryRow(db, service.TableName(), service.ColumnNames(), *cond)
	err = sacn(func(s vanilla.Scanner) (err error) {
		err = service.Scan(s)
		err = errors.Wrapf(err, "service Scan")
		return
	})

	err = errors.Wrapf(err, "query multi rows")
	t.Error(err)

}
func TestQueryMultiRows(t *testing.T) {

	uuid := "s:1"

	cond := vanilla.NewCond(
		"WHERE uuid = ?",
		uuid,
	)

	service := servicev2.Service_tangled{}
	sacn := vanilla.QueryRows(db, service.TableName(), service.ColumnNames(), *cond)
	err := sacn(func(s vanilla.Scanner) (err error) {
		err = service.Scan(s)
		err = errors.Wrapf(err, "service Scan")
		return
	})

	steps := make([]servicev2.ServiceStep_tangled, 0, 10)
	step := servicev2.ServiceStep_tangled{}
	err = vanilla.QueryRows(db, step.TableName(), step.ColumnNames(), *cond)(func(s vanilla.Scanner) (err error) {
		err = step.Scan(s)
		err = errors.Wrapf(err, "step Scan")
		if err == nil {
			steps = append(steps, step)
		}
		return
	})

	Do(&err, func() (err error) {
		b, err := json.Marshal(service)
		if err != nil {
			println(err.Error())
		}
		println(string(b))

		err = json.Unmarshal(b, &service)
		if err != nil {
			println(err.Error())
		}
		return
	})
	Do(&err, func() (err error) {
		b, err := json.Marshal(steps)
		if err != nil {
			println(err.Error())
		}
		println(string(b))

		err = json.Unmarshal(b, &steps)
		if err != nil {
			println(err.Error())
		}
		return
	})

	rsp := servicev2.HttpRsp_Service{Service: service, Steps: steps}

	Do(&err, func() (err error) {
		b, err := json.Marshal(rsp)
		if err != nil {
			println(err.Error())
		}
		println(string(b))

		err = json.Unmarshal(b, &rsp)
		if err != nil {
			println(err.Error())
		}
		return
	})

	err = errors.Wrapf(err, "query multi rows")
	t.Error(err)

}

func TestGetService(t *testing.T) {

	uuid := "s:1"

	condition := "WHERE uuid = ?"
	args := []interface{}{
		uuid,
	}
	service := servicev2.Service_tangled{}
	Do(&err, func() (err error) {
		args, stmt, err := vanilla.Select(service.TableName(), service.ColumnNames(), vanilla.Condition{Condition: condition, Args: args}).Prepare(db)
		Do(&err, func() (err error) {
			defer stmt.Close()
			rows, err := stmt.Query(args...)
			err = errors.Wrapf(err, "stmt query; service")
			Do(&err, func() (err error) {
				for rows.Next() {
					err = service.Scan(rows)
					err = errors.Wrapf(err, "scan row; service")
					if err != nil {
						return
					}
				}
				return
			})
			return
		})
		return
	})

	steps := make([]servicev2.ServiceStep_tangled, 0, __INIT_RECORD_CAPACITY__)
	Do(&err, func() (err error) {
		step := servicev2.ServiceStep_tangled{}
		args, stmt, err := vanilla.Select(step.TableName(), step.ColumnNames(), vanilla.Condition{Condition: condition, Args: args}).Prepare(db)
		Do(&err, func() (err error) {
			defer stmt.Close()
			rows, err := stmt.Query(args...)
			err = errors.Wrapf(err, "stmt query; step")
			Do(&err, func() (err error) {
				defer rows.Close()
				for rows.Next() {
					// step := servicev2.ServiceStep_tangled{}
					err = step.Scan(rows)
					err = errors.Wrapf(err, "scan row; step")
					if err != nil {
						break
					}

					steps = append(steps, step)
				}
				return
			})
			return
		})
		return
	})

	Do(&err, func() (err error) {
		b, err := json.Marshal(service)
		if err != nil {
			println(err.Error())
		}
		println(string(b))

		newservice := servicev2.Service_tangled{}

		err = json.Unmarshal(b, &newservice)
		if err != nil {
			println(err.Error())
		}
		return
	})
	Do(&err, func() (err error) {
		b, err := json.Marshal(steps)
		if err != nil {
			println(err.Error())
		}
		println(string(b))

		err = json.Unmarshal(b, &steps)
		if err != nil {
			println(err.Error())
		}
		return
	})

	// newservice := servicev2.Service_tangled{}

	// err = json.Unmarshal(b, &newservice)
	// if err != nil {
	// 	println(err.Error())
	// }

	t.Error(err)
}

// 여러 테이블의 결과를 한번에 쿼리 하는 것은 기본 지원되는 사양이 아닌 것으로 보임
func TestGetServiceMultiResultSet(t *testing.T) {

	uuid := "s:1"

	condition := "WHERE uuid = ?"
	args := []interface{}{
		uuid,
	}
	service := servicev2.Service_tangled{}
	steps := make([]servicev2.ServiceStep_tangled, 0, __INIT_RECORD_CAPACITY__)
	step := servicev2.ServiceStep_tangled{}
	Do(&err, func() (err error) {
		// args, stmt, err := vanilla.Get(service.TableName(), service.ColumnNames(), vanilla.Condition{Condition: condition, Args: args}).
		// 	Combine(vanilla.Get(step.TableName(), step.ColumnNames(), vanilla.Condition{Condition: condition, Args: args})).
		// 	Prepare(db)

		q := vanilla.Select(service.TableName(), service.ColumnNames(), vanilla.Condition{Condition: condition, Args: args}).
			Combine(vanilla.Select(step.TableName(), step.ColumnNames(), vanilla.Condition{Condition: condition, Args: args}))

		// db.Query(q.Query(), q.Args()...)

		Do(&err, func() (err error) {
			// defer stmt.Close()
			// rows, err := stmt.Query(args...)

			rows, err := db.Query(q.Query(), q.Args()...)
			err = errors.Wrapf(err, "stmt query")
			Do(&err, func() (err error) {
				defer rows.Close()

				rows.NextResultSet()
				// if rows.Next() {
				//scan service
				for rows.Next() {
					err = service.Scan(rows)
					err = errors.Wrapf(err, "scan service")

					if err != nil {
						return
					}
				}

				rows.NextResultSet()

				for rows.Next() {
					step := servicev2.ServiceStep_tangled{}
					//scan step
					err = step.Scan(rows)
					err = errors.Wrapf(err, "scan step")
					if err != nil {
						return
					}

					steps = append(steps, step)
				}
				// }
				return
			})
			return
		})
		return
	})

	Do(&err, func() (err error) {
		b, err := json.Marshal(service)
		if err != nil {
			println(err.Error())
		}
		println(string(b))

		newservice := servicev2.Service_tangled{}

		err = json.Unmarshal(b, &newservice)
		if err != nil {
			println(err.Error())
		}
		return
	})
	Do(&err, func() (err error) {
		b, err := json.Marshal(steps)
		if err != nil {
			println(err.Error())
		}
		println(string(b))

		err = json.Unmarshal(b, &steps)
		if err != nil {
			println(err.Error())
		}
		return
	})

	t.Error(err)
}
