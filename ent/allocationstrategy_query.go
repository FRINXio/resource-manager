// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/ent/predicate"
	"github.com/net-auto/resourceManager/ent/resourcepool"
)

// AllocationStrategyQuery is the builder for querying AllocationStrategy entities.
type AllocationStrategyQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.AllocationStrategy
	withPools  *ResourcePoolQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the AllocationStrategyQuery builder.
func (asq *AllocationStrategyQuery) Where(ps ...predicate.AllocationStrategy) *AllocationStrategyQuery {
	asq.predicates = append(asq.predicates, ps...)
	return asq
}

// Limit adds a limit step to the query.
func (asq *AllocationStrategyQuery) Limit(limit int) *AllocationStrategyQuery {
	asq.limit = &limit
	return asq
}

// Offset adds an offset step to the query.
func (asq *AllocationStrategyQuery) Offset(offset int) *AllocationStrategyQuery {
	asq.offset = &offset
	return asq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (asq *AllocationStrategyQuery) Unique(unique bool) *AllocationStrategyQuery {
	asq.unique = &unique
	return asq
}

// Order adds an order step to the query.
func (asq *AllocationStrategyQuery) Order(o ...OrderFunc) *AllocationStrategyQuery {
	asq.order = append(asq.order, o...)
	return asq
}

// QueryPools chains the current query on the "pools" edge.
func (asq *AllocationStrategyQuery) QueryPools() *ResourcePoolQuery {
	query := &ResourcePoolQuery{config: asq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := asq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := asq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(allocationstrategy.Table, allocationstrategy.FieldID, selector),
			sqlgraph.To(resourcepool.Table, resourcepool.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, allocationstrategy.PoolsTable, allocationstrategy.PoolsColumn),
		)
		fromU = sqlgraph.SetNeighbors(asq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first AllocationStrategy entity from the query.
// Returns a *NotFoundError when no AllocationStrategy was found.
func (asq *AllocationStrategyQuery) First(ctx context.Context) (*AllocationStrategy, error) {
	nodes, err := asq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{allocationstrategy.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (asq *AllocationStrategyQuery) FirstX(ctx context.Context) *AllocationStrategy {
	node, err := asq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first AllocationStrategy ID from the query.
// Returns a *NotFoundError when no AllocationStrategy ID was found.
func (asq *AllocationStrategyQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = asq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{allocationstrategy.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (asq *AllocationStrategyQuery) FirstIDX(ctx context.Context) int {
	id, err := asq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single AllocationStrategy entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one AllocationStrategy entity is found.
// Returns a *NotFoundError when no AllocationStrategy entities are found.
func (asq *AllocationStrategyQuery) Only(ctx context.Context) (*AllocationStrategy, error) {
	nodes, err := asq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{allocationstrategy.Label}
	default:
		return nil, &NotSingularError{allocationstrategy.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (asq *AllocationStrategyQuery) OnlyX(ctx context.Context) *AllocationStrategy {
	node, err := asq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only AllocationStrategy ID in the query.
// Returns a *NotSingularError when more than one AllocationStrategy ID is found.
// Returns a *NotFoundError when no entities are found.
func (asq *AllocationStrategyQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = asq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{allocationstrategy.Label}
	default:
		err = &NotSingularError{allocationstrategy.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (asq *AllocationStrategyQuery) OnlyIDX(ctx context.Context) int {
	id, err := asq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of AllocationStrategies.
func (asq *AllocationStrategyQuery) All(ctx context.Context) ([]*AllocationStrategy, error) {
	if err := asq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return asq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (asq *AllocationStrategyQuery) AllX(ctx context.Context) []*AllocationStrategy {
	nodes, err := asq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of AllocationStrategy IDs.
func (asq *AllocationStrategyQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := asq.Select(allocationstrategy.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (asq *AllocationStrategyQuery) IDsX(ctx context.Context) []int {
	ids, err := asq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (asq *AllocationStrategyQuery) Count(ctx context.Context) (int, error) {
	if err := asq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return asq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (asq *AllocationStrategyQuery) CountX(ctx context.Context) int {
	count, err := asq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (asq *AllocationStrategyQuery) Exist(ctx context.Context) (bool, error) {
	if err := asq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return asq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (asq *AllocationStrategyQuery) ExistX(ctx context.Context) bool {
	exist, err := asq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the AllocationStrategyQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (asq *AllocationStrategyQuery) Clone() *AllocationStrategyQuery {
	if asq == nil {
		return nil
	}
	return &AllocationStrategyQuery{
		config:     asq.config,
		limit:      asq.limit,
		offset:     asq.offset,
		order:      append([]OrderFunc{}, asq.order...),
		predicates: append([]predicate.AllocationStrategy{}, asq.predicates...),
		withPools:  asq.withPools.Clone(),
		// clone intermediate query.
		sql:    asq.sql.Clone(),
		path:   asq.path,
		unique: asq.unique,
	}
}

// WithPools tells the query-builder to eager-load the nodes that are connected to
// the "pools" edge. The optional arguments are used to configure the query builder of the edge.
func (asq *AllocationStrategyQuery) WithPools(opts ...func(*ResourcePoolQuery)) *AllocationStrategyQuery {
	query := &ResourcePoolQuery{config: asq.config}
	for _, opt := range opts {
		opt(query)
	}
	asq.withPools = query
	return asq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.AllocationStrategy.Query().
//		GroupBy(allocationstrategy.FieldName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (asq *AllocationStrategyQuery) GroupBy(field string, fields ...string) *AllocationStrategyGroupBy {
	grbuild := &AllocationStrategyGroupBy{config: asq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := asq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return asq.sqlQuery(ctx), nil
	}
	grbuild.label = allocationstrategy.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//	}
//
//	client.AllocationStrategy.Query().
//		Select(allocationstrategy.FieldName).
//		Scan(ctx, &v)
func (asq *AllocationStrategyQuery) Select(fields ...string) *AllocationStrategySelect {
	asq.fields = append(asq.fields, fields...)
	selbuild := &AllocationStrategySelect{AllocationStrategyQuery: asq}
	selbuild.label = allocationstrategy.Label
	selbuild.flds, selbuild.scan = &asq.fields, selbuild.Scan
	return selbuild
}

func (asq *AllocationStrategyQuery) prepareQuery(ctx context.Context) error {
	for _, f := range asq.fields {
		if !allocationstrategy.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if asq.path != nil {
		prev, err := asq.path(ctx)
		if err != nil {
			return err
		}
		asq.sql = prev
	}
	if allocationstrategy.Policy == nil {
		return errors.New("ent: uninitialized allocationstrategy.Policy (forgotten import ent/runtime?)")
	}
	if err := allocationstrategy.Policy.EvalQuery(ctx, asq); err != nil {
		return err
	}
	return nil
}

func (asq *AllocationStrategyQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*AllocationStrategy, error) {
	var (
		nodes       = []*AllocationStrategy{}
		_spec       = asq.querySpec()
		loadedTypes = [1]bool{
			asq.withPools != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*AllocationStrategy).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &AllocationStrategy{config: asq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, asq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := asq.withPools; query != nil {
		if err := asq.loadPools(ctx, query, nodes,
			func(n *AllocationStrategy) { n.Edges.Pools = []*ResourcePool{} },
			func(n *AllocationStrategy, e *ResourcePool) { n.Edges.Pools = append(n.Edges.Pools, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (asq *AllocationStrategyQuery) loadPools(ctx context.Context, query *ResourcePoolQuery, nodes []*AllocationStrategy, init func(*AllocationStrategy), assign func(*AllocationStrategy, *ResourcePool)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*AllocationStrategy)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.ResourcePool(func(s *sql.Selector) {
		s.Where(sql.InValues(allocationstrategy.PoolsColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.resource_pool_allocation_strategy
		if fk == nil {
			return fmt.Errorf(`foreign-key "resource_pool_allocation_strategy" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "resource_pool_allocation_strategy" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (asq *AllocationStrategyQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := asq.querySpec()
	_spec.Node.Columns = asq.fields
	if len(asq.fields) > 0 {
		_spec.Unique = asq.unique != nil && *asq.unique
	}
	return sqlgraph.CountNodes(ctx, asq.driver, _spec)
}

func (asq *AllocationStrategyQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := asq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (asq *AllocationStrategyQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   allocationstrategy.Table,
			Columns: allocationstrategy.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: allocationstrategy.FieldID,
			},
		},
		From:   asq.sql,
		Unique: true,
	}
	if unique := asq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := asq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, allocationstrategy.FieldID)
		for i := range fields {
			if fields[i] != allocationstrategy.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := asq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := asq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := asq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := asq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (asq *AllocationStrategyQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(asq.driver.Dialect())
	t1 := builder.Table(allocationstrategy.Table)
	columns := asq.fields
	if len(columns) == 0 {
		columns = allocationstrategy.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if asq.sql != nil {
		selector = asq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if asq.unique != nil && *asq.unique {
		selector.Distinct()
	}
	for _, p := range asq.predicates {
		p(selector)
	}
	for _, p := range asq.order {
		p(selector)
	}
	if offset := asq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := asq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// AllocationStrategyGroupBy is the group-by builder for AllocationStrategy entities.
type AllocationStrategyGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (asgb *AllocationStrategyGroupBy) Aggregate(fns ...AggregateFunc) *AllocationStrategyGroupBy {
	asgb.fns = append(asgb.fns, fns...)
	return asgb
}

// Scan applies the group-by query and scans the result into the given value.
func (asgb *AllocationStrategyGroupBy) Scan(ctx context.Context, v any) error {
	query, err := asgb.path(ctx)
	if err != nil {
		return err
	}
	asgb.sql = query
	return asgb.sqlScan(ctx, v)
}

func (asgb *AllocationStrategyGroupBy) sqlScan(ctx context.Context, v any) error {
	for _, f := range asgb.fields {
		if !allocationstrategy.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := asgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := asgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (asgb *AllocationStrategyGroupBy) sqlQuery() *sql.Selector {
	selector := asgb.sql.Select()
	aggregation := make([]string, 0, len(asgb.fns))
	for _, fn := range asgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(asgb.fields)+len(asgb.fns))
		for _, f := range asgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(asgb.fields...)...)
}

// AllocationStrategySelect is the builder for selecting fields of AllocationStrategy entities.
type AllocationStrategySelect struct {
	*AllocationStrategyQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (ass *AllocationStrategySelect) Scan(ctx context.Context, v any) error {
	if err := ass.prepareQuery(ctx); err != nil {
		return err
	}
	ass.sql = ass.AllocationStrategyQuery.sqlQuery(ctx)
	return ass.sqlScan(ctx, v)
}

func (ass *AllocationStrategySelect) sqlScan(ctx context.Context, v any) error {
	rows := &sql.Rows{}
	query, args := ass.sql.Query()
	if err := ass.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
