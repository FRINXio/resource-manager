// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/symphony/pkg/ent-contrib/entgql"
	"time"

	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
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
			NotEmpty(),
	}
}

// Edges of the ResourceType.
func (ResourceType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("property_types", PropertyType.Type).
			Annotations(entgql.Bind()),
		edge.To("pools", ResourcePool.Type).
			Annotations(entgql.Bind()),
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
		edge.To("pools", ResourcePool.Type).
			Annotations(entgql.Bind()),
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
		field.Text("description").
			Optional().
			Nillable(),
		field.Enum("lang").
			Values("py", "js").
			Default("js"),
		field.Text("script").
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
const ResourcePoolDealocationImmediately = 0

// Fields of the ResourcePool.
func (ResourcePool) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique(),
		field.Text("description").
			Optional().
			Nillable(),
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
		edge.To("claims", Resource.Type).
			Annotations(entgql.Bind()),
		edge.To("poolProperties", PoolProperties.Type).
			Unique().
			Annotations(entgql.Bind()),
		edge.To("allocation_strategy", AllocationStrategy.Type).
			Unique().
			Annotations(entgql.Bind()),
		edge.From("parent_resource", Resource.Type).
			Ref("nested_pool").
			Comment("pool hierarchies can use this link between resoruce and pool").
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
		edge.To("properties", Property.Type).
			Annotations(entgql.Bind()),
		edge.To("nested_pool", ResourcePool.Type).
			Comment("pool hierarchies can use this link between resoruce and pool").
			Unique().
			Annotations(entgql.Bind()),
	}
}

// PoolProperties hold information on the current pool
type PoolProperties struct {
	ent.Schema
}

// Edges of PoolProperties
func (PoolProperties) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("pool", ResourcePool.Type).
			Ref("poolProperties").
			Unique(),
		edge.To("resourceType", ResourceType.Type).
			Annotations(entgql.Bind()),
		edge.To("properties", Property.Type).
			Annotations(entgql.Bind()),
	}
}
