// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/net-auto/resourceManager/ent/poolproperties"
	"github.com/net-auto/resourceManager/ent/resourcepool"
)

// PoolProperties is the model entity for the PoolProperties schema.
type PoolProperties struct {
	config
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PoolPropertiesQuery when eager-loading is set.
	Edges                         PoolPropertiesEdges `json:"edges"`
	resource_pool_pool_properties *int
}

// PoolPropertiesEdges holds the relations/edges for other nodes in the graph.
type PoolPropertiesEdges struct {
	// Pool holds the value of the pool edge.
	Pool *ResourcePool `json:"pool,omitempty"`
	// ResourceType holds the value of the resourceType edge.
	ResourceType []*ResourceType `json:"resourceType,omitempty"`
	// Properties holds the value of the properties edge.
	Properties []*Property `json:"properties,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
	// totalCount holds the count of the edges above.
	totalCount [3]map[string]int

	namedResourceType map[string][]*ResourceType
	namedProperties   map[string][]*Property
}

// PoolOrErr returns the Pool value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PoolPropertiesEdges) PoolOrErr() (*ResourcePool, error) {
	if e.loadedTypes[0] {
		if e.Pool == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: resourcepool.Label}
		}
		return e.Pool, nil
	}
	return nil, &NotLoadedError{edge: "pool"}
}

// ResourceTypeOrErr returns the ResourceType value or an error if the edge
// was not loaded in eager-loading.
func (e PoolPropertiesEdges) ResourceTypeOrErr() ([]*ResourceType, error) {
	if e.loadedTypes[1] {
		return e.ResourceType, nil
	}
	return nil, &NotLoadedError{edge: "resourceType"}
}

// PropertiesOrErr returns the Properties value or an error if the edge
// was not loaded in eager-loading.
func (e PoolPropertiesEdges) PropertiesOrErr() ([]*Property, error) {
	if e.loadedTypes[2] {
		return e.Properties, nil
	}
	return nil, &NotLoadedError{edge: "properties"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*PoolProperties) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case poolproperties.FieldID:
			values[i] = new(sql.NullInt64)
		case poolproperties.ForeignKeys[0]: // resource_pool_pool_properties
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type PoolProperties", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the PoolProperties fields.
func (pp *PoolProperties) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case poolproperties.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			pp.ID = int(value.Int64)
		case poolproperties.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field resource_pool_pool_properties", value)
			} else if value.Valid {
				pp.resource_pool_pool_properties = new(int)
				*pp.resource_pool_pool_properties = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryPool queries the "pool" edge of the PoolProperties entity.
func (pp *PoolProperties) QueryPool() *ResourcePoolQuery {
	return (&PoolPropertiesClient{config: pp.config}).QueryPool(pp)
}

// QueryResourceType queries the "resourceType" edge of the PoolProperties entity.
func (pp *PoolProperties) QueryResourceType() *ResourceTypeQuery {
	return (&PoolPropertiesClient{config: pp.config}).QueryResourceType(pp)
}

// QueryProperties queries the "properties" edge of the PoolProperties entity.
func (pp *PoolProperties) QueryProperties() *PropertyQuery {
	return (&PoolPropertiesClient{config: pp.config}).QueryProperties(pp)
}

// Update returns a builder for updating this PoolProperties.
// Note that you need to call PoolProperties.Unwrap() before calling this method if this PoolProperties
// was returned from a transaction, and the transaction was committed or rolled back.
func (pp *PoolProperties) Update() *PoolPropertiesUpdateOne {
	return (&PoolPropertiesClient{config: pp.config}).UpdateOne(pp)
}

// Unwrap unwraps the PoolProperties entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pp *PoolProperties) Unwrap() *PoolProperties {
	_tx, ok := pp.config.driver.(*txDriver)
	if !ok {
		panic("ent: PoolProperties is not a transactional entity")
	}
	pp.config.driver = _tx.drv
	return pp
}

// String implements the fmt.Stringer.
func (pp *PoolProperties) String() string {
	var builder strings.Builder
	builder.WriteString("PoolProperties(")
	builder.WriteString(fmt.Sprintf("id=%v", pp.ID))
	builder.WriteByte(')')
	return builder.String()
}

// NamedResourceType returns the ResourceType named value or an error if the edge was not
// loaded in eager-loading with this name.
func (pp *PoolProperties) NamedResourceType(name string) ([]*ResourceType, error) {
	if pp.Edges.namedResourceType == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := pp.Edges.namedResourceType[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (pp *PoolProperties) appendNamedResourceType(name string, edges ...*ResourceType) {
	if pp.Edges.namedResourceType == nil {
		pp.Edges.namedResourceType = make(map[string][]*ResourceType)
	}
	if len(edges) == 0 {
		pp.Edges.namedResourceType[name] = []*ResourceType{}
	} else {
		pp.Edges.namedResourceType[name] = append(pp.Edges.namedResourceType[name], edges...)
	}
}

// NamedProperties returns the Properties named value or an error if the edge was not
// loaded in eager-loading with this name.
func (pp *PoolProperties) NamedProperties(name string) ([]*Property, error) {
	if pp.Edges.namedProperties == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := pp.Edges.namedProperties[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (pp *PoolProperties) appendNamedProperties(name string, edges ...*Property) {
	if pp.Edges.namedProperties == nil {
		pp.Edges.namedProperties = make(map[string][]*Property)
	}
	if len(edges) == 0 {
		pp.Edges.namedProperties[name] = []*Property{}
	} else {
		pp.Edges.namedProperties[name] = append(pp.Edges.namedProperties[name], edges...)
	}
}

// PoolPropertiesSlice is a parsable slice of PoolProperties.
type PoolPropertiesSlice []*PoolProperties

func (pp PoolPropertiesSlice) config(cfg config) {
	for _i := range pp {
		pp[_i].config = cfg
	}
}
