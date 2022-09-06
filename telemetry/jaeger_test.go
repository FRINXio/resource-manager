// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry_test

import (
	"github.com/net-auto/resourceManager/telemetry"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewJaegerExporter(t *testing.T) {
	err := os.Setenv("JAEGER_AGENT_ENDPOINT", "localhost:6831")
	require.NoError(t, err)
	defer func() {
		err := os.Unsetenv("JAEGER_AGENT_ENDPOINT")
		require.NoError(t, err)
	}()
	exporter, err := telemetry.GetTraceExporter("jaeger",
		telemetry.TraceExporterOptions{ServiceName: t.Name()},
	)
	require.NoError(t, err)
	require.NotNil(t, exporter)
}
