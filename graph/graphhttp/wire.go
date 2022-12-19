// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wireinject
// +build wireinject

package graphhttp

import (
	"github.com/net-auto/resourceManager/logging/log"
	"github.com/net-auto/resourceManager/server"
	"github.com/net-auto/resourceManager/server/xserver"
	"github.com/net-auto/resourceManager/telemetry"
	"net/http"

	"github.com/google/wire"
	"github.com/gorilla/mux"
	"github.com/net-auto/resourceManager/viewer"
	"gocloud.dev/server/health"
)

// Config defines the http server config.
type Config struct {
	Tenancy      viewer.Tenancy
	Logger       log.Logger
	Telemetry    telemetry.Config
	HealthChecks []health.Checker
}

// NewServer creates a server from config.
func NewServer(cfg Config) (*server.Server, func(), error) {
	wire.Build(
		xserver.ServiceSet,
		wire.FieldsOf(new(Config), "Logger", "Telemetry", "HealthChecks"),
		newRouterConfig,
		newRouter,
		wire.Bind(new(http.Handler), new(*mux.Router)),
	)
	return nil, nil, nil
}

func newRouterConfig(config Config) (cfg routerConfig, err error) {
	cfg = routerConfig{logger: config.Logger}
	cfg.viewer.tenancy = config.Tenancy
	return cfg, nil
}
