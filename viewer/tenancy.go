// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copied from inventory and modified to use our ent.Client

package viewer

import (
	"context"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/migrate"
	pools "github.com/net-auto/resourceManager/pools/allocating_strategies"
	"github.com/net-auto/resourceManager/psql"
	"go.uber.org/zap"
	"gocloud.dev/server/health"
	"net/url"

	"gocloud.dev/server/health/sqlhealth"
)

// Tenancy provides tenant client for key.
type Tenancy interface {
	health.Checker
	ClientFor(context.Context, string, *zap.Logger) (*ent.Client, error)
}

// FixedTenancy returns a fixed client.
type FixedTenancy struct {
	health.Checker
	client *ent.Client
}

// NewFixedTenancy creates fixed tenancy from client.
func NewFixedTenancy(client *ent.Client, checker health.Checker) *FixedTenancy {
	return &FixedTenancy{checker, client}
}

// ClientFor implements Tenancy interface.
func (f FixedTenancy) ClientFor(context.Context, string, *zap.Logger) (*ent.Client, error) {
	return f.Client(), nil
}

// Client returns the client stored in fixed tenancy.
func (f FixedTenancy) Client() *ent.Client {
	return f.client
}

func PsqlClient(ctx context.Context, u *url.URL, maxConns int, logger *zap.Logger) (*ent.Client, health.Checker, error) {
	db, _, err := psql.Provide(ctx, u)
	if err != nil {
		return nil, nil, fmt.Errorf("opening psql database: %w", err)
	}
	db.SetMaxOpenConns(maxConns)
	checker := sqlhealth.New(db)

	drv := ent.Driver(entsql.OpenDB(dialect.Postgres, db))
	client := ent.NewClient(drv)

	logger.Debug("Invoking db migration")
	if err := dbMigrate(ctx, client, logger); err != nil {
		return nil, nil, err
	}

	logger.Debug("Loading built-in resource types")
	if err := pools.LoadBuiltinTypes(ctx, client); err != nil {
		return nil, nil, err
	}

	return client, checker, nil
}

func dbMigrate(ctx context.Context, client *ent.Client, logger *zap.Logger) error {
	if err := client.Schema.Create(ctx,
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		logger.Error("db migrate", zap.Error(err))
		return fmt.Errorf("running db migration: %w", err)
	}
	return nil
}
