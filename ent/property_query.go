// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/net-auto/resourceManager/ent/predicate"
	"github.com/net-auto/resourceManager/ent/property"
	"github.com/net-auto/resourceManager/ent/propertytype"
	"github.com/net-auto/resourceManager/ent/resource"
)

// PropertyQuery is the builder for querying Property entities.
type PropertyQuery struct {
	config
	limit         *int
	offset        *int
	unique        *bool
	order         []OrderFunc
	fields        []string
	predicates    []predicate.Property
	withType      *PropertyTypeQuery
	withResources *ResourceQuery
	withFKs       bool
	modifiers     []func(*sql.Selector)
	loadTotal     []func(context.Context, []*Property) error
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the PropertyQuery builder.
func (pq *PropertyQuery) Where(ps ...predicate.Property) *PropertyQuery {
	pq.predicates = append(pq.predicates, ps...)
	return pq
}

// Limit adds a limit step to the query.
func (pq *PropertyQuery) Limit(limit int) *PropertyQuery {
	pq.limit = &limit
	return pq
}

// Offset adds an offset step to the query.
func (pq *PropertyQuery) Offset(offset int) *PropertyQuery {
	pq.offset = &offset
	return pq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (pq *PropertyQuery) Unique(unique bool) *PropertyQuery {
	pq.unique = &unique
	return pq
}

// Order adds an order step to the query.
func (pq *PropertyQuery) Order(o ...OrderFunc) *PropertyQuery {
	pq.order = append(pq.order, o...)
	return pq
}

// QueryType chains the current query on the "type" edge.
func (pq *PropertyQuery) QueryType() *PropertyTypeQuery {
	query := &PropertyTypeQuery{config: pq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := pq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := pq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(property.Table, property.FieldID, selector),
			sqlgraph.To(propertytype.Table, propertytype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, property.TypeTable, property.TypeColumn),
		)
		fromU = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryResources chains the current query on the "resources" edge.
func (pq *PropertyQuery) QueryResources() *ResourceQuery {
	query := &ResourceQuery{config: pq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := pq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := pq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(property.Table, property.FieldID, selector),
			sqlgraph.To(resource.Table, resource.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, property.ResourcesTable, property.ResourcesColumn),
		)
		fromU = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Property entity from the query.
// Returns a *NotFoundError when no Property was found.
func (pq *PropertyQuery) First(ctx context.Context) (*Property, error) {
	nodes, err := pq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{property.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (pq *PropertyQuery) FirstX(ctx context.Context) *Property {
	node, err := pq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Property ID from the query.
// Returns a *NotFoundError when no Property ID was found.
func (pq *PropertyQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = pq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{property.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (pq *PropertyQuery) FirstIDX(ctx context.Context) int {
	id, err := pq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Property entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Property entity is found.
// Returns a *NotFoundError when no Property entities are found.
func (pq *PropertyQuery) Only(ctx context.Context) (*Property, error) {
	nodes, err := pq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{property.Label}
	default:
		return nil, &NotSingularError{property.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (pq *PropertyQuery) OnlyX(ctx context.Context) *Property {
	node, err := pq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Property ID in the query.
// Returns a *NotSingularError when more than one Property ID is found.
// Returns a *NotFoundError when no entities are found.
func (pq *PropertyQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = pq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{property.Label}
	default:
		err = &NotSingularError{property.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (pq *PropertyQuery) OnlyIDX(ctx context.Context) int {
	id, err := pq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Properties.
func (pq *PropertyQuery) All(ctx context.Context) ([]*Property, error) {
	if err := pq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return pq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (pq *PropertyQuery) AllX(ctx context.Context) []*Property {
	nodes, err := pq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Property IDs.
func (pq *PropertyQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := pq.Select(property.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (pq *PropertyQuery) IDsX(ctx context.Context) []int {
	ids, err := pq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (pq *PropertyQuery) Count(ctx context.Context) (int, error) {
	if err := pq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return pq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (pq *PropertyQuery) CountX(ctx context.Context) int {
	count, err := pq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (pq *PropertyQuery) Exist(ctx context.Context) (bool, error) {
	if err := pq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return pq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (pq *PropertyQuery) ExistX(ctx context.Context) bool {
	exist, err := pq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the PropertyQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (pq *PropertyQuery) Clone() *PropertyQuery {
	if pq == nil {
		return nil
	}
	return &PropertyQuery{
		config:        pq.config,
		limit:         pq.limit,
		offset:        pq.offset,
		order:         append([]OrderFunc{}, pq.order...),
		predicates:    append([]predicate.Property{}, pq.predicates...),
		withType:      pq.withType.Clone(),
		withResources: pq.withResources.Clone(),
		// clone intermediate query.
		sql:    pq.sql.Clone(),
		path:   pq.path,
		unique: pq.unique,
	}
}

// WithType tells the query-builder to eager-load the nodes that are connected to
// the "type" edge. The optional arguments are used to configure the query builder of the edge.
func (pq *PropertyQuery) WithType(opts ...func(*PropertyTypeQuery)) *PropertyQuery {
	query := &PropertyTypeQuery{config: pq.config}
	for _, opt := range opts {
		opt(query)
	}
	pq.withType = query
	return pq
}

// WithResources tells the query-builder to eager-load the nodes that are connected to
// the "resources" edge. The optional arguments are used to configure the query builder of the edge.
func (pq *PropertyQuery) WithResources(opts ...func(*ResourceQuery)) *PropertyQuery {
	query := &ResourceQuery{config: pq.config}
	for _, opt := range opts {
		opt(query)
	}
	pq.withResources = query
	return pq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		IntVal int `json:"intValue" gqlgen:"intValue"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Property.Query().
//		GroupBy(property.FieldIntVal).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (pq *PropertyQuery) GroupBy(field string, fields ...string) *PropertyGroupBy {
	grbuild := &PropertyGroupBy{config: pq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := pq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return pq.sqlQuery(ctx), nil
	}
	grbuild.label = property.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		IntVal int `json:"intValue" gqlgen:"intValue"`
//	}
//
//	client.Property.Query().
//		Select(property.FieldIntVal).
//		Scan(ctx, &v)
//
func (pq *PropertyQuery) Select(fields ...string) *PropertySelect {
	pq.fields = append(pq.fields, fields...)
	selbuild := &PropertySelect{PropertyQuery: pq}
	selbuild.label = property.Label
	selbuild.flds, selbuild.scan = &pq.fields, selbuild.Scan
	return selbuild
}

func (pq *PropertyQuery) prepareQuery(ctx context.Context) error {
	for _, f := range pq.fields {
		if !property.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if pq.path != nil {
		prev, err := pq.path(ctx)
		if err != nil {
			return err
		}
		pq.sql = prev
	}
	if property.Policy == nil {
		return errors.New("ent: uninitialized property.Policy (forgotten import ent/runtime?)")
	}
	if err := property.Policy.EvalQuery(ctx, pq); err != nil {
		return err
	}
	return nil
}

func (pq *PropertyQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Property, error) {
	var (
		nodes       = []*Property{}
		withFKs     = pq.withFKs
		_spec       = pq.querySpec()
		loadedTypes = [2]bool{
			pq.withType != nil,
			pq.withResources != nil,
		}
	)
	if pq.withType != nil || pq.withResources != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, property.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Property).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Property{config: pq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(pq.modifiers) > 0 {
		_spec.Modifiers = pq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, pq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := pq.withType; query != nil {
		if err := pq.loadType(ctx, query, nodes, nil,
			func(n *Property, e *PropertyType) { n.Edges.Type = e }); err != nil {
			return nil, err
		}
	}
	if query := pq.withResources; query != nil {
		if err := pq.loadResources(ctx, query, nodes, nil,
			func(n *Property, e *Resource) { n.Edges.Resources = e }); err != nil {
			return nil, err
		}
	}
	for i := range pq.loadTotal {
		if err := pq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (pq *PropertyQuery) loadType(ctx context.Context, query *PropertyTypeQuery, nodes []*Property, init func(*Property), assign func(*Property, *PropertyType)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*Property)
	for i := range nodes {
		if nodes[i].property_type == nil {
			continue
		}
		fk := *nodes[i].property_type
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	query.Where(propertytype.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "property_type" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (pq *PropertyQuery) loadResources(ctx context.Context, query *ResourceQuery, nodes []*Property, init func(*Property), assign func(*Property, *Resource)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*Property)
	for i := range nodes {
		if nodes[i].resource_properties == nil {
			continue
		}
		fk := *nodes[i].resource_properties
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	query.Where(resource.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "resource_properties" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (pq *PropertyQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := pq.querySpec()
	if len(pq.modifiers) > 0 {
		_spec.Modifiers = pq.modifiers
	}
	_spec.Node.Columns = pq.fields
	if len(pq.fields) > 0 {
		_spec.Unique = pq.unique != nil && *pq.unique
	}
	return sqlgraph.CountNodes(ctx, pq.driver, _spec)
}

func (pq *PropertyQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := pq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (pq *PropertyQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   property.Table,
			Columns: property.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: property.FieldID,
			},
		},
		From:   pq.sql,
		Unique: true,
	}
	if unique := pq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := pq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, property.FieldID)
		for i := range fields {
			if fields[i] != property.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := pq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := pq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := pq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := pq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (pq *PropertyQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(pq.driver.Dialect())
	t1 := builder.Table(property.Table)
	columns := pq.fields
	if len(columns) == 0 {
		columns = property.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if pq.sql != nil {
		selector = pq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if pq.unique != nil && *pq.unique {
		selector.Distinct()
	}
	for _, p := range pq.predicates {
		p(selector)
	}
	for _, p := range pq.order {
		p(selector)
	}
	if offset := pq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := pq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// PropertyGroupBy is the group-by builder for Property entities.
type PropertyGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (pgb *PropertyGroupBy) Aggregate(fns ...AggregateFunc) *PropertyGroupBy {
	pgb.fns = append(pgb.fns, fns...)
	return pgb
}

// Scan applies the group-by query and scans the result into the given value.
func (pgb *PropertyGroupBy) Scan(ctx context.Context, v any) error {
	query, err := pgb.path(ctx)
	if err != nil {
		return err
	}
	pgb.sql = query
	return pgb.sqlScan(ctx, v)
}

func (pgb *PropertyGroupBy) sqlScan(ctx context.Context, v any) error {
	for _, f := range pgb.fields {
		if !property.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := pgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := pgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (pgb *PropertyGroupBy) sqlQuery() *sql.Selector {
	selector := pgb.sql.Select()
	aggregation := make([]string, 0, len(pgb.fns))
	for _, fn := range pgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(pgb.fields)+len(pgb.fns))
		for _, f := range pgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(pgb.fields...)...)
}

// PropertySelect is the builder for selecting fields of Property entities.
type PropertySelect struct {
	*PropertyQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (ps *PropertySelect) Scan(ctx context.Context, v any) error {
	if err := ps.prepareQuery(ctx); err != nil {
		return err
	}
	ps.sql = ps.PropertyQuery.sqlQuery(ctx)
	return ps.sqlScan(ctx, v)
}

func (ps *PropertySelect) sqlScan(ctx context.Context, v any) error {
	rows := &sql.Rows{}
	query, args := ps.sql.Query()
	if err := ps.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
