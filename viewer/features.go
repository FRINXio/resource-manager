// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewer

import (
	"net/http"
	"time"

	"go.opencensus.io/plugin/ochttp"
	"gocloud.dev/runtimevar"
	"gocloud.dev/runtimevar/httpvar"
)

// The following strings are features that can be allowed on for specific tenants\users
// To query if feature is allowed use FeatureSet.Enabled
const (
	FeatureMandatoryPropertiesOnWorkOrderClose = "mandatory_properties_on_work_order_close"
	FeatureExecuteAutomationFlows              = "execute_automation_flows"
)

// How frequently the features variable is updated by http calls
const featuresUpdateFreq = 5 * time.Second

// FeatureSet holds the list of features of the viewer
type TenantFeatures map[string][]string

// FeatureSet holds the list of features of the viewer
type FeatureSet map[string]struct{}

// NewFeatureSet create FeatureSet from a list of features.
func NewFeatureSet(features ...string) FeatureSet {
	set := make(FeatureSet, len(features))
	for _, feature := range features {
		set[feature] = struct{}{}
	}
	return set
}

// SyncFeatures syncs feature flags to variable periodically via http
func SyncFeatures(cfg Config) (*runtimevar.Variable, func(), error) {
	decoder := runtimevar.NewDecoder(TenantFeatures{}, runtimevar.JSONDecode)
	v, err := httpvar.OpenVariable(
		&http.Client{Transport: &ochttp.Transport{}},
		cfg.FeaturesURL.String(), decoder,
		&httpvar.Options{WaitDuration: featuresUpdateFreq},
	)
	if err != nil {
		return nil, nil, err
	}
	return v, func() {
		_ = v.Close()
	}, nil
}
