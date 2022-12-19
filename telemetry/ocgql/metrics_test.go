// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocgql_test

import (
	ocgql2 "github.com/net-auto/resourceManager/telemetry/ocgql"
	"sort"
	"sync"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/testserver"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

func TestMetrics(t *testing.T) {
	err := view.Register(ocgql2.DefaultViews...)
	require.NoError(t, err)
	defer view.Unregister(ocgql2.DefaultViews...)

	h := testserver.New()
	h.AddTransport(transport.POST{})
	h.AddTransport(transport.Websocket{})
	h.Use(ocgql2.Metrics{})
	h.Use(extension.FixedComplexityLimit(100))
	h.SetCalculatedComplexity(50)

	c := client.New(h)
	err = c.Post(`query { name }`, &struct{ Name string }{})
	require.NoError(t, err)

	sk := c.Websocket(`subscription { name }`)
	defer func() {
		err := sk.Close()
		require.NoError(t, err)
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := sk.Next(&struct{ Name string }{})
		require.NoError(t, err)
	}()
	h.SendNextSubscriptionMessage()
	wg.Wait()

	counters := []struct {
		name string
		rows []view.Row
	}{
		{
			name: ocgql2.RequestTotalView.Name,
			rows: []view.Row{
				{
					Tags: []tag.Tag{{Key: ocgql2.Operation, Value: "query"}},
					Data: &view.CountData{Value: 1},
				},
				{
					Tags: []tag.Tag{{Key: ocgql2.Operation, Value: "subscription"}},
					Data: &view.CountData{Value: 1},
				},
			},
		},
		{
			name: ocgql2.ResponseTotalView.Name,
			rows: []view.Row{
				{
					Tags: []tag.Tag{
						{Key: ocgql2.Operation, Value: "query"},
						{Key: ocgql2.Errors, Value: "0"},
					},
					Data: &view.CountData{Value: 1},
				},
				{
					Tags: []tag.Tag{
						{Key: ocgql2.Operation, Value: "subscription"},
						{Key: ocgql2.Errors, Value: "0"},
					},
					Data: &view.CountData{Value: 1},
				},
			},
		},
		{
			name: ocgql2.ResolveTotalView.Name,
			rows: []view.Row{
				{
					Tags: []tag.Tag{
						{Key: ocgql2.Object, Value: "Query"},
						{Key: ocgql2.Field, Value: "name"},
						{Key: ocgql2.Errors, Value: "0"},
					},
					Data: &view.CountData{Value: 1},
				},
			},
		},
		{
			name: ocgql2.DeprecatedResolveTotalView.Name,
			rows: []view.Row{},
		},
	}
	for _, v := range counters {
		rows, err := view.RetrieveData(v.name)
		require.NoError(t, err)
		sort.Slice(rows, func(i, j int) bool {
			var leftOp, rightOp string
			for _, t := range rows[i].Tags {
				if t.Key == ocgql2.Operation {
					leftOp = t.Value
					break
				}
			}
			for _, t := range rows[j].Tags {
				if t.Key == ocgql2.Operation {
					rightOp = t.Value
					break
				}
			}
			return leftOp <= rightOp
		})
		for i := range rows {
			view.ClearStart(rows[i].Data)
			require.Equal(t, v.rows[i].Data, rows[i].Data)
			require.ElementsMatch(t, v.rows[i].Tags, rows[i].Tags)
		}
	}

	distributions := []struct {
		name string
	}{
		{name: ocgql2.RequestLatencyView.Name},
		{name: ocgql2.ResolveLatencyView.Name},
		{name: ocgql2.RequestComplexityView.Name},
	}
	for _, v := range distributions {
		rows, err := view.RetrieveData(v.name)
		require.NoError(t, err)
		require.Len(t, rows, 1)
		data, ok := rows[0].Data.(*view.DistributionData)
		require.True(t, ok)
		require.Greater(t, data.Sum(), float64(0))
	}

	rows, err := view.RetrieveData(ocgql2.NumSubscriptionsView.Name)
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.Equal(t, &view.LastValueData{Value: 1}, rows[0].Data)
}
