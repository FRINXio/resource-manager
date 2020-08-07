// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/facebookincubator/symphony/graph/graphql/plugin/txgen"
)

func main() {
	verbose := flag.Bool("v", false, "show logs")
	flag.Parse()

	var (
		output bytes.Buffer
		err    error
	)
	defer func() {
		if err != nil {
			io.Copy(os.Stderr, &output)
			os.Exit(1)
		}
	}()
	if !*verbose {
		log.SetOutput(&output)
	}

	var cfg *config.Config
	if cfg, err = config.LoadConfigFromDefaultLocations(); err != nil {
		log.Println("cannot load config file", err)
		return
	}
	if err = api.Generate(cfg, api.AddPlugin(
		txgen.New(config.PackageConfig{}),
	)); err != nil {
		log.Println("cannot generate code", err)
		return
	}
}
