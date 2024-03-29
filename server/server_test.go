// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server_test

import (
	"context"
	"errors"
	"github.com/net-auto/resourceManager/server"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"
	"gocloud.dev/server/requestlog"
)

func TestListenAndServe(t *testing.T) {
	var td testDriver
	td.On("ListenAndServe", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			h, _ := args.Get(1).(http.Handler)
			assert.NotNil(t, h)
		}).
		Return(nil).
		Once()
	defer td.AssertExpectations(t)

	s := server.New(http.NotFoundHandler(), &server.Options{Driver: &td})
	err := s.ListenAndServe(":8080")
	require.NoError(t, err)
}

func TestMiddleware(t *testing.T) {
	var tl testLogger
	tl.On("Log", mock.Anything).
		Run(func(args mock.Arguments) {
			ent, _ := args.Get(0).(*requestlog.Entry)
			require.NotNil(t, ent)
			assert.NotEmpty(t, ent.TraceID)
			assert.NotEmpty(t, ent.SpanID)
		}).
		Once()
	defer tl.AssertExpectations(t)

	var td testDriver
	td.On("ListenAndServe", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			handler, _ := args.Get(1).(http.Handler)
			require.NotNil(t, handler)
			tests := []struct {
				target string
				expect func(*testing.T, *httptest.ResponseRecorder)
			}{
				{
					target: "/",
					expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
						assert.Equal(t, http.StatusInternalServerError, rec.Code)
						assert.NotEmpty(t, rec.Header().Get("X-Correlation-ID"))
					},
				},
				{
					target: "/debug/pprof/",
					expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
						assert.Equal(t, http.StatusOK, rec.Code)
						assert.NotZero(t, rec.Body.Len())
					},
				},
			}
			for _, tt := range tests {
				req := httptest.NewRequest(http.MethodGet, tt.target, nil)
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, req)
				tt.expect(t, rec)
			}
		}).
		Return(nil).
		Once()
	defer td.AssertExpectations(t)

	var te testExporter
	te.On("ExportSpan", mock.AnythingOfType("*trace.SpanData")).
		Run(func(args mock.Arguments) {
			s, _ := args.Get(0).(*trace.SpanData)
			assert.Equal(t, "/", s.Name)
		}).
		Once()
	te.On("ExportView", mock.Anything).
		Maybe()
	defer trace.UnregisterExporter(&te)
	defer te.AssertExpectations(t)

	var tr testRecovery
	tr.On("Recover", mock.Anything, mock.Anything).
		Return(nil).
		Once()
	defer tr.AssertExpectations(t)

	handler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic(errors.New("bad handler"))
	})
	s := server.New(handler, &server.Options{
		RequestLogger:         &tl,
		TraceExporter:         &te,
		DefaultSamplingPolicy: trace.AlwaysSample(),
		EnableProfiling:       true,
		RecoveryHandler:       tr.Recover,
		Driver:                &td,
	})
	err := s.ListenAndServe(":8080")
	assert.NoError(t, err)
}

type testDriver struct {
	mock.Mock
}

func (td *testDriver) ListenAndServe(addr string, h http.Handler) error {
	args := td.Called(addr, h)
	return args.Error(0)
}

func (td *testDriver) Shutdown(ctx context.Context) error {
	args := td.Called(ctx)
	return args.Error(0)
}

type testLogger struct {
	mock.Mock
}

func (tl *testLogger) Log(ent *requestlog.Entry) {
	tl.Called(ent)
}

type testExporter struct {
	mock.Mock
}

func (te *testExporter) ExportSpan(s *trace.SpanData) {
	te.Called(s)
}

type testRecovery struct {
	mock.Mock
}

func (tr *testRecovery) Recover(ctx context.Context, p interface{}) error {
	args := tr.Called(ctx, p)
	return args.Error(0)
}
