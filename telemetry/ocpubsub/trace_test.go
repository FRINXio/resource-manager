// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocpubsub_test

import (
	"context"
	"github.com/net-auto/resourceManager/telemetry/ocpubsub"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/mempubsub"
)

type testExporter struct {
	spans []*trace.SpanData
}

func (t *testExporter) ExportSpan(s *trace.SpanData) {
	t.spans = append(t.spans, s)
}

func TestTraces(t *testing.T) {
	te := &testExporter{}
	trace.RegisterExporter(te)
	defer trace.UnregisterExporter(te)

	pstopic := mempubsub.NewTopic()
	topic := ocpubsub.TraceTopic{Topic: pstopic}
	subscription := mempubsub.NewSubscription(pstopic, time.Second)

	ctx, span := trace.StartSpan(context.Background(),
		t.Name(), trace.WithSampler(trace.AlwaysSample()),
	)
	body := []byte("foobar")
	err := topic.Send(ctx, &pubsub.Message{Body: body})
	require.NoError(t, err)
	span.End()

	msg, err := subscription.Receive(ctx)
	require.NoError(t, err)
	require.Equal(t, []byte("foobar"), msg.Body)
	require.Len(t, te.spans, 3)
	msg.Ack()

	sc, ok := ocpubsub.SpanContextFromMessage(msg)
	require.True(t, ok)
	require.Equal(t, span.SpanContext(), sc)
}
