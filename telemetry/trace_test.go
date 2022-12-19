// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry_test

import (
	"github.com/net-auto/resourceManager/telemetry"
	"sort"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"
)

type traceIniter struct {
	mock.Mock
}

func (ti *traceIniter) Init(opts telemetry.TraceExporterOptions) (trace.Exporter, error) {
	args := ti.Called(opts)
	exporter, _ := args.Get(0).(trace.Exporter)
	return exporter, args.Error(1)
}

func TestGetTraceExporter(t *testing.T) {
	_, err := telemetry.GetTraceExporter("noexist",
		telemetry.TraceExporterOptions{},
	)
	require.EqualError(t, err, `trace exporter "noexist" not found`)
	var ti traceIniter
	ti.On("Init", mock.Anything).Return(nil, nil).Once()
	defer ti.AssertExpectations(t)
	require.NotPanics(t, func() {
		telemetry.MustRegisterTraceExporter(t.Name(), ti.Init)
	})
	defer telemetry.UnregisterTraceExporter(t.Name())
	_, err = telemetry.GetTraceExporter(t.Name(),
		telemetry.TraceExporterOptions{},
	)
	require.NoError(t, err)
}

func TestAvailableTraceExporters(t *testing.T) {
	var ti traceIniter
	defer ti.AssertExpectations(t)
	suffixes := []string{"foo", "bar", "baz"}
	for _, suffix := range suffixes {
		err := telemetry.RegisterTraceExporter(t.Name()+suffix, ti.Init)
		require.NoError(t, err)
	}
	defer func() {
		for _, suffix := range suffixes {
			telemetry.UnregisterTraceExporter(t.Name() + suffix)
		}
	}()
	require.Panics(t, func() {
		telemetry.MustRegisterTraceExporter(t.Name()+suffixes[0], ti.Init)
	})
	exporters := telemetry.AvailableTraceExporters()
	require.True(t, sort.IsSorted(sort.StringSlice(exporters)))
	for _, suffix := range suffixes {
		require.Contains(t, exporters, t.Name()+suffix)
	}
}

func TestWithoutNameSampler(t *testing.T) {
	sampler := telemetry.WithoutNameSampler("foo", "bar")
	decision := sampler(trace.SamplingParameters{Name: "foo"})
	require.False(t, decision.Sample)
	decision = sampler(trace.SamplingParameters{Name: "bar"})
	require.False(t, decision.Sample)
	decision = sampler(trace.SamplingParameters{Name: "baz"})
	require.True(t, decision.Sample)
}
