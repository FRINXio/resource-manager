// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log_test

import (
	"context"
	"errors"
	"github.com/net-auto/resourceManager/logging/log"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type testExporter struct {
	spans []*trace.SpanData
}

func (e *testExporter) ExportSpan(s *trace.SpanData) {
	e.spans = append(e.spans, s)
}

func newSpanCore(span *trace.Span) zapcore.Core {
	return log.NewSpanCore(span, zap.InfoLevel)
}

func TestSpanCoreCheck(t *testing.T) {
	tests := []struct {
		name    string
		sampler trace.Sampler
		level   zapcore.Level
		expect  func(*testing.T, []*trace.SpanData)
	}{
		{
			name:    "suppressed-level",
			sampler: trace.AlwaysSample(),
			level:   zap.DebugLevel,
			expect: func(t *testing.T, spans []*trace.SpanData) {
				require.Len(t, spans, 1)
				require.Empty(t, spans[0].Attributes)
				require.Empty(t, spans[0].Annotations)
			},
		},
		{
			name:    "suppressed-sampler",
			sampler: trace.NeverSample(),
			level:   zap.ErrorLevel,
			expect: func(t *testing.T, spans []*trace.SpanData) {
				require.Empty(t, spans)
			},
		},
		{
			name:    "emitted",
			sampler: trace.AlwaysSample(),
			level:   zap.InfoLevel,
			expect: func(t *testing.T, spans []*trace.SpanData) {
				require.Len(t, spans, 1)
				require.Len(t, spans[0].Annotations, 1)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			exporter := &testExporter{}
			trace.RegisterExporter(exporter)

			_, span := trace.StartSpan(context.Background(), "test",
				trace.WithSampler(tc.sampler))
			logger := zap.New(newSpanCore(span))
			if ce := logger.Check(tc.level, tc.level.String()+" message"); ce != nil {
				ce.Write()
			}
			span.End()

			tc.expect(t, exporter.spans)
		})
	}
}

func TestSpanCoreWith(t *testing.T) {
	exporter := &testExporter{}
	trace.RegisterExporter(exporter)
	_, span := trace.StartSpan(context.Background(), "test",
		trace.WithSampler(trace.AlwaysSample()))

	root := zap.New(newSpanCore(span))
	root = root.With(zap.String("root", "root"))
	left := root.With(zap.String("left", "left"))
	right := root.With(zap.String("right", "right"))
	leaf := left.With(zap.String("leaf", "leaf"))

	loggers := []*zap.Logger{root, left, right, leaf}
	for _, logger := range loggers {
		logger.Info("")
	}

	span.End()
	spans := exporter.spans
	require.Len(t, spans, 1)

	annotations := spans[0].Annotations
	require.Len(t, annotations, len(loggers))

	require.Len(t, annotations[0].Attributes, 2)
	require.Equal(t, "root", annotations[0].Attributes["root"])

	require.Len(t, annotations[1].Attributes, 3)
	require.Equal(t, "root", annotations[1].Attributes["root"])
	require.Equal(t, "left", annotations[1].Attributes["left"])

	require.Len(t, annotations[2].Attributes, 3)
	require.Equal(t, "root", annotations[2].Attributes["root"])
	require.Equal(t, "right", annotations[2].Attributes["right"])

	require.Len(t, annotations[3].Attributes, 4)
	require.Equal(t, "root", annotations[3].Attributes["root"])
	require.Equal(t, "left", annotations[3].Attributes["left"])
	require.Equal(t, "leaf", annotations[3].Attributes["leaf"])
}

type loggable struct{ bool }

func (l loggable) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if !l.bool {
		return errors.New("can't marshal")
	}
	enc.AddString("loggable", "yes")
	return nil
}

