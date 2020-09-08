// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/facebook/ent/schema/field"
	"github.com/net-auto/resourceManager/ent/predicate"
	"github.com/net-auto/resourceManager/ent/resourcepool"
)

// ResourcePoolDelete is the builder for deleting a ResourcePool entity.
type ResourcePoolDelete struct {
	config
	hooks      []Hook
	mutation   *ResourcePoolMutation
	predicates []predicate.ResourcePool
}

// Where adds a new predicate to the delete builder.
func (rpd *ResourcePoolDelete) Where(ps ...predicate.ResourcePool) *ResourcePoolDelete {
	rpd.predicates = append(rpd.predicates, ps...)
	return rpd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (rpd *ResourcePoolDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(rpd.hooks) == 0 {
		affected, err = rpd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ResourcePoolMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			rpd.mutation = mutation
			affected, err = rpd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(rpd.hooks) - 1; i >= 0; i-- {
			mut = rpd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, rpd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (rpd *ResourcePoolDelete) ExecX(ctx context.Context) int {
	n, err := rpd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (rpd *ResourcePoolDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: resourcepool.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: resourcepool.FieldID,
			},
		},
	}
	if ps := rpd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, rpd.driver, _spec)
}

// ResourcePoolDeleteOne is the builder for deleting a single ResourcePool entity.
type ResourcePoolDeleteOne struct {
	rpd *ResourcePoolDelete
}

// Exec executes the deletion query.
func (rpdo *ResourcePoolDeleteOne) Exec(ctx context.Context) error {
	n, err := rpdo.rpd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{resourcepool.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (rpdo *ResourcePoolDeleteOne) ExecX(ctx context.Context) {
	rpdo.rpd.ExecX(ctx)
}
