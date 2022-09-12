// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// AllocationStrategiesColumns holds the columns for the "allocation_strategies" table.
	AllocationStrategiesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "description", Type: field.TypeString, Nullable: true, Size: 2147483647},
		{Name: "lang", Type: field.TypeEnum, Enums: []string{"py", "js", "go"}, Default: "js"},
		{Name: "script", Type: field.TypeString, Size: 2147483647},
	}
	// AllocationStrategiesTable holds the schema information for the "allocation_strategies" table.
	AllocationStrategiesTable = &schema.Table{
		Name:       "allocation_strategies",
		Columns:    AllocationStrategiesColumns,
		PrimaryKey: []*schema.Column{AllocationStrategiesColumns[0]},
	}
	// PoolPropertiesColumns holds the columns for the "pool_properties" table.
	PoolPropertiesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "resource_pool_pool_properties", Type: field.TypeInt, Unique: true, Nullable: true},
	}
	// PoolPropertiesTable holds the schema information for the "pool_properties" table.
	PoolPropertiesTable = &schema.Table{
		Name:       "pool_properties",
		Columns:    PoolPropertiesColumns,
		PrimaryKey: []*schema.Column{PoolPropertiesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "pool_properties_resource_pools_poolProperties",
				Columns:    []*schema.Column{PoolPropertiesColumns[1]},
				RefColumns: []*schema.Column{ResourcePoolsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "poolproperties_resource_pool_pool_properties",
				Unique:  false,
				Columns: []*schema.Column{PoolPropertiesColumns[1]},
			},
		},
	}
	// PropertiesColumns holds the columns for the "properties" table.
	PropertiesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "int_val", Type: field.TypeInt, Nullable: true},
		{Name: "bool_val", Type: field.TypeBool, Nullable: true},
		{Name: "float_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "latitude_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "longitude_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "range_from_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "range_to_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "string_val", Type: field.TypeString, Nullable: true},
		{Name: "pool_properties_properties", Type: field.TypeInt, Nullable: true},
		{Name: "property_type", Type: field.TypeInt},
		{Name: "resource_properties", Type: field.TypeInt, Nullable: true},
	}
	// PropertiesTable holds the schema information for the "properties" table.
	PropertiesTable = &schema.Table{
		Name:       "properties",
		Columns:    PropertiesColumns,
		PrimaryKey: []*schema.Column{PropertiesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "properties_pool_properties_properties",
				Columns:    []*schema.Column{PropertiesColumns[9]},
				RefColumns: []*schema.Column{PoolPropertiesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "properties_property_types_type",
				Columns:    []*schema.Column{PropertiesColumns[10]},
				RefColumns: []*schema.Column{PropertyTypesColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "properties_resources_properties",
				Columns:    []*schema.Column{PropertiesColumns[11]},
				RefColumns: []*schema.Column{ResourcesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "property_resource_properties",
				Unique:  false,
				Columns: []*schema.Column{PropertiesColumns[11]},
			},
			{
				Name:    "property_property_type",
				Unique:  false,
				Columns: []*schema.Column{PropertiesColumns[10]},
			},
			{
				Name:    "property_int_val",
				Unique:  false,
				Columns: []*schema.Column{PropertiesColumns[1]},
			},
		},
	}
	// PropertyTypesColumns holds the columns for the "property_types" table.
	PropertyTypesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "type", Type: field.TypeEnum, Enums: []string{"string", "int", "bool", "float", "date", "enum", "range", "email", "gps_location", "datetime_local", "node"}},
		{Name: "name", Type: field.TypeString},
		{Name: "external_id", Type: field.TypeString, Unique: true, Nullable: true},
		{Name: "index", Type: field.TypeInt, Nullable: true},
		{Name: "category", Type: field.TypeString, Nullable: true},
		{Name: "int_val", Type: field.TypeInt, Nullable: true},
		{Name: "bool_val", Type: field.TypeBool, Nullable: true},
		{Name: "float_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "latitude_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "longitude_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "string_val", Type: field.TypeString, Nullable: true, Size: 2147483647},
		{Name: "range_from_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "range_to_val", Type: field.TypeFloat64, Nullable: true},
		{Name: "is_instance_property", Type: field.TypeBool, Default: true},
		{Name: "editable", Type: field.TypeBool, Default: true},
		{Name: "mandatory", Type: field.TypeBool, Default: false},
		{Name: "deleted", Type: field.TypeBool, Default: false},
		{Name: "node_type", Type: field.TypeString, Nullable: true},
		{Name: "allocation_strategy_pool_property_types", Type: field.TypeInt, Nullable: true},
		{Name: "resource_type_property_types", Type: field.TypeInt, Nullable: true},
	}
	// PropertyTypesTable holds the schema information for the "property_types" table.
	PropertyTypesTable = &schema.Table{
		Name:       "property_types",
		Columns:    PropertyTypesColumns,
		PrimaryKey: []*schema.Column{PropertyTypesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "property_types_allocation_strategies_pool_property_types",
				Columns:    []*schema.Column{PropertyTypesColumns[19]},
				RefColumns: []*schema.Column{AllocationStrategiesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "property_types_resource_types_property_types",
				Columns:    []*schema.Column{PropertyTypesColumns[20]},
				RefColumns: []*schema.Column{ResourceTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// ResourcesColumns holds the columns for the "resources" table.
	ResourcesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"free", "claimed", "retired", "bench"}},
		{Name: "description", Type: field.TypeString, Nullable: true, Size: 2147483647},
		{Name: "alternate_id", Type: field.TypeJSON, Nullable: true},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "resource_pool_claims", Type: field.TypeInt, Nullable: true},
	}
	// ResourcesTable holds the schema information for the "resources" table.
	ResourcesTable = &schema.Table{
		Name:       "resources",
		Columns:    ResourcesColumns,
		PrimaryKey: []*schema.Column{ResourcesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "resources_resource_pools_claims",
				Columns:    []*schema.Column{ResourcesColumns[5]},
				RefColumns: []*schema.Column{ResourcePoolsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "resource_resource_pool_claims",
				Unique:  false,
				Columns: []*schema.Column{ResourcesColumns[5]},
			},
		},
	}
	// ResourcePoolsColumns holds the columns for the "resource_pools" table.
	ResourcePoolsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "description", Type: field.TypeString, Nullable: true, Size: 2147483647},
		{Name: "pool_type", Type: field.TypeEnum, Enums: []string{"singleton", "set", "allocating"}},
		{Name: "dealocation_safety_period", Type: field.TypeInt, Default: 0},
		{Name: "resource_nested_pool", Type: field.TypeInt, Unique: true, Nullable: true},
		{Name: "resource_pool_allocation_strategy", Type: field.TypeInt, Nullable: true},
		{Name: "resource_type_pools", Type: field.TypeInt, Nullable: true},
	}
	// ResourcePoolsTable holds the schema information for the "resource_pools" table.
	ResourcePoolsTable = &schema.Table{
		Name:       "resource_pools",
		Columns:    ResourcePoolsColumns,
		PrimaryKey: []*schema.Column{ResourcePoolsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "resource_pools_resources_nested_pool",
				Columns:    []*schema.Column{ResourcePoolsColumns[5]},
				RefColumns: []*schema.Column{ResourcesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "resource_pools_allocation_strategies_allocation_strategy",
				Columns:    []*schema.Column{ResourcePoolsColumns[6]},
				RefColumns: []*schema.Column{AllocationStrategiesColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "resource_pools_resource_types_pools",
				Columns:    []*schema.Column{ResourcePoolsColumns[7]},
				RefColumns: []*schema.Column{ResourceTypesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "resourcepool_resource_pool_allocation_strategy",
				Unique:  false,
				Columns: []*schema.Column{ResourcePoolsColumns[6]},
			},
		},
	}
	// ResourceTypesColumns holds the columns for the "resource_types" table.
	ResourceTypesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString},
	}
	// ResourceTypesTable holds the schema information for the "resource_types" table.
	ResourceTypesTable = &schema.Table{
		Name:       "resource_types",
		Columns:    ResourceTypesColumns,
		PrimaryKey: []*schema.Column{ResourceTypesColumns[0]},
	}
	// TagsColumns holds the columns for the "tags" table.
	TagsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "tag", Type: field.TypeString, Unique: true},
	}
	// TagsTable holds the schema information for the "tags" table.
	TagsTable = &schema.Table{
		Name:       "tags",
		Columns:    TagsColumns,
		PrimaryKey: []*schema.Column{TagsColumns[0]},
	}
	// PoolPropertiesResourceTypeColumns holds the columns for the "pool_properties_resourceType" table.
	PoolPropertiesResourceTypeColumns = []*schema.Column{
		{Name: "pool_properties_id", Type: field.TypeInt},
		{Name: "resource_type_id", Type: field.TypeInt},
	}
	// PoolPropertiesResourceTypeTable holds the schema information for the "pool_properties_resourceType" table.
	PoolPropertiesResourceTypeTable = &schema.Table{
		Name:       "pool_properties_resourceType",
		Columns:    PoolPropertiesResourceTypeColumns,
		PrimaryKey: []*schema.Column{PoolPropertiesResourceTypeColumns[0], PoolPropertiesResourceTypeColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "pool_properties_resourceType_pool_properties_id",
				Columns:    []*schema.Column{PoolPropertiesResourceTypeColumns[0]},
				RefColumns: []*schema.Column{PoolPropertiesColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "pool_properties_resourceType_resource_type_id",
				Columns:    []*schema.Column{PoolPropertiesResourceTypeColumns[1]},
				RefColumns: []*schema.Column{ResourceTypesColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// TagPoolsColumns holds the columns for the "tag_pools" table.
	TagPoolsColumns = []*schema.Column{
		{Name: "tag_id", Type: field.TypeInt},
		{Name: "resource_pool_id", Type: field.TypeInt},
	}
	// TagPoolsTable holds the schema information for the "tag_pools" table.
	TagPoolsTable = &schema.Table{
		Name:       "tag_pools",
		Columns:    TagPoolsColumns,
		PrimaryKey: []*schema.Column{TagPoolsColumns[0], TagPoolsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "tag_pools_tag_id",
				Columns:    []*schema.Column{TagPoolsColumns[0]},
				RefColumns: []*schema.Column{TagsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "tag_pools_resource_pool_id",
				Columns:    []*schema.Column{TagPoolsColumns[1]},
				RefColumns: []*schema.Column{ResourcePoolsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		AllocationStrategiesTable,
		PoolPropertiesTable,
		PropertiesTable,
		PropertyTypesTable,
		ResourcesTable,
		ResourcePoolsTable,
		ResourceTypesTable,
		TagsTable,
		PoolPropertiesResourceTypeTable,
		TagPoolsTable,
	}
)

func init() {
	PoolPropertiesTable.ForeignKeys[0].RefTable = ResourcePoolsTable
	PropertiesTable.ForeignKeys[0].RefTable = PoolPropertiesTable
	PropertiesTable.ForeignKeys[1].RefTable = PropertyTypesTable
	PropertiesTable.ForeignKeys[2].RefTable = ResourcesTable
	PropertyTypesTable.ForeignKeys[0].RefTable = AllocationStrategiesTable
	PropertyTypesTable.ForeignKeys[1].RefTable = ResourceTypesTable
	ResourcesTable.ForeignKeys[0].RefTable = ResourcePoolsTable
	ResourcePoolsTable.ForeignKeys[0].RefTable = ResourcesTable
	ResourcePoolsTable.ForeignKeys[1].RefTable = AllocationStrategiesTable
	ResourcePoolsTable.ForeignKeys[2].RefTable = ResourceTypesTable
	PoolPropertiesResourceTypeTable.ForeignKeys[0].RefTable = PoolPropertiesTable
	PoolPropertiesResourceTypeTable.ForeignKeys[1].RefTable = ResourceTypesTable
	TagPoolsTable.ForeignKeys[0].RefTable = TagsTable
	TagPoolsTable.ForeignKeys[1].RefTable = ResourcePoolsTable
}
