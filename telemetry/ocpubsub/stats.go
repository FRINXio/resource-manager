// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocpubsub

import (
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

// The following measures are supported for use in custom views.
var (
	MessagesSentTotal = stats.Int64(
		"pubsub/messages_sent_total",
		"Number of messages sent",
		stats.UnitDimensionless,
	)
	MessagesErrorTotal = stats.Int64(
		"pubsub/messages_errors_total",
		"Number of send errors",
		stats.UnitDimensionless,
	)
	MessagesSentLatency = stats.Float64(
		"pubsub/messages_sent_latency",
		"Latency of messages sent",
		stats.UnitMilliseconds,
	)
	MessagesSentBytes = stats.Int64(
		"pubsub/message_sent_bytes_total",
		"Number of bytes sent",
		stats.UnitBytes,
	)
	MessagesReceivedTotal = stats.Int64(
		"pubsub/messages_received_total",
		"Number of messages received",
		stats.UnitDimensionless,
	)
	MessagesReceivedBytes = stats.Int64(
		"pubsub/message_received_bytes_total",
		"Number of bytes received",
		stats.UnitBytes,
	)
	MessagesReceivedErrorTotal = stats.Int64(
		"pubsub/messages_received_errors_total",
		"Number of receive errors",
		stats.UnitDimensionless,
	)
)

// Package ocpubsub provides some convenience views for server measures.
// You still need to register these views for data to actually be collected.
var (
	MessagesSentTotalView = &view.View{
		Name:        "pubsub/messages_sent_total",
		Description: "Number of messages sent",
		Measure:     MessagesSentTotal,
		Aggregation: view.Count(),
	}
	MessagesErrorTotalView = &view.View{
		Name:        "pubsub/messages_errors_total",
		Description: "Number of send errors",
		Measure:     MessagesErrorTotal,
		Aggregation: view.Count(),
	}
	MessagesSentLatencyView = &view.View{
		Name:        "pubsub/messages_sent_latency_milliseconds",
		Description: "Latency distribution of pubsub messages",
		Measure:     MessagesSentLatency,
		Aggregation: ocgrpc.DefaultMillisecondsDistribution,
	}
	MessagesSentBytesView = &view.View{
		Name:        "pubsub/messages_sent_bytes_total",
		Description: "Number of bytes sent",
		Measure:     MessagesSentBytes,
		Aggregation: ocgrpc.DefaultBytesDistribution,
	}
	MessagesReceivedTotalView = &view.View{
		Name:        "pubsub/messages_received_total",
		Description: "Number of messages received",
		Measure:     MessagesReceivedTotal,
		Aggregation: view.Count(),
	}
	MessagesReceivedBytesView = &view.View{
		Name:        "pubsub/messages_received_bytes_total",
		Description: "Number of bytes received",
		Measure:     MessagesReceivedBytes,
		Aggregation: ocgrpc.DefaultBytesDistribution,
	}
	MessagesReceivedErrorTotalView = &view.View{
		Name:        "pubsub/messages_received_errors_total",
		Description: "Number of receive errors",
		Measure:     MessagesReceivedErrorTotal,
		Aggregation: view.Count(),
	}
)

// DefaultViews are the default views provided by this package.
var DefaultViews = []*view.View{
	MessagesSentTotalView,
	MessagesErrorTotalView,
	MessagesSentLatencyView,
	MessagesSentBytesView,
	MessagesReceivedTotalView,
	MessagesReceivedBytesView,
	MessagesReceivedErrorTotalView,
}
