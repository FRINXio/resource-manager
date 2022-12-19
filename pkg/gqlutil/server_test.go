// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gqlutil_test

import (
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/net-auto/resourceManager/pkg/gqlutil"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	srv := gqlutil.NewServer(&graphql.ExecutableSchemaMock{})
	require.NotNil(t, srv)
}
