// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copied over from inventory, since ent.go cannot handle schema split across multiple packages e.g. resource manager + inventory

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// PropertyType defines the property type schema.
type PropertyType struct {
	ent.Schema
}

// Fields returns property type fields.
func (PropertyType) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("type").
			Values(
				"string",
				"int",
				"bool",
				"float",
				"date",
				"enum",
				"range",
				"email",
				"gps_location",
				"datetime_local",
				"node",
			),
		field.String("name"),
		field.String("external_id").
			Unique().
			Optional(),
		field.Int("index").
			Optional(),
		field.String("category").
			Optional(),
		field.Int("int_val").
			StructTag(`json:"intValue" gqlgen:"intValue"`).
			Optional().
			Nillable(),
		field.Bool("bool_val").
			StructTag(`json:"booleanValue" gqlgen:"booleanValue"`).
			Optional().
			Nillable(),
		field.Float("float_val").
			StructTag(`json:"floatValue" gqlgen:"floatValue"`).
			Optional().
			Nillable(),
		field.Float("latitude_val").
			StructTag(`json:"latitudeValue" gqlgen:"latitudeValue"`).
			Optional().
			Nillable(),
		field.Float("longitude_val").
			StructTag(`json:"longitudeValue" gqlgen:"longitudeValue"`).
			Optional().
			Nillable(),
		field.Text("string_val").
			StructTag(`json:"stringValue" gqlgen:"stringValue"`).
			Optional().
			Nillable(),
		field.Float("range_from_val").
			StructTag(`json:"rangeFromValue" gqlgen:"rangeFromValue"`).
			Optional().
			Nillable(),
		field.Float("range_to_val").
			StructTag(`json:"rangeToValue" gqlgen:"rangeToValue"`).
			Optional().
			Nillable(),
		field.Bool("is_instance_property").
			StructTag(`gqlgen:"isInstanceProperty"`).
			Default(true),
		field.Bool("editable").
			StructTag(`gqlgen:"isEditable"`).
			Default(true),
		field.Bool("mandatory").
			StructTag(`gqlgen:"isMandatory"`).
			Default(false),
		field.Bool("deleted").
			StructTag(`gqlgen:"isDeleted"`).
			Default(false),
		field.String("nodeType").Optional(),
	}
}

// Edges returns property type edges.
func (PropertyType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("properties", Property.Type).
			Ref("type"),
		edge.From("resource_type", ResourceType.Type).
			Ref("property_types").Unique(),
	}
}

// Property defines the property schema.
type Property struct {
	ent.Schema
}

// Fields returns property fields.
func (Property) Fields() []ent.Field {
	return []ent.Field{
		field.Int("int_val").
			StructTag(`json:"intValue" gqlgen:"intValue"`).
			Optional().
			Nillable(),
		field.Bool("bool_val").
			StructTag(`json:"booleanValue" gqlgen:"booleanValue"`).
			Optional().
			Nillable(),
		field.Float("float_val").
			StructTag(`json:"floatValue" gqlgen:"floatValue"`).
			Optional().
			Nillable(),
		field.Float("latitude_val").
			StructTag(`json:"latitudeValue" gqlgen:"latitudeValue"`).
			Optional().
			Nillable(),
		field.Float("longitude_val").
			StructTag(`json:"longitudeValue" gqlgen:"longitudeValue"`).
			Optional().
			Nillable(),
		field.Float("range_from_val").
			StructTag(`json:"rangeFromValue" gqlgen:"rangeFromValue"`).
			Optional().
			Nillable(),
		field.Float("range_to_val").
			StructTag(`json:"rangeToValue" gqlgen:"rangeToValue"`).
			Optional().
			Nillable(),
		field.String("string_val").
			StructTag(`json:"stringValue" gqlgen:"stringValue"`).
			Optional().
			Nillable(),
	}
}

// Edges returns property edges.
func (Property) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("type", PropertyType.Type).
			Unique().
			Required().
			StructTag(`gqlgen:"propertyType"`),
	}
}
