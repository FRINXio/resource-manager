// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/facebook/ent/schema/field"
	"github.com/net-auto/resourceManager/ent/poolproperties"
	"github.com/net-auto/resourceManager/ent/predicate"
	"github.com/net-auto/resourceManager/ent/propertytype"
	"github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/net-auto/resourceManager/ent/resourcetype"
)

// ResourceTypeUpdate is the builder for updating ResourceType entities.
type ResourceTypeUpdate struct {
	config
	hooks    []Hook
	mutation *ResourceTypeMutation
}

// Where adds a new predicate for the builder.
func (rtu *ResourceTypeUpdate) Where(ps ...predicate.ResourceType) *ResourceTypeUpdate {
	rtu.mutation.predicates = append(rtu.mutation.predicates, ps...)
	return rtu
}

// SetName sets the name field.
func (rtu *ResourceTypeUpdate) SetName(s string) *ResourceTypeUpdate {
	rtu.mutation.SetName(s)
	return rtu
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (rtu *ResourceTypeUpdate) AddPropertyTypeIDs(ids ...int) *ResourceTypeUpdate {
	rtu.mutation.AddPropertyTypeIDs(ids...)
	return rtu
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (rtu *ResourceTypeUpdate) AddPropertyTypes(p ...*PropertyType) *ResourceTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return rtu.AddPropertyTypeIDs(ids...)
}

// AddPoolIDs adds the pools edge to ResourcePool by ids.
func (rtu *ResourceTypeUpdate) AddPoolIDs(ids ...int) *ResourceTypeUpdate {
	rtu.mutation.AddPoolIDs(ids...)
	return rtu
}

// AddPools adds the pools edges to ResourcePool.
func (rtu *ResourceTypeUpdate) AddPools(r ...*ResourcePool) *ResourceTypeUpdate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return rtu.AddPoolIDs(ids...)
}

// AddPoolPropertyIDs adds the pool_properties edge to PoolProperties by ids.
func (rtu *ResourceTypeUpdate) AddPoolPropertyIDs(ids ...int) *ResourceTypeUpdate {
	rtu.mutation.AddPoolPropertyIDs(ids...)
	return rtu
}

// AddPoolProperties adds the pool_properties edges to PoolProperties.
func (rtu *ResourceTypeUpdate) AddPoolProperties(p ...*PoolProperties) *ResourceTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return rtu.AddPoolPropertyIDs(ids...)
}

// Mutation returns the ResourceTypeMutation object of the builder.
func (rtu *ResourceTypeUpdate) Mutation() *ResourceTypeMutation {
	return rtu.mutation
}

