package main_test

import (
	_ "embed"
	"fmt"
	"testing"

	env "github.com/NexClipper/sudory/pkg/server/model/service/gen_test_data/env"
)

var (
	//go:embed init_service.sql
	init_service_sql string
)

func TestInitTables(t *testing.T) {

	// SQL: CREATE DATABASE sudory_schema_test_v2;
	db := env.NewSqlDB("sudory_schema_test_v2")

	rst, err := db.Exec(init_service_sql)
	if err != nil {
		t.Fatal(err)
	}

	affected, err := rst.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}

	lastInserted, err := rst.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(fmt.Sprintf("RowsAffected=%v", affected))
	t.Log(fmt.Sprintf("LastInsertId=%v", lastInserted))
}
