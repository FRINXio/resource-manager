// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package psql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"contrib.go.opencensus.io/integrations/ocsql"
	slog "github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/telemetry"
	"github.com/jackc/pgx"
	psql "github.com/jackc/pgx/stdlib"
)

// Config is a DSN string wrapper.
type Config struct {
	dsn string
}

// String returns the textual representation of a config.
func (c *Config) String() string {
	return c.dsn
}

// MarshalText marshals the Config to text.
func (c *Config) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// ReplaceDbName replaces name of DB in dsn url string
// This cannot be done voa ParseConnectionString + SerializeConnectionString from psql lib
// since the lib only supports parsing, but not serializing
func ReplaceDbName(dsnUrl string, newDbName string) (string, error) {
	_, err := pgx.ParseConnectionString(dsnUrl)
	if err != nil {
		return "", fmt.Errorf("cannot parse dsn: %w", err)
	}

	dsnPrefix := "postgresql://"
	if strings.Index(dsnUrl, dsnPrefix) != 0 {
		return "", fmt.Errorf("cannot parse dsn: Must start with 'postgres://'")
	}
	dsnUrl = dsnUrl[len(dsnPrefix):]

	indexQ := strings.Index(dsnUrl, "?")
	indexS := -1

	if indexQ > 0 {
		indexS = strings.LastIndex(dsnUrl[0:indexQ], "/")
	} else {
		indexS = strings.LastIndex(dsnUrl, "/")
	}

	newDsn := ""

	if indexQ > 0 && indexS > 0 {
		oldDb := dsnUrl[indexS : indexQ+1]
		newDsn = strings.Replace(dsnUrl, oldDb, "/"+newDbName+"?", 1)
	} else if indexQ < 0 && indexS > 0 {
		oldDb := dsnUrl[indexS : ]
		if oldDb == "/" {
			newDsn = dsnUrl + newDbName
		} else {
			newDsn = strings.Replace(dsnUrl, oldDb, "/" + newDbName, 1)
		}
	} else {
		newDsn = dsnUrl + "/" + newDbName
	}

	return dsnPrefix + newDsn, nil

}

// UnmarshalText unmarshalls text to a Config.
func (c *Config) UnmarshalText(text []byte) error {
	// Parse the DSN to make sure it's valid
	_, err := pgx.ParseConnectionString(string(text))
	if err != nil {
		return fmt.Errorf("cannot parse dsn: %w", err)
	}
	c.dsn = string(text)
	return nil
}

// Set sets the Config for the flag.Value interface.
func (c *Config) Set(value string) error {
	return c.UnmarshalText([]byte(value))
}

// Open new connection and start stats recorder.
func Open(dsn string) *sql.DB {
	return sql.OpenDB(connector{dsn})
}

// RecordStats records database statistics for provided sql.DB.
func RecordStats(db *sql.DB) func() {
	return ocsql.RecordStats(db, 10*time.Second)
}

// SetLogger is used to set the logger for critical errors.
func SetLogger(cfg *Config, logger slog.Logger) {
	//const lvl = zap.ErrorLevel
	//l, _ := zap.NewStdLogAt(
	//	logger.Background().
	//		WithOptions(zap.AddStacktrace(lvl)).
	//		With(zap.String("pkg", "mysql")),
	//	lvl,
	//)

	// FIXME set the logger
	logger.Background().Error("Unable to set logger for psql, FIXME")
}

//type PqToStdLogger struct {
//	logger *log.Logger
//}
//
//func (p PqToStdLogger) Log(level pgx.LogLevel, msg string, data map[string]interface{}) {
//	p.logger.Printf(msg, data);
//}

type connector struct {
	dsn string
}

func (c connector) Connect(context.Context) (driver.Conn, error) {
	return c.Driver().Open(c.dsn)
}

func (connector) Driver() driver.Driver {
	return ocsql.Wrap(psql.GetDefaultDriver(),
		ocsql.WithAllTraceOptions(),
		ocsql.WithRowsClose(false),
		ocsql.WithRowsNext(false),
		ocsql.WithDisableErrSkip(true),
		ocsql.WithSampler(
			telemetry.WithoutNameSampler("sql:prepare"),
		),
	)
}

// DefaultViews are predefined views for opencensus metrics.
var DefaultViews = ocsql.DefaultViews

// Provider is a wire provider that produces *sql.DB from config.
func Provider(config *Config) (*sql.DB, func()) {
	db := Open(config.String())
	return db, RecordStats(db)
}
