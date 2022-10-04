package vanilla_test

// import (
// 	"database/sql"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/NexClipper/sudory/pkg/server/database/vanilla"

// 	// . "github.com/NexClipper/sudory/pkg/server/macro"
// 	clusterv2 "github.com/NexClipper/sudory/pkg/server/model/cluster/v2"
// 	servicev2 "github.com/NexClipper/sudory/pkg/server/model/service/v2"
// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/pkg/errors"
// )

// var (
// 	db_conn = func() *sql.DB {
// 		db, err := sql.Open("mysql", "root:verysecret@tcp(localhost:3306)/sudory?parseTime=true&multiStatements=true")
// 		if err != nil {
// 			panic(err)
// 		}
// 		return db
// 	}
// )

// func TestQueryRowCluster(t *testing.T) {

// 	uuid := "7dcf4f5eb62641c083614dccab3a044b"

// 	// cond := vanilla.NewCond(
// 	// 	"WHERE uuid = ?",
// 	// 	uuid,
// 	// )

// 	cluster := clusterv2.Cluster{}
// 	// sacn := vanilla.QueryRow(db, cluster.TableName(), cluster.ColumnNames(), *cond)
// 	// err = sacn(func(s vanilla.Scanner) (err error) {
// 	// 	err = cluster.Scan(s)
// 	// 	err = errors.Wrapf(err, "cluster Scan")
// 	// 	return
// 	// })

// 	eq_uuid := vanilla.And(vanilla.Equal("uuid", uuid)).Parse()

// 	stmt := vanilla.Stmt.Select(cluster.TableName(), cluster.ColumnNames(), eq_uuid, nil, nil)
// 	err := stmt.QueryRow(db_conn())(func(s vanilla.Scanner) (err error) {
// 		err = cluster.Scan(s)
// 		err = errors.Wrapf(err, "cluster Scan")
// 		return
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	println(fmt.Sprintf("%+v", cluster))

// 	err = errors.Wrapf(err, "query rows")
// 	t.Error(err)
// }

// func TestQueryRow(t *testing.T) {

// 	uuid := "b9809f912c54483cb324b6e5e8b058a8"

// 	eq_uuid := vanilla.And(vanilla.Equal("uuid", uuid)).Parse()

// 	service := servicev2.Service_tangled{}
// 	stmt := vanilla.Stmt.Select(service.TableName(), service.ColumnNames(), eq_uuid, nil, nil)
// 	err := stmt.QueryRow(db_conn())(func(s vanilla.Scanner) (err error) {
// 		err = service.Scan(s)
// 		err = errors.Wrapf(err, "service Scan")
// 		return
// 	})

// 	println(fmt.Sprintf("%+v", service))

// 	err = errors.Wrapf(err, "query row")
// 	t.Error(err)

// }

// func TestQueryRowIsNill(t *testing.T) {

// 	isnull := vanilla.IsNull("deleted").Parse()

// 	var n int
// 	stmt := vanilla.Stmt.Select("session", []string{"COUNT(1)"}, isnull, nil, nil)
// 	err := stmt.QueryRow(db_conn())(func(s vanilla.Scanner) (err error) {
// 		err = s.Scan(&n)
// 		err = errors.Wrapf(err, "session Scan")
// 		return
// 	})

// 	println(fmt.Sprintf("serssion count %+v", n))

// 	err = errors.Wrapf(err, "query row")
// 	t.Error(err)

// }

// func TestQueryMultiRows(t *testing.T) {

// 	uuid := "b9809f912c54483cb324b6e5e8b058a8"

// 	eq_uuid := vanilla.And(vanilla.Equal("uuid", uuid)).Parse()

// 	steps := make([]servicev2.ServiceStep_tangled, 0, 10)
// 	step := servicev2.ServiceStep_tangled{}
// 	stmt := vanilla.Stmt.Select(step.TableName(), step.ColumnNames(), eq_uuid, nil, nil)
// 	err := stmt.QueryRows(db_conn())(func(scan vanilla.Scanner, _ int) (err error) {
// 		err = step.Scan(scan)
// 		err = errors.Wrapf(err, "step Scan")
// 		if err != nil {
// 			return
// 		}
// 		steps = append(steps, step)
// 		return
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	println(fmt.Sprintf("%+v", steps))

