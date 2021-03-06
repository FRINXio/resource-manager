// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (as *AllocationStrategyQuery) CollectFields(ctx context.Context, satisfies ...string) *AllocationStrategyQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		as = as.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return as
}

func (as *AllocationStrategyQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *AllocationStrategyQuery {
	return as
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (pp *PoolPropertiesQuery) CollectFields(ctx context.Context, satisfies ...string) *PoolPropertiesQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		pp = pp.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return pp
}

func (pp *PoolPropertiesQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *PoolPropertiesQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "properties":
			pp = pp.WithProperties(func(query *PropertyQuery) {
				query.collectField(ctx, field)
			})
		case "resourceType":
			pp = pp.WithResourceType(func(query *ResourceTypeQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return pp
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (pr *PropertyQuery) CollectFields(ctx context.Context, satisfies ...string) *PropertyQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		pr = pr.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return pr
}

func (pr *PropertyQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *PropertyQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "type":
			pr = pr.WithType(func(query *PropertyTypeQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return pr
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (pt *PropertyTypeQuery) CollectFields(ctx context.Context, satisfies ...string) *PropertyTypeQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		pt = pt.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return pt
}

func (pt *PropertyTypeQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *PropertyTypeQuery {
	return pt
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (r *ResourceQuery) CollectFields(ctx context.Context, satisfies ...string) *ResourceQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		r = r.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return r
}

func (r *ResourceQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ResourceQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "nested_pool":
			r = r.WithNestedPool(func(query *ResourcePoolQuery) {
				query.collectField(ctx, field)
			})
		case "properties":
			r = r.WithProperties(func(query *PropertyQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return r
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (rp *ResourcePoolQuery) CollectFields(ctx context.Context, satisfies ...string) *ResourcePoolQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		rp = rp.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return rp
}

func (rp *ResourcePoolQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ResourcePoolQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "allocation_strategy":
			rp = rp.WithAllocationStrategy(func(query *AllocationStrategyQuery) {
				query.collectField(ctx, field)
			})
		case "claims":
			rp = rp.WithClaims(func(query *ResourceQuery) {
				query.collectField(ctx, field)
			})
		case "poolProperties":
			rp = rp.WithPoolProperties(func(query *PoolPropertiesQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return rp
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (rt *ResourceTypeQuery) CollectFields(ctx context.Context, satisfies ...string) *ResourceTypeQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		rt = rt.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return rt
}

func (rt *ResourceTypeQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *ResourceTypeQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "pools":
			rt = rt.WithPools(func(query *ResourcePoolQuery) {
				query.collectField(ctx, field)
			})
		case "property_types":
			rt = rt.WithPropertyTypes(func(query *PropertyTypeQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return rt
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (t *TagQuery) CollectFields(ctx context.Context, satisfies ...string) *TagQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		t = t.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return t
}

func (t *TagQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *TagQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "pools":
			t = t.WithPools(func(query *ResourcePoolQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return t
}
