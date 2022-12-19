// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocgql_test

import (
	"context"
	"errors"
	"github.com/net-auto/resourceManager/telemetry/ocgql"
	"sync"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/testserver"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/cenkalti/backoff/v4"
	"github.com/stretchr/testify/suite"
	"go.opencensus.io/trace"
)

type tracerTestSuite struct {
	suite.Suite
	sample bool
	client *client.Client
	server *testserver.TestServer
	spans  sync.Map
}

func (s *tracerTestSuite) SetupSuite() {
	srv := testserver.New()
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{})
	srv.Use(extension.FixedComplexityLimit(1000))
	srv.SetCalculatedComplexity(100)
	srv.Use(&ocgql.Tracer{
		AllowRoot: true,
		GetOpAttrs: func(context.Context) []trace.Attribute {
			return []trace.Attribute{
				trace.StringAttribute("username", "tester"),
			}
		},
		Field: true,
		Sampler: func(trace.SamplingParameters) trace.SamplingDecision {
			return trace.SamplingDecision{Sample: s.sample}
		},
	})
	srv.AroundOperations(
		func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
			ctx, span := trace.StartSpan(ctx, "test",
				trace.WithSampler(trace.AlwaysSample()),
			)
			defer span.End()
			return next(ctx)
		},
	)
	s.client = client.New(srv)
	s.server = srv
	trace.RegisterExporter(s)
}

func (s *tracerTestSuite) TearDownSuite() {
	trace.UnregisterExporter(s)
}

func (s *tracerTestSuite) SetupTest() {
	s.sample = true
}

func (s *tracerTestSuite) TearDownTest() {
	s.spans = sync.Map{}
}

func (s *tracerTestSuite) ExportSpan(span *trace.SpanData) {
	s.spans.Store(span.Name, span)
}

func (s *tracerTestSuite) GetSpan(name string) *trace.SpanData {
	span, ok := s.spans.Load(name)
	if !ok {
		return nil
	}
	return span.(*trace.SpanData)
}

func TestTracer(t *testing.T) {
	suite.Run(t, &tracerTestSuite{})
}

func (s *tracerTestSuite) TestOperation() {
	const (
		query = "query($id: Int!) { name: find(id: $id) }"
		id    = "42"
	)
	err := s.post(query, client.Var("id", id))
	s.Require().NoError(err)

	span := s.GetSpan("query")
	s.Require().NotNil(span)
	s.Require().Equal("tester", span.Attributes["username"])
	s.Require().Equal(query, span.Attributes["graphql.query"])
	s.Require().Equal(id, span.Attributes["graphql.vars.id"])
	s.Require().EqualValues(100, span.Attributes["graphql.complexity.value"])
	s.Require().EqualValues(1000, span.Attributes["graphql.complexity.limit"])
	s.Require().EqualValues(trace.StatusCodeOK, span.Code)
	s.Require().Empty(span.Message)

	span = s.GetSpan("name")
	s.Require().NotNil(span)
	for _, attr := range []string{"path", "name", "alias"} {
		s.Require().Equal("name", span.Attributes["graphql.field."+attr])
	}
	s.Require().EqualValues(trace.StatusCodeOK, span.Code)
	s.Require().Empty(span.Message)
}

func (s *tracerTestSuite) TestWithoutSampling() {
	s.sample = false
	err := s.post("query { name }")
	s.Require().NoError(err)
	span := s.GetSpan("query")
	s.Require().Nil(span)
}

func (s *tracerTestSuite) TestNamedOperation() {
	err := s.post("query foobar { name }")
	s.Require().NoError(err)

	span := s.GetSpan("foobar")
	s.Require().NotNil(span)
	s.Require().EqualValues(trace.StatusCodeOK, span.Code)
}

func (s *tracerTestSuite) TestSubscription() {
	sk := s.client.Websocket(`subscription { name }`)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := sk.Next(&struct{ Name string }{})
		s.Require().NoError(err)
	}()
	s.server.SendNextSubscriptionMessage()
	wg.Wait()
	err := sk.Close()
	s.Require().NoError(err)

	ctx, cancel := context.WithTimeout(
		context.Background(), time.Second,
	)
	defer cancel()
	err = backoff.Retry(
		func() error {
			if s.GetSpan("subscription") == nil ||
				s.GetSpan("subscription.response") == nil {
				return errors.New("span not found")
			}
			return nil
		},
		backoff.WithContext(
			backoff.NewConstantBackOff(10*time.Millisecond), ctx,
		),
	)
	s.Require().NoError(err)
}

func (s *tracerTestSuite) post(query string, opts ...client.Option) error {
	return s.client.Post(query, &struct{ Name string }{}, opts...)
}
