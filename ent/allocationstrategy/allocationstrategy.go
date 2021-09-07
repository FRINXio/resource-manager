// Code generated by entc, DO NOT EDIT.

package allocationstrategy

import (
	"fmt"
	"io"
	"strconv"

	"github.com/facebook/ent"
)

const (
	// Label holds the string label denoting the allocationstrategy type in the database.
	Label = "allocation_strategy"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldLang holds the string denoting the lang field in the database.
	FieldLang = "lang"
	// FieldScript holds the string denoting the script field in the database.
	FieldScript = "script"

	// EdgePools holds the string denoting the pools edge name in mutations.
	EdgePools = "pools"

	// Table holds the table name of the allocationstrategy in the database.
	Table = "allocation_strategies"
	// PoolsTable is the table the holds the pools relation/edge.
	PoolsTable = "resource_pools"
	// PoolsInverseTable is the table name for the ResourcePool entity.
	// It exists in this package in order to avoid circular dependency with the "resourcepool" package.
	PoolsInverseTable = "resource_pools"
	// PoolsColumn is the table column denoting the pools relation/edge.
	PoolsColumn = "resource_pool_allocation_strategy"
)

// Columns holds all SQL columns for allocationstrategy fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldDescription,
	FieldLang,
	FieldScript,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/net-auto/resourceManager/ent/runtime"
//
var (
	Hooks  [1]ent.Hook
	Policy ent.Policy
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// ScriptValidator is a validator for the "script" field. It is called by the builders before save.
	ScriptValidator func(string) error
)

// Lang defines the type for the lang enum field.
type Lang string

// LangJs is the default Lang.
const DefaultLang = LangJs

// Lang values.
const (
	LangPy Lang = "py"
	LangJs Lang = "js"
	LangGo Lang = "go"
)

func (l Lang) String() string {
	return string(l)
}

// LangValidator is a validator for the "lang" field enum values. It is called by the builders before save.
func LangValidator(l Lang) error {
	switch l {
	case LangPy, LangJs, LangGo:
		return nil
	default:
		return fmt.Errorf("allocationstrategy: invalid enum value for lang field: %q", l)
	}
}

// MarshalGQL implements graphql.Marshaler interface.
func (l Lang) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(l.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (l *Lang) UnmarshalGQL(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("enum %T must be a string", val)
	}
	*l = Lang(str)
	if err := LangValidator(*l); err != nil {
		return fmt.Errorf("%s is not a valid Lang", str)
	}
	return nil
}
