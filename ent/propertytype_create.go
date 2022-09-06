// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/net-auto/resourceManager/ent/property"
	"github.com/net-auto/resourceManager/ent/propertytype"
	"github.com/net-auto/resourceManager/ent/resourcetype"
)

// PropertyTypeCreate is the builder for creating a PropertyType entity.
type PropertyTypeCreate struct {
	config
	mutation *PropertyTypeMutation
	hooks    []Hook
}

// SetType sets the "type" field.
func (ptc *PropertyTypeCreate) SetType(pr propertytype.Type) *PropertyTypeCreate {
	ptc.mutation.SetType(pr)
	return ptc
}

// SetName sets the "name" field.
func (ptc *PropertyTypeCreate) SetName(s string) *PropertyTypeCreate {
	ptc.mutation.SetName(s)
	return ptc
}

// SetExternalID sets the "external_id" field.
func (ptc *PropertyTypeCreate) SetExternalID(s string) *PropertyTypeCreate {
	ptc.mutation.SetExternalID(s)
	return ptc
}

// SetNillableExternalID sets the "external_id" field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableExternalID(s *string) *PropertyTypeCreate {
	if s != nil {
		ptc.SetExternalID(*s)
	}
	return ptc
}

// SetIndex sets the "index" field.
func (ptc *PropertyTypeCreate) SetIndex(i int) *PropertyTypeCreate {
	ptc.mutation.SetIndex(i)
	return ptc
}

// SetNillableIndex sets the "index" field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableIndex(i *int) *PropertyTypeCreate {
	if i != nil {
		ptc.SetIndex(*i)
	}
	return ptc
}

// SetCategory sets the "category" field.
func (ptc *PropertyTypeCreate) SetCategory(s string) *PropertyTypeCreate {
	ptc.mutation.SetCategory(s)
	return ptc
}

// SetNillableCategory sets the "category" field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableCategory(s *string) *PropertyTypeCreate {
	if s != nil {
		ptc.SetCategory(*s)
	}
	return ptc
}

// SetIntVal sets the "int_val" field.
func (ptc *PropertyTypeCreate) SetIntVal(i int) *PropertyTypeCreate {
	ptc.mutation.SetIntVal(i)
	return ptc
}

// SetNillableIntVal sets the "int_val" field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableIntVal(i *int) *PropertyTypeCreate {
	if i != nil {
		ptc.SetIntVal(*i)
	}
	return ptc
}

// SetBoolVal sets the "bool_val" field.
func (ptc *PropertyTypeCreate) SetBoolVal(b bool) *PropertyTypeCreate {
	ptc.mutation.SetBoolVal(b)
	return ptc
}

// SetNillableBoolVal sets the "bool_val" field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableBoolVal(b *bool) *PropertyTypeCreate {
	if b != nil {
		ptc.SetBoolVal(*b)
	}
	return ptc
}

// SetFloatVal sets the "float_val" field.
func (ptc *PropertyTypeCreate) SetFloatVal(f float64) *PropertyTypeCreate {
	ptc.mutation.SetFloatVal(f)
	return ptc
}

// SetNillableFloatVal sets the "float_val" field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableFloatVal(f *float64) *PropertyTypeCreate {
	if f != nil {
		ptc.SetFloatVal(*f)
	}
	return ptc
}

// SetLatitudeVal sets the "latitude_val" field.
func (ptc *PropertyTypeCreate) SetLatitudeVal(f float64) *PropertyTypeCreate {
	ptc.mutation.SetLatitudeVal(f)
	return ptc
}

// SetNillableLatitudeVal sets the "latitude_val" field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableLatitudeVal(f *float64) *PropertyTypeCreate {
	if f != nil {
		ptc.SetLatitudeVal(*f)
	}
	return ptc
}

// SetLongitudeVal sets the "longitude_val" field.
func (ptc *PropertyTypeCreate) SetLongitudeVal(f float64) *PropertyTypeCreate {
	ptc.mutation.SetLongitudeVal(f)
	return ptc
}

