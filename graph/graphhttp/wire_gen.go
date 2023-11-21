// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package graphhttp

import (
	"github.com/net-auto/resourceManager/logging/log"
	"github.com/net-auto/resourceManager/server"
	"github.com/net-auto/resourceManager/server/xserver"
	"github.com/net-auto/resourceManager/telemetry"
	"github.com/net-auto/resourceManager/viewer"
	"gocloud.dev/server/health"
)

// Injectors from wire.go:

// NewServer creates a server from config.
func NewServer(cfg Config) (*server.Server, func(), error) {
	graphhttpRouterConfig, err := newRouterConfig(cfg)
	if err != nil {
		return nil, nil, err
	}
	router, err := newRouter(graphhttpRouterConfig)
	if err != nil {
		return nil, nil, err
	}
	logger := cfg.Logger
	zapLogger := xserver.NewRequestLogger(logger)
	v := cfg.HealthChecks
	config := cfg.Telemetry
	exporter, cleanup, err := telemetry.ProvideTraceExporter(config)
	if err != nil {
		return nil, nil, err
	}
	profilingEnabler := _wireProfilingEnablerValue
	sampler := telemetry.ProvideTraceSampler(config)
	handlerFunc := xserver.NewRecoveryHandler(logger)
	defaultDriver := _wireDefaultDriverValue
	options := &server.Options{
		RequestLogger:         zapLogger,
		HealthChecks:          v,
		TraceExporter:         exporter,
		EnableProfiling:       profilingEnabler,
		DefaultSamplingPolicy: sampler,
		RecoveryHandler:       handlerFunc,
		Driver:                defaultDriver,
	}
	serverServer := server.New(router, options)
	return serverServer, func() {
		cleanup()
	}, nil
}

var (
	_wireProfilingEnablerValue = server.ProfilingEnabler(true)
	_wireDefaultDriverValue    = &server.DefaultDriver{}
)

// wire.go:

// Config defines the http server config.
type Config struct {
	Tenancy      viewer.Tenancy
	Logger       log.Logger
	Telemetry    telemetry.Config
	HealthChecks []health.Checker
}

func newRouterConfig(config Config) (cfg routerConfig, err error) {
	cfg = routerConfig{logger: config.Logger}
	cfg.viewer.tenancy = config.Tenancy
	return cfg, nil
}
