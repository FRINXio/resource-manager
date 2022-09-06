// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log_test

import (
	"bytes"
	log2 "github.com/net-auto/resourceManager/logging/log"
	stdlog "log"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestProvider(t *testing.T) {
	var buf bytes.Buffer
	stdlog.SetOutput(&buf)
	logger, restorer, err := log2.ProvideLogger(log2.Config{})
	require.NoError(t, err)
	defer restorer()
	require.Equal(t, logger.Background(), log2.ProvideZapLogger(logger))
	require.Equal(t, logger.Background(), zap.L())
	stdlog.Println("suppressed message")
	require.Zero(t, buf.Len())
}
