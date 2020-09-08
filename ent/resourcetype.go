// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"github.com/facebook/ent/dialect/sql"
	"github.com/net-auto/resourceManager/ent/resourcetype"
)

// ResourceType is the model entity for the ResourceType schema.
type ResourceType struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ResourceTypeQuery when eager-loading is set.
	Edges ResourceTypeEdges `json:"edges"`
}

// ResourceTypeEdges holds the relations/edges for other nodes in the graph.
type ResourceTypeEdges struct {
	// PropertyTypes holds the value of the property_types edge.
	PropertyTypes []*PropertyType
	// Pools holds the value of the pools edge.
	Pools []*ResourcePool
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// PropertyTypesOrErr returns the PropertyTypes value or an error if the edge
// was not loaded in eager-loading.
func (e ResourceTypeEdges) PropertyTypesOrErr() ([]*PropertyType, error) {
	if e.loadedTypes[0] {
		return e.PropertyTypes, nil
	}
	return nil, &NotLoadedError{edge: "property_types"}
}

// PoolsOrErr returns the Pools value or an error if the edge
// was not loaded in eager-loading.
func (e ResourceTypeEdges) PoolsOrErr() ([]*ResourcePool, error) {
	if e.loadedTypes[1] {
		return e.Pools, nil
	}
	return nil, &NotLoadedError{edge: "pools"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ResourceType) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullString{}, // name
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ResourceType fields.
func (rt *ResourceType) assignValues(values ...interface{}) error {
	if m, n := len(values), len(resourcetype.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	rt.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[0])
	} else if value.Valid {
		rt.Name = value.String
	}
	return nil
}

// QueryPropertyTypes queries the property_types edge of the ResourceType.
func (rt *ResourceType) QueryPropertyTypes() *PropertyTypeQuery {
	return (&ResourceTypeClient{config: rt.config}).QueryPropertyTypes(rt)
}

// QueryPools queries the pools edge of the ResourceType.
func (rt *ResourceType) QueryPools() *ResourcePoolQuery {
	return (&ResourceTypeClient{config: rt.config}).QueryPools(rt)
}

// Update returns a builder for updating this ResourceType.
// Note that, you need to call ResourceType.Unwrap() before calling this method, if this ResourceType
// was returned from a transaction, and the transaction was committed or rolled back.
func (rt *ResourceType) Update() *ResourceTypeUpdateOne {
	return (&ResourceTypeClient{config: rt.config}).UpdateOne(rt)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (rt *ResourceType) Unwrap() *ResourceType {
	tx, ok := rt.config.driver.(*txDriver)
	if !ok {
		panic("ent: ResourceType is not a transactional entity")
	}
	rt.config.driver = tx.drv
	return rt
}

// String implements the fmt.Stringer.
func (rt *ResourceType) String() string {
	var builder strings.Builder
	builder.WriteString("ResourceType(")
	builder.WriteString(fmt.Sprintf("id=%v", rt.ID))
	builder.WriteString(", name=")
	builder.WriteString(rt.Name)
	builder.WriteByte(')')
	return builder.String()
}

// ResourceTypes is a parsable slice of ResourceType.
type ResourceTypes []*ResourceType

func (rt ResourceTypes) config(cfg config) {
	for _i := range rt {
		rt[_i].config = cfg
	}
}
