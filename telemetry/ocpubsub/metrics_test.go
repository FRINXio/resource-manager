// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocpubsub_test

import (
	"context"
	ocpubsub2 "github.com/net-auto/resourceManager/telemetry/ocpubsub"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opencensus.io/stats/view"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/mempubsub"
)

func TestMetrics(t *testing.T) {
	views := ocpubsub2.DefaultViews
	err := view.Register(views...)
	require.NoError(t, err)
	defer view.Unregister(views...)

	pstopic := mempubsub.NewTopic()
	topic := ocpubsub2.MetricsTopic{Topic: pstopic}
	subscription := ocpubsub2.MetricsSubscription{
		Subscription: mempubsub.NewSubscription(
			pstopic, time.Second,
		),
	}

	ctx := context.Background()
	err = topic.Send(ctx, &pubsub.Message{Body: []byte("foobar")})
	require.NoError(t, err)
	msg, err := subscription.Receive(ctx)
	require.NoError(t, err)
	require.Equal(t, []byte("foobar"), msg.Body)
	msg.Ack()

	for _, v := range views {
		v := view.Find(v.Name)
		require.NotNil(t, v)

		rows, err := view.RetrieveData(v.Name)
		require.NoError(t, err)
		if strings.Contains(v.Name, "_errors_") {
			continue
		}
		require.NotEmpty(t, rows, "no data on view %q", v.Name)

		switch data := rows[0].Data.(type) {
		case *view.CountData:
			require.Greater(t, data.Value, int64(0))
		case *view.DistributionData:
			require.Greater(t, data.Sum(), float64(0))
		default:
			require.Failf(t, "unknown data type", "type=%T", data)
		}
	}
}
