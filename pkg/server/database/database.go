package database

import (
	"bytes"
	"database/sql"
	"time"

	"github.com/NexClipper/sudory/pkg/server/config"
	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
)

func New(cfgdb config.Database) (*sql.DB, error) {
	db, err := sql.Open(cfgdb.Type, FormatDSN(cfgdb))
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}

	db.SetMaxOpenConns(cfgdb.MaxOpenConns)
	db.SetMaxIdleConns(cfgdb.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfgdb.MaxConnLifeTime) * time.Second)

	return db, nil
}

func FormatDSN(cfgdb config.Database) string {
	const (
		defaultDbConnParams = "charset=utf8mb4&parseTime=True&loc=Local"
	)

	var buf bytes.Buffer

	if len(cfgdb.Username) > 0 {
		buf.WriteString(cfgdb.Username)
		if len(cfgdb.Password) > 0 {
			buf.WriteByte(':')
			buf.WriteString(cfgdb.Password)
		}
		buf.WriteByte('@')
	}

	if len(cfgdb.Protocol) > 0 {
		buf.WriteString(cfgdb.Protocol)
		if len(cfgdb.Host) > 0 && len(cfgdb.Port) > 0 {
			buf.WriteByte('(')
			buf.WriteString(cfgdb.Host + ":" + cfgdb.Port)
			buf.WriteByte(')')
		}
	}

	buf.WriteByte('/')
	buf.WriteString(cfgdb.DBName)

	if len(defaultDbConnParams) > 0 {
		buf.WriteByte('?')
		buf.WriteString(defaultDbConnParams)
	}

	return buf.String()
}
