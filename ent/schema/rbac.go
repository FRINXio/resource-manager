package schema

import (
	"context"
	"fmt"
	"github.com/facebook/ent"
	"strings"
)

type Identity struct {
	Tenant string
	User   string
	Roles  []string
	Groups []string
}

var fullAccessIdentity = Identity{}

var identityKey string;

var adminRoles []string
var adminGroups []string

type RbacConfig struct {
	AdminRoles []string `name:"rbac.admin-roles" sep:"," help:"List of user roles granting admin access"`;
	AdminGroups []string `name:"rbac.admin-groups" sep:"," help:"List of user groups granting admin access"`;
}

type generalRBACPolicy struct {
	mutationsAllowed bool
}

// RBAC grants all access to queries, but only admins can invoke mutations
// admins are identified from roles/groups configured via InitializeAdminRoles, InitializeAdminGroups
var RBAC = generalRBACPolicy{false}
// ALWAYS_ALLOWED grants all access to queries and mutations
var ALWAYS_ALLOWED = generalRBACPolicy{true}

// WithIdentity pushes user identity into context
func WithIdentity(ctx context.Context, tenant string, user string, roles string, groups string) context.Context {
	return WithIdentityParsed(ctx, tenant, user, splitRoles(roles), splitRoles(groups))
}

// WithIdentity pushes user identity into context
func WithIdentityParsed(ctx context.Context, tenant string, user string, roles []string, groups []string) context.Context {
	identity := &Identity{tenant, user, roles, groups}
	return context.WithValue(ctx, identityKey, identity)
}

// WithIdentity pushes user identity into context
func WithFullAccessIdentity(ctx context.Context) context.Context {
	return context.WithValue(ctx, identityKey, &fullAccessIdentity)
}

func splitRoles(roles string) []string {
	split := strings.Split(roles, ",")
	trimmed := make([]string, len(split))
	for i, role := range split {
		trimmed[i] = strings.TrimSpace(role)
	}
	return trimmed
}

// GetIdentity loads user identity from context
func GetIdentity(ctx context.Context) (*Identity, error) {
	identity, err := ctx.Value(identityKey).(*Identity)
	if !err {
		return nil, fmt.Errorf("Unable to get identity from context")
	}
	if identity == nil {
		return nil, fmt.Errorf("Identity information not provided.")
	}
	return identity, nil
}

// InitializeAdminRoles globally configures a list of user roles granting admin access
func InitializeAdminRoles(rolesAsString string) {
	InitializeAdminRolesFromSlice(splitRoles(rolesAsString))
}

// InitializeAdminRoles globally configures a list of user roles granting admin access
func InitializeAdminRolesFromSlice(roles []string) {
	adminRoles = roles
}

// InitializeAdminRoles globally configures a list of user groups granting admin access
func InitializeAdminGroups(groupsAsString string) {
	InitializeAdminGroupsFromSlice(splitRoles(groupsAsString))
}

// InitializeAdminRoles globally configures a list of user groups granting admin access
func InitializeAdminGroupsFromSlice(groups []string) {
	adminGroups = groups
}

// EvalQuery always grants access
func (generalRBACPolicy) EvalQuery(ctx context.Context, m ent.Query) error {
	// All queries allowed
	return nil
}

// EvalMutation grants access to user possessing admin roles or groups
func (generalRBACPolicy) EvalMutation(ctx context.Context, m ent.Mutation) error {
	// Mutations allowed only for admins
	identity, err := GetIdentity(ctx)
	if err != nil {
		return fmt.Errorf("Identity context not found. User unauthorized to mutate: %q", m.Type())
	}
	if isAllowed(identity) {
		return nil
	} else {
		return fmt.Errorf("User: %q unauthorized to mutate: %q. Must be role: %s or group: %s",
			identity, m.Type(), adminRoles, adminGroups)
	}
}

func isAllowed(identity *Identity) bool {
	if (identity == &fullAccessIdentity) {
		// Full access
		return true
	}

	if len(adminRoles) == 0 && len(adminGroups) == 0 {
		// TODO log warning about incorrect RBAC setting !
		// return true if RBAC required roles/groups not set
		return true
	}

	if len(adminRoles) != 0 {
		for _, role := range identity.Roles {
			for _, roleRequired := range adminRoles {
				if strings.ToUpper(role) == strings.ToUpper(roleRequired) {
					return true
				}
			}
		}
	}

	if len(adminGroups) != 0 {
		for _, group := range identity.Groups {
			for _, groupRequired := range adminGroups {
				if strings.ToUpper(group) == strings.ToUpper(groupRequired) {
					return true
				}
			}
		}
	}

	return false
}