// SetNillableLongitudeVal sets the "longitude_val" field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableLongitudeVal(f *float64) *PropertyTypeCreate {
	if f != nil {
		ptc.SetLongitudeVal(*f)
	}
	return ptc
}

// SetStringVal sets the "string_val" field.
func (ptc *PropertyTypeCreate) SetStringVal(s string) *PropertyTypeCreate {
	ptc.mutation.SetStringVal(s)
	return ptc
}

// SetNillableStringVal sets the "string_val" field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableStringVal(s *string) *PropertyTypeCreate {
	if s != nil {
		ptc.SetStringVal(*s)
	}
	return ptc
}

// SetRangeFromVal sets the "range_from_val" field.
func (ptc *PropertyTypeCreate) SetRangeFromVal(f float64) *PropertyTypeCreate {
	ptc.mutation.SetRangeFromVal(f)
	return ptc
}

// SetNillableRangeFromVal sets the "range_from_val" field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableRangeFromVal(f *float64) *PropertyTypeCreate {
	if f != nil {
		ptc.SetRangeFromVal(*f)
	}
	return ptc
}

// SetRangeToVal sets the "range_to_val" field.
func (ptc *PropertyTypeCreate) SetRangeToVal(f float64) *PropertyTypeCreate {
	ptc.mutation.SetRangeToVal(f)
	return ptc
}

// SetNillableRangeToVal sets the "range_to_val" field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableRangeToVal(f *float64) *PropertyTypeCreate {
	if f != nil {
		ptc.SetRangeToVal(*f)
	}
	return ptc
}

// SetIsInstanceProperty sets the "is_instance_property" field.
func (ptc *PropertyTypeCreate) SetIsInstanceProperty(b bool) *PropertyTypeCreate {
	ptc.mutation.SetIsInstanceProperty(b)
	return ptc
}

// SetNillableIsInstanceProperty sets the "is_instance_property" field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableIsInstanceProperty(b *bool) *PropertyTypeCreate {
	if b != nil {
		ptc.SetIsInstanceProperty(*b)
	}
	return ptc
}

// SetEditable sets the "editable" field.
func (ptc *PropertyTypeCreate) SetEditable(b bool) *PropertyTypeCreate {
	ptc.mutation.SetEditable(b)
	return ptc
}

// SetNillableEditable sets the "editable" field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableEditable(b *bool) *PropertyTypeCreate {
	if b != nil {
		ptc.SetEditable(*b)
	}
	return ptc
}

// SetMandatory sets the "mandatory" field.
func (ptc *PropertyTypeCreate) SetMandatory(b bool) *PropertyTypeCreate {
	ptc.mutation.SetMandatory(b)
	return ptc
}

// SetNillableMandatory sets the "mandatory" field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableMandatory(b *bool) *PropertyTypeCreate {
	if b != nil {
		ptc.SetMandatory(*b)
	}
	return ptc
}

// SetDeleted sets the "deleted" field.
func (ptc *PropertyTypeCreate) SetDeleted(b bool) *PropertyTypeCreate {
	ptc.mutation.SetDeleted(b)
	return ptc
}

// SetNillableDeleted sets the "deleted" field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableDeleted(b *bool) *PropertyTypeCreate {
	if b != nil {
		ptc.SetDeleted(*b)
	}
	return ptc
}

// SetNodeType sets the "nodeType" field.
func (ptc *PropertyTypeCreate) SetNodeType(s string) *PropertyTypeCreate {
	ptc.mutation.SetNodeType(s)
	return ptc
}

// SetNillableNodeType sets the "nodeType" field if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableNodeType(s *string) *PropertyTypeCreate {
	if s != nil {
		ptc.SetNodeType(*s)
	}
	return ptc
}

// AddPropertyIDs adds the "properties" edge to the Property entity by IDs.
func (ptc *PropertyTypeCreate) AddPropertyIDs(ids ...int) *PropertyTypeCreate {
	ptc.mutation.AddPropertyIDs(ids...)
	return ptc
}

