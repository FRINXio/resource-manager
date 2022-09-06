// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry

import (
	"fmt"
	"sort"
	"sync"

	"github.com/scylladb/go-set/strset"
	"go.opencensus.io/trace"
)

// TraceExporterOptions defines a set of options shared between trace exporters.
type TraceExporterOptions struct {
	ServiceName string            `name:"service" env:"TELEMETRY_TRACE_SERVICE" default:"${service_name}" help:"Traced service name."`
	Tags        map[string]string `env:"TELEMETRY_TRACE_TAGS" help:"Fixed set of tags to add to every trace."`
}

// TraceExporterInitFunc is the function that is called to initialize a trace exporter.
type TraceExporterInitFunc func(TraceExporterOptions) (trace.Exporter, error)

var traceExporters sync.Map

// RegisterTraceExporter registers a trace exporter.
func RegisterTraceExporter(name string, f TraceExporterInitFunc) error {
	if _, loaded := traceExporters.LoadOrStore(name, f); loaded {
		return fmt.Errorf("trace exporter %q already registered", name)
	}
	return nil
}

// MustRegisterTraceExporter registers a trace exporter and panics on error.
func MustRegisterTraceExporter(name string, f TraceExporterInitFunc) {
	if err := RegisterTraceExporter(name, f); err != nil {
		panic(err)
	}
}

// UnregisterTraceExporter unregisters a trace exporter.
func UnregisterTraceExporter(name string) {
	traceExporters.Delete(name)
}

// GetTraceExporter gets the specified trace exporter passing in the options to the exporter init function.
func GetTraceExporter(name string, opts TraceExporterOptions) (trace.Exporter, error) {
	f, ok := traceExporters.Load(name)
	if !ok {
		return nil, fmt.Errorf("trace exporter %q not found", name)
	}
	return f.(TraceExporterInitFunc)(opts)
}

func availableExporters(exporters *sync.Map) []string {
	var names []string
	exporters.Range(func(key, _ interface{}) bool {
		names = append(names, key.(string))
		return true
	})
	sort.Strings(names)
	return names
}

// AvailableTraceExporters gets the names of registered trace exporters.
func AvailableTraceExporters() []string {
	return availableExporters(&traceExporters)
}

// WithoutNameSampler returns a trace sampler filtering out a set of span names.
func WithoutNameSampler(name string, names ...string) trace.Sampler {
	if len(names) == 0 {
		return func(p trace.SamplingParameters) trace.SamplingDecision {
			return trace.SamplingDecision{
				Sample: p.Name != name,
			}
		}
	}
	set := strset.New(names...)
	set.Add(name)
	return func(p trace.SamplingParameters) trace.SamplingDecision {
		return trace.SamplingDecision{Sample: !set.Has(p.Name)}
	}
}
