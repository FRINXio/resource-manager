// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build wireinject

package main

import (
	"context"
	"contrib.go.opencensus.io/integrations/ocsql"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/server/metrics"
	"github.com/facebookincubator/symphony/pkg/server/xserver"
	"github.com/facebookincubator/symphony/pkg/telemetry"
	"github.com/google/wire"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/graph/graphhttp"
	"github.com/net-auto/resourceManager/viewer"
	"go.opencensus.io/stats/view"
	"gocloud.dev/server/health"
)

func newApplication(ctx context.Context, flags *cliFlags) (*application, func(), error) {
	wire.Build(
		wire.FieldsOf(new(*cliFlags),
			"ListenAddress",
			"MetricsAddress",
			"LogConfig",
			"TelemetryConfig",
		),
		log.Provider,
		wire.Struct(
			new(application),
			"*",
		),
		newTenancy,
		newHealthChecks,
		newPsqlTenancy,
		metrics.Provider,
		telemetry.ProvideViewExporter,
		provideViews,
		graphhttp.NewServer,
		wire.Struct(new(graphhttp.Config), "*"),
	)
	return nil, nil, nil
}

func newTenancy(ctx context.Context, tenancy *viewer.PsqlTenancy) (viewer.Tenancy, error) {
	initFunc := func(client *ent.Client) {

	}
	return viewer.NewCacheTenancy(tenancy, initFunc), nil
}

func newHealthChecks(tenancy *viewer.PsqlTenancy) []health.Checker {
	return []health.Checker{tenancy}
}

func newPsqlTenancy(ctx context.Context, flags *cliFlags) (*viewer.PsqlTenancy, error) {
	return viewer.NewPsqlTenancy(ctx, flags.DatabaseURL, flags.TenancyConfig.TenantMaxConn)
}

func provideViews() []*view.View {
	views := xserver.DefaultViews()
	views = append(views, ocsql.DefaultViews...)
	return views
}
