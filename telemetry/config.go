// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/google/wire"
	"go.opencensus.io/trace"
)

// Config is a struct containing configurable telemetry settings.
type Config struct {
	Trace struct {
		ExporterName         string   `name:"exporter" env:"TELEMETRY_TRACE_EXPORTER" default:"nop" enum:"${trace_exporters}" help:"Exporter to use when exporting telemetry trace data."`
		SamplingProbability  float64  `name:"sampling_probability" env:"TELEMETRY_TRACE_SAMPLING_PROBABILITY" default:"1.0" help:"Sampling probability for trace creation."`
		ExcludeSpanNames     []string `name:"exclude_span_names" env:"TELEMETRY_TRACE_EXCLUDE_SPAN_NAMES" help:"List of span names to exclude on creation."`
		TraceExporterOptions `embed:""`
	} `prefix:"telemetry.trace." embed:""`
	View struct {
		ExporterName        string `name:"exporter" env:"TELEMETRY_VIEW_EXPORTER" default:"prometheus" enum:"${view_exporters}" help:"Exporter to use when exporting telemetry metrics data."`
		ViewExporterOptions `embed:""`
	} `prefix:"telemetry.view." embed:""`
}

// Apply implements kong.Option interface.
func (Config) Apply(k *kong.Kong) error {
	vars := kong.Vars{
		"service_name": func() string {
			exec, _ := os.Executable()
			return filepath.Base(exec)
		}(),
		"trace_exporters": strings.Join(
			AvailableTraceExporters(), ",",
		),
		"view_exporters": strings.Join(
			AvailableViewExporters(), ",",
		),
	}
	return vars.Apply(k)
}

// ProvideTraceExporter is a wire provider that produces trace exporter from config.
func ProvideTraceExporter(config Config) (trace.Exporter, func(), error) {
	exporter, err := GetTraceExporter(
		config.Trace.ExporterName,
		config.Trace.TraceExporterOptions,
	)
	if err != nil {
		return nil, nil, err
	}
	if flusher, ok := exporter.(interface{ Flush() }); ok {
		return exporter, flusher.Flush, nil
	}
	return exporter, func() {}, nil
}

// ProvideTraceSampler is a wire provider that produces trace sampler from config.
func ProvideTraceSampler(config Config) trace.Sampler {
	sampler := trace.ProbabilitySampler(
		config.Trace.SamplingProbability,
	)
	if len(config.Trace.ExcludeSpanNames) == 0 {
		return sampler
	}
	filter := WithoutNameSampler(
		config.Trace.ExcludeSpanNames[0],
		config.Trace.ExcludeSpanNames[1:]...,
	)
	return func(p trace.SamplingParameters) trace.SamplingDecision {
		if decision := filter(p); !decision.Sample {
			return decision
		}
		return sampler(p)
	}
}

// ProvideViewExporter is a wire provider that produces view exporter from config.
func ProvideViewExporter(config Config) (ViewExporter, error) {
	return GetViewExporter(
		config.View.ExporterName,
		config.View.ViewExporterOptions,
	)
}

// Provider is a wire provider that produces telemetry exports from config.
var Provider = wire.NewSet(
	ProvideTraceExporter,
	ProvideTraceSampler,
	ProvideViewExporter,
)
