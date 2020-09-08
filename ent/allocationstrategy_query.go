// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"math"

	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/facebook/ent/schema/field"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/ent/predicate"
	"github.com/net-auto/resourceManager/ent/resourcepool"
)

// AllocationStrategyQuery is the builder for querying AllocationStrategy entities.
type AllocationStrategyQuery struct {
	config
	limit      *int
	offset     *int
	order      []OrderFunc
	unique     []string
	predicates []predicate.AllocationStrategy
	// eager-loading edges.
	withPools *ResourcePoolQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
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

// Order adds an order step to the query.
func (asq *AllocationStrategyQuery) Order(o ...OrderFunc) *AllocationStrategyQuery {
	asq.order = append(asq.order, o...)
	return asq
}

// QueryPools chains the current query on the pools edge.
func (asq *AllocationStrategyQuery) QueryPools() *ResourcePoolQuery {
	query := &ResourcePoolQuery{config: asq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := asq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(allocationstrategy.Table, allocationstrategy.FieldID, asq.sqlQuery()),
			sqlgraph.To(resourcepool.Table, resourcepool.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, allocationstrategy.PoolsTable, allocationstrategy.PoolsColumn),
		)
		fromU = sqlgraph.SetNeighbors(asq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first AllocationStrategy entity in the query. Returns *NotFoundError when no allocationstrategy was found.
func (asq *AllocationStrategyQuery) First(ctx context.Context) (*AllocationStrategy, error) {
	asSlice, err := asq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(asSlice) == 0 {
		return nil, &NotFoundError{allocationstrategy.Label}
	}
	return asSlice[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (asq *AllocationStrategyQuery) FirstX(ctx context.Context) *AllocationStrategy {
	as, err := asq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return as
}

// FirstID returns the first AllocationStrategy id in the query. Returns *NotFoundError when no id was found.
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

// FirstXID is like FirstID, but panics if an error occurs.
func (asq *AllocationStrategyQuery) FirstXID(ctx context.Context) int {
	id, err := asq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only AllocationStrategy entity in the query, returns an error if not exactly one entity was returned.
func (asq *AllocationStrategyQuery) Only(ctx context.Context) (*AllocationStrategy, error) {
	asSlice, err := asq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(asSlice) {
	case 1:
		return asSlice[0], nil
	case 0:
		return nil, &NotFoundError{allocationstrategy.Label}
	default:
		return nil, &NotSingularError{allocationstrategy.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (asq *AllocationStrategyQuery) OnlyX(ctx context.Context) *AllocationStrategy {
	as, err := asq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return as
}

// OnlyID returns the only AllocationStrategy id in the query, returns an error if not exactly one id was returned.
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
	asSlice, err := asq.All(ctx)
	if err != nil {
		panic(err)
	}
	return asSlice
}

// IDs executes the query and returns a list of AllocationStrategy ids.
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

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (asq *AllocationStrategyQuery) Clone() *AllocationStrategyQuery {
	return &AllocationStrategyQuery{
		config:     asq.config,
		limit:      asq.limit,
		offset:     asq.offset,
		order:      append([]OrderFunc{}, asq.order...),
		unique:     append([]string{}, asq.unique...),
		predicates: append([]predicate.AllocationStrategy{}, asq.predicates...),
		// clone intermediate query.
		sql:  asq.sql.Clone(),
		path: asq.path,
	}
}

//  WithPools tells the query-builder to eager-loads the nodes that are connected to
// the "pools" edge. The optional arguments used to configure the query builder of the edge.
func (asq *AllocationStrategyQuery) WithPools(opts ...func(*ResourcePoolQuery)) *AllocationStrategyQuery {
	query := &ResourcePoolQuery{config: asq.config}
	for _, opt := range opts {
		opt(query)
	}
	asq.withPools = query
	return asq
}

// GroupBy used to group vertices by one or more fields/columns.
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
//
func (asq *AllocationStrategyQuery) GroupBy(field string, fields ...string) *AllocationStrategyGroupBy {
	group := &AllocationStrategyGroupBy{config: asq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := asq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return asq.sqlQuery(), nil
	}
	return group
}

// Select one or more fields from the given query.
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
//
func (asq *AllocationStrategyQuery) Select(field string, fields ...string) *AllocationStrategySelect {
	selector := &AllocationStrategySelect{config: asq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := asq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return asq.sqlQuery(), nil
	}
	return selector
}

func (asq *AllocationStrategyQuery) prepareQuery(ctx context.Context) error {
	if asq.path != nil {
		prev, err := asq.path(ctx)
		if err != nil {
			return err
		}
		asq.sql = prev
	}
	return nil
}

func (asq *AllocationStrategyQuery) sqlAll(ctx context.Context) ([]*AllocationStrategy, error) {
	var (
		nodes       = []*AllocationStrategy{}
		_spec       = asq.querySpec()
		loadedTypes = [1]bool{
			asq.withPools != nil,
		}
	)
	_spec.ScanValues = func() []interface{} {
		node := &AllocationStrategy{config: asq.config}
		nodes = append(nodes, node)
		values := node.scanValues()
		return values
	}
	_spec.Assign = func(values ...interface{}) error {
		if len(nodes) == 0 {
			return fmt.Errorf("ent: Assign called without calling ScanValues")
		}
		node := nodes[len(nodes)-1]
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(values...)
	}
	if err := sqlgraph.QueryNodes(ctx, asq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if query := asq.withPools; query != nil {
		fks := make([]driver.Value, 0, len(nodes))
		nodeids := make(map[int]*AllocationStrategy)
		for i := range nodes {
			fks = append(fks, nodes[i].ID)
			nodeids[nodes[i].ID] = nodes[i]
		}
		query.withFKs = true
		query.Where(predicate.ResourcePool(func(s *sql.Selector) {
			s.Where(sql.InValues(allocationstrategy.PoolsColumn, fks...))
		}))
		neighbors, err := query.All(ctx)
		if err != nil {
			return nil, err
		}
		for _, n := range neighbors {
			fk := n.resource_pool_allocation_strategy
			if fk == nil {
				return nil, fmt.Errorf(`foreign-key "resource_pool_allocation_strategy" is nil for node %v`, n.ID)
			}
			node, ok := nodeids[*fk]
			if !ok {
				return nil, fmt.Errorf(`unexpected foreign-key "resource_pool_allocation_strategy" returned %v for node %v`, *fk, n.ID)
			}
			node.Edges.Pools = append(node.Edges.Pools, n)
		}
	}

	return nodes, nil
}

func (asq *AllocationStrategyQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := asq.querySpec()
	return sqlgraph.CountNodes(ctx, asq.driver, _spec)
}

func (asq *AllocationStrategyQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := asq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
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

func (asq *AllocationStrategyQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(asq.driver.Dialect())
	t1 := builder.Table(allocationstrategy.Table)
	selector := builder.Select(t1.Columns(allocationstrategy.Columns...)...).From(t1)
	if asq.sql != nil {
		selector = asq.sql
		selector.Select(selector.Columns(allocationstrategy.Columns...)...)
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

// AllocationStrategyGroupBy is the builder for group-by AllocationStrategy entities.
type AllocationStrategyGroupBy struct {
	config
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

// Scan applies the group-by query and scan the result into the given value.
func (asgb *AllocationStrategyGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := asgb.path(ctx)
	if err != nil {
		return err
	}
	asgb.sql = query
	return asgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (asgb *AllocationStrategyGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := asgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (asgb *AllocationStrategyGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(asgb.fields) > 1 {
		return nil, errors.New("ent: AllocationStrategyGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := asgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (asgb *AllocationStrategyGroupBy) StringsX(ctx context.Context) []string {
	v, err := asgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from group-by. It is only allowed when querying group-by with one field.
func (asgb *AllocationStrategyGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = asgb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{allocationstrategy.Label}
	default:
		err = fmt.Errorf("ent: AllocationStrategyGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (asgb *AllocationStrategyGroupBy) StringX(ctx context.Context) string {
	v, err := asgb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (asgb *AllocationStrategyGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(asgb.fields) > 1 {
		return nil, errors.New("ent: AllocationStrategyGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := asgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (asgb *AllocationStrategyGroupBy) IntsX(ctx context.Context) []int {
	v, err := asgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from group-by. It is only allowed when querying group-by with one field.
func (asgb *AllocationStrategyGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = asgb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{allocationstrategy.Label}
	default:
		err = fmt.Errorf("ent: AllocationStrategyGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (asgb *AllocationStrategyGroupBy) IntX(ctx context.Context) int {
	v, err := asgb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (asgb *AllocationStrategyGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(asgb.fields) > 1 {
		return nil, errors.New("ent: AllocationStrategyGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := asgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (asgb *AllocationStrategyGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := asgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from group-by. It is only allowed when querying group-by with one field.
func (asgb *AllocationStrategyGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = asgb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{allocationstrategy.Label}
	default:
		err = fmt.Errorf("ent: AllocationStrategyGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (asgb *AllocationStrategyGroupBy) Float64X(ctx context.Context) float64 {
	v, err := asgb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (asgb *AllocationStrategyGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(asgb.fields) > 1 {
		return nil, errors.New("ent: AllocationStrategyGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := asgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (asgb *AllocationStrategyGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := asgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from group-by. It is only allowed when querying group-by with one field.
func (asgb *AllocationStrategyGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = asgb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{allocationstrategy.Label}
	default:
		err = fmt.Errorf("ent: AllocationStrategyGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (asgb *AllocationStrategyGroupBy) BoolX(ctx context.Context) bool {
	v, err := asgb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (asgb *AllocationStrategyGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := asgb.sqlQuery().Query()
	if err := asgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (asgb *AllocationStrategyGroupBy) sqlQuery() *sql.Selector {
	selector := asgb.sql
	columns := make([]string, 0, len(asgb.fields)+len(asgb.fns))
	columns = append(columns, asgb.fields...)
	for _, fn := range asgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(asgb.fields...)
}

// AllocationStrategySelect is the builder for select fields of AllocationStrategy entities.
type AllocationStrategySelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (ass *AllocationStrategySelect) Scan(ctx context.Context, v interface{}) error {
	query, err := ass.path(ctx)
	if err != nil {
		return err
	}
	ass.sql = query
	return ass.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ass *AllocationStrategySelect) ScanX(ctx context.Context, v interface{}) {
	if err := ass.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (ass *AllocationStrategySelect) Strings(ctx context.Context) ([]string, error) {
	if len(ass.fields) > 1 {
		return nil, errors.New("ent: AllocationStrategySelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := ass.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ass *AllocationStrategySelect) StringsX(ctx context.Context) []string {
	v, err := ass.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from selector. It is only allowed when selecting one field.
func (ass *AllocationStrategySelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = ass.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{allocationstrategy.Label}
	default:
		err = fmt.Errorf("ent: AllocationStrategySelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (ass *AllocationStrategySelect) StringX(ctx context.Context) string {
	v, err := ass.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (ass *AllocationStrategySelect) Ints(ctx context.Context) ([]int, error) {
	if len(ass.fields) > 1 {
		return nil, errors.New("ent: AllocationStrategySelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := ass.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ass *AllocationStrategySelect) IntsX(ctx context.Context) []int {
	v, err := ass.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from selector. It is only allowed when selecting one field.
func (ass *AllocationStrategySelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = ass.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{allocationstrategy.Label}
	default:
		err = fmt.Errorf("ent: AllocationStrategySelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (ass *AllocationStrategySelect) IntX(ctx context.Context) int {
	v, err := ass.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (ass *AllocationStrategySelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(ass.fields) > 1 {
		return nil, errors.New("ent: AllocationStrategySelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := ass.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ass *AllocationStrategySelect) Float64sX(ctx context.Context) []float64 {
	v, err := ass.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from selector. It is only allowed when selecting one field.
func (ass *AllocationStrategySelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = ass.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{allocationstrategy.Label}
	default:
		err = fmt.Errorf("ent: AllocationStrategySelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (ass *AllocationStrategySelect) Float64X(ctx context.Context) float64 {
	v, err := ass.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (ass *AllocationStrategySelect) Bools(ctx context.Context) ([]bool, error) {
	if len(ass.fields) > 1 {
		return nil, errors.New("ent: AllocationStrategySelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := ass.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ass *AllocationStrategySelect) BoolsX(ctx context.Context) []bool {
	v, err := ass.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from selector. It is only allowed when selecting one field.
func (ass *AllocationStrategySelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = ass.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{allocationstrategy.Label}
	default:
		err = fmt.Errorf("ent: AllocationStrategySelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (ass *AllocationStrategySelect) BoolX(ctx context.Context) bool {
	v, err := ass.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ass *AllocationStrategySelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := ass.sqlQuery().Query()
	if err := ass.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ass *AllocationStrategySelect) sqlQuery() sql.Querier {
	selector := ass.sql
	selector.Select(selector.Columns(ass.fields...)...)
	return selector
}
