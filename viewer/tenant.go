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

	"github.com/facebookincubator/symphony/pkg/viewer"
)

type (
	// TenantService is a tenant service.
	TenantService struct{ db *sql.DB }

	Tenant struct {
		Id   string
		Name string
	}
)

// NewTenantService create a new tenant service.
func NewTenantService(db *sql.DB) *TenantService {
	return &TenantService{db}
}

// Create a tenant by name.
func (s TenantService) Create(ctx context.Context, name string) (tenant *Tenant, err error) {
	var tx *sql.Tx
	if tx, err = s.db.BeginTx(ctx, nil); err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
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

	return create(tx, ctx, name, s)
}

func create(tx *sql.Tx, ctx context.Context, name string, s TenantService) (*Tenant, error) {
	if name == "" {
		return nil, errors.Errorf("missing tenant name")
	}

	switch exist, err := s.exist(tx, ctx, name); {
	case err != nil:
		return nil, err
	case exist:
		return nil, errors.Errorf("tenant %q exists", name)
	}

	if _, err := tx.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_bin", viewer.DBName(name))); err != nil {
		return nil, err
	}
	return &Tenant{Id: name, Name: name}, nil
}

// Delete tenant by name.
func (s TenantService) Delete(ctx context.Context, name string) (err error) {
	var tx *sql.Tx
	if tx, err = s.db.BeginTx(ctx, nil); err != nil {
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

	return s.delete(tx, ctx, name)
}

func (s TenantService) delete(tx *sql.Tx, ctx context.Context, name string) error {
	if name == "" {
		return errors.Errorf("missing tenant name")
	}
	switch exist, err := s.exist(tx, ctx, name); {
	case err != nil:
		return err
	case !exist:
		return errors.Errorf("missing tenant %s", name)
	}
	if _, err := tx.ExecContext(ctx,
		fmt.Sprintf("DROP DATABASE `%s`", viewer.DBName(name)),
	); err != nil {
		return err
	}
	return nil
}

func (s TenantService) Exist(ctx context.Context, name string) (exists bool, err error) {
	var tx *sql.Tx
	if tx, err = s.db.BeginTx(ctx, nil); err != nil {
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

	return s.exist(tx, ctx, name)
}

func (s TenantService) exist(tx *sql.Tx, ctx context.Context, name string) (bool, error) {
	rows, err := tx.QueryContext(ctx,
		"SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", viewer.DBName(name),
	)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if !rows.Next() {
		return false, sql.ErrNoRows
	}
	var n int
	if err := rows.Scan(&n); err != nil {
		return false, fmt.Errorf("scanning count: %w", err)
	}
	return n > 0, rows.Err()
}
