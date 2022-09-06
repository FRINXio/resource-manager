// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logtest

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
)

// TestLogger is a logger to be used for testing.
type TestLogger struct {
	logger *zap.Logger
}

// TestingT is a subset of the API provided by all *testing.T and *testing.B objects.
type TestingT = zaptest.TestingT

// NewTestLogger creates a new testing logger.
func NewTestLogger(t TestingT) *TestLogger {
	logger := zaptest.NewLogger(t, zaptest.WrapOptions(zap.AddCaller()))
	return &TestLogger{logger}
}

// Background returns a context-unaware logger.
func (l TestLogger) Background() *zap.Logger {
	return l.logger
}

// For ignores context and returns background logger.
func (l TestLogger) For(context.Context) *zap.Logger {
	return l.Background()
}

// WithObserver creates an observer for emitted logs.
func (l TestLogger) WithObserver(enabler zapcore.LevelEnabler) (TestLogger, *observer.ObservedLogs) {
	core, o := observer.New(enabler)
	logger := l.logger.WithOptions(
		zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(c, core)
		}),
	)
	return TestLogger{logger}, o
}
