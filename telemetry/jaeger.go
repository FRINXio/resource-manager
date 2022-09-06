// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry

import (
	"fmt"
	"os"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/opentracing/opentracing-go"
	jaggeropentracing "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

func init() {
	MustRegisterTraceExporter("jaeger", NewJaegerExporter)
}

// NewJaegerExporter creates a new opencensus trace exporter.
func NewJaegerExporter(opts TraceExporterOptions) (trace.Exporter, error) {
	o := jaeger.Options{
		Endpoint:          os.Getenv("JAEGER_ENDPOINT"),
		CollectorEndpoint: os.Getenv("JAEGER_COLLECTOR_ENDPOINT"),
		AgentEndpoint:     os.Getenv("JAEGER_AGENT_ENDPOINT"),
		Username:          os.Getenv("JAEGER_USER"),
		Password:          os.Getenv("JAEGER_PASSWORD"),
		BufferMaxCount:    16 << 10,
		Process: jaeger.Process{
			ServiceName: opts.ServiceName,
			Tags:        make([]jaeger.Tag, 0, len(opts.Tags)),
		},
		OnError: func(err error) {
			zap.L().Error("cannot upload jaeger stats data", zap.Error(err))
		},
	}
	for key, value := range opts.Tags {
		o.Process.Tags = append(o.Process.Tags, jaeger.StringTag(key, value))
	}
	exporter, err := jaeger.NewExporter(o)
	if err != nil {
		return nil, fmt.Errorf("cannot create jaeger exporter: %w", err)
	}
	return exporter, nil
}

// ProvideJaegerTracer creates a new opentracing tracer.
func ProvideJaegerTracer(cfg *Config) (opentracing.Tracer, func(), error) {
	config := jaegercfg.Configuration{
		ServiceName: cfg.Trace.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaggeropentracing.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			QueueSize:          16 << 10,
			LocalAgentHostPort: os.Getenv("JAEGER_AGENT_ENDPOINT"),
			CollectorEndpoint:  os.Getenv("JAEGER_COLLECTOR_ENDPOINT"),
			User:               os.Getenv("JAEGER_USER"),
			Password:           os.Getenv("JAEGER_PASSWORD"),
		},
		Tags: make([]opentracing.Tag, 0, len(cfg.Trace.Tags)),
	}
	for key, value := range cfg.Trace.Tags {
		config.Tags = append(config.Tags, opentracing.Tag{
			Key:   key,
			Value: value,
		})
	}
	tracer, closer, err := config.NewTracer(jaegercfg.Logger(jaggeropentracing.StdLogger))
	return tracer, func() {
		_ = closer.Close()
	}, err
}
