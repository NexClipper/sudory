package database

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewMigrate(src, dest string) (m *migrate.Migrate, err error) {
	m, err = migrate.New(src, dest)
	err = errors.Wrapf(err, "failed to new migrate")
	return
}
