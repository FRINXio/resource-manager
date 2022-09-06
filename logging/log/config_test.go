// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log_test

import (
	"github.com/net-auto/resourceManager/logging/log"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConfigParse(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		var cfg log.Config
		parser, err := kong.New(&cfg)
		require.NoError(t, err)
		_, err = parser.Parse(nil)
		require.NoError(t, err)
		require.Equal(t, zap.InfoLevel, cfg.Level)
		require.Equal(t, "console", cfg.Format)
	})
	t.Run("OK", func(t *testing.T) {
		var cfg log.Config
		parser, err := kong.New(&cfg)
		require.NoError(t, err)
		_, err = parser.Parse([]string{
			"--log.level", "error",
			"--log.format", "json",
		})
		require.NoError(t, err)
		require.Equal(t, zap.ErrorLevel, cfg.Level)
		require.Equal(t, "json", cfg.Format)
	})
	t.Run("BadLevel", func(t *testing.T) {
		var cfg log.Config
		parser, err := kong.New(&cfg)
		require.NoError(t, err)
		_, err = parser.Parse([]string{
			"--log.level", "foo",
		})
		require.Error(t, err)
	})
	t.Run("BadFormat", func(t *testing.T) {
		var cfg log.Config
		parser, err := kong.New(&cfg)
		require.NoError(t, err)
		_, err = parser.Parse([]string{
			"--log.format", "bar",
		})
		require.Error(t, err)
	})
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  log.Config
		wantErr bool
	}{
		{
			name: "Production",
			config: log.Config{
				Level:  zap.InfoLevel,
				Format: "json",
			},
		},
		{
			name: "Development",
			config: log.Config{
				Level:  zap.DebugLevel,
				Format: "console",
			},
		},
		{
			name: "Nop",
		},
		{
			name: "BadFormat",
			config: log.Config{
				Format: "fmt",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			logger, err := log.New(tc.config)
			if !tc.wantErr {
				require.NotNil(t, logger)
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestMustNew(t *testing.T) {
	var config log.Config
	require.NotPanics(t, func() { _ = log.MustNew(config) })
	config.Format = "baz"
	require.Panics(t, func() { _ = log.MustNew(config) })
}