// ClearPropertyTypes clears all "property_types" edges to type PropertyType.
func (rtu *ResourceTypeUpdate) ClearPropertyTypes() *ResourceTypeUpdate {
	rtu.mutation.ClearPropertyTypes()
	return rtu
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (rtu *ResourceTypeUpdate) RemovePropertyTypeIDs(ids ...int) *ResourceTypeUpdate {
	rtu.mutation.RemovePropertyTypeIDs(ids...)
	return rtu
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (rtu *ResourceTypeUpdate) RemovePropertyTypes(p ...*PropertyType) *ResourceTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return rtu.RemovePropertyTypeIDs(ids...)
}

// ClearPools clears all "pools" edges to type ResourcePool.
func (rtu *ResourceTypeUpdate) ClearPools() *ResourceTypeUpdate {
	rtu.mutation.ClearPools()
	return rtu
}

// RemovePoolIDs removes the pools edge to ResourcePool by ids.
func (rtu *ResourceTypeUpdate) RemovePoolIDs(ids ...int) *ResourceTypeUpdate {
	rtu.mutation.RemovePoolIDs(ids...)
	return rtu
}

// RemovePools removes pools edges to ResourcePool.
func (rtu *ResourceTypeUpdate) RemovePools(r ...*ResourcePool) *ResourceTypeUpdate {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return rtu.RemovePoolIDs(ids...)
}

// ClearPoolProperties clears all "pool_properties" edges to type PoolProperties.
func (rtu *ResourceTypeUpdate) ClearPoolProperties() *ResourceTypeUpdate {
	rtu.mutation.ClearPoolProperties()
	return rtu
}

// RemovePoolPropertyIDs removes the pool_properties edge to PoolProperties by ids.
func (rtu *ResourceTypeUpdate) RemovePoolPropertyIDs(ids ...int) *ResourceTypeUpdate {
	rtu.mutation.RemovePoolPropertyIDs(ids...)
	return rtu
}

// RemovePoolProperties removes pool_properties edges to PoolProperties.
func (rtu *ResourceTypeUpdate) RemovePoolProperties(p ...*PoolProperties) *ResourceTypeUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return rtu.RemovePoolPropertyIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (rtu *ResourceTypeUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(rtu.hooks) == 0 {
		if err = rtu.check(); err != nil {
			return 0, err
		}
		affected, err = rtu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ResourceTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = rtu.check(); err != nil {
				return 0, err
			}
			rtu.mutation = mutation
			affected, err = rtu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(rtu.hooks) - 1; i >= 0; i-- {
			mut = rtu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, rtu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (rtu *ResourceTypeUpdate) SaveX(ctx context.Context) int {
	affected, err := rtu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (rtu *ResourceTypeUpdate) Exec(ctx context.Context) error {
	_, err := rtu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rtu *ResourceTypeUpdate) ExecX(ctx context.Context) {
	if err := rtu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rtu *ResourceTypeUpdate) check() error {
	if v, ok := rtu.mutation.Name(); ok {
		if err := resourcetype.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf("ent: validator failed for field \"name\": %w", err)}
		}
	}
	return nil
}

func (rtu *ResourceTypeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   resourcetype.Table,
			Columns: resourcetype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: resourcetype.FieldID,
			},
		},
	}
	if ps := rtu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := rtu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: resourcetype.FieldName,
		})
	}
	if rtu.mutation.PropertyTypesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   resourcetype.PropertyTypesTable,
			Columns: []string{resourcetype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: propertytype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := rtu.mutation.RemovedPropertyTypesIDs(); len(nodes) > 0 && !rtu.mutation.PropertyTypesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   resourcetype.PropertyTypesTable,
			Columns: []string{resourcetype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: propertytype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := rtu.mutation.PropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   resourcetype.PropertyTypesTable,
			Columns: []string{resourcetype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: propertytype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if rtu.mutation.PoolsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   resourcetype.PoolsTable,
			Columns: []string{resourcetype.PoolsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: resourcepool.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := rtu.mutation.RemovedPoolsIDs(); len(nodes) > 0 && !rtu.mutation.PoolsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   resourcetype.PoolsTable,
			Columns: []string{resourcetype.PoolsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: resourcepool.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := rtu.mutation.PoolsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   resourcetype.PoolsTable,
			Columns: []string{resourcetype.PoolsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: resourcepool.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if rtu.mutation.PoolPropertiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   resourcetype.PoolPropertiesTable,
			Columns: resourcetype.PoolPropertiesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: poolproperties.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := rtu.mutation.RemovedPoolPropertiesIDs(); len(nodes) > 0 && !rtu.mutation.PoolPropertiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   resourcetype.PoolPropertiesTable,
			Columns: resourcetype.PoolPropertiesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: poolproperties.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := rtu.mutation.PoolPropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   resourcetype.PoolPropertiesTable,
			Columns: resourcetype.PoolPropertiesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: poolproperties.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, rtu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{resourcetype.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// ResourceTypeUpdateOne is the builder for updating a single ResourceType entity.
type ResourceTypeUpdateOne struct {
	config
	hooks    []Hook
	mutation *ResourceTypeMutation
}

// SetName sets the name field.
func (rtuo *ResourceTypeUpdateOne) SetName(s string) *ResourceTypeUpdateOne {
	rtuo.mutation.SetName(s)
	return rtuo
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (rtuo *ResourceTypeUpdateOne) AddPropertyTypeIDs(ids ...int) *ResourceTypeUpdateOne {
	rtuo.mutation.AddPropertyTypeIDs(ids...)
	return rtuo
}

// AddPropertyTypes adds the property_types edges to PropertyType.
func (rtuo *ResourceTypeUpdateOne) AddPropertyTypes(p ...*PropertyType) *ResourceTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return rtuo.AddPropertyTypeIDs(ids...)
}

// AddPoolIDs adds the pools edge to ResourcePool by ids.
func (rtuo *ResourceTypeUpdateOne) AddPoolIDs(ids ...int) *ResourceTypeUpdateOne {
	rtuo.mutation.AddPoolIDs(ids...)
	return rtuo
}

// AddPools adds the pools edges to ResourcePool.
func (rtuo *ResourceTypeUpdateOne) AddPools(r ...*ResourcePool) *ResourceTypeUpdateOne {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return rtuo.AddPoolIDs(ids...)
}

// AddPoolPropertyIDs adds the pool_properties edge to PoolProperties by ids.
func (rtuo *ResourceTypeUpdateOne) AddPoolPropertyIDs(ids ...int) *ResourceTypeUpdateOne {
	rtuo.mutation.AddPoolPropertyIDs(ids...)
	return rtuo
}

// AddPoolProperties adds the pool_properties edges to PoolProperties.
func (rtuo *ResourceTypeUpdateOne) AddPoolProperties(p ...*PoolProperties) *ResourceTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return rtuo.AddPoolPropertyIDs(ids...)
}

// Mutation returns the ResourceTypeMutation object of the builder.
func (rtuo *ResourceTypeUpdateOne) Mutation() *ResourceTypeMutation {
	return rtuo.mutation
}

// ClearPropertyTypes clears all "property_types" edges to type PropertyType.
func (rtuo *ResourceTypeUpdateOne) ClearPropertyTypes() *ResourceTypeUpdateOne {
	rtuo.mutation.ClearPropertyTypes()
	return rtuo
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (rtuo *ResourceTypeUpdateOne) RemovePropertyTypeIDs(ids ...int) *ResourceTypeUpdateOne {
	rtuo.mutation.RemovePropertyTypeIDs(ids...)
	return rtuo
}

// RemovePropertyTypes removes property_types edges to PropertyType.
func (rtuo *ResourceTypeUpdateOne) RemovePropertyTypes(p ...*PropertyType) *ResourceTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return rtuo.RemovePropertyTypeIDs(ids...)
}

// ClearPools clears all "pools" edges to type ResourcePool.
func (rtuo *ResourceTypeUpdateOne) ClearPools() *ResourceTypeUpdateOne {
	rtuo.mutation.ClearPools()
	return rtuo
}

// RemovePoolIDs removes the pools edge to ResourcePool by ids.
func (rtuo *ResourceTypeUpdateOne) RemovePoolIDs(ids ...int) *ResourceTypeUpdateOne {
	rtuo.mutation.RemovePoolIDs(ids...)
	return rtuo
}

// RemovePools removes pools edges to ResourcePool.
func (rtuo *ResourceTypeUpdateOne) RemovePools(r ...*ResourcePool) *ResourceTypeUpdateOne {
	ids := make([]int, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return rtuo.RemovePoolIDs(ids...)
}

// ClearPoolProperties clears all "pool_properties" edges to type PoolProperties.
func (rtuo *ResourceTypeUpdateOne) ClearPoolProperties() *ResourceTypeUpdateOne {
	rtuo.mutation.ClearPoolProperties()
	return rtuo
}

// RemovePoolPropertyIDs removes the pool_properties edge to PoolProperties by ids.
func (rtuo *ResourceTypeUpdateOne) RemovePoolPropertyIDs(ids ...int) *ResourceTypeUpdateOne {
	rtuo.mutation.RemovePoolPropertyIDs(ids...)
	return rtuo
}

// RemovePoolProperties removes pool_properties edges to PoolProperties.
func (rtuo *ResourceTypeUpdateOne) RemovePoolProperties(p ...*PoolProperties) *ResourceTypeUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return rtuo.RemovePoolPropertyIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (rtuo *ResourceTypeUpdateOne) Save(ctx context.Context) (*ResourceType, error) {
	var (
		err  error
		node *ResourceType
	)
	if len(rtuo.hooks) == 0 {
		if err = rtuo.check(); err != nil {
			return nil, err
		}
		node, err = rtuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ResourceTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = rtuo.check(); err != nil {
				return nil, err
			}
			rtuo.mutation = mutation
			node, err = rtuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(rtuo.hooks) - 1; i >= 0; i-- {
			mut = rtuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, rtuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (rtuo *ResourceTypeUpdateOne) SaveX(ctx context.Context) *ResourceType {
	node, err := rtuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (rtuo *ResourceTypeUpdateOne) Exec(ctx context.Context) error {
	_, err := rtuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rtuo *ResourceTypeUpdateOne) ExecX(ctx context.Context) {
	if err := rtuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rtuo *ResourceTypeUpdateOne) check() error {
	if v, ok := rtuo.mutation.Name(); ok {
		if err := resourcetype.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf("ent: validator failed for field \"name\": %w", err)}
		}
	}
	return nil
}

func (rtuo *ResourceTypeUpdateOne) sqlSave(ctx context.Context) (_node *ResourceType, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   resourcetype.Table,
			Columns: resourcetype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: resourcetype.FieldID,
			},
		},
	}
	id, ok := rtuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing ResourceType.ID for update")}
	}
	_spec.Node.ID.Value = id
	if value, ok := rtuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: resourcetype.FieldName,
		})
	}
	if rtuo.mutation.PropertyTypesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   resourcetype.PropertyTypesTable,
			Columns: []string{resourcetype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: propertytype.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := rtuo.mutation.RemovedPropertyTypesIDs(); len(nodes) > 0 && !rtuo.mutation.PropertyTypesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   resourcetype.PropertyTypesTable,
			Columns: []string{resourcetype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: propertytype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := rtuo.mutation.PropertyTypesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   resourcetype.PropertyTypesTable,
			Columns: []string{resourcetype.PropertyTypesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: propertytype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if rtuo.mutation.PoolsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   resourcetype.PoolsTable,
			Columns: []string{resourcetype.PoolsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: resourcepool.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := rtuo.mutation.RemovedPoolsIDs(); len(nodes) > 0 && !rtuo.mutation.PoolsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   resourcetype.PoolsTable,
			Columns: []string{resourcetype.PoolsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: resourcepool.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := rtuo.mutation.PoolsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   resourcetype.PoolsTable,
			Columns: []string{resourcetype.PoolsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: resourcepool.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if rtuo.mutation.PoolPropertiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   resourcetype.PoolPropertiesTable,
			Columns: resourcetype.PoolPropertiesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: poolproperties.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := rtuo.mutation.RemovedPoolPropertiesIDs(); len(nodes) > 0 && !rtuo.mutation.PoolPropertiesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   resourcetype.PoolPropertiesTable,
			Columns: resourcetype.PoolPropertiesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: poolproperties.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := rtuo.mutation.PoolPropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   resourcetype.PoolPropertiesTable,
			Columns: resourcetype.PoolPropertiesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: poolproperties.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &ResourceType{config: rtuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues()
	if err = sqlgraph.UpdateNode(ctx, rtuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{resourcetype.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return _node, nil
}
