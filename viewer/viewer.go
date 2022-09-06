// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copied from inventory and modified to use our ent.Client and our tenancy.go

package viewer

import (
	"context"
	"encoding/json"
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/schema"
	logger "github.com/net-auto/resourceManager/logging"
	log2 "github.com/net-auto/resourceManager/logging/log"
	"go.uber.org/atomic"
	"go.uber.org/zap/zapcore"
	"sort"

	"net/http"
	"net/url"
	"strings"

	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

const (
	// TenantHeader is the http tenant header.
	TenantHeader = "x-tenant-id"
	UserHeader   = "from"
	RoleHeader   = "x-auth-user-roles"
	GroupHeader  = "x-auth-user-groups"
)

// Attributes recorded on the span of the requests.
const (
	TenantAttribute    = "viewer.tenant"
	UserAttribute      = "viewer.user"
	RoleAttribute      = "viewer.role"
	UserAgentAttribute = "viewer.user_agent"
)

// The following tags are applied to context recorded by this package.
var (
	KeyTenant    = tag.MustNewKey(TenantAttribute)
	KeyUser      = tag.MustNewKey(UserAttribute)
	KeyRole      = tag.MustNewKey(RoleAttribute)
	KeyUserAgent = tag.MustNewKey(UserAgentAttribute)
)

// Config configures the viewer package.
type Config struct {
	TenantMaxConn int      `name:"tenancy.db_max_conn" env:"TENANCY_DB_MAX_CONN" default:"1" help:"Max connections to database per tenant."`
	FeaturesURL   *url.URL `name:"features.url" env:"FEATURES_URL" placeholder:"URL" help:"URL to fetch tenant features."`
}

// TenancyHandler adds viewer / tenancy into incoming requests.
func TenancyHandler(h http.Handler, tenancy Tenancy, symphLogger log2.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenant := r.Header.Get(TenantHeader)
		if tenant == "" {
			logger.Warn(r.Context(), "request missing tenant header")
			http.Error(w, "missing tenant header", http.StatusBadRequest)
			return
		}

		user := r.Header.Get(UserHeader)
		if user == "" {
			logger.Warn(r.Context(), "request missing user header")
			http.Error(w, "missing user header", http.StatusBadRequest)
			return
		}

		roles := r.Header.Get(RoleHeader)
		groups := r.Header.Get(GroupHeader)

		logger.Debug(r.Context(), "getting tenancy client for %+v", zap.String("tenant", tenant))

		client, err := tenancy.ClientFor(r.Context(), tenant, symphLogger.For(r.Context()))
		if err != nil {
			const msg = "cannot get tenancy client"
			logger.Warn(r.Context(), msg+"error: %+v", zap.Error(err))
			http.Error(w, msg, http.StatusServiceUnavailable)
			return
		}

		var ctx = ent.NewContext(r.Context(), client)

		ctx = schema.WithIdentity(ctx, tenant, user, roles, groups)
		identity, _ := schema.GetIdentity(ctx)
		marshal, err := json.Marshal(identity)
		if err != nil {
			logger.Warn(r.Context(), "Cannot marshall identity %+v", zap.Error(err))
			http.Error(w, "marshalling identity", http.StatusServiceUnavailable)
			return
		}
		ctx = log2.NewFieldsContext(ctx, zap.String("identity", string(marshal)))
		trace.FromContext(ctx).AddAttributes(traceAttrs(identity)...)
		ctx, _ = tag.New(ctx, tags(r, identity)...)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

var (
	GroupAttribute = "viewer.group"
	KeyGroup       = tag.MustNewKey(GroupAttribute)
)

func traceAttrs(v *schema.Identity) []trace.Attribute {
	return []trace.Attribute{
		trace.StringAttribute(TenantAttribute, v.Tenant),
		trace.StringAttribute(UserAttribute, v.User),
		trace.StringAttribute(RoleAttribute, strings.Join(v.Roles, ",")),
		trace.StringAttribute(GroupAttribute, strings.Join(v.Groups, ",")),
	}
}

func tags(r *http.Request, v *schema.Identity) []tag.Mutator {
	var userAgent string
	if parts := strings.SplitN(r.UserAgent(), " ", 2); len(parts) > 0 {
		userAgent = parts[0]
	}
	return []tag.Mutator{
		tag.Upsert(KeyTenant, v.Tenant),
		tag.Upsert(KeyUser, v.User),
		tag.Upsert(KeyRole, strings.Join(v.Roles, ",")),
		tag.Upsert(KeyGroup, strings.Join(v.Groups, ",")),
		tag.Upsert(KeyUserAgent, userAgent),
	}
}

// Option enables viewer customization.
type Option func(Viewer)

type viewer struct {
	tenant   string
	features FeatureSet
}

type User struct {
	Role   string
	AuthID string
}

// Tenant is the tenant of the viewer.
func (v *viewer) Tenant() string {
	return v.tenant
}

// Features is the features applied for the viewer.
func (v *viewer) Features() FeatureSet {
	return v.features
}

// User returns the ent user of the viewer.
func (v *UserViewer) User() *User {
	u, _ := v.user.Load().(*User)
	return u
}

// Name implements Viewer.Name by getting user's Auth ID.
func (v *UserViewer) Name() string {
	return v.User().AuthID
}

// Name implements Viewer.Name by getting user's Role.
func (v *UserViewer) Role() string {
	return v.User().Role
}

// TraceAttrs returns a set of trace attributes of viewer.
func TraceAttrs(v Viewer) []trace.Attribute {
	return []trace.Attribute{
		trace.StringAttribute(TenantAttribute, v.Tenant()),
		trace.StringAttribute(UserAttribute, v.Name()),
		trace.StringAttribute(RoleAttribute, v.Role()),
	}
}

// Viewer is the interface to hold additional per request information.
type Viewer interface {
	zapcore.ObjectMarshaler
	Tenant() string
	Features() FeatureSet
	Name() string
	Role() string
}

// UserViewer is a viewer that holds a user ent.
type UserViewer struct {
	viewer
	user atomic.Value
}

// String returns the textual representation of a feature set.
func (f FeatureSet) String() string {
	features := make([]string, 0, len(f))
	for feature := range f {
		features = append(features, feature)
	}
	sort.Strings(features)
	return strings.Join(features, ",")
}

// Enabled check if feature is in FeatureSet.
func (f FeatureSet) Enabled(feature string) bool {
	_, ok := f[feature]
	return ok
}

type ctxKey struct{}

// FromContext returns the Viewer stored in a context, or nil if there isn't one.
func FromContext(ctx context.Context) Viewer {
	v, _ := ctx.Value(ctxKey{}).(Viewer)
	return v
}

// NewContext returns a new context with the given Viewer attached.
func NewContext(parent context.Context, v Viewer) context.Context {
	return context.WithValue(parent, ctxKey{}, v)
}
