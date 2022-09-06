// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gqlutil

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/net-auto/resourceManager/telemetry/ocgql"
	"github.com/net-auto/resourceManager/viewer"
	"go.opencensus.io/trace"
)

// NewServer creates a default graphql server.
func NewServer(es graphql.ExecutableSchema) *handler.Server {
	srv := handler.New(es)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New(1000))
	srv.Use(extension.Introspection{})
	srv.Use(ocgql.Metrics{})
	srv.Use(&ocgql.Tracer{
		GetOpAttrs: func(ctx context.Context) []trace.Attribute {
			if v := viewer.FromContext(ctx); v != nil {
				return viewer.TraceAttrs(v)
			}
			return nil
		},
	})
	return srv
}
