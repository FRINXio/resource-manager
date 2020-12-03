// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package psql

import (
	"context"
	"database/sql"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"net/url"
	"time"

	"contrib.go.opencensus.io/integrations/ocsql"
	cdkpsql "gocloud.dev/postgres"
	"gocloud.dev/postgres/awspostgres"
)

var defaultURLMux = cdkpsql.URLMux{}

func init() {
	traceOpts := []ocsql.TraceOption{
		ocsql.WithDisableErrSkip(true),
		ocsql.WithSampler(trace.NeverSample()),
	}
	defaultURLMux.RegisterPostgres(cdkpsql.Scheme, &cdkpsql.URLOpener{})

	defaultURLMux.RegisterPostgres(awspostgres.Scheme, &awspostgres.URLOpener{
		TraceOpts: traceOpts,
	})
}

// Open opens the database identified by the URL string given.
func Open(ctx context.Context, urlstr string) (*sql.DB, error) {
	return defaultURLMux.OpenPostgres(ctx, urlstr)
}

// OpenURL opens the database identified by the URL given.
func OpenURL(ctx context.Context, u *url.URL) (*sql.DB, error) {
	db, err := defaultURLMux.OpenPostgresURL(ctx, u)
	if db != nil {
		db.SetConnMaxIdleTime(30 * time.Second)
		db.SetConnMaxLifetime(5 * time.Minute)
		db.SetMaxIdleConns(1)
		db.SetMaxOpenConns(5)
	}
	return db, err
}

// Provide is a wire provider that produces *sql.DB from url.
func Provide(ctx context.Context, u *url.URL) (*sql.DB, func(), error) {
	db, err := OpenURL(ctx, u)
	if err != nil {
		return nil, nil, err
	}
	return db, ocsql.RecordStats(db, 10*time.Second), nil
}

// logger forwards mysql logs to zap global logger.
type logger struct{}

// Print implements mysql.Logger interface.
func (logger) Print(args ...interface{}) {
	zap.S().With("pkg", "psql").Error(args...)
}
