// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copied from inventory and modified to use our ent.Client

package viewer

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"runtime"
	"strings"
	"sync"

	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/net-auto/resourceManager/ent/migrate"
	pools "github.com/net-auto/resourceManager/pools/allocating_strategies"
	"github.com/net-auto/resourceManager/psql"
	"go.uber.org/zap"

	"github.com/facebook/ent/dialect"
	entsql "github.com/facebook/ent/dialect/sql"
	"github.com/net-auto/resourceManager/ent"

	"gocloud.dev/server/health"
	"gocloud.dev/server/health/sqlhealth"
)

// Tenancy provides tenant client for key.
type Tenancy interface {
	ClientFor(context.Context, string, *zap.Logger) (*ent.Client, *dialect.Driver, error)
}

// FixedTenancy returns a fixed client.
type FixedTenancy struct {
	client *ent.Client
}

// NewFixedTenancy creates fixed tenancy from client.
func NewFixedTenancy(client *ent.Client) FixedTenancy {
	return FixedTenancy{client}
}

// ClientFor implements Tenancy interface.
func (f FixedTenancy) ClientFor(context.Context, string, *log.Logger) (*ent.Client, error) {
	return f.Client(), nil
}

// Client returns the client stored in fixed tenancy.
func (f FixedTenancy) Client() *ent.Client {
	return f.client
}

// CacheTenancy is a tenancy wrapper cashing underlying clients.
type CacheTenancy struct {
	tenancy  Tenancy
	initFunc func(*ent.Client)
	clients  sync.Map
	mu       sync.Mutex
}

// NewCacheTenancy creates a tenancy cache.
func NewCacheTenancy(tenancy Tenancy, initFunc func(*ent.Client)) *CacheTenancy {
	return &CacheTenancy{
		tenancy:  tenancy,
		initFunc: initFunc,
	}
}

type clientAndDriver struct {
	client *ent.Client
	driver *dialect.Driver
}

// ClientFor implements Tenancy interface.
func (c *CacheTenancy) ClientFor(ctx context.Context, name string, logger *zap.Logger) (*ent.Client, *dialect.Driver, error) {
	if cAndD, ok := c.clients.Load(name); ok {
		cAndDTyped := cAndD.(*clientAndDriver)
		return cAndDTyped.client, cAndDTyped.driver, nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if cAndD, ok := c.clients.Load(name); ok {
		cAndDTyped := cAndD.(*clientAndDriver)
		return cAndDTyped.client, cAndDTyped.driver, nil
	}
	client, driver, err := c.tenancy.ClientFor(ctx, name, logger)
	if err != nil {
		return client, driver, err
	}
	if c.initFunc != nil {
		c.initFunc(client)
	}
	c.clients.Store(name, clientAndDriver{client, driver})
	return client, driver, nil
}

// CheckHealth implements health.Checker interface.
func (c *CacheTenancy) CheckHealth() error {
	if checker, ok := c.tenancy.(health.Checker); ok {
		return checker.CheckHealth()
	}
	return nil
}

// PsqlTenancy provides logical database per tenant.
type PsqlTenancy struct {
	health.Checker
	url      url.URL
	maxConns int
	mu       sync.Mutex
	closers  []func()
}

// NewPsqlTenancy creates psql tenancy for data source.
func NewPsqlTenancy(ctx context.Context, u *url.URL, maxConns int) (*PsqlTenancy, error) {
	db, err := psql.OpenURL(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("opening postgres database: %w", err)
	}
	checker := sqlhealth.New(db)
	tenancy := &PsqlTenancy{
		Checker:  checker,
		url:      *u,
		maxConns: maxConns,
		closers: []func(){
			checker.Stop,
		},
	}
	runtime.SetFinalizer(tenancy, func(tenancy *PsqlTenancy) {
		for _, closer := range tenancy.closers {
			closer()
		}
	})
	return tenancy, nil
}

// ClientFor implements Tenancy interface.
func (m *PsqlTenancy) ClientFor(ctx context.Context, name string, logger *zap.Logger) (*ent.Client, error) {
	if err := m.createTenanDb(ctx, logger, name); err != nil {
		return nil, err
	}

	u := m.url
	dbName := DBName(name)
	u.Path = "/" + dbName
	db, closer, err := psql.Provide(ctx, &u)
	if err != nil {
		return nil, fmt.Errorf("opening psql database: %w", err)
	}
	db.SetMaxOpenConns(m.maxConns)
	m.mu.Lock()
	m.closers = append(m.closers, closer)
	m.mu.Unlock()
	driver := entsql.OpenDB(dialect.Postgres, db)
	driverOption := ent.Driver(driver)

	client := ent.NewClient(driverOption)

	logger.Debug("Invoking db migration for tenant", zap.String("tenant", name))
	if err := m.migrate(ctx, client, logger); err != nil {
		return nil, err
	}

	logger.Debug("Loading built-in resource types for tenant", zap.String("tenant", name))
	if err := pools.LoadBuiltinTypes(ctx, client); err != nil {
		return nil, err
	}

	return client, nil
}

func (m *PsqlTenancy) createTenanDb(ctx context.Context, logger *zap.Logger, name string) error {
	dbRoot, closer, err := psql.Provide(ctx, &m.url)
	defer closer()
	if err != nil {
		return fmt.Errorf("opening root psql database: %w", err)
	}

	if exists, err := ExistTenantDb(ctx, name, dbRoot); !exists && err == nil {
		logger.Info("Creating db for new tenant", zap.String("tenant", name))
		if _, err := CreateTenantDb(ctx, name, dbRoot); err != nil {
			return err
		}
	} else if exists && err == nil {
		// Do nothing, db already in place
	} else {
		logger.Error("Creating db for new tenant failed", zap.String("tenant", name))
		return err
	}

	return nil
}

const dbPrefix = "rm_tenant_"

// DBName returns the prefixed database name in order to avoid collision with Postgres internal databases.
func DBName(name string) string {
	return dbPrefix + name
}

// FromDBName returns the source name of the tenant.
func FromDBName(name string) string {
	return strings.TrimPrefix(name, dbPrefix)
}

type queryer interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
}

type tenancyCtxKey struct{}

// TenancyFromContext returns the Tenancy stored in a context, or nil if there isn't one.
func TenancyFromContext(ctx context.Context) Tenancy {
	t, _ := ctx.Value(tenancyCtxKey{}).(Tenancy)
	return t
}

// NewTenancyContext returns a new context with the given Tenancy attached.
func NewTenancyContext(parent context.Context, tenancy Tenancy) context.Context {
	return context.WithValue(parent, tenancyCtxKey{}, tenancy)
}

func (m *PsqlTenancy) migrate(ctx context.Context, client *ent.Client, logger *zap.Logger) error {
	if err := client.Schema.Create(ctx,
		migrate.WithFixture(false),
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		logger.Error("tenancy migrate", zap.Error(err))
		return fmt.Errorf("running tenancy migration: %w", err)
	}
	return nil
}
