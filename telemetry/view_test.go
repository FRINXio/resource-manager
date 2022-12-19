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
)

type viewIniter struct {
	mock.Mock
}

func (vi *viewIniter) Init(opts telemetry.ViewExporterOptions) (telemetry.ViewExporter, error) {
	args := vi.Called(opts)
	exporter, _ := args.Get(0).(telemetry.ViewExporter)
	return exporter, args.Error(1)
}

func TestGetViewExporter(t *testing.T) {
	_, err := telemetry.GetViewExporter("noexist",
		telemetry.ViewExporterOptions{},
	)
	require.EqualError(t, err, `view exporter "noexist" not found`)
	var vi viewIniter
	vi.On("Init", mock.Anything).
		Return(nil, nil).Once()
	defer vi.AssertExpectations(t)
	require.NotPanics(t, func() {
		telemetry.MustRegisterViewExporter(t.Name(), vi.Init)
	})
	defer telemetry.UnregisterViewExporter(t.Name())
	_, err = telemetry.GetViewExporter(t.Name(),
		telemetry.ViewExporterOptions{},
	)
	require.NoError(t, err)
}

func TestAvailableViewExporters(t *testing.T) {
	var vi viewIniter
	defer vi.AssertExpectations(t)
	suffixes := []string{"foo", "bar", "baz"}
	for _, suffix := range suffixes {
		err := telemetry.RegisterViewExporter(t.Name()+suffix, vi.Init)
		require.NoError(t, err)
	}
	defer func() {
		for _, suffix := range suffixes {
			telemetry.UnregisterViewExporter(t.Name() + suffix)
		}
	}()
	require.Panics(t, func() {
		telemetry.MustRegisterViewExporter(t.Name()+suffixes[0], vi.Init)
	})
	exporters := telemetry.AvailableViewExporters()
	require.True(t, sort.IsSorted(sort.StringSlice(exporters)))
	for _, suffix := range suffixes {
		require.Contains(t, exporters, t.Name()+suffix)
	}
}
