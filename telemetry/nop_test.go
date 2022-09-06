// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry_test

import (
	telemetry2 "github.com/net-auto/resourceManager/telemetry"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNopExporter(t *testing.T) {
	_, err := telemetry2.GetTraceExporter("nop",
		telemetry2.TraceExporterOptions{},
	)
	require.NoError(t, err)
	_, err = telemetry2.GetViewExporter("nop",
		telemetry2.ViewExporterOptions{},
	)
	require.NoError(t, err)
}
