// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/facebook/ent/dialect"
	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/schema"
	"github.com/facebookincubator/ent-contrib/entgql"
	"github.com/hashicorp/go-multierror"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
	"github.com/net-auto/resourceManager/ent/poolproperties"
	"github.com/net-auto/resourceManager/ent/property"
	"github.com/net-auto/resourceManager/ent/propertytype"
	"github.com/net-auto/resourceManager/ent/resource"
	"github.com/net-auto/resourceManager/ent/resourcepool"
	"github.com/net-auto/resourceManager/ent/resourcetype"
	"github.com/net-auto/resourceManager/ent/tag"
	"golang.org/x/sync/semaphore"
)

// Noder wraps the basic Node method.
type Noder interface {
	Node(context.Context) (*Node, error)
}

// Node in the graph.
type Node struct {
	ID     int      `json:"id,omitemty"`      // node id.
	Type   string   `json:"type,omitempty"`   // node type.
	Fields []*Field `json:"fields,omitempty"` // node fields.
	Edges  []*Edge  `json:"edges,omitempty"`  // node edges.
}

// Field of a node.
type Field struct {
	Type  string `json:"type,omitempty"`  // field type.
	Name  string `json:"name,omitempty"`  // field name (as in struct).
	Value string `json:"value,omitempty"` // stringified value.
}

// Edges between two nodes.
type Edge struct {
	Type string `json:"type,omitempty"` // edge type.
	Name string `json:"name,omitempty"` // edge name.
	IDs  []int  `json:"ids,omitempty"`  // node ids (where this edge point to).
}

