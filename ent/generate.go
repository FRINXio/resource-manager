// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ent

// This little code using symphony is here just to make sure symphony as a dependency gets downloaded before go:generate is triggered
//  go:generate (the last command) uses some files from symphony so the dependency needs to be downloaded before its triggerred
//  simple solution is to use something from symphony before go:generate like this useless function
import "github.com/facebookincubator/symphony/pkg/ctxutil"

func noop() {
	ctxutil.DoneCtx()
}

//go:generate echo ""
//go:generate echo "------> Generating ent.go entities from ent/schema folder"
//go:generate go run github.com/facebook/ent/cmd/entc generate --storage=sql --template ./template --template $GOPATH/pkg/mod/github.com/facebookincubator/ent-contrib@v0.0.0-20201018112627-5709d2185a62/entgql/template --template $GOPATH/pkg/mod/github.com/!f!r!i!n!xio/symphony@v0.0.0-20201029142040-3dd67baee7b6/pkg/ent/template --feature privacy,schema/snapshot --header "// Code generated by entc, DO NOT EDIT." ./schema
