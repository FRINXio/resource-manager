// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copied from inventory and modified to use our ent.Client and our tenancy.go

package viewer

import (
	fb_ent "github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/user"
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
func TenancyHandler(h http.Handler, tenancy Tenancy, logger log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logger.For(r.Context())

		tenant := r.Header.Get(TenantHeader)
		if tenant == "" {
			logger.Warn("request missing tenant header")
			http.Error(w, "missing tenant header", http.StatusBadRequest)
			return
		}

		user := r.Header.Get(UserHeader)
		if user == "" {
			logger.Warn("request missing user header")
			http.Error(w, "missing user header", http.StatusBadRequest)
			return
		}

		roles := r.Header.Get(RoleHeader)
		groups := r.Header.Get(GroupHeader)

		logger.Info("getting tenancy client for %s", zap.String("tenant", tenant));

		client, err := tenancy.ClientFor(r.Context(), tenant)
		if err != nil {
			const msg = "cannot get tenancy client"
			logger.Warn(msg, zap.Error(err))
			http.Error(w, msg, http.StatusServiceUnavailable)
			return
		}

		var (
			ctx  = ent.NewContext(r.Context(), client)
			opts = make([]fb_viewer.Option, 0, 1)
			v    fb_viewer.Viewer
		)

		if username := r.Header.Get(UserHeader); username != "" {
			u, err := GetOrCreateUser(
				username,
				"OWNER")
			if err != nil {
				logger.Warn("cannot get user ent", zap.Error(err))
				http.Error(w, "getting user entity", http.StatusServiceUnavailable)
				return
			}
			ctx = schema.WithIdentity(ctx, tenant, username, roles, groups)
			v = fb_viewer.NewUser(tenant, u, opts...)
		} else {
			logger.Warn("request missing identity header")
			http.Error(w, "missing identity header", http.StatusBadRequest)
			return
		}

		ctx = log.NewFieldsContext(ctx, zap.Object("viewer", v))
		trace.FromContext(ctx).AddAttributes(traceAttrs(v)...)
		ctx, _ = tag.New(ctx, tags(r, v)...)
		ctx = fb_viewer.NewContext(ctx, v)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetOrCreateUser creates or returns existing user with given authID and role.
// we do not manage users locally, so just return the inputs
func GetOrCreateUser(authID string, role user.Role) (*fb_ent.User, error) {
	return &fb_ent.User{AuthID: authID, Role: role}, nil
}

func traceAttrs(v fb_viewer.Viewer) []trace.Attribute {
	return []trace.Attribute{
		trace.StringAttribute(fb_viewer.TenantAttribute, v.Tenant()),
		trace.StringAttribute(fb_viewer.UserAttribute, v.Name()),
		trace.StringAttribute(fb_viewer.RoleAttribute, v.Role().String()),
	}
}

func tags(r *http.Request, v fb_viewer.Viewer) []tag.Mutator {
	var userAgent string
	if parts := strings.SplitN(r.UserAgent(), " ", 2); len(parts) > 0 {
		userAgent = parts[0]
	}
	return []tag.Mutator{
		tag.Upsert(fb_viewer.KeyTenant, v.Tenant()),
		tag.Upsert(fb_viewer.KeyUser, v.Name()),
		tag.Upsert(fb_viewer.KeyRole, v.Role().String()),
		tag.Upsert(fb_viewer.KeyUserAgent, userAgent),
	}
}
