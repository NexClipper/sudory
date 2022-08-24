package env

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func NewSqlDB(database string) *sql.DB {
	var (
		driver = "mysql"
		dsn    = fmt.Sprintf("root:verysecret@tcp(127.0.0.1:3306)/%v?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true", database)
	)

	db, err := sql.Open(driver, dsn)
	if err != nil {
		panic(err)
	}

	return db
}