func (as *AllocationStrategy) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     as.ID,
		Type:   "AllocationStrategy",
		Fields: make([]*Field, 4),
		Edges:  make([]*Edge, 1),
	}
	var buf []byte
	if buf, err = json.Marshal(as.Name); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(as.Description); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "string",
		Name:  "description",
		Value: string(buf),
	}
	if buf, err = json.Marshal(as.Lang); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "allocationstrategy.Lang",
		Name:  "lang",
		Value: string(buf),
	}
	if buf, err = json.Marshal(as.Script); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "script",
		Value: string(buf),
	}
	node.Edges[0] = &Edge{
		Type: "ResourcePool",
		Name: "pools",
	}
	node.Edges[0].IDs, err = as.QueryPools().
		Select(resourcepool.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (pp *PoolProperties) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     pp.ID,
		Type:   "PoolProperties",
		Fields: make([]*Field, 0),
		Edges:  make([]*Edge, 3),
	}
	node.Edges[0] = &Edge{
		Type: "ResourcePool",
		Name: "pool",
	}
	node.Edges[0].IDs, err = pp.QueryPool().
		Select(resourcepool.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		Type: "ResourceType",
		Name: "resourceType",
	}
	node.Edges[1].IDs, err = pp.QueryResourceType().
		Select(resourcetype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		Type: "Property",
		Name: "properties",
	}
	node.Edges[2].IDs, err = pp.QueryProperties().
		Select(property.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (pr *Property) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     pr.ID,
		Type:   "Property",
		Fields: make([]*Field, 8),
		Edges:  make([]*Edge, 2),
	}
	var buf []byte
	if buf, err = json.Marshal(pr.IntVal); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "int",
		Name:  "int_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.BoolVal); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "bool",
		Name:  "bool_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.FloatVal); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "float64",
		Name:  "float_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.LatitudeVal); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "float64",
		Name:  "latitude_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.LongitudeVal); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "float64",
		Name:  "longitude_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.RangeFromVal); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "float64",
		Name:  "range_from_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.RangeToVal); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "float64",
		Name:  "range_to_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.StringVal); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "string",
		Name:  "string_val",
		Value: string(buf),
	}
	node.Edges[0] = &Edge{
		Type: "PropertyType",
		Name: "type",
	}
	node.Edges[0].IDs, err = pr.QueryType().
		Select(propertytype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		Type: "Resource",
		Name: "resources",
	}
	node.Edges[1].IDs, err = pr.QueryResources().
		Select(resource.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (pt *PropertyType) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     pt.ID,
		Type:   "PropertyType",
		Fields: make([]*Field, 18),
		Edges:  make([]*Edge, 2),
	}
	var buf []byte
	if buf, err = json.Marshal(pt.Type); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "propertytype.Type",
		Name:  "type",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Name); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.ExternalID); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "external_id",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Index); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "int",
		Name:  "index",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Category); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "category",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.IntVal); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "int",
		Name:  "int_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.BoolVal); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "bool",
		Name:  "bool_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.FloatVal); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "float64",
		Name:  "float_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.LatitudeVal); err != nil {
		return nil, err
	}
	node.Fields[8] = &Field{
		Type:  "float64",
		Name:  "latitude_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.LongitudeVal); err != nil {
		return nil, err
	}
	node.Fields[9] = &Field{
		Type:  "float64",
		Name:  "longitude_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.StringVal); err != nil {
		return nil, err
	}
	node.Fields[10] = &Field{
		Type:  "string",
		Name:  "string_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.RangeFromVal); err != nil {
		return nil, err
	}
	node.Fields[11] = &Field{
		Type:  "float64",
		Name:  "range_from_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.RangeToVal); err != nil {
		return nil, err
	}
	node.Fields[12] = &Field{
		Type:  "float64",
		Name:  "range_to_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.IsInstanceProperty); err != nil {
		return nil, err
	}
	node.Fields[13] = &Field{
		Type:  "bool",
		Name:  "is_instance_property",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Editable); err != nil {
		return nil, err
	}
	node.Fields[14] = &Field{
		Type:  "bool",
		Name:  "editable",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Mandatory); err != nil {
		return nil, err
	}
	node.Fields[15] = &Field{
		Type:  "bool",
		Name:  "mandatory",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Deleted); err != nil {
		return nil, err
	}
	node.Fields[16] = &Field{
		Type:  "bool",
		Name:  "deleted",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.NodeType); err != nil {
		return nil, err
	}
	node.Fields[17] = &Field{
		Type:  "string",
		Name:  "nodeType",
		Value: string(buf),
	}
	node.Edges[0] = &Edge{
		Type: "Property",
		Name: "properties",
	}
	node.Edges[0].IDs, err = pt.QueryProperties().
		Select(property.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		Type: "ResourceType",
		Name: "resource_type",
	}
	node.Edges[1].IDs, err = pt.QueryResourceType().
		Select(resourcetype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (r *Resource) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     r.ID,
		Type:   "Resource",
		Fields: make([]*Field, 4),
		Edges:  make([]*Edge, 3),
	}
	var buf []byte
	if buf, err = json.Marshal(r.Status); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "resource.Status",
		Name:  "status",
		Value: string(buf),
	}
	if buf, err = json.Marshal(r.Description); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "string",
		Name:  "description",
		Value: string(buf),
	}
	if buf, err = json.Marshal(r.AlternateID); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "map[string]interface {}",
		Name:  "alternate_id",
		Value: string(buf),
	}
	if buf, err = json.Marshal(r.UpdatedAt); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "time.Time",
		Name:  "updated_at",
		Value: string(buf),
	}
	node.Edges[0] = &Edge{
		Type: "ResourcePool",
		Name: "pool",
	}
	node.Edges[0].IDs, err = r.QueryPool().
		Select(resourcepool.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		Type: "Property",
		Name: "properties",
	}
	node.Edges[1].IDs, err = r.QueryProperties().
		Select(property.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		Type: "ResourcePool",
		Name: "nested_pool",
	}
	node.Edges[2].IDs, err = r.QueryNestedPool().
		Select(resourcepool.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (rp *ResourcePool) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     rp.ID,
		Type:   "ResourcePool",
		Fields: make([]*Field, 4),
		Edges:  make([]*Edge, 6),
	}
	var buf []byte
	if buf, err = json.Marshal(rp.Name); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(rp.Description); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "string",
		Name:  "description",
		Value: string(buf),
	}
	if buf, err = json.Marshal(rp.PoolType); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "resourcepool.PoolType",
		Name:  "pool_type",
		Value: string(buf),
	}
	if buf, err = json.Marshal(rp.DealocationSafetyPeriod); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "int",
		Name:  "dealocation_safety_period",
		Value: string(buf),
	}
	node.Edges[0] = &Edge{
		Type: "ResourceType",
		Name: "resource_type",
	}
	node.Edges[0].IDs, err = rp.QueryResourceType().
		Select(resourcetype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		Type: "Tag",
		Name: "tags",
	}
	node.Edges[1].IDs, err = rp.QueryTags().
		Select(tag.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		Type: "Resource",
		Name: "claims",
	}
	node.Edges[2].IDs, err = rp.QueryClaims().
		Select(resource.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		Type: "PoolProperties",
		Name: "poolProperties",
	}
	node.Edges[3].IDs, err = rp.QueryPoolProperties().
		Select(poolproperties.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[4] = &Edge{
		Type: "AllocationStrategy",
		Name: "allocation_strategy",
	}
	node.Edges[4].IDs, err = rp.QueryAllocationStrategy().
		Select(allocationstrategy.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[5] = &Edge{
		Type: "Resource",
		Name: "parent_resource",
	}
	node.Edges[5].IDs, err = rp.QueryParentResource().
		Select(resource.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (rt *ResourceType) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     rt.ID,
		Type:   "ResourceType",
		Fields: make([]*Field, 1),
		Edges:  make([]*Edge, 3),
	}
	var buf []byte
	if buf, err = json.Marshal(rt.Name); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	node.Edges[0] = &Edge{
		Type: "PropertyType",
		Name: "property_types",
	}
	node.Edges[0].IDs, err = rt.QueryPropertyTypes().
		Select(propertytype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		Type: "ResourcePool",
		Name: "pools",
	}
	node.Edges[1].IDs, err = rt.QueryPools().
		Select(resourcepool.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		Type: "PoolProperties",
		Name: "pool_properties",
	}
	node.Edges[2].IDs, err = rt.QueryPoolProperties().
		Select(poolproperties.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (t *Tag) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     t.ID,
		Type:   "Tag",
		Fields: make([]*Field, 1),
		Edges:  make([]*Edge, 1),
	}
	var buf []byte
	if buf, err = json.Marshal(t.Tag); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "string",
		Name:  "tag",
		Value: string(buf),
	}
	node.Edges[0] = &Edge{
		Type: "ResourcePool",
		Name: "pools",
	}
	node.Edges[0].IDs, err = t.QueryPools().
		Select(resourcepool.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (c *Client) Node(ctx context.Context, id int) (*Node, error) {
	n, err := c.Noder(ctx, id)
	if err != nil {
		return nil, err
	}
	return n.Node(ctx)
}

var errNodeInvalidID = &NotFoundError{"node"}

// NodeOption allows configuring the Noder execution using functional options.
type NodeOption func(*NodeOptions)

// WithNodeType sets the Type of the node (i.e. the table to query).
// If was not provided, the table will be derived from the universal-id
// configuration as described in: https://entgo.io/docs/migrate/#universal-ids.
func WithNodeType(t string) NodeOption {
	return func(o *NodeOptions) {
		o.Type = t
	}
}

// NodeOptions holds the configuration for Noder execution.
type NodeOptions struct {
	// Type of the node (schema table).
	Type string
}

// Noder returns a Node by its id. If the NodeType was not provided, it will
// be derived from the id value according to the universal-id configuration.
//
//		c.Noder(ctx, id)
//		c.Noder(ctx, id, ent.WithNodeType(pet.Table))
//
func (c *Client) Noder(ctx context.Context, id int, opts ...NodeOption) (_ Noder, err error) {
	defer func() {
		if IsNotFound(err) {
			err = multierror.Append(err, entgql.ErrNodeNotFound(id))
		}
	}()
	options := &NodeOptions{}
	for _, opt := range opts {
		opt(options)
	}
	if options.Type == "" {
		options.Type, err = c.tables.nodeType(ctx, c.driver, id)
		if err != nil {
			return nil, err
		}
	}
	return c.noder(ctx, options.Type, id)
}

func (c *Client) noder(ctx context.Context, tbl string, id int) (Noder, error) {
	switch tbl {
	case allocationstrategy.Table:
		n, err := c.AllocationStrategy.Query().
			Where(allocationstrategy.ID(id)).
			CollectFields(ctx, "AllocationStrategy").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case poolproperties.Table:
		n, err := c.PoolProperties.Query().
			Where(poolproperties.ID(id)).
			CollectFields(ctx, "PoolProperties").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case property.Table:
		n, err := c.Property.Query().
			Where(property.ID(id)).
			CollectFields(ctx, "Property").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case propertytype.Table:
		n, err := c.PropertyType.Query().
			Where(propertytype.ID(id)).
			CollectFields(ctx, "PropertyType").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case resource.Table:
		n, err := c.Resource.Query().
			Where(resource.ID(id)).
			CollectFields(ctx, "Resource").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case resourcepool.Table:
		n, err := c.ResourcePool.Query().
			Where(resourcepool.ID(id)).
			CollectFields(ctx, "ResourcePool").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case resourcetype.Table:
		n, err := c.ResourceType.Query().
			Where(resourcetype.ID(id)).
			CollectFields(ctx, "ResourceType").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case tag.Table:
		n, err := c.Tag.Query().
			Where(tag.ID(id)).
			CollectFields(ctx, "Tag").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	default:
		return nil, fmt.Errorf("cannot resolve Noder from table %q: %w", tbl, errNodeInvalidID)
	}
}

type tables struct {
	once  sync.Once
	sem   *semaphore.Weighted
	value atomic.Value
}

func (t *tables) nodeType(ctx context.Context, drv dialect.Driver, id int) (string, error) {
	tables, err := t.Load(ctx, drv)
	if err != nil {
		return "", err
	}
	idx := id / (1<<32 - 1)
	if idx < 0 || idx >= len(tables) {
		return "", fmt.Errorf("cannot resolve table from id %v: %w", id, errNodeInvalidID)
	}
	return tables[idx], nil
}

func (t *tables) Load(ctx context.Context, drv dialect.Driver) ([]string, error) {
	if tables := t.value.Load(); tables != nil {
		return tables.([]string), nil
	}
	t.once.Do(func() { t.sem = semaphore.NewWeighted(1) })
	if err := t.sem.Acquire(ctx, 1); err != nil {
		return nil, err
	}
	defer t.sem.Release(1)
	if tables := t.value.Load(); tables != nil {
		return tables.([]string), nil
	}
	tables, err := t.load(ctx, drv)
	if err == nil {
		t.value.Store(tables)
	}
	return tables, err
}

func (tables) load(ctx context.Context, drv dialect.Driver) ([]string, error) {
	rows := &sql.Rows{}
	query, args := sql.Dialect(drv.Dialect()).
		Select("type").
		From(sql.Table(schema.TypeTable)).
		OrderBy(sql.Asc("id")).
		Query()
	if err := drv.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []string
	return tables, sql.ScanSlice(rows, &tables)
}