// 	err = errors.Wrapf(err, "query rows")
// 	t.Error(err)

// }

// // 여러 테이블의 결과를 한번에 쿼리 하는 것은 기본 지원되는 사양이 아닌 것으로 보임
// func TestGetServiceMultiResultSet(t *testing.T) {

// 	// uuid := "b9809f912c54483cb324b6e5e8b058a8"

// 	// eq_uuid := vanilla.Equal("uuid", uuid).Parse()

// 	// service := servicev2.Service_tangled{}
// 	// stmt1 := vanilla.Stmt.Select(service.TableName(), service.ColumnNames(), eq_uuid, nil, nil)

// 	// steps := make([]servicev2.ServiceStep_tangled, 0, 10)
// 	// step := servicev2.ServiceStep_tangled{}
// 	// stmt2 := vanilla.Stmt.Select(step.TableName(), step.ColumnNames(), eq_uuid, nil, nil)

// 	// query := strings.Join([]string{
// 	// 	stmt1.Query(), stmt2.Query(),
// 	// }, ";")
// 	// args := append(stmt1.Args(), stmt2.Args()...)

// 	rows, err := db_conn().Query(`

// --	insert into managed_channel (uuid, name, event_category, created) values('5', '1', 0, NOW());
// --	insert into managed_channel (uuid, name, event_category, created) values('6', '2', 0, NOW());

// 	SELECT uuid FROM managed_channel;

// 	SELECT uuid FROM managed_channel;

// `)
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }

// 	// rows, err := stmt.Query()

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer rows.Close()

// 	var uuid string
// 	for rows.Next() {
// 		err = rows.Scan(&uuid)
// 		err = errors.Wrapf(err, "sql.Row.Scan")
// 		if err != nil {
// 			break
// 		}
// 		println(uuid)
// 	}

// 	if ok := rows.NextResultSet(); ok {
// 		t.Log(ok)
// 	}

// 	for rows.Next() {
// 		err = rows.Scan(&uuid)
// 		err = errors.Wrapf(err, "sql.Row.Scan")
// 		if err != nil {
// 			break
// 		}
// 		println(uuid)
// 	}

// 	// err := vanilla.QueryRows(db_conn(), query, args)(func(s vanilla.Scanner, i int) (err error) {
// 	// 	err = step.Scan(s)
// 	// 	err = errors.Wrapf(err, "step Scan")
// 	// 	if err != nil {
// 	// 		steps = append(steps, step)
// 	// 	}
// 	// 	return
// 	// })
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// println(fmt.Sprintf("%+v", servicev2.HttpRsp_Service{Service_tangled: service, Steps: steps}))

// 	t.Error(err)
// }

// func TestUpdateServcie(t *testing.T) {

// 	uuid := "b9809f912c54483cb324b6e5e8b058a8"

// 	eq_uuid := vanilla.Equal("uuid", uuid).Parse()

// 	service := servicev2.Service_essential{}

// 	KV := map[string]interface{}{
// 		"summary": "hello, vanilla!",
// 	}

// 	stmt := vanilla.Stmt.Update(service.TableName(), KV, eq_uuid)
// 	affected, err := stmt.Exec(db_conn())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	println("affected", affected)

// 	t.Error(err)
// }

// func TestInsertServcie(t *testing.T) {

// 	service1 := servicev2.Service{}
// 	service1.Uuid = "b9809f912c54483cb324b6e5e8b058a8"
// 	service1.Created = time.Now()
// 	service1.Name = "created by vanilla1"
// 	service1.Summary = *vanilla.NewNullString("summary of created by vanilla1")
// 	service2 := servicev2.Service{}
// 	service2.Uuid = "97f80887254e4636a2255c188a06dc36"
// 	service2.Created = time.Now()
// 	service2.Name = "created by vanilla2"
// 	service2.Summary = *vanilla.NewNullString("summary of created by vanilla2")

// 	stmt, err := vanilla.Stmt.Insert(service1.TableName(), service1.ColumnNames(), append(service1.Values(), service2.Values()...))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	affected, err := stmt.Exec(db_conn())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	println("affected", affected)

// 	t.Error(err)
// }
