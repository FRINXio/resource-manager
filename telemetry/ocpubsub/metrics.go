// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocpubsub

import (
	"context"
	"time"

	"go.opencensus.io/stats"
	"gocloud.dev/pubsub"
)

// MetricsTopic is an topic that collects metrics for the outgoing messages.
type MetricsTopic struct {
	// Topic to collect metrics for.
	Topic
}

// Send delegates the actual send to underlying topic and
// record metrics for the outgoing message.
func (t MetricsTopic) Send(ctx context.Context, msg *pubsub.Message) error {
	stats.Record(ctx, MessagesSentTotal.M(1))

	start := time.Now()
	err := t.Topic.Send(ctx, msg)
	elapsed := time.Since(start)

	latency := float64(elapsed) / float64(time.Millisecond)
	measurements := make([]stats.Measurement, 0, 2)
	measurements = append(measurements,
		MessagesSentLatency.M(latency),
	)
	if err != nil {
		measurements = append(measurements,
			MessagesErrorTotal.M(1),
		)
	} else {
		measurements = append(measurements,
			MessagesSentBytes.M(msgLen(msg)),
		)
	}
	stats.Record(ctx, measurements...)

	return err
}

// MetricsSubscription is a subscription that collects metrics for the incoming messages.
type MetricsSubscription struct {
	// Subscription to collects metrics for.
	Subscription
}

// Receive delegates the actual receives to underlying subscription and
// record metrics for the incoming message.
func (s MetricsSubscription) Receive(ctx context.Context) (msg *pubsub.Message, err error) {
	defer func() {
		if err == nil {
			stats.Record(ctx,
				MessagesReceivedTotal.M(1),
				MessagesReceivedBytes.M(msgLen(msg)),
			)
		} else {
			stats.Record(ctx,
				MessagesReceivedErrorTotal.M(1),
			)
		}
	}()

	return s.Subscription.Receive(ctx)
}

func msgLen(msg *pubsub.Message) int64 {
	if msg == nil {
		return 0
	}
	length := len(msg.Body)
	for k, v := range msg.Metadata {
		length += len(k) + len(v)
	}
	return int64(length)
}
