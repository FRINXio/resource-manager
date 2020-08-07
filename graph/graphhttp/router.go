// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphhttp

import (
	"fmt"
	"net/http"

	"github.com/facebookincubator/symphony/pkg/actions"
	"github.com/facebookincubator/symphony/pkg/actions/executor"
	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/log"
	fb_viewer "github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/net-auto/resourceManager/graph/graphql"
	"github.com/net-auto/resourceManager/viewer"

	"github.com/gorilla/mux"
)

type routerConfig struct {
	viewer struct {
		tenancy viewer.Tenancy
		authurl string
	}
	logger  log.Logger
	actions struct{ registry *executor.Registry }
}

func newRouter(cfg routerConfig) (*mux.Router, func(), error) {
	router := mux.NewRouter()
	router.Use(
		func(h http.Handler) http.Handler {
			return viewer.TenancyHandler(h, cfg.viewer.tenancy, cfg.logger)
		},
		func(h http.Handler) http.Handler {
			return fb_viewer.UserHandler(h, cfg.logger)
		},
		func(h http.Handler) http.Handler {
			return authz.Handler(h, cfg.logger)
		},
		func(h http.Handler) http.Handler {
			return actions.Handler(h, cfg.logger, cfg.actions.registry)
		},
	)

	handler, cleanup, err := graphql.NewHandler(
		graphql.HandlerConfig{
			Logger: cfg.logger,
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("creating graphql handler: %w", err)
	}
	router.PathPrefix("/").
		Handler(handler).
		Name("root")

	return router, cleanup, nil
}
