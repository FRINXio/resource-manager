// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocpubsub

import (
	"context"

	"gocloud.dev/pubsub"
)

type (
	// Topic is the interface extracted from gocloud.dev/pubsub.Topic.
	Topic interface {
		As(interface{}) bool
		ErrorAs(error, interface{}) bool
		Send(context.Context, *pubsub.Message) error
		Shutdown(context.Context) error
	}

	// Subscription is the interface extracted from gocloud.dev/pubsub.Subscription.
	Subscription interface {
		As(interface{}) bool
		ErrorAs(error, interface{}) bool
		Receive(context.Context) (*pubsub.Message, error)
		Shutdown(context.Context) error
	}
)

// WrapTopic wraps a given topic with a topic that
// collects telemetry for the outgoing messages.
func WrapTopic(topic Topic) Topic {
	if topic == nil {
		panic("topic is nil")
	}
	return MetricsTopic{
		TraceTopic{topic},
	}
}

// WrapSubscription wraps a given subscription with a subscription that
// collects telemetry for the incoming messages.
func WrapSubscription(sub Subscription) Subscription {
	if sub == nil {
		panic("subscription is nil")
	}
	return MetricsSubscription{sub}
}
