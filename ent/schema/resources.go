// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/index"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
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
		edge.From("pool_properties", PoolProperties.Type).
			Ref("resourceType"),
	}
}

func (ResourceType) Policy() ent.Policy {
	return RBAC
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
			Values("py", "js", "go").
			Default("js"),
		field.Text("script").
			NotEmpty(),
	}
}

func (AllocationStrategy) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("pools", ResourcePool.Type).
			Ref("allocation_strategy"),
		edge.To("pool_property_types", PropertyType.Type).
			Annotations(entgql.Bind()),
	}
}

func (AllocationStrategy) Policy() ent.Policy {
	return RBAC
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

func (ResourcePool) Indexes() []ent.Index {
	return []ent.Index{
		index.
			Edges("allocation_strategy"),
	}
}

func (ResourcePool) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
	}
}

func (ResourcePool) Policy() ent.Policy {
	return RBAC
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
		field.Text("description").
			Optional().
			Nillable(),
		field.JSON("alternate_id", make(map[string]interface{})).
			Optional(),
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

func (Resource) Indexes() []ent.Index {
	return []ent.Index{
		index.
			Edges("pool"),
	}
}

func (Resource) Policy() ent.Policy {
	return ALWAYS_ALLOWED
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

func (PoolProperties) Indexes() []ent.Index {
	return []ent.Index{
		index.
			Edges("pool"),
	}
}

func (PoolProperties) Policy() ent.Policy {
	return RBAC
}
