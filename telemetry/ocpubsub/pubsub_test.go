// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocpubsub_test

import (
	"github.com/net-auto/resourceManager/telemetry/ocpubsub"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/mempubsub"
)

func TestCDKInterfaces(t *testing.T) {
	tests := []struct {
		name          string
		interfaceType interface{}
		structType    interface{}
	}{
		{
			name:          "Topic",
			interfaceType: (*ocpubsub.Topic)(nil),
			structType:    (*pubsub.Topic)(nil),
		},
		{
			name:          "Subscription",
			interfaceType: (*ocpubsub.Subscription)(nil),
			structType:    (*pubsub.Subscription)(nil),
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Implements(t, tc.interfaceType, tc.structType)
			iface := reflect.TypeOf(tc.interfaceType).Elem()
			concrete := reflect.TypeOf(tc.structType)
			numMethod := iface.NumMethod()
			require.Equal(t, numMethod, concrete.NumMethod())
			for i := 0; i < numMethod; i++ {
				_, ok := concrete.MethodByName(iface.Method(i).Name)
				require.True(t, ok)
			}
		})
	}
}

func TestWrapper(t *testing.T) {
	pstopic := mempubsub.NewTopic()
	t.Run("Topic", func(t *testing.T) {
		topic := ocpubsub.WrapTopic(pstopic)
		require.NotNil(t, topic)
		require.Panics(t, func() { ocpubsub.WrapTopic(nil) })
	})
	t.Run("Subscription", func(t *testing.T) {
		subscription := mempubsub.NewSubscription(pstopic, time.Second)
		require.NotNil(t, subscription)
		require.Panics(t, func() { ocpubsub.WrapSubscription(nil) })
	})
}
