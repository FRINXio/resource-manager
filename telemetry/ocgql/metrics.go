// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocgql

import (
	"context"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/vektah/gqlparser/v2/ast"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

// Metrics for opencensus.
type Metrics struct{}

var _ interface {
	graphql.HandlerExtension
	graphql.OperationInterceptor
	graphql.ResponseInterceptor
	graphql.FieldInterceptor
} = Metrics{}

// ExtensionName returns the metrics extension name.
func (Metrics) ExtensionName() string {
	return "OpenCensusMetrics"
}

// Validate the executable graphql schema.
func (Metrics) Validate(graphql.ExecutableSchema) error {
	return nil
}

// numSubscriptions tracks the current number of subscriptions.
var numSubscriptions int64

// InterceptOperation measures graphql operation execution.
func (Metrics) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	if op := graphql.GetOperationContext(ctx).Operation; op != nil {
		if op.Operation == ast.Subscription {
			stats.Record(ctx, NumSubscriptions.M(
				atomic.AddInt64(&numSubscriptions, 1),
			))
			go func() {
				<-ctx.Done()
				stats.Record(ctx, NumSubscriptions.M(
					atomic.AddInt64(&numSubscriptions, -1),
				))
			}()
		}
		ctx, _ = tag.New(ctx,
			tag.Upsert(Operation, string(op.Operation)),
		)
	}
	stats.Record(ctx, RequestTotal.M(1))
	return next(ctx)
}

// InterceptResponse measures graphql response execution.
func (Metrics) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	rsp := next(ctx)
	end := graphql.Now()

	measurements := []stats.Measurement{
		ResponseTotal.M(1),
	}
	if oc := graphql.GetOperationContext(ctx); oc.Operation == nil || oc.Operation.Operation != ast.Subscription {
		measurements = append(measurements,
			RequestLatency.M(
				float64(end.Sub(oc.Stats.OperationStart))/float64(time.Millisecond),
			),
		)
	}
	if complexity := extension.GetComplexityStats(ctx); complexity != nil {
		measurements = append(measurements,
			RequestComplexity.M(
				int64(complexity.Complexity),
			),
		)
	}
	ctx, _ = tag.New(ctx,
		tag.Upsert(Errors, strconv.Itoa(len(graphql.GetErrors(ctx)))),
	)
	stats.Record(ctx, measurements...)
	return rsp
}

// InterceptField measures graphql field execution.
func (Metrics) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	fc := graphql.GetFieldContext(ctx)
	deprecated := fc.Field.Definition.Directives.ForName("deprecated")
	ctx, _ = tag.New(ctx,
		tag.Upsert(Object, fc.Object),
		tag.Upsert(Field, fc.Field.Name),
		tag.Upsert(GraphQLObject, fc.Object),
		tag.Upsert(GraphQLField, fc.Field.Name),
		tag.Upsert(GraphQLDeprecated, strconv.FormatBool(deprecated != nil)),
	)

	start := graphql.Now()
	res, err := next(ctx)
	end := graphql.Now()

	ctx, _ = tag.New(ctx,
		tag.Upsert(Errors, strconv.Itoa(
			len(graphql.GetFieldErrors(ctx, fc)),
		)),
	)
	measurements := []stats.Measurement{
		ResolveLatency.M(
			float64(end.Sub(start)) / float64(time.Millisecond),
		),
		ResolveTotal.M(1),
	}
	if fc.Field.Definition.Directives.ForName("deprecated") != nil {
		measurements = append(measurements,
			DeprecatedResolveTotal.M(1),
		)
	}
	stats.Record(ctx, measurements...)
	return res, err
}
