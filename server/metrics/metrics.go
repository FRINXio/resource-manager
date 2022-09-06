// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metrics

import (
	"context"
	"errors"
	telemetry2 "github.com/net-auto/resourceManager/telemetry"
	"net"
	"net/http"
	"time"

	"github.com/google/wire"
	"go.opencensus.io/stats/view"
	"go.uber.org/zap"
)

const (
	readTimeout     = 8 * time.Second
	writeTimeout    = 8 * time.Second
	shutdownTimeout = 5 * time.Second
)

// Metrics exposes metrics over http.
type Metrics struct {
	log      *zap.Logger
	exporter telemetry2.ViewExporter
	views    []*view.View
}

// Config configures metrics.
type Config struct {
	Log      *zap.Logger
	Exporter telemetry2.ViewExporter
	Views    []*view.View
}

// Addr represents an endpoint address.
type Addr string

func New(cfg Config) *Metrics {
	if cfg.Log == nil {
		cfg.Log = zap.NewNop()
	}
	if cfg.Exporter == nil {
		cfg.Exporter = telemetry2.NopExporter{}
	}
	return &Metrics{
		log:      cfg.Log.Named("metrics"),
		exporter: cfg.Exporter,
		views:    append([]*view.View(nil), cfg.Views...),
	}
}

// Start will register views and start the metrics server.
func (m *Metrics) Start(addr Addr) (_ *http.Server, err error) {
	ln, err := net.Listen("tcp", string(addr))
	if err != nil {
		m.log.Error("cannot create listener", zap.Error(err))
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = ln.Close()
		}
	}()
	if err := view.Register(m.views...); err != nil {
		m.log.Error("cannot register views", zap.Error(err))
		return nil, err
	}
	view.RegisterExporter(m.exporter)
	defer func() {
		if err != nil {
			view.Unregister(m.views...)
			view.UnregisterExporter(m.exporter)
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/metrics", m.exporter)

	srv := &http.Server{
		Addr:         ln.Addr().String(),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Handler:      mux,
	}

	go func() {
		log := m.log.With(zap.Stringer("address", ln.Addr()))
		log.Info("listening for connections")
		if err := srv.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
			log.Error("cannot run server", zap.Error(err))
		}
	}()

	return srv, nil
}

// Shutdown terminates the metrics server.
func (m *Metrics) Shutdown(srv *http.Server) error {
	m.log.Info("stopping server")

	view.Unregister(m.views...)
	view.UnregisterExporter(m.exporter)

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		m.log.Error("cannot shutdown server", zap.Error(err))
		return err
	}

	m.log.Info("server gracefully stopped")
	return nil
}

// Serve starts the metrics server and blocks until context is done.
func (m *Metrics) Serve(ctx context.Context, addr Addr) error {
	srv, err := m.Start(addr)
	if err != nil {
		return err
	}
	defer m.Shutdown(srv)
	<-ctx.Done()
	return nil
}

// Provider is a wire provider of this package.
var Provider = wire.NewSet(
	New,
	wire.Struct(new(Config), "*"),
)
