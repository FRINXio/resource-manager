// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocpubsub

import (
	"context"
	"encoding/base64"

	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
	"gocloud.dev/pubsub"
)

const spanMetaKey = "trace-span-context"

// TraceTopic is a topic that propagates span context.
type TraceTopic struct {
	// Topic is the underlying topic that actually sends messages.
	Topic
}

// Send propagates context and delegates the actual send to underlying topic.
func (t TraceTopic) Send(ctx context.Context, msg *pubsub.Message) error {
	b := propagation.Binary(trace.FromContext(ctx).SpanContext())
	if b != nil {
		if msg.Metadata == nil {
			msg.Metadata = map[string]string{
				spanMetaKey: base64.StdEncoding.EncodeToString(b),
			}
		} else if _, ok := msg.Metadata[spanMetaKey]; !ok {
			msg.Metadata[spanMetaKey] = base64.StdEncoding.EncodeToString(b)
		}
	}
	return t.Topic.Send(ctx, msg)
}

// SpanContextFromMessage returns the span context stored in message.
func SpanContextFromMessage(msg *pubsub.Message) (trace.SpanContext, bool) {
	if msg == nil || msg.Metadata == nil {
		return trace.SpanContext{}, false
	}
	s, ok := msg.Metadata[spanMetaKey]
	if !ok {
		return trace.SpanContext{}, false
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return trace.SpanContext{}, false
	}
	return propagation.FromBinary(b)
}
