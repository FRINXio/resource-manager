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
	"github.com/net-auto/resourceManager/ent/predicate"
	"github.com/net-auto/resourceManager/ent/property"
	"github.com/net-auto/resourceManager/ent/propertytype"
	"github.com/net-auto/resourceManager/ent/resourcetype"
)

// PropertyTypeQuery is the builder for querying PropertyType entities.
type PropertyTypeQuery struct {
	config
	limit               *int
	offset              *int
	unique              *bool
	order               []OrderFunc
	fields              []string
	predicates          []predicate.PropertyType
	withProperties      *PropertyQuery
	withResourceType    *ResourceTypeQuery
	withFKs             bool
	modifiers           []func(*sql.Selector)
	loadTotal           []func(context.Context, []*PropertyType) error
	withNamedProperties map[string]*PropertyQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the PropertyTypeQuery builder.
func (ptq *PropertyTypeQuery) Where(ps ...predicate.PropertyType) *PropertyTypeQuery {
	ptq.predicates = append(ptq.predicates, ps...)
	return ptq
}

// Limit adds a limit step to the query.
func (ptq *PropertyTypeQuery) Limit(limit int) *PropertyTypeQuery {
	ptq.limit = &limit
	return ptq
}

// Offset adds an offset step to the query.
func (ptq *PropertyTypeQuery) Offset(offset int) *PropertyTypeQuery {
	ptq.offset = &offset
	return ptq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (ptq *PropertyTypeQuery) Unique(unique bool) *PropertyTypeQuery {
	ptq.unique = &unique
	return ptq
}

// Order adds an order step to the query.
func (ptq *PropertyTypeQuery) Order(o ...OrderFunc) *PropertyTypeQuery {
	ptq.order = append(ptq.order, o...)
	return ptq
}

// QueryProperties chains the current query on the "properties" edge.
func (ptq *PropertyTypeQuery) QueryProperties() *PropertyQuery {
	query := &PropertyQuery{config: ptq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := ptq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(propertytype.Table, propertytype.FieldID, selector),
			sqlgraph.To(property.Table, property.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, propertytype.PropertiesTable, propertytype.PropertiesColumn),
		)
		fromU = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryResourceType chains the current query on the "resource_type" edge.
func (ptq *PropertyTypeQuery) QueryResourceType() *ResourceTypeQuery {
	query := &ResourceTypeQuery{config: ptq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := ptq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(propertytype.Table, propertytype.FieldID, selector),
			sqlgraph.To(resourcetype.Table, resourcetype.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, propertytype.ResourceTypeTable, propertytype.ResourceTypeColumn),
		)
		fromU = sqlgraph.SetNeighbors(ptq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first PropertyType entity from the query.
// Returns a *NotFoundError when no PropertyType was found.
func (ptq *PropertyTypeQuery) First(ctx context.Context) (*PropertyType, error) {
	nodes, err := ptq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{propertytype.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (ptq *PropertyTypeQuery) FirstX(ctx context.Context) *PropertyType {
	node, err := ptq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first PropertyType ID from the query.
// Returns a *NotFoundError when no PropertyType ID was found.
func (ptq *PropertyTypeQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ptq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{propertytype.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (ptq *PropertyTypeQuery) FirstIDX(ctx context.Context) int {
	id, err := ptq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single PropertyType entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one PropertyType entity is found.
// Returns a *NotFoundError when no PropertyType entities are found.
func (ptq *PropertyTypeQuery) Only(ctx context.Context) (*PropertyType, error) {
	nodes, err := ptq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{propertytype.Label}
	default:
		return nil, &NotSingularError{propertytype.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (ptq *PropertyTypeQuery) OnlyX(ctx context.Context) *PropertyType {
	node, err := ptq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only PropertyType ID in the query.
// Returns a *NotSingularError when more than one PropertyType ID is found.
// Returns a *NotFoundError when no entities are found.
func (ptq *PropertyTypeQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ptq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{propertytype.Label}
	default:
		err = &NotSingularError{propertytype.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (ptq *PropertyTypeQuery) OnlyIDX(ctx context.Context) int {
	id, err := ptq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of PropertyTypes.
func (ptq *PropertyTypeQuery) All(ctx context.Context) ([]*PropertyType, error) {
	if err := ptq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return ptq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (ptq *PropertyTypeQuery) AllX(ctx context.Context) []*PropertyType {
	nodes, err := ptq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of PropertyType IDs.
func (ptq *PropertyTypeQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := ptq.Select(propertytype.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (ptq *PropertyTypeQuery) IDsX(ctx context.Context) []int {
	ids, err := ptq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (ptq *PropertyTypeQuery) Count(ctx context.Context) (int, error) {
	if err := ptq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return ptq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (ptq *PropertyTypeQuery) CountX(ctx context.Context) int {
	count, err := ptq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (ptq *PropertyTypeQuery) Exist(ctx context.Context) (bool, error) {
	if err := ptq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return ptq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (ptq *PropertyTypeQuery) ExistX(ctx context.Context) bool {
	exist, err := ptq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the PropertyTypeQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (ptq *PropertyTypeQuery) Clone() *PropertyTypeQuery {
	if ptq == nil {
		return nil
	}
	return &PropertyTypeQuery{
		config:           ptq.config,
		limit:            ptq.limit,
		offset:           ptq.offset,
		order:            append([]OrderFunc{}, ptq.order...),
		predicates:       append([]predicate.PropertyType{}, ptq.predicates...),
		withProperties:   ptq.withProperties.Clone(),
		withResourceType: ptq.withResourceType.Clone(),
		// clone intermediate query.
		sql:    ptq.sql.Clone(),
		path:   ptq.path,
		unique: ptq.unique,
	}
}

// WithProperties tells the query-builder to eager-load the nodes that are connected to
// the "properties" edge. The optional arguments are used to configure the query builder of the edge.
func (ptq *PropertyTypeQuery) WithProperties(opts ...func(*PropertyQuery)) *PropertyTypeQuery {
	query := &PropertyQuery{config: ptq.config}
	for _, opt := range opts {
		opt(query)
	}
	ptq.withProperties = query
	return ptq
}

// WithResourceType tells the query-builder to eager-load the nodes that are connected to
// the "resource_type" edge. The optional arguments are used to configure the query builder of the edge.
func (ptq *PropertyTypeQuery) WithResourceType(opts ...func(*ResourceTypeQuery)) *PropertyTypeQuery {
	query := &ResourceTypeQuery{config: ptq.config}
	for _, opt := range opts {
		opt(query)
	}
	ptq.withResourceType = query
	return ptq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Type propertytype.Type `json:"type,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.PropertyType.Query().
//		GroupBy(propertytype.FieldType).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (ptq *PropertyTypeQuery) GroupBy(field string, fields ...string) *PropertyTypeGroupBy {
	grbuild := &PropertyTypeGroupBy{config: ptq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := ptq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return ptq.sqlQuery(ctx), nil
	}
	grbuild.label = propertytype.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Type propertytype.Type `json:"type,omitempty"`
//	}
//
//	client.PropertyType.Query().
//		Select(propertytype.FieldType).
//		Scan(ctx, &v)
func (ptq *PropertyTypeQuery) Select(fields ...string) *PropertyTypeSelect {
	ptq.fields = append(ptq.fields, fields...)
	selbuild := &PropertyTypeSelect{PropertyTypeQuery: ptq}
	selbuild.label = propertytype.Label
	selbuild.flds, selbuild.scan = &ptq.fields, selbuild.Scan
	return selbuild
}

func (ptq *PropertyTypeQuery) prepareQuery(ctx context.Context) error {
	for _, f := range ptq.fields {
		if !propertytype.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if ptq.path != nil {
		prev, err := ptq.path(ctx)
		if err != nil {
			return err
		}
		ptq.sql = prev
	}
	if propertytype.Policy == nil {
		return errors.New("ent: uninitialized propertytype.Policy (forgotten import ent/runtime?)")
	}
	if err := propertytype.Policy.EvalQuery(ctx, ptq); err != nil {
		return err
	}
	return nil
}

func (ptq *PropertyTypeQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*PropertyType, error) {
	var (
		nodes       = []*PropertyType{}
		withFKs     = ptq.withFKs
		_spec       = ptq.querySpec()
		loadedTypes = [2]bool{
			ptq.withProperties != nil,
			ptq.withResourceType != nil,
		}
	)
	if ptq.withResourceType != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, propertytype.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*PropertyType).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &PropertyType{config: ptq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(ptq.modifiers) > 0 {
		_spec.Modifiers = ptq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, ptq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := ptq.withProperties; query != nil {
		if err := ptq.loadProperties(ctx, query, nodes,
			func(n *PropertyType) { n.Edges.Properties = []*Property{} },
			func(n *PropertyType, e *Property) { n.Edges.Properties = append(n.Edges.Properties, e) }); err != nil {
			return nil, err
		}
	}
	if query := ptq.withResourceType; query != nil {
		if err := ptq.loadResourceType(ctx, query, nodes, nil,
			func(n *PropertyType, e *ResourceType) { n.Edges.ResourceType = e }); err != nil {
			return nil, err
		}
	}
	for name, query := range ptq.withNamedProperties {
		if err := ptq.loadProperties(ctx, query, nodes,
			func(n *PropertyType) { n.appendNamedProperties(name) },
			func(n *PropertyType, e *Property) { n.appendNamedProperties(name, e) }); err != nil {
			return nil, err
		}
	}
	for i := range ptq.loadTotal {
		if err := ptq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (ptq *PropertyTypeQuery) loadProperties(ctx context.Context, query *PropertyQuery, nodes []*PropertyType, init func(*PropertyType), assign func(*PropertyType, *Property)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*PropertyType)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.Property(func(s *sql.Selector) {
		s.Where(sql.InValues(propertytype.PropertiesColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.property_type
		if fk == nil {
			return fmt.Errorf(`foreign-key "property_type" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "property_type" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (ptq *PropertyTypeQuery) loadResourceType(ctx context.Context, query *ResourceTypeQuery, nodes []*PropertyType, init func(*PropertyType), assign func(*PropertyType, *ResourceType)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*PropertyType)
	for i := range nodes {
		if nodes[i].resource_type_property_types == nil {
			continue
		}
		fk := *nodes[i].resource_type_property_types
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	query.Where(resourcetype.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "resource_type_property_types" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (ptq *PropertyTypeQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := ptq.querySpec()
	if len(ptq.modifiers) > 0 {
		_spec.Modifiers = ptq.modifiers
	}
	_spec.Node.Columns = ptq.fields
	if len(ptq.fields) > 0 {
		_spec.Unique = ptq.unique != nil && *ptq.unique
	}
	return sqlgraph.CountNodes(ctx, ptq.driver, _spec)
}

func (ptq *PropertyTypeQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := ptq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (ptq *PropertyTypeQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   propertytype.Table,
			Columns: propertytype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: propertytype.FieldID,
			},
		},
		From:   ptq.sql,
		Unique: true,
	}
	if unique := ptq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := ptq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, propertytype.FieldID)
		for i := range fields {
			if fields[i] != propertytype.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := ptq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := ptq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := ptq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := ptq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (ptq *PropertyTypeQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(ptq.driver.Dialect())
	t1 := builder.Table(propertytype.Table)
	columns := ptq.fields
	if len(columns) == 0 {
		columns = propertytype.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if ptq.sql != nil {
		selector = ptq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if ptq.unique != nil && *ptq.unique {
		selector.Distinct()
	}
	for _, p := range ptq.predicates {
		p(selector)
	}
	for _, p := range ptq.order {
		p(selector)
	}
	if offset := ptq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := ptq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// WithNamedProperties tells the query-builder to eager-load the nodes that are connected to the "properties"
// edge with the given name. The optional arguments are used to configure the query builder of the edge.
func (ptq *PropertyTypeQuery) WithNamedProperties(name string, opts ...func(*PropertyQuery)) *PropertyTypeQuery {
	query := &PropertyQuery{config: ptq.config}
	for _, opt := range opts {
		opt(query)
	}
	if ptq.withNamedProperties == nil {
		ptq.withNamedProperties = make(map[string]*PropertyQuery)
	}
	ptq.withNamedProperties[name] = query
	return ptq
}

// PropertyTypeGroupBy is the group-by builder for PropertyType entities.
type PropertyTypeGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (ptgb *PropertyTypeGroupBy) Aggregate(fns ...AggregateFunc) *PropertyTypeGroupBy {
	ptgb.fns = append(ptgb.fns, fns...)
	return ptgb
}

// Scan applies the group-by query and scans the result into the given value.
func (ptgb *PropertyTypeGroupBy) Scan(ctx context.Context, v any) error {
	query, err := ptgb.path(ctx)
	if err != nil {
		return err
	}
	ptgb.sql = query
	return ptgb.sqlScan(ctx, v)
}

func (ptgb *PropertyTypeGroupBy) sqlScan(ctx context.Context, v any) error {
	for _, f := range ptgb.fields {
		if !propertytype.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := ptgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ptgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ptgb *PropertyTypeGroupBy) sqlQuery() *sql.Selector {
	selector := ptgb.sql.Select()
	aggregation := make([]string, 0, len(ptgb.fns))
	for _, fn := range ptgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(ptgb.fields)+len(ptgb.fns))
		for _, f := range ptgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(ptgb.fields...)...)
}

// PropertyTypeSelect is the builder for selecting fields of PropertyType entities.
type PropertyTypeSelect struct {
	*PropertyTypeQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (pts *PropertyTypeSelect) Scan(ctx context.Context, v any) error {
	if err := pts.prepareQuery(ctx); err != nil {
		return err
	}
	pts.sql = pts.PropertyTypeQuery.sqlQuery(ctx)
	return pts.sqlScan(ctx, v)
}

func (pts *PropertyTypeSelect) sqlScan(ctx context.Context, v any) error {
	rows := &sql.Rows{}
	query, args := pts.sql.Query()
	if err := pts.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
