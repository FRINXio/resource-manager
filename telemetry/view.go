// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry

import (
	"fmt"
	"net/http"
	"sync"

	"go.opencensus.io/stats/view"
)

// ViewExporter groups view.Exporter and http.Handler interfaces.
type ViewExporter interface {
	view.Exporter
	http.Handler
}

// ViewExporterOptions defines a set of options shared between view exporters.
type ViewExporterOptions struct {
	Labels map[string]string `env:"TELEMETRY_VIEW_LABELS" help:"Fixed set of labels to add to every view."`
}

// ViewExporterInitFunc is the function that is called to initialize a view exporter.
type ViewExporterInitFunc func(ViewExporterOptions) (ViewExporter, error)

var viewExporters sync.Map

// RegisterViewExporter registers a view exporter.
func RegisterViewExporter(name string, f ViewExporterInitFunc) error {
	if _, loaded := viewExporters.LoadOrStore(name, f); loaded {
		return fmt.Errorf("view exporter %q already registered", name)
	}
	return nil
}

// MustRegisterViewExporter registers a view exporter and panics on error.
func MustRegisterViewExporter(name string, f ViewExporterInitFunc) {
	if err := RegisterViewExporter(name, f); err != nil {
		panic(err)
	}
}

// UnregisterViewExporter unregisters a view exporter.
func UnregisterViewExporter(name string) {
	viewExporters.Delete(name)
}

// AvailableViewExporters gets the names of registered view exporters.
func AvailableViewExporters() []string {
	return availableExporters(&viewExporters)
}

// GetViewExporter gets the specified view exporter passing in the options to the exporter init function.
func GetViewExporter(name string, opts ViewExporterOptions) (ViewExporter, error) {
	f, ok := viewExporters.Load(name)
	if !ok {
		return nil, fmt.Errorf("view exporter %q not found", name)
	}
	return f.(ViewExporterInitFunc)(opts)
}
