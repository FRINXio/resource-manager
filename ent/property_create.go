// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/facebook/ent/schema/field"
	"github.com/net-auto/resourceManager/ent/property"
	"github.com/net-auto/resourceManager/ent/propertytype"
)

// PropertyCreate is the builder for creating a Property entity.
type PropertyCreate struct {
	config
	mutation *PropertyMutation
	hooks    []Hook
}

// SetIntVal sets the int_val field.
func (pc *PropertyCreate) SetIntVal(i int) *PropertyCreate {
	pc.mutation.SetIntVal(i)
	return pc
}

// SetNillableIntVal sets the int_val field if the given value is not nil.
func (pc *PropertyCreate) SetNillableIntVal(i *int) *PropertyCreate {
	if i != nil {
		pc.SetIntVal(*i)
	}
	return pc
}

// SetBoolVal sets the bool_val field.
func (pc *PropertyCreate) SetBoolVal(b bool) *PropertyCreate {
	pc.mutation.SetBoolVal(b)
	return pc
}

// SetNillableBoolVal sets the bool_val field if the given value is not nil.
func (pc *PropertyCreate) SetNillableBoolVal(b *bool) *PropertyCreate {
	if b != nil {
		pc.SetBoolVal(*b)
	}
	return pc
}

// SetFloatVal sets the float_val field.
func (pc *PropertyCreate) SetFloatVal(f float64) *PropertyCreate {
	pc.mutation.SetFloatVal(f)
	return pc
}

// SetNillableFloatVal sets the float_val field if the given value is not nil.
func (pc *PropertyCreate) SetNillableFloatVal(f *float64) *PropertyCreate {
	if f != nil {
		pc.SetFloatVal(*f)
	}
	return pc
}

// SetLatitudeVal sets the latitude_val field.
func (pc *PropertyCreate) SetLatitudeVal(f float64) *PropertyCreate {
	pc.mutation.SetLatitudeVal(f)
	return pc
}

// SetNillableLatitudeVal sets the latitude_val field if the given value is not nil.
func (pc *PropertyCreate) SetNillableLatitudeVal(f *float64) *PropertyCreate {
	if f != nil {
		pc.SetLatitudeVal(*f)
	}
	return pc
}

// SetLongitudeVal sets the longitude_val field.
func (pc *PropertyCreate) SetLongitudeVal(f float64) *PropertyCreate {
	pc.mutation.SetLongitudeVal(f)
	return pc
}

// SetNillableLongitudeVal sets the longitude_val field if the given value is not nil.
func (pc *PropertyCreate) SetNillableLongitudeVal(f *float64) *PropertyCreate {
	if f != nil {
		pc.SetLongitudeVal(*f)
	}
	return pc
}

// SetRangeFromVal sets the range_from_val field.
func (pc *PropertyCreate) SetRangeFromVal(f float64) *PropertyCreate {
	pc.mutation.SetRangeFromVal(f)
	return pc
}

// SetNillableRangeFromVal sets the range_from_val field if the given value is not nil.
func (pc *PropertyCreate) SetNillableRangeFromVal(f *float64) *PropertyCreate {
	if f != nil {
		pc.SetRangeFromVal(*f)
	}
	return pc
}

// SetRangeToVal sets the range_to_val field.
func (pc *PropertyCreate) SetRangeToVal(f float64) *PropertyCreate {
	pc.mutation.SetRangeToVal(f)
	return pc
}

// SetNillableRangeToVal sets the range_to_val field if the given value is not nil.
func (pc *PropertyCreate) SetNillableRangeToVal(f *float64) *PropertyCreate {
	if f != nil {
		pc.SetRangeToVal(*f)
	}
	return pc
}

// SetStringVal sets the string_val field.
func (pc *PropertyCreate) SetStringVal(s string) *PropertyCreate {
	pc.mutation.SetStringVal(s)
	return pc
}

// SetNillableStringVal sets the string_val field if the given value is not nil.
func (pc *PropertyCreate) SetNillableStringVal(s *string) *PropertyCreate {
	if s != nil {
		pc.SetStringVal(*s)
	}
	return pc
}

