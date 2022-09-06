// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocgql

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.opencensus.io/trace"
)

// Tracer for opencensus.
type Tracer struct {
	// AllowRoot, if set to true, will allow the creation of root spans in
	// absence of existing spans.
	// Default is to not trace calls if no existing parent span is found.
	AllowRoot bool

	// GetOpAttrs allows to add additional attributes per operation.
	GetOpAttrs func(context.Context) []trace.Attribute

	// Field, if set to true, will enable recording of field spans.
	Field bool

	// Sampler to use when creating spans.
	Sampler trace.Sampler
}

var _ interface {
	graphql.HandlerExtension
	graphql.OperationInterceptor
	graphql.ResponseInterceptor
	graphql.FieldInterceptor
} = &Tracer{}

// ExtensionName returns the metrics extension name.
func (Tracer) ExtensionName() string {
	return "OpenCensusTracing"
}

// Validate the executable graphql schema.
func (Tracer) Validate(graphql.ExecutableSchema) error {
	return nil
}

func (t *Tracer) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	if !t.AllowRoot && trace.FromContext(ctx) == nil {
		return next(ctx)
	}
	oc := graphql.GetOperationContext(ctx)
	if !isSubscription(oc) {
		return next(ctx)
	}

	ctx, span := t.startSpan(ctx, spanNameFromOperation(oc))
	span.AddAttributes(t.operationAttrs(ctx, oc)...)
	go func() {
		<-ctx.Done()
		span.End()
	}()
	return next(ctx)
}

func (t *Tracer) startSpan(ctx context.Context, name string) (context.Context, *trace.Span) {
	return trace.StartSpan(ctx, name, trace.WithSampler(t.Sampler))
}

func spanNameFromOperation(oc *graphql.OperationContext) string {
	if oc.OperationName != "" {
		return oc.OperationName
	}
	if op := oc.Operation; op != nil {
		if op.Name != "" {
			return op.Name
		}
		return string(op.Operation)
	}
	return "operation"
}

func isSubscription(oc *graphql.OperationContext) bool {
	return oc.Operation != nil && oc.Operation.Operation == ast.Subscription
}

// InterceptResponse traces graphql response execution.
func (t *Tracer) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	if !t.AllowRoot && trace.FromContext(ctx) == nil {
		return next(ctx)
	}
	oc := graphql.GetOperationContext(ctx)
	name := spanNameFromOperation(oc)
	ctx, span := t.startSpan(ctx, name)
	defer span.End()
	if !span.IsRecordingEvents() {
		return next(ctx)
	}

	if isSubscription(oc) {
		span.SetName(name + ".response")
	} else {
		span.AddAttributes(t.operationAttrs(ctx, oc)...)
	}
	defer setSpanStatus(ctx, span, graphql.GetErrors)

	return next(ctx)
}

func setSpanStatus(ctx context.Context, span *trace.Span, getErrors func(context.Context) gqlerror.List) {
	if errs := getErrors(ctx); errs != nil {
		span.SetStatus(trace.Status{
			Code:    trace.StatusCodeUnknown,
			Message: errs.Error(),
		})
	}
}

func (t *Tracer) operationAttrs(ctx context.Context, oc *graphql.OperationContext) []trace.Attribute {
	var attrs []trace.Attribute
	if t.GetOpAttrs != nil {
		attrs = append(attrs, t.GetOpAttrs(ctx)...)
	}
	if op := oc.Operation; op != nil {
		attrs = append(attrs,
			trace.StringAttribute("graphql.operation", string(op.Operation)),
		)
	}
	attrs = append(attrs,
		trace.StringAttribute("graphql.query", oc.RawQuery),
	)
	for name, value := range oc.Variables {
		attrs = append(attrs,
			trace.StringAttribute("graphql.vars."+name, fmt.Sprintf("%+v", value)),
		)
	}
	if stats := extension.GetComplexityStats(ctx); stats != nil {
		attrs = append(attrs,
			trace.Int64Attribute("graphql.complexity.value", int64(stats.Complexity)),
			trace.Int64Attribute("graphql.complexity.limit", int64(stats.ComplexityLimit)),
		)
	}
	return attrs
}

// InterceptField traces graphql field execution.
func (t *Tracer) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	if !t.Field || (!t.AllowRoot && trace.FromContext(ctx) == nil) {
		return next(ctx)
	}
	fc := graphql.GetFieldContext(ctx)
	ctx, span := t.startSpan(ctx, spanNameFromField(fc.Field))
	defer span.End()
	if !span.IsRecordingEvents() {
		return next(ctx)
	}

	span.AddAttributes(fieldAttrs(ctx)...)
	defer setSpanStatus(ctx, span, func(ctx context.Context) gqlerror.List {
		return graphql.GetFieldErrors(ctx, fc)
	})

	return next(ctx)
}

func fieldAttrs(ctx context.Context) []trace.Attribute {
	fc := graphql.GetFieldContext(ctx)
	attrs := []trace.Attribute{
		trace.StringAttribute("graphql.field.path", graphql.GetPath(ctx).String()),
		trace.StringAttribute("graphql.field.name", fc.Field.Name),
		trace.StringAttribute("graphql.field.alias", fc.Field.Alias),
	}
	if object := fc.Field.ObjectDefinition; object != nil {
		attrs = append(attrs,
			trace.StringAttribute("graphql.field.object", object.Name),
		)
	}
	for _, arg := range fc.Field.Arguments {
		attrs = append(attrs,
			trace.StringAttribute("graphql.field.args."+arg.Name, arg.Value.String()),
		)
	}
	return attrs
}

func spanNameFromField(field graphql.CollectedField) string {
	if object := field.ObjectDefinition; object != nil {
		return object.Name + "." + field.Name
	}
	return field.Name
}
