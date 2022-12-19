// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metrics_test

import (
	"github.com/net-auto/resourceManager/server/metrics"
	"github.com/net-auto/resourceManager/telemetry/mocks"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestMetrics(t *testing.T) {
	var e mocks.ViewExporter
	e.On("ExportView", mock.Anything).
		Maybe()
	e.On("ServeHTTP", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			w, ok := args.Get(0).(http.ResponseWriter)
			require.True(t, ok)
			w.WriteHeader(http.StatusOK)
		}).
		Once()

	m := metrics.New(metrics.Config{
		Log:      zaptest.NewLogger(t),
		Exporter: &e,
	})
	srv, err := m.Start("localhost:0")
	require.NoError(t, err)

	rsp, err := http.Get("http://" + srv.Addr + "/metrics")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rsp.StatusCode)
	rsp.Body.Close()

	err = m.Shutdown(srv)
	require.NoError(t, err)
}