func TestSpanCoreWrite(t *testing.T) {
	exporter := &testExporter{}
	trace.RegisterExporter(exporter)

	_, span := trace.StartSpan(context.Background(), "test",
		trace.WithSampler(trace.AlwaysSample()))
	logger := zap.New(newSpanCore(span))
	logger.Info("field dump",
		zap.Bool("b", true),
		zap.Float32("f32", math.Pi),
		zap.Float64("f64", math.E),
		zap.Int("i", 0),
		zap.Int8("i8", -8),
		zap.Int16("i16", -16),
		zap.Int32("i32", -32),
		zap.Int64("i64", -64),
		zap.Uintptr("ptr", 0xbadbeef),
		zap.Uint("u", 0),
		zap.Uint8("u8", 8),
		zap.Uint16("u16", 16),
		zap.Uint32("u32", 32),
		zap.Uint64("u64", 64),
		zap.Complex64("c64", 1+1i),
		zap.Complex128("c128", 2+2i),
		zap.Duration("duration", time.Second+time.Second/2),
		zap.Time("date", time.Now()),
		zap.Binary("bin", []byte{5, 4, 3}),
		zap.ByteString("bytes", []byte{1, 2, 3}),
		zap.Reflect("numbers", []int{1, 2, 3}),
		zap.Bools("bools", []bool{true, true, false}),
		zap.Object("obj", loggable{true}),
	)
	span.End()

	spans := exporter.spans
	require.Len(t, spans, 1)
	annotations := spans[0].Annotations
	require.Len(t, annotations, 1)
	annotation := annotations[0]
	require.Equal(t, true, annotation.Attributes["b"])
	require.EqualValues(t, math.Float32bits(math.Pi), annotation.Attributes["f32"])
	require.EqualValues(t, math.Float64bits(math.E), annotation.Attributes["f64"])
	require.EqualValues(t, 0, annotation.Attributes["i"])
	require.EqualValues(t, -8, annotation.Attributes["i8"])
	require.EqualValues(t, -16, annotation.Attributes["i16"])
	require.EqualValues(t, -32, annotation.Attributes["i32"])
	require.EqualValues(t, -64, annotation.Attributes["i64"])
	require.EqualValues(t, 0xbadbeef, annotation.Attributes["ptr"])
	require.EqualValues(t, 0, annotation.Attributes["u"])
	require.EqualValues(t, 8, annotation.Attributes["u8"])
	require.EqualValues(t, 16, annotation.Attributes["u16"])
	require.EqualValues(t, 32, annotation.Attributes["u32"])
	require.EqualValues(t, 64, annotation.Attributes["u64"])
	require.Equal(t, "(1+1i)", annotation.Attributes["c64"])
	require.Equal(t, "(2+2i)", annotation.Attributes["c128"])
	require.Equal(t, "1.5s", annotation.Attributes["duration"])
	require.NotEmpty(t, annotation.Attributes["date"])
	require.Equal(t, "BQQD", annotation.Attributes["bin"])
	require.Equal(t, "\x01\x02\x03", annotation.Attributes["bytes"])
	require.Equal(t, "[1 2 3]", annotation.Attributes["numbers"])
	require.Equal(t, "[true true false]", annotation.Attributes["bools"])
	require.Equal(t, "map[loggable:yes]", annotation.Attributes["obj"])
}

func TestSpanCoreOnPanic(t *testing.T) {
	exporter := &testExporter{}
	trace.RegisterExporter(exporter)
	_, span := trace.StartSpan(context.Background(), "test",
		trace.WithSampler(trace.AlwaysSample()))

	logger := zap.New(newSpanCore(span), zap.AddStacktrace(zap.PanicLevel))
	require.Panics(t, func() { logger.Panic("oh no!") })
	span.End()

	spans := exporter.spans
	require.Len(t, spans, 1)
	require.Equal(t, true, spans[0].Attributes["error"])
	annotations := spans[0].Annotations
	require.Len(t, annotations, 1)
	var annotation *trace.Annotation
	require.Condition(t, func() bool {
		for i := range annotations {
			if annotations[i].Message == "oh no!" {
				annotation = &annotations[i]
				return true
			}
		}
		return false
	})
	require.NotNil(t, annotation.Attributes)
	require.Equal(t, "panic", annotation.Attributes["level"])
	require.NotEmpty(t, annotation.Attributes["stack"])
}

func TestSpanCoreSync(t *testing.T) {
	core := newSpanCore(nil)
	require.NoError(t, core.Sync())
}
