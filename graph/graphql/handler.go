// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/facebookincubator/ent-contrib/entgql"
	"github.com/facebookincubator/symphony/graph/graphql/complexity"
	"github.com/facebookincubator/symphony/pkg/gqlutil"

	// "github.com/net-auto/resourceManager/graph/graphql/directive"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/telemetry/ocgql"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/privacy"
	"github.com/net-auto/resourceManager/graph/graphql/generated"
	"github.com/net-auto/resourceManager/graph/graphql/resolver"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.opencensus.io/plugin/ochttp"
	"go.uber.org/zap"
)

// HandlerConfig configures graphql handler.
type HandlerConfig struct {
	Client *ent.Client
	Logger log.Logger
}

func init() {
	// TODO directive
	// views := append(
	// 	ocgql.DefaultServerViews,
	// 	directive.ServerDeprecatedCountByObjectInputField,
	// )
	views := ocgql.DefaultViews

	for _, v := range views {
		v.TagKeys = append(v.TagKeys,
			viewer.KeyTenant,
			viewer.KeyUser,
			viewer.KeyUserAgent,
		)
	}
}

// Modified OpenTxFromContext from ent
func OpenTxFromContext(ctx context.Context) (context.Context, driver.Tx, error) {
	client := ent.FromContext(ctx)
	if client == nil {
		return nil, nil, errors.New("no client attached to context")
	}
	txOptions := sql.TxOptions{}
	opCtx := graphql.GetOperationContext(ctx)
	if opCtx != nil && opCtx.Operation.Operation == ast.Mutation {
		for i := 0; i < len(opCtx.Operation.SelectionSet); i++ {
			field, ok := opCtx.Operation.SelectionSet[i].(*ast.Field)
			if ok {
				if field.Name == "ClaimResource" {
					txOptions = sql.TxOptions{
						Isolation: sql.LevelSerializable,
					}
					break
				}
			}
		}
	}
	tx, err := client.BeginTx(ctx, &txOptions)
	if err != nil {
		return nil, nil, err
	}
	ctx = ent.NewTxContext(ctx, tx)
	ctx = ent.NewContext(ctx, tx.Client())
	return ctx, tx, nil
}

// NewHandler creates a graphql http handler.
func NewHandler(cfg HandlerConfig) (http.Handler, error) {
	rsv := resolver.New(
		resolver.Config{
			Logger: cfg.Logger,
		},
	)

	router := mux.NewRouter()
	router.Use(func(handler http.Handler) http.Handler {
		timeouter := http.TimeoutHandler(handler, 3*time.Minute, "request timed out")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := timeouter
			if websocket.IsWebSocketUpgrade(r) {
				h = handler
			}
			h.ServeHTTP(w, r)
		})
	})

	srv := gqlutil.NewServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: rsv,
				// TODO directive
				// Directives: directive.New(cfg.Logger),
				Directives: generated.DirectiveRoot{},
			},
		),
	)

	srv.Use(entgql.Transactioner{
		TxOpener: entgql.TxOpenerFunc(OpenTxFromContext),
	})

	srv.SetErrorPresenter(errorPresenter(cfg.Logger))
	srv.SetRecoverFunc(gqlutil.RecoverFunc(cfg.Logger))
	srv.Use(extension.FixedComplexityLimit(complexity.Infinite))

	router.Path("/graphiql").
		Handler(
			ochttp.WithRouteTag(
				playground.Handler(
					"GraphQL playground",
					"/graph/query",
				),
				"graphiql",
			),
		)
	router.Path("/query").
		Handler(
			ochttp.WithRouteTag(
				gziphandler.GzipHandler(srv),
				"query",
			),
		)

	return router, nil
}

func errorPresenter(logger log.Logger) graphql.ErrorPresenterFunc {
	return func(ctx context.Context, err error) (gqlerr *gqlerror.Error) {
		defer func() {
			if errors.Is(err, privacy.Deny) {
				gqlerr.Message = "Permission denied"
			}
		}()
		if errors.As(err, &gqlerr) {
			if gqlerr.Path == nil {
				gqlerr.Path = graphql.GetPath(ctx)
			}
			return gqlerr
		}
		logger.For(ctx).
			Error("graphql internal failure",
				zap.Error(err),
			)
		return &gqlerror.Error{
			Message: "Sorry, something went wrong",
			Path:    graphql.GetPath(ctx),
		}
	}
}