// AddProperties adds the "properties" edges to the Property entity.
func (ptc *PropertyTypeCreate) AddProperties(p ...*Property) *PropertyTypeCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ptc.AddPropertyIDs(ids...)
}

// SetResourceTypeID sets the "resource_type" edge to the ResourceType entity by ID.
func (ptc *PropertyTypeCreate) SetResourceTypeID(id int) *PropertyTypeCreate {
	ptc.mutation.SetResourceTypeID(id)
	return ptc
}

// SetNillableResourceTypeID sets the "resource_type" edge to the ResourceType entity by ID if the given value is not nil.
func (ptc *PropertyTypeCreate) SetNillableResourceTypeID(id *int) *PropertyTypeCreate {
	if id != nil {
		ptc = ptc.SetResourceTypeID(*id)
	}
	return ptc
}

// SetResourceType sets the "resource_type" edge to the ResourceType entity.
func (ptc *PropertyTypeCreate) SetResourceType(r *ResourceType) *PropertyTypeCreate {
	return ptc.SetResourceTypeID(r.ID)
}

// Mutation returns the PropertyTypeMutation object of the builder.
func (ptc *PropertyTypeCreate) Mutation() *PropertyTypeMutation {
	return ptc.mutation
}

// Save creates the PropertyType in the database.
func (ptc *PropertyTypeCreate) Save(ctx context.Context) (*PropertyType, error) {
	var (
		err  error
		node *PropertyType
	)
	if err := ptc.defaults(); err != nil {
		return nil, err
	}
	if len(ptc.hooks) == 0 {
		if err = ptc.check(); err != nil {
			return nil, err
		}
		node, err = ptc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PropertyTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = ptc.check(); err != nil {
				return nil, err
			}
			ptc.mutation = mutation
			if node, err = ptc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(ptc.hooks) - 1; i >= 0; i-- {
			if ptc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = ptc.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, ptc.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*PropertyType)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from PropertyTypeMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (ptc *PropertyTypeCreate) SaveX(ctx context.Context) *PropertyType {
	v, err := ptc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ptc *PropertyTypeCreate) Exec(ctx context.Context) error {
	_, err := ptc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ptc *PropertyTypeCreate) ExecX(ctx context.Context) {
	if err := ptc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ptc *PropertyTypeCreate) defaults() error {
	if _, ok := ptc.mutation.IsInstanceProperty(); !ok {
		v := propertytype.DefaultIsInstanceProperty
		ptc.mutation.SetIsInstanceProperty(v)
	}
	if _, ok := ptc.mutation.Editable(); !ok {
		v := propertytype.DefaultEditable
		ptc.mutation.SetEditable(v)
	}
	if _, ok := ptc.mutation.Mandatory(); !ok {
		v := propertytype.DefaultMandatory
		ptc.mutation.SetMandatory(v)
	}
	if _, ok := ptc.mutation.Deleted(); !ok {
		v := propertytype.DefaultDeleted
		ptc.mutation.SetDeleted(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (ptc *PropertyTypeCreate) check() error {
	if _, ok := ptc.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New(`ent: missing required field "PropertyType.type"`)}
	}
	if v, ok := ptc.mutation.GetType(); ok {
		if err := propertytype.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "PropertyType.type": %w`, err)}
		}
	}
	if _, ok := ptc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "PropertyType.name"`)}
	}
	if _, ok := ptc.mutation.IsInstanceProperty(); !ok {
		return &ValidationError{Name: "is_instance_property", err: errors.New(`ent: missing required field "PropertyType.is_instance_property"`)}
	}
	if _, ok := ptc.mutation.Editable(); !ok {
		return &ValidationError{Name: "editable", err: errors.New(`ent: missing required field "PropertyType.editable"`)}
	}
	if _, ok := ptc.mutation.Mandatory(); !ok {
		return &ValidationError{Name: "mandatory", err: errors.New(`ent: missing required field "PropertyType.mandatory"`)}
	}
	if _, ok := ptc.mutation.Deleted(); !ok {
		return &ValidationError{Name: "deleted", err: errors.New(`ent: missing required field "PropertyType.deleted"`)}
	}
	return nil
}

func (ptc *PropertyTypeCreate) sqlSave(ctx context.Context) (*PropertyType, error) {
	_node, _spec := ptc.createSpec()
	if err := sqlgraph.CreateNode(ctx, ptc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (ptc *PropertyTypeCreate) createSpec() (*PropertyType, *sqlgraph.CreateSpec) {
	var (
		_node = &PropertyType{config: ptc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: propertytype.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: propertytype.FieldID,
			},
		}
	)
	if value, ok := ptc.mutation.GetType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: propertytype.FieldType,
		})
		_node.Type = value
	}
	if value, ok := ptc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldName,
		})
		_node.Name = value
	}
	if value, ok := ptc.mutation.ExternalID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldExternalID,
		})
		_node.ExternalID = value
	}
	if value, ok := ptc.mutation.Index(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: propertytype.FieldIndex,
		})
		_node.Index = value
	}
	if value, ok := ptc.mutation.Category(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldCategory,
		})
		_node.Category = value
	}
	if value, ok := ptc.mutation.IntVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: propertytype.FieldIntVal,
		})
		_node.IntVal = &value
	}
	if value, ok := ptc.mutation.BoolVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldBoolVal,
		})
		_node.BoolVal = &value
	}
	if value, ok := ptc.mutation.FloatVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldFloatVal,
		})
		_node.FloatVal = &value
	}
	if value, ok := ptc.mutation.LatitudeVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldLatitudeVal,
		})
		_node.LatitudeVal = &value
	}
	if value, ok := ptc.mutation.LongitudeVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldLongitudeVal,
		})
		_node.LongitudeVal = &value
	}
	if value, ok := ptc.mutation.StringVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldStringVal,
		})
		_node.StringVal = &value
	}
	if value, ok := ptc.mutation.RangeFromVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldRangeFromVal,
		})
		_node.RangeFromVal = &value
	}
	if value, ok := ptc.mutation.RangeToVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: propertytype.FieldRangeToVal,
		})
		_node.RangeToVal = &value
	}
	if value, ok := ptc.mutation.IsInstanceProperty(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldIsInstanceProperty,
		})
		_node.IsInstanceProperty = value
	}
	if value, ok := ptc.mutation.Editable(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldEditable,
		})
		_node.Editable = value
	}
	if value, ok := ptc.mutation.Mandatory(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldMandatory,
		})
		_node.Mandatory = value
	}
	if value, ok := ptc.mutation.Deleted(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: propertytype.FieldDeleted,
		})
		_node.Deleted = value
	}
	if value, ok := ptc.mutation.NodeType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: propertytype.FieldNodeType,
		})
		_node.NodeType = value
	}
	if nodes := ptc.mutation.PropertiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   propertytype.PropertiesTable,
			Columns: []string{propertytype.PropertiesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: property.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ptc.mutation.ResourceTypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   propertytype.ResourceTypeTable,
			Columns: []string{propertytype.ResourceTypeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: resourcetype.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.resource_type_property_types = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// PropertyTypeCreateBulk is the builder for creating many PropertyType entities in bulk.
type PropertyTypeCreateBulk struct {
	config
	builders []*PropertyTypeCreate
}

// Save creates the PropertyType entities in the database.
func (ptcb *PropertyTypeCreateBulk) Save(ctx context.Context) ([]*PropertyType, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ptcb.builders))
	nodes := make([]*PropertyType, len(ptcb.builders))
	mutators := make([]Mutator, len(ptcb.builders))
	for i := range ptcb.builders {
		func(i int, root context.Context) {
			builder := ptcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*PropertyTypeMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ptcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ptcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ptcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ptcb *PropertyTypeCreateBulk) SaveX(ctx context.Context) []*PropertyType {
	v, err := ptcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ptcb *PropertyTypeCreateBulk) Exec(ctx context.Context) error {
	_, err := ptcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ptcb *PropertyTypeCreateBulk) ExecX(ctx context.Context) {
	if err := ptcb.Exec(ctx); err != nil {
		panic(err)
	}
}