// SetTypeID sets the type edge to PropertyType by id.
func (pc *PropertyCreate) SetTypeID(id int) *PropertyCreate {
	pc.mutation.SetTypeID(id)
	return pc
}

// SetType sets the type edge to PropertyType.
func (pc *PropertyCreate) SetType(p *PropertyType) *PropertyCreate {
	return pc.SetTypeID(p.ID)
}

// Mutation returns the PropertyMutation object of the builder.
func (pc *PropertyCreate) Mutation() *PropertyMutation {
	return pc.mutation
}

// Save creates the Property in the database.
func (pc *PropertyCreate) Save(ctx context.Context) (*Property, error) {
	if err := pc.preSave(); err != nil {
		return nil, err
	}
	var (
		err  error
		node *Property
	)
	if len(pc.hooks) == 0 {
		node, err = pc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PropertyMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			pc.mutation = mutation
			node, err = pc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(pc.hooks) - 1; i >= 0; i-- {
			mut = pc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, pc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (pc *PropertyCreate) SaveX(ctx context.Context) *Property {
	v, err := pc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (pc *PropertyCreate) preSave() error {
	if _, ok := pc.mutation.TypeID(); !ok {
		return &ValidationError{Name: "type", err: errors.New("ent: missing required edge \"type\"")}
	}
	return nil
}

func (pc *PropertyCreate) sqlSave(ctx context.Context) (*Property, error) {
	pr, _spec := pc.createSpec()
	if err := sqlgraph.CreateNode(ctx, pc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	pr.ID = int(id)
	return pr, nil
}

func (pc *PropertyCreate) createSpec() (*Property, *sqlgraph.CreateSpec) {
	var (
		pr    = &Property{config: pc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: property.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: property.FieldID,
			},
		}
	)
	if value, ok := pc.mutation.IntVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: property.FieldIntVal,
		})
		pr.IntVal = &value
	}
	if value, ok := pc.mutation.BoolVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: property.FieldBoolVal,
		})
		pr.BoolVal = &value
	}
	if value, ok := pc.mutation.FloatVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: property.FieldFloatVal,
		})
		pr.FloatVal = &value
	}
	if value, ok := pc.mutation.LatitudeVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: property.FieldLatitudeVal,
		})
		pr.LatitudeVal = &value
	}
	if value, ok := pc.mutation.LongitudeVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: property.FieldLongitudeVal,
		})
		pr.LongitudeVal = &value
	}
	if value, ok := pc.mutation.RangeFromVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: property.FieldRangeFromVal,
		})
		pr.RangeFromVal = &value
	}
	if value, ok := pc.mutation.RangeToVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeFloat64,
			Value:  value,
			Column: property.FieldRangeToVal,
		})
		pr.RangeToVal = &value
	}
	if value, ok := pc.mutation.StringVal(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: property.FieldStringVal,
		})
		pr.StringVal = &value
	}
	if nodes := pc.mutation.TypeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   property.TypeTable,
			Columns: []string{property.TypeColumn},
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	return pr, _spec
}

// PropertyCreateBulk is the builder for creating a bulk of Property entities.
type PropertyCreateBulk struct {
	config
	builders []*PropertyCreate
}

// Save creates the Property entities in the database.
func (pcb *PropertyCreateBulk) Save(ctx context.Context) ([]*Property, error) {
	specs := make([]*sqlgraph.CreateSpec, len(pcb.builders))
	nodes := make([]*Property, len(pcb.builders))
	mutators := make([]Mutator, len(pcb.builders))
	for i := range pcb.builders {
		func(i int, root context.Context) {
			builder := pcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				if err := builder.preSave(); err != nil {
					return nil, err
				}
				mutation, ok := m.(*PropertyMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, pcb.builders[i+1].mutation)
				} else {
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, pcb.driver, &sqlgraph.BatchCreateSpec{Nodes: specs}); err != nil {
						if cerr, ok := isSQLConstraintError(err); ok {
							err = cerr
						}
					}
				}
				mutation.done = true
				if err != nil {
					return nil, err
				}
				id := specs[i].ID.Value.(int64)
				nodes[i].ID = int(id)
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, pcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX calls Save and panics if Save returns an error.
func (pcb *PropertyCreateBulk) SaveX(ctx context.Context) []*Property {
	v, err := pcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}
