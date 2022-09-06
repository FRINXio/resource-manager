// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gqlutil_test

import (
	"context"
	"database/sql"
	"io"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/testserver"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/net-auto/resourceManager/pkg/gqlutil"
	"github.com/stretchr/testify/require"
)

func TestDBInjector(t *testing.T) {
	newServer := func(db *sql.DB) *testserver.TestServer {
		srv := testserver.New()
		srv.AddTransport(transport.POST{})
		srv.Use(gqlutil.DBInjector{DB: db})
		return srv
	}

	t.Run("Query", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer func() {
			err := mock.ExpectationsWereMet()
			require.NoError(t, err)
		}()

		srv := newServer(db)
		srv.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
			require.Equal(t, db, gqlutil.DBFromContext(ctx))
			return next(ctx)
		})

		c := client.New(srv)
		err = c.Post(`query { name }`, &struct{ Name string }{})
		require.NoError(t, err)
	})
	t.Run("Mutation", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			mock.ExpectBegin()
			mock.ExpectCommit()
			defer func() {
				err := mock.ExpectationsWereMet()
				require.NoError(t, err)
			}()

			srv := newServer(db)
			srv.AroundResponses(func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
				require.NotNil(t, gqlutil.TxFromContext(ctx))
				return &graphql.Response{Data: []byte(`{"name":"test"}`)}
			})

			c := client.New(srv)
			err = c.Post(`mutation { name }`, &struct{ Name string }{})
			require.NoError(t, err)
		})
		t.Run("Err", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			mock.ExpectBegin()
			mock.ExpectRollback()
			defer func() {
				err := mock.ExpectationsWereMet()
				require.NoError(t, err)
			}()

			srv := testserver.NewError()
			srv.AddTransport(transport.POST{})
			srv.Use(gqlutil.DBInjector{DB: db})

			c := client.New(srv)
			err = c.Post(`mutation { name }`, &struct{ Name string }{})
			require.Error(t, err)
		})
		t.Run("Panic", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			mock.ExpectBegin()
			mock.ExpectRollback()
			defer func() {
				err := mock.ExpectationsWereMet()
				require.NoError(t, err)
			}()

			srv := newServer(db)
			srv.SetRecoverFunc(func(_ context.Context, err interface{}) error {
				return err.(error)
			})
			srv.AroundResponses(func(ctx context.Context, _ graphql.ResponseHandler) *graphql.Response {
				panic(graphql.ErrorOnPath(ctx, io.ErrUnexpectedEOF))
			})

			c := client.New(srv)
			err = c.Post(`mutation { name }`, &struct{ Name string }{})
			require.Error(t, err)
			require.Contains(t, err.Error(), io.ErrUnexpectedEOF.Error())
		})
		t.Run("NoTx", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
			defer func() {
				err := mock.ExpectationsWereMet()
				require.NoError(t, err)
			}()

			srv := newServer(db)
			c := client.New(srv)
			err = c.Post(`mutation { name }`, &struct{ Name string }{})
			require.Error(t, err)
			require.Contains(t, err.Error(), sql.ErrConnDone.Error())
		})
	})
	t.Run("NoDB", func(t *testing.T) {
		srv := testserver.New()
		require.PanicsWithError(t, "DB is nil", func() {
			srv.Use(gqlutil.DBInjector{})
		})
	})
}

func TestDBImplementors(t *testing.T) {
	require.Implements(t, (*gqlutil.BeginTxExecQueryer)(nil), &sql.DB{})
	require.Implements(t, (*gqlutil.BeginTxExecQueryer)(nil), &sql.Conn{})
	require.Implements(t, (*gqlutil.ExecQueryer)(nil), &sql.Tx{})
}
