// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/alecthomas/kong"
	"github.com/net-auto/resourceManager/ent/schema"
	logger "github.com/net-auto/resourceManager/logging"
	"github.com/net-auto/resourceManager/server"
	"github.com/net-auto/resourceManager/server/metrics"
	"github.com/net-auto/resourceManager/telemetry"
	stdlog "log"
	"net/url"
	"os"
	"syscall"

	_ "github.com/net-auto/resourceManager/ent/runtime"
	"github.com/net-auto/resourceManager/pkg/ctxgroup"
	"github.com/net-auto/resourceManager/pkg/ctxutil"
	"go.uber.org/zap"
)

type cliFlags struct {
	ConfigFile       kong.ConfigFlag   `type:"existingfile" placeholder:"PATH" help:"Configuration file path."`
	ListenAddress    string            `name:"web.listen-address" default:":http" help:"Web address to listen on."`
	MetricsAddress   metrics.Addr      `name:"metrics.listen-address" default:":9464" help:"Metrics address to listen on."`
	DatabaseURL      *url.URL          `name:"db.url" env:"DB_URL" required:"" placeholder:"URL" help:"Database URL."`
	MaxDbConnections int               `name:"tenancy.db_max_conn" env:"TENANCY_DB_MAX_CONNECTIONS" default:"20" help:"Database max connections."`
	TelemetryConfig  telemetry.Config  `embed:""`
	RbacConfig       schema.RbacConfig `embed:""`
	LogPath          string            `name:"logPath" env:"RM_LOG_PATH" default:"./rm.log" help:"Path to logfile." type:"path"`
	LogLevel         string            `name:"loglevel" env:"RM_LOG_LEVEL" default:"info" help:"Logging level - fatal, error, warning, info, debug or trace." type:"string"`
	LogWithColors    bool              `name:"logWithColors" default:"false" help:"Force colors in log." type:"bool"`
}

func main() {
	var cf cliFlags
	kong.Parse(&cf, cf.TelemetryConfig)

	ctx := ctxutil.WithSignal(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	app, cleanup, err := newApplication(ctx, &cf)
	if err != nil {
		stdlog.Fatal(err)
	}
	defer cleanup()

	defer logger.Close()

	logger.Info(ctx, "initializing RBAC with %+v", zap.Reflect("rbacConfig", cf.RbacConfig))
	initializeRbacSettings(cf)

	logger.Info(ctx, "starting application %+v", zap.String("httpEndpoint", cf.ListenAddress))

	err = app.run(ctx)
	logger.Info(ctx, "terminating application %+v", zap.Error(err))
}

// initializeRbacSettings configures which roles and groups grant users admin access ... globally
func initializeRbacSettings(cf cliFlags) {
	schema.InitializeAdminRolesFromSlice(cf.RbacConfig.AdminRoles)
	schema.InitializeAdminGroupsFromSlice(cf.RbacConfig.AdminGroups)
}

type application struct {
	*zap.Logger
	server      *server.Server
	addr        string
	metrics     *metrics.Metrics
	metricsAddr metrics.Addr
}

func (app *application) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	g := ctxgroup.WithContext(ctx)
	g.Go(func(context.Context) error {
		err := app.server.ListenAndServe(app.addr)
		app.Debug("http server terminated", zap.Error(err))
		return err
	})
	g.Go(func(ctx context.Context) error {
		return app.metrics.Serve(ctx, app.metricsAddr)
	})
	g.Go(func(ctx context.Context) error {
		defer cancel()
		<-ctx.Done()
		return nil
	})
	<-ctx.Done()

	g.Go(func(context.Context) error {
		app.Debug("start http server termination")
		err := app.server.Shutdown(context.Background())
		app.Debug("end http server termination", zap.Error(err))
		return err
	})
	return g.Wait()
}
