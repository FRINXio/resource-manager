// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build wireinject

package main

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/server"
	fb_viewer "github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/psql"
	"github.com/net-auto/resourceManager/graph/graphhttp"
	"github.com/net-auto/resourceManager/viewer"

	"github.com/google/wire"
	"gocloud.dev/server/health"
)

func newApplication(ctx context.Context, flags *cliFlags) (*application, func(), error) {
	wire.Build(
		wire.FieldsOf(new(*cliFlags),
			"PsqlConfig",
			"LogConfig",
			"TelemetryConfig",
			"TenancyConfig",
		),
		log.Provider,
		newApp,
		newTenancy,
		newHealthChecks,
		newPsqlTenancy,
		graphhttp.NewServer,
		wire.Struct(new(graphhttp.Config), "*"),
	)
	return nil, nil, nil
}

func newApp(logger log.Logger, httpServer *server.Server, flags *cliFlags) *application {
	var app application
	app.Logger = logger.Background()
	app.http.Server = httpServer
	app.http.addr = flags.HTTPAddr
	return &app
}

func newTenancy(tenancy *viewer.PsqlTenancy) (viewer.Tenancy, error) {
	initFunc := func(*ent.Client) {
		// NOOP
	}
	return viewer.NewCacheTenancy(tenancy, initFunc), nil
}

func newHealthChecks(tenancy *viewer.PsqlTenancy) []health.Checker {
	return []health.Checker{tenancy}
}

func newPsqlTenancy(config psql.Config, tenancyConfig fb_viewer.Config, logger log.Logger) (*viewer.PsqlTenancy, error) {
	tenancy, err := viewer.NewPsqlTenancy(config.String(), tenancyConfig.TenantMaxConn)
	if err != nil {
		return nil, fmt.Errorf("creating psql tenancy: %w", err)
	}
	tenancy.SetLogger(logger)
	psql.SetLogger(&config, logger)
	return tenancy, nil
}
