// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocgql

import (
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// The following measures are supported for use in custom views.
var (
	RequestTotal = stats.Int64(
		"graphql/requests_total",
		"Total number of requests",
		stats.UnitDimensionless,
	)
	RequestLatency = stats.Float64(
		"graphql/request_latency_milliseconds",
		"Latency of requests",
		stats.UnitMilliseconds,
	)
	ResponseTotal = stats.Int64(
		"graphql/responses_total",
		"Total number of responses",
		stats.UnitDimensionless,
	)
	ResolveTotal = stats.Int64(
		"graphql/resolves_total",
		"Total number of resolves",
		stats.UnitDimensionless,
	)
	ResolveLatency = stats.Float64(
		"graphql/resolve_latency_milliseconds",
		"Latency of resolves",
		stats.UnitMilliseconds,
	)
	DeprecatedResolveTotal = stats.Int64(
		"graphql/deprecated_resolves_total",
		"Total number of deprecated resolves",
		stats.UnitDimensionless,
	)
	RequestComplexity = stats.Int64(
		"graphql/request_complexities",
		"Complexity of requests",
		stats.UnitDimensionless,
	)
	NumSubscriptions = stats.Int64(
		"graphql/current_subscriptions",
		"Current number of subscriptions",
		stats.UnitDimensionless,
	)
)

// The following tags are applied to stats recorded by this package.
var (
	// Operation is the operation (query, mutation or subscription).
	Operation = tag.MustNewKey("operation")

	// Object is the object being resolved.
	Object = tag.MustNewKey("object")

	// Field is the field being resolved.
	Field = tag.MustNewKey("field")

	// Errors is the number of request/resolve errors.
	Errors = tag.MustNewKey("errors")
)

// The following tags are deprecated.
var (
	GraphQLObject     = tag.MustNewKey("graphql.object")
	GraphQLField      = tag.MustNewKey("graphql.field")
	GraphQLDeprecated = tag.MustNewKey("graphql.deprecated")
)

// Package ocgql provides some convenience views for measures.
// You still need to register these views for data to actually be collected.
var (
	RequestTotalView = &view.View{
		Name:        RequestTotal.Name(),
		Description: RequestTotal.Description(),
		TagKeys:     []tag.Key{Operation},
		Measure:     RequestTotal,
		Aggregation: view.Count(),
	}
	RequestLatencyView = &view.View{
		Name:        RequestLatency.Name(),
		Description: RequestLatency.Description(),
		Measure:     RequestLatency,
		Aggregation: ochttp.DefaultLatencyDistribution,
	}
	ResponseTotalView = &view.View{
		Name:        ResponseTotal.Name(),
		Description: ResponseTotal.Description(),
		TagKeys:     []tag.Key{Operation, Errors},
		Measure:     ResponseTotal,
		Aggregation: view.Count(),
	}
	ResolveTotalView = &view.View{
		Name:        ResolveTotal.Name(),
		Description: ResolveTotal.Description(),
		TagKeys:     []tag.Key{Object, Field, Errors},
		Measure:     ResolveTotal,
		Aggregation: view.Count(),
	}
	ResolveLatencyView = &view.View{
		Name:        ResolveLatency.Name(),
		Description: ResolveLatency.Description(),
		Measure:     ResolveLatency,
		Aggregation: ocgrpc.DefaultMillisecondsDistribution,
	}
	DeprecatedResolveTotalView = &view.View{
		Name:        DeprecatedResolveTotal.Name(),
		Description: DeprecatedResolveTotal.Description(),
		TagKeys:     []tag.Key{Object, Field},
		Measure:     DeprecatedResolveTotal,
		Aggregation: view.Count(),
	}
	RequestComplexityView = &view.View{
		Name:        RequestComplexity.Name(),
		Description: RequestComplexity.Description(),
		Measure:     RequestComplexity,
		Aggregation: view.Distribution(10, 100, 1000, 5000, 10000, 50000, 100000, 1000000),
	}
	NumSubscriptionsView = &view.View{
		Name:        NumSubscriptions.Name(),
		Description: NumSubscriptions.Description(),
		Measure:     NumSubscriptions,
		Aggregation: view.LastValue(),
	}
)

// The following views are deprecated.
var (
	ResolveCountByObjectField = &view.View{
		Name:        "graphql/server/resolve_count_by_object_field",
		Description: "Count of GraphQL resolves by object and field",
		TagKeys:     []tag.Key{GraphQLObject, GraphQLField, GraphQLDeprecated},
		Measure:     ResolveTotal,
		Aggregation: view.Count(),
	}
)

// DefaultViews are the default views provided by this package.
var DefaultViews = []*view.View{
	RequestTotalView,
	RequestLatencyView,
	ResponseTotalView,
	ResolveTotalView,
	ResolveLatencyView,
	DeprecatedResolveTotalView,
	RequestComplexityView,
	NumSubscriptionsView,
	ResolveCountByObjectField,
}
