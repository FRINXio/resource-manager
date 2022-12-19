// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config is a struct containing configurable settings for the logger.
type Config struct {
	Level  zapcore.Level `name:"log.level" env:"LOG_LEVEL" default:"info" help:"Only log messages with the given severity or above."`
	Format string        `name:"log.format" env:"LOG_FORMAT" default:"console" enum:"console,json" help:"Output format of log messages."`
}

// empty returns true when the config is equal to its zero value.
func (c Config) empty() bool {
	return c == Config{}
}

// New returns a new leveled contextual logger.
func New(config Config) (Logger, error) {
	if config.empty() {
		return NewNopLogger(), nil
	}
	var cfg zap.Config
	switch config.Format {
	case "json":
		cfg = zap.NewProductionConfig()
	case "console":
		cfg = zap.NewDevelopmentConfig()
	default:
		return nil, fmt.Errorf("unrecognized format: %q", config.Format)
	}
	cfg.Level = zap.NewAtomicLevelAt(config.Level)
	logger, err := cfg.Build(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		return nil, fmt.Errorf("building logger: %w", err)
	}
	return NewDefaultLogger(logger), nil
}

// MustNew returns a new leveled contextual logger, and panic on error.
func MustNew(config Config) Logger {
	logger, err := New(config)
	if err != nil {
		panic(err)
	}
	return logger
}
