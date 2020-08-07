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
	"net/http"
	"strings"

	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

// TenancyHandler adds viewer / tenancy into incoming requests.
func TenancyHandler(h http.Handler, tenancy Tenancy, logger log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logger.For(r.Context())
		tenant := r.Header.Get(fb_viewer.TenantHeader)
		if tenant == "" {
			logger.Warn("request missing tenant header")
			http.Error(w, "missing tenant header", http.StatusBadRequest)
			return
		}
		logger = logger.With(zap.String("tenant", tenant))

		role := user.Role(r.Header.Get(fb_viewer.RoleHeader))
		if err := user.RoleValidator(role); err != nil {
			logger.Warn("request contains invalid role",
				zap.Stringer("role", role),
				zap.Error(err),
			)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

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
		if features := r.Header.Get(fb_viewer.FeaturesHeader); features != "" {
			opts = append(opts, fb_viewer.WithFeatures(strings.Split(features, ",")...))
		}
		if username := r.Header.Get(fb_viewer.UserHeader); username != "" {
			u, err := GetOrCreateUser(
				username,
				role)
			if err != nil {
				logger.Warn("cannot get user ent", zap.Error(err))
				http.Error(w, "getting user entity", http.StatusServiceUnavailable)
				return
			}
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
