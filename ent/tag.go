// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/net-auto/resourceManager/ent/tag"
)

// Tag is the model entity for the Tag schema.
type Tag struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Tag holds the value of the "tag" field.
	Tag string `json:"tag,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the TagQuery when eager-loading is set.
	Edges TagEdges `json:"edges"`
}

// TagEdges holds the relations/edges for other nodes in the graph.
type TagEdges struct {
	// Pools holds the value of the pools edge.
	Pools []*ResourcePool `json:"pools,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
	// totalCount holds the count of the edges above.
	totalCount [1]map[string]int

	namedPools map[string][]*ResourcePool
}

// PoolsOrErr returns the Pools value or an error if the edge
// was not loaded in eager-loading.
func (e TagEdges) PoolsOrErr() ([]*ResourcePool, error) {
	if e.loadedTypes[0] {
		return e.Pools, nil
	}
	return nil, &NotLoadedError{edge: "pools"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Tag) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case tag.FieldID:
			values[i] = new(sql.NullInt64)
		case tag.FieldTag:
			values[i] = new(sql.NullString)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Tag", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Tag fields.
func (t *Tag) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case tag.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			t.ID = int(value.Int64)
		case tag.FieldTag:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field tag", values[i])
			} else if value.Valid {
				t.Tag = value.String
			}
		}
	}
	return nil
}

// QueryPools queries the "pools" edge of the Tag entity.
func (t *Tag) QueryPools() *ResourcePoolQuery {
	return (&TagClient{config: t.config}).QueryPools(t)
}

// Update returns a builder for updating this Tag.
// Note that you need to call Tag.Unwrap() before calling this method if this Tag
// was returned from a transaction, and the transaction was committed or rolled back.
func (t *Tag) Update() *TagUpdateOne {
	return (&TagClient{config: t.config}).UpdateOne(t)
}

// Unwrap unwraps the Tag entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (t *Tag) Unwrap() *Tag {
	_tx, ok := t.config.driver.(*txDriver)
	if !ok {
		panic("ent: Tag is not a transactional entity")
	}
	t.config.driver = _tx.drv
	return t
}

// String implements the fmt.Stringer.
func (t *Tag) String() string {
	var builder strings.Builder
	builder.WriteString("Tag(")
	builder.WriteString(fmt.Sprintf("id=%v, ", t.ID))
	builder.WriteString("tag=")
	builder.WriteString(t.Tag)
	builder.WriteByte(')')
	return builder.String()
}

// NamedPools returns the Pools named value or an error if the edge was not
// loaded in eager-loading with this name.
func (t *Tag) NamedPools(name string) ([]*ResourcePool, error) {
	if t.Edges.namedPools == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := t.Edges.namedPools[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (t *Tag) appendNamedPools(name string, edges ...*ResourcePool) {
	if t.Edges.namedPools == nil {
		t.Edges.namedPools = make(map[string][]*ResourcePool)
	}
	if len(edges) == 0 {
		t.Edges.namedPools[name] = []*ResourcePool{}
	} else {
		t.Edges.namedPools[name] = append(t.Edges.namedPools[name], edges...)
	}
}

// Tags is a parsable slice of Tag.
type Tags []*Tag

func (t Tags) config(cfg config) {
	for _i := range t {
		t[_i].config = cfg
	}
}
