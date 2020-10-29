// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copied from inventory and modified to use our ent.Client and our tenancy.go

package viewer

import (
	logger "github.com/net-auto/resourceManager/logging"
	"encoding/json"
	"github.com/facebookincubator/symphony/pkg/log"
	fb_viewer "github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/schema"
	"net/http"
	"strings"

	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

const (
	// TenantHeader is the http tenant header.
	TenantHeader = "x-tenant-id"
	UserHeader = "from"
	RoleHeader = "x-auth-user-roles"
	GroupHeader = "x-auth-user-groups"
)

// TenancyHandler adds viewer / tenancy into incoming requests.
func TenancyHandler(h http.Handler, tenancy Tenancy, _ log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenant := r.Header.Get(TenantHeader)
		if tenant == "" {
			logger.Warn(r.Context(),"request missing tenant header")
			http.Error(w, "missing tenant header", http.StatusBadRequest)
			return
		}

		user := r.Header.Get(UserHeader)
		if user == "" {
			logger.Warn(r.Context(),"request missing user header")
			http.Error(w, "missing user header", http.StatusBadRequest)
			return
		}

		roles := r.Header.Get(RoleHeader)
		groups := r.Header.Get(GroupHeader)

		logger.Debug(r.Context(),"getting tenancy client for %+v", zap.String("tenant", tenant))

		client, err := tenancy.ClientFor(r.Context(), tenant, logger)
		if err != nil {
			const msg = "cannot get tenancy client"
			logger.Warn(r.Context(), msg + "error: %+v", zap.Error(err))
			http.Error(w, msg, http.StatusServiceUnavailable)
			return
		}

		var ctx  = ent.NewContext(r.Context(), client)

		ctx = schema.WithIdentity(ctx, tenant, user, roles, groups)
		identity, _ := schema.GetIdentity(ctx)
		marshal, err := json.Marshal(identity)
		if err != nil {
			logger.Warn(r.Context(),"Cannot marshall identity %+v", zap.Error(err))
			http.Error(w, "marshalling identity", http.StatusServiceUnavailable)
			return
		}
		ctx = log.NewFieldsContext(ctx, zap.String("identity", string(marshal)))
		trace.FromContext(ctx).AddAttributes(traceAttrs(identity)...)
		ctx, _ = tag.New(ctx, tags(r, identity)...)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

var (
	GroupAttribute = "viewer.group"
	KeyGroup    = tag.MustNewKey(GroupAttribute)
)

func traceAttrs(v *schema.Identity) []trace.Attribute {
	return []trace.Attribute{
		trace.StringAttribute(fb_viewer.TenantAttribute, v.Tenant),
		trace.StringAttribute(fb_viewer.UserAttribute, v.User),
		trace.StringAttribute(fb_viewer.RoleAttribute, strings.Join(v.Roles, ",")),
		trace.StringAttribute(GroupAttribute, strings.Join(v.Groups, ",")),
	}
}

func tags(r *http.Request, v *schema.Identity) []tag.Mutator {
	var userAgent string
	if parts := strings.SplitN(r.UserAgent(), " ", 2); len(parts) > 0 {
		userAgent = parts[0]
	}
	return []tag.Mutator{
		tag.Upsert(fb_viewer.KeyTenant, v.Tenant),
		tag.Upsert(fb_viewer.KeyUser, v.User),
		tag.Upsert(fb_viewer.KeyRole, strings.Join(v.Roles, ",")),
		tag.Upsert(KeyGroup, strings.Join(v.Groups, ",")),
		tag.Upsert(fb_viewer.KeyUserAgent, userAgent),
	}
}
