// Code generated by entc, DO NOT EDIT.

package resourcepool

import (
	"fmt"
	"io"
	"strconv"
)

const (
	// Label holds the string label denoting the resourcepool type in the database.
	Label = "resource_pool"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldPoolType holds the string denoting the pool_type field in the database.
	FieldPoolType = "pool_type"
	// FieldDealocationSafetyPeriod holds the string denoting the dealocation_safety_period field in the database.
	FieldDealocationSafetyPeriod = "dealocation_safety_period"

	// EdgeResourceType holds the string denoting the resource_type edge name in mutations.
	EdgeResourceType = "resource_type"
	// EdgeTags holds the string denoting the tags edge name in mutations.
	EdgeTags = "tags"
	// EdgeClaims holds the string denoting the claims edge name in mutations.
	EdgeClaims = "claims"
	// EdgeAllocationStrategy holds the string denoting the allocation_strategy edge name in mutations.
	EdgeAllocationStrategy = "allocation_strategy"
	// EdgeParentResource holds the string denoting the parent_resource edge name in mutations.
	EdgeParentResource = "parent_resource"

	// Table holds the table name of the resourcepool in the database.
	Table = "resource_pools"
	// ResourceTypeTable is the table the holds the resource_type relation/edge.
	ResourceTypeTable = "resource_pools"
	// ResourceTypeInverseTable is the table name for the ResourceType entity.
	// It exists in this package in order to avoid circular dependency with the "resourcetype" package.
	ResourceTypeInverseTable = "resource_types"
	// ResourceTypeColumn is the table column denoting the resource_type relation/edge.
	ResourceTypeColumn = "resource_type_pools"
	// TagsTable is the table the holds the tags relation/edge. The primary key declared below.
	TagsTable = "tag_pools"
	// TagsInverseTable is the table name for the Tag entity.
	// It exists in this package in order to avoid circular dependency with the "tag" package.
	TagsInverseTable = "tags"
	// ClaimsTable is the table the holds the claims relation/edge.
	ClaimsTable = "resources"
	// ClaimsInverseTable is the table name for the Resource entity.
	// It exists in this package in order to avoid circular dependency with the "resource" package.
	ClaimsInverseTable = "resources"
	// ClaimsColumn is the table column denoting the claims relation/edge.
	ClaimsColumn = "resource_pool_claims"
	// AllocationStrategyTable is the table the holds the allocation_strategy relation/edge.
	AllocationStrategyTable = "resource_pools"
	// AllocationStrategyInverseTable is the table name for the AllocationStrategy entity.
	// It exists in this package in order to avoid circular dependency with the "allocationstrategy" package.
	AllocationStrategyInverseTable = "allocation_strategies"
	// AllocationStrategyColumn is the table column denoting the allocation_strategy relation/edge.
	AllocationStrategyColumn = "resource_pool_allocation_strategy"
	// ParentResourceTable is the table the holds the parent_resource relation/edge.
	ParentResourceTable = "resource_pools"
	// ParentResourceInverseTable is the table name for the Resource entity.
	// It exists in this package in order to avoid circular dependency with the "resource" package.
	ParentResourceInverseTable = "resources"
	// ParentResourceColumn is the table column denoting the parent_resource relation/edge.
	ParentResourceColumn = "resource_nested_pool"
)

// Columns holds all SQL columns for resourcepool fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldDescription,
	FieldPoolType,
	FieldDealocationSafetyPeriod,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the ResourcePool type.
var ForeignKeys = []string{
	"resource_nested_pool",
	"resource_pool_allocation_strategy",
	"resource_type_pools",
}

var (
	// TagsPrimaryKey and TagsColumn2 are the table columns denoting the
	// primary key for the tags relation (M2M).
	TagsPrimaryKey = []string{"tag_id", "resource_pool_id"}
)

var (
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultDealocationSafetyPeriod holds the default value on creation for the dealocation_safety_period field.
	DefaultDealocationSafetyPeriod int
)

// PoolType defines the type for the pool_type enum field.
type PoolType string

// PoolType values.
const (
	PoolTypeSingleton  PoolType = "singleton"
	PoolTypeSet        PoolType = "set"
	PoolTypeAllocating PoolType = "allocating"
)

func (pt PoolType) String() string {
	return string(pt)
}

// PoolTypeValidator is a validator for the "pool_type" field enum values. It is called by the builders before save.
func PoolTypeValidator(pt PoolType) error {
	switch pt {
	case PoolTypeSingleton, PoolTypeSet, PoolTypeAllocating:
		return nil
	default:
		return fmt.Errorf("resourcepool: invalid enum value for pool_type field: %q", pt)
	}
}

// MarshalGQL implements graphql.Marshaler interface.
func (pt PoolType) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(pt.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (pt *PoolType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enum %T must be a string", v)
	}
	*pt = PoolType(str)
	if err := PoolTypeValidator(*pt); err != nil {
		return fmt.Errorf("%s is not a valid PoolType", str)
	}
	return nil
}
