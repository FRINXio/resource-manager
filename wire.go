// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build wireinject

package main

import (
	"context"
	"contrib.go.opencensus.io/integrations/ocsql"
	symphLogger "github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/server/metrics"
	"github.com/facebookincubator/symphony/pkg/server/xserver"
	"github.com/facebookincubator/symphony/pkg/telemetry"
	"github.com/google/wire"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/graph/graphhttp"
	logger "github.com/net-auto/resourceManager/logging"
	"github.com/net-auto/resourceManager/viewer"
	"go.opencensus.io/stats/view"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gocloud.dev/server/health"
)

func newApplication(ctx context.Context, flags *cliFlags) (*application, func(), error) {
	wire.Build(
		wire.FieldsOf(new(*cliFlags),
			"ListenAddress",
			"MetricsAddress",
			"TelemetryConfig",
		),
		provideLogger,
		symphLogger.ProvideZapLogger,
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

// Adapter exposing RM logger into symphony logger
type rmLoggerSymphonyAdapter struct {
	bg *zap.Logger
}

// Background returns a context-unaware logger.
func (l rmLoggerSymphonyAdapter) Background() *zap.Logger {
	return l.bg
}

// For returns a context-aware logger.
func (l rmLoggerSymphonyAdapter) For(ctx context.Context) *zap.Logger {
	return l.Background().With(symphLogger.FieldsFromContext(ctx)...)
}

// provideLogger initializes RM logger and exposes it to the wiring system as symphony logger
func provideLogger(ctx context.Context, cf *cliFlags) symphLogger.Logger {
	logger.Init(cf.LogPath, cf.LogLevel, cf.LogWithColors)
	var zapLevel zapcore.Level
	switch cf.LogLevel {
	case "trace":
	case "debug":
		zapLevel = zapcore.DebugLevel
		break;
	case "warn":
		zapLevel = zapcore.WarnLevel
		break;
	case "error":
		zapLevel = zapcore.ErrorLevel
		break;
	case "panic":
		zapLevel = zapcore.PanicLevel
		break;
	case "fatal":
		zapLevel = zapcore.FatalLevel
		break;
	case "info":
	default:
		zapLevel = zapcore.InfoLevel
		break;
	}
	zapLevel.Enabled(zap.DebugLevel)
	zapEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:       "msg",
		LevelKey:         "",
		TimeKey:          "",
		NameKey:          "",
		CallerKey:        "",
		FunctionKey:      "",
		StacktraceKey:    "",
		LineEnding:       "",
		EncodeLevel:      nil,
		EncodeTime:       nil,
		EncodeDuration:   nil,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: " ",
	})
	core := zapcore.NewCore(zapEncoder, zapcore.AddSync(logger.GetLogger().Writer()), zapLevel)
	return rmLoggerSymphonyAdapter{zap.New(core)}
}

func provideViews() []*view.View {
	views := xserver.DefaultViews()
	views = append(views, ocsql.DefaultViews...)
	return views
}
