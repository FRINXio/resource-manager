// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log_test

import (
	"context"
	log2 "github.com/net-auto/resourceManager/logging/log"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggerFieldContext(t *testing.T) {
	core, o := observer.New(zap.InfoLevel)
	logger := log2.NewDefaultLogger(zap.New(core))

	ctx := log2.NewFieldsContext(context.Background(), zap.String("name", "test"))
	ctx = log2.NewFieldsContext(ctx, zap.String("lang", "go"))
	logger.For(ctx).Info("test message", zap.Int("speed", 42))

	require.Equal(t, 1, o.
		FilterMessage("test message").
		FilterField(zap.String("name", "test")).
		FilterField(zap.String("lang", "go")).
		Len(),
	)
}
