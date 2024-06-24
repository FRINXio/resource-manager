// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"contrib.go.opencensus.io/integrations/ocsql"
	"github.com/google/wire"
	"github.com/net-auto/resourceManager/graph/graphhttp"
	logger "github.com/net-auto/resourceManager/logging"
	"github.com/net-auto/resourceManager/logging/log"
	"github.com/net-auto/resourceManager/server/metrics"
	"github.com/net-auto/resourceManager/server/xserver"
	"github.com/net-auto/resourceManager/telemetry"
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
		log.ProvideZapLogger,
		wire.Struct(
			new(application),
			"*",
		),
		newFixedTenancy,
		newHealthChecks,
		metrics.Provider,
		telemetry.ProvideViewExporter,
		provideViews,
		graphhttp.NewServer,
		wire.Struct(new(graphhttp.Config), "*"),
	)
	return nil, nil, nil
}

func newHealthChecks(tenancy viewer.Tenancy) []health.Checker {
	return []health.Checker{tenancy}
}

func newFixedTenancy(ctx context.Context, flags *cliFlags, logger *zap.Logger) (viewer.Tenancy, error) {
	client, checker, err := viewer.PsqlClient(ctx, flags.DatabaseURL, flags.MaxDbConnections, logger)
	if err != nil {
		return nil, err
	}
	return viewer.NewFixedTenancy(client, checker), nil
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
	return l.Background().With(log.FieldsFromContext(ctx)...)
}

// provideLogger initializes RM logger and exposes it to the wiring system as symphony logger
func provideLogger(ctx context.Context, cf *cliFlags) log.Logger {
	logger.Init(cf.LogPath, cf.LogLevel, cf.LogWithColors)
	var zapLevel zapcore.Level
	switch cf.LogLevel {
	case "trace":
	case "debug":
		zapLevel = zapcore.DebugLevel
		break
	case "warn":
		zapLevel = zapcore.WarnLevel
		break
	case "error":
		zapLevel = zapcore.ErrorLevel
		break
	case "panic":
		zapLevel = zapcore.PanicLevel
		break
	case "fatal":
		zapLevel = zapcore.FatalLevel
		break
	case "info":
	default:
		zapLevel = zapcore.InfoLevel
		break
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
