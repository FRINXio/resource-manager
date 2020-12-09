// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Adapted from inventory's grpc tenancy manager service

package viewer

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

// Create a tenant by name.
func CreateTenantDb(ctx context.Context, name string, db *sql.DB) (tenant string, err error) {
	return create(db, ctx, name)
}

func create(db *sql.DB, ctx context.Context, name string) (string, error) {
	if name == "" {
		return "", errors.Errorf("missing tenant name")
	}

	switch exist, err := ExistTenantDb(ctx, name, db); {
	case err != nil:
		return "", err
	case exist:
		return "", errors.Errorf("tenant %q exists", name)
	}

	if _, err := db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", DBName(name))); err != nil {
		return "", err
	}
	return name, nil
}

// FIXME: Who calls this?
// Delete tenant by name.
func DeleteTenantDb(ctx context.Context, name string, db *sql.DB) (err error) {
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
		if err != nil {
			_ = tx.Rollback()
			return
		}
		if err = tx.Commit(); err != nil {
			err = fmt.Errorf("committing transaction: %w", err)
		}
	}()

	return delete(tx, ctx, name)
}

func delete(tx *sql.Tx, ctx context.Context, name string) error {
	if name == "" {
		return errors.Errorf("missing tenant name")
	}
	switch exist, err := exist(tx, ctx, name); {
	case err != nil:
		return err
	case !exist:
		return errors.Errorf("missing tenant %s", name)
	}
	if _, err := tx.ExecContext(ctx,
		fmt.Sprintf("DROP DATABASE `%s`", DBName(name)),
	); err != nil {
		return err
	}
	return nil
}

func ExistTenantDb(ctx context.Context, name string, db *sql.DB) (exists bool, err error) {
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return false, fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
		if err != nil {
			_ = tx.Rollback()
			return
		}
		if err = tx.Commit(); err != nil {
			err = fmt.Errorf("committing transaction: %w", err)
		}
	}()

	return exist(tx, ctx, name)
}

func exist(tx *sql.Tx, ctx context.Context, name string) (bool, error) {
	rows, err := tx.QueryContext(
		ctx,
		"SELECT COUNT(*) FROM pg_database WHERE datname = $1",
		DBName(name),
	)

	if err != nil {
		return false, err
	}
	defer rows.Close()
	if !rows.Next() {
		return false, nil
	}
	var n int
	if err := rows.Scan(&n); err != nil {
		return false, fmt.Errorf("scanning count: %w", err)
	}
	return n > 0, rows.Err()
}
