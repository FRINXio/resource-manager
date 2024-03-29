// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/ent/predicate"
)

// AllocationStrategyDelete is the builder for deleting a AllocationStrategy entity.
type AllocationStrategyDelete struct {
	config
	hooks    []Hook
	mutation *AllocationStrategyMutation
}

// Where appends a list predicates to the AllocationStrategyDelete builder.
func (asd *AllocationStrategyDelete) Where(ps ...predicate.AllocationStrategy) *AllocationStrategyDelete {
	asd.mutation.Where(ps...)
	return asd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (asd *AllocationStrategyDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(asd.hooks) == 0 {
		affected, err = asd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AllocationStrategyMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			asd.mutation = mutation
			affected, err = asd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(asd.hooks) - 1; i >= 0; i-- {
			if asd.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = asd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, asd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (asd *AllocationStrategyDelete) ExecX(ctx context.Context) int {
	n, err := asd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (asd *AllocationStrategyDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: allocationstrategy.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: allocationstrategy.FieldID,
			},
		},
	}
	if ps := asd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, asd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	return affected, err
}

// AllocationStrategyDeleteOne is the builder for deleting a single AllocationStrategy entity.
type AllocationStrategyDeleteOne struct {
	asd *AllocationStrategyDelete
}

// Exec executes the deletion query.
func (asdo *AllocationStrategyDeleteOne) Exec(ctx context.Context) error {
	n, err := asdo.asd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{allocationstrategy.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (asdo *AllocationStrategyDeleteOne) ExecX(ctx context.Context) {
	asdo.asd.ExecX(ctx)
}
