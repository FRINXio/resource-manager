// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/alecthomas/kong"
	"github.com/net-auto/resourceManager/ent/schema"
	logger "github.com/net-auto/resourceManager/logging"
	stdlog "log"
	"os"
	"syscall"

	"github.com/facebookincubator/symphony/pkg/ctxgroup"
	"github.com/facebookincubator/symphony/pkg/ctxutil"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/mysql"
	"github.com/facebookincubator/symphony/pkg/server"
	"github.com/facebookincubator/symphony/pkg/telemetry"
	"github.com/facebookincubator/symphony/pkg/viewer"

	_ "github.com/net-auto/resourceManager/ent/runtime"
	"go.uber.org/zap"
)

type cliFlags struct {
	HTTPAddr           string            `name:"web.listen-address" default:":http" help:"Web address to listen on."`
	MySQLConfig        mysql.Config      `name:"mysql.dsn" env:"MYSQL_DSN" required:"" placeholder:"STRING" help:"MySQL data source name."`
	LogConfig          log.Config        `embed:""`
	TelemetryConfig    telemetry.Config  `embed:""`
	TenancyConfig      viewer.Config     `embed:""`
	RbacConfig         schema.RbacConfig `embed:""`
	LogPath 		   string            `name:"logPath" env:"RM_LOG_PATH" default:"./rm.log" help:"Path to logfile." type:"path"`
	LogLevel 		   string            `name:"loglevel" env:"RM_LOG_LEVEL" default:"info" help:"Logging level - fatal, error, warning, info, debug or trace." type:"string"`
	LogWithColors 	   bool              `name:"logWithColors" default:"false" help:"Force colors in log." type:"bool"`
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

	logger.Init(cf.LogPath, cf.LogLevel, cf.LogWithColors)
	defer logger.Close()

	logger.Info(ctx,"initializing RBAC with %+v", zap.Reflect("rbacConfig", cf.RbacConfig))
	initializeRbacSettings(cf)

	logger.Info(ctx,"starting application %+v", zap.String("httpEndpoint", cf.HTTPAddr))

	err = app.run(ctx)
	logger.Info(ctx,"terminating application %+v", zap.Error(err))
}

// initializeRbacSettings configures which roles and groups grant users admin access ... globally
func initializeRbacSettings(cf cliFlags) {
	schema.InitializeAdminRolesFromSlice(cf.RbacConfig.AdminRoles)
	schema.InitializeAdminGroupsFromSlice(cf.RbacConfig.AdminGroups)
}

type application struct {
	*zap.Logger
	http struct {
		*server.Server
		addr string
	}
}

func (app *application) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	g := ctxgroup.WithContext(ctx)
	g.Go(func(context.Context) error {
		err := app.http.ListenAndServe(app.http.addr)
		app.Debug("http server terminated", zap.Error(err))
		return err
	})
	g.Go(func(ctx context.Context) error {
		defer cancel()
		<-ctx.Done()
		return nil
	})
	<-ctx.Done()

	app.Warn("start application termination",
		zap.NamedError("reason", ctx.Err()),
	)
	defer app.Debug("end application termination")

	g.Go(func(context.Context) error {
		app.Debug("start http server termination")
		err := app.http.Shutdown(context.Background())
		app.Debug("end http server termination", zap.Error(err))
		return err
	})
	return g.Wait()
}
