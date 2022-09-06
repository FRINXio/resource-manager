// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gqlutil_test

import (
	"context"
	"testing"

	"github.com/net-auto/resourceManager/logging/log/logtest"
	"github.com/net-auto/resourceManager/pkg/gqlutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRecoverFunc(t *testing.T) {
	logger, observer := logtest.NewTestLogger(t).WithObserver(zap.ErrorLevel)
	err := gqlutil.RecoverFunc(logger)(context.Background(), "oh no")
	require.Equal(t, 1, observer.
		FilterMessage("graphql panic recovery").
		FilterField(zap.String("error", "oh no")).Len())
	require.EqualError(t, err, "internal system error")
}
