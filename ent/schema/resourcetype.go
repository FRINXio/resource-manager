// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
	"time"

	// "github.com/facebookincubator/symphony/pkg/authz"
	// "github.com/net-auto/resourceManager/ent/privacy"
	// "github.com/net-auto/resourceManager/ent/privacy"
	// "github.com/facebookincubator/symphony/pkg/viewer"
)

// ResourceType holds the schema definition for the ResourceType entity.
type ResourceType struct {
	ent.Schema
}

// Fields of the ResourceType.
func (ResourceType) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique(),
	}
}

// Edges of the ResourceType.
func (ResourceType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("property_types", PropertyType.Type),
		edge.To("pools", ResourcePool.Type),
	}
}

// Policy returns resource type policy.
func (ResourceType) Policy() ent.Policy {
	// TODO setup RBAC policies for entities (RBAC based on user's role) such as:
	// return authz.NewPolicy(
	// 	authz.WithMutationRules(
	// 		authz.ResourceTypeWritePolicyRule(),
	// 	),
	// )
	return nil
}

type Tag struct {
	ent.Schema
}

func (Tag) Fields() []ent.Field {
	return []ent.Field{
		field.String("tag").
			NotEmpty().
			Unique(),
	}
}

func (Tag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("pools", ResourcePool.Type),
	}
}

type AllocationStrategy struct {
	ent.Schema
}

// Fields of the ResourcePool.
func (AllocationStrategy) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique(),
		field.Enum("lang").
			Values("py", "js").
			Default("js"),
		field.String("script").
			NotEmpty(),
	}
}

func (AllocationStrategy) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("pools", ResourcePool.Type).
			Ref("allocation_strategy"),
	}
}

// ResourcePool holds the schema definition for the Resource pool entity.
type ResourcePool struct {
	ent.Schema
}

const ResourcePoolDealocationRetire = -1
const ResourcePoolDealocationImmediatelly = 0

// Fields of the ResourcePool.
func (ResourcePool) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique(),
		field.Enum("pool_type").
			Values("singleton", "set", "allocating"),
		field.Int("dealocation_safety_period").
			Default(0).
			Comment("How long to keep resources unavailable after dealocation (in seconds)." +
				" -1 release never, 0 release immediately"),
	}
}

// Edges of the ResourcePool.
func (ResourcePool) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("resource_type", ResourceType.Type).
			Ref("pools").
			Unique(),
		edge.From("tags", Tag.Type).
			Ref("pools"),
		edge.To("claims", Resource.Type),
		edge.To("allocation_strategy", AllocationStrategy.Type).
			Unique(),
	}
}

// Resource holds the schema definition for the Resource entity.
type Resource struct {
	ent.Schema
}

// Fields of the Resource.
func (Resource) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").
			Values("free", "claimed", "retired", "bench"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Resource.
func (Resource) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("pool", ResourcePool.Type).
			Ref("claims").
			Unique(),
		edge.To("properties", Property.Type),
	}
}
