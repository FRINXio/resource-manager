// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"github.com/net-auto/resourceManager/ent"
	"github.com/net-auto/resourceManager/ent/allocationstrategy"
)

// Input parameters for creating an allocation pool
type CreateAllocatingPoolInput struct {
	AllocationStrategyID        int                    `json:"allocationStrategyId"`
	Description                 *string                `json:"description"`
	PoolDealocationSafetyPeriod int                    `json:"poolDealocationSafetyPeriod"`
	PoolName                    string                 `json:"poolName"`
	PoolProperties              map[string]interface{} `json:"poolProperties"`
	PoolPropertyTypes           map[string]interface{} `json:"poolPropertyTypes"`
	ResourceTypeID              int                    `json:"resourceTypeId"`
	Tags                        []string               `json:"tags"`
}

// Output of creating an allocating pool
type CreateAllocatingPoolPayload struct {
	Pool *ent.ResourcePool `json:"pool"`
}

// Input parameters for creating a new allocation strategy
type CreateAllocationStrategyInput struct {
	Name        string                  `json:"name"`
	Description *string                 `json:"description"`
	Script      string                  `json:"script"`
	Lang        allocationstrategy.Lang `json:"lang"`
}

// Output of creating a new allocation strategy
type CreateAllocationStrategyPayload struct {
	Strategy *ent.AllocationStrategy `json:"strategy"`
}

// Input parameters for creating a nested allocation pool
type CreateNestedAllocatingPoolInput struct {
	AllocationStrategyID        int      `json:"allocationStrategyId"`
	Description                 *string  `json:"description"`
	ParentResourceID            int      `json:"parentResourceId"`
	PoolDealocationSafetyPeriod int      `json:"poolDealocationSafetyPeriod"`
	PoolName                    string   `json:"poolName"`
	ResourceTypeID              int      `json:"resourceTypeId"`
	Tags                        []string `json:"tags"`
}

// Output of creating a nested allocating pool
type CreateNestedAllocatingPoolPayload struct {
	Pool *ent.ResourcePool `json:"pool"`
}

// Input parameters for creating a nested set pool
type CreateNestedSetPoolInput struct {
	Description                 *string                  `json:"description"`
	ParentResourceID            int                      `json:"parentResourceId"`
	PoolDealocationSafetyPeriod int                      `json:"poolDealocationSafetyPeriod"`
	PoolName                    string                   `json:"poolName"`
	PoolValues                  []map[string]interface{} `json:"poolValues"`
	ResourceTypeID              int                      `json:"resourceTypeId"`
	Tags                        []string                 `json:"tags"`
}

// Output of creating a nested set pool
type CreateNestedSetPoolPayload struct {
	Pool *ent.ResourcePool `json:"pool"`
}

// Input parameters for creating a nested singleton pool
type CreateNestedSingletonPoolInput struct {
	Description      *string                  `json:"description"`
	ParentResourceID int                      `json:"parentResourceId"`
	PoolName         string                   `json:"poolName"`
	PoolValues       []map[string]interface{} `json:"poolValues"`
	ResourceTypeID   int                      `json:"resourceTypeId"`
	Tags             []string                 `json:"tags"`
}

// Output of creating a nested singleton pool
type CreateNestedSingletonPoolPayload struct {
	Pool *ent.ResourcePool `json:"pool"`
}

// Creating a new resource-type
type CreateResourceTypeInput struct {
	// name of the resource type AND property type (should they be different?)
	ResourceName string `json:"resourceName"`
	// resourceProperties: Map! - for key "init" the value is the initial value of the property type (like 7)
	//                          - for key "type" the value is the name of the type like "int"
	ResourceProperties map[string]interface{} `json:"resourceProperties"`
}

// Output of creating a new resource-type
type CreateResourceTypePayload struct {
	ResourceType *ent.ResourceType `json:"resourceType"`
}

// Input parameters for creating a set pool
type CreateSetPoolInput struct {
	Description                 *string                  `json:"description"`
	PoolDealocationSafetyPeriod int                      `json:"poolDealocationSafetyPeriod"`
	PoolName                    string                   `json:"poolName"`
	PoolValues                  []map[string]interface{} `json:"poolValues"`
	ResourceTypeID              int                      `json:"resourceTypeId"`
	Tags                        []string                 `json:"tags"`
}

// Output of creating set pool
type CreateSetPoolPayload struct {
	Pool *ent.ResourcePool `json:"pool"`
}

// Input parameters for creating a singleton pool
type CreateSingletonPoolInput struct {
	Description    *string                  `json:"description"`
	PoolName       string                   `json:"poolName"`
	PoolValues     []map[string]interface{} `json:"poolValues"`
	ResourceTypeID int                      `json:"resourceTypeId"`
	Tags           []string                 `json:"tags"`
}

// Output of creating a singleton pool
type CreateSingletonPoolPayload struct {
	Pool *ent.ResourcePool `json:"pool"`
}

// Input parameters for creating a new tag
type CreateTagInput struct {
	TagText string `json:"tagText"`
}

// Output of creating a tag
type CreateTagPayload struct {
	Tag *ent.Tag `json:"tag"`
}

// Input parameters for deleting an existing allocation strategy
type DeleteAllocationStrategyInput struct {
	AllocationStrategyID int `json:"allocationStrategyId"`
}

// Output of deleting an existing allocation strategy
type DeleteAllocationStrategyPayload struct {
	Strategy *ent.AllocationStrategy `json:"strategy"`
}

// Input entity for deleting a pool
type DeleteResourcePoolInput struct {
	ResourcePoolID int `json:"resourcePoolId"`
}

// Output entity for deleting a pool
type DeleteResourcePoolPayload struct {
	ResourcePoolID int `json:"resourcePoolId"`
}

// Input parameters for deleting an existing resource-type
type DeleteResourceTypeInput struct {
	ResourceTypeID int `json:"resourceTypeId"`
}

// Output of deleting a resource-type
type DeleteResourceTypePayload struct {
	ResourceTypeID int `json:"resourceTypeId"`
}

// Input parameters for deleting an existing tag
type DeleteTagInput struct {
	TagID int `json:"tagId"`
}

// Output of deleting a tag
type DeleteTagPayload struct {
	TagID int `json:"tagId"`
}

// Entity representing capacity of a pool
type PoolCapacityPayload struct {
	FreeCapacity     float64 `json:"freeCapacity"`
	UtilizedCapacity float64 `json:"utilizedCapacity"`
}

// Alternative representation of identity of a resource (i.e. alternative to resource ID)
type ResourceInput struct {
	Properties map[string]interface{} `json:"Properties"`
	Status     string                 `json:"Status"`
	UpdatedAt  string                 `json:"UpdatedAt"`
}

// Convenience entity representing the identity of a pool in some calls
type ResourcePoolInput struct {
	ResourcePoolName string                 `json:"ResourcePoolName"`
	PoolProperties   map[string]interface{} `json:"poolProperties"`
}

// Helper entities for tag search
type TagAnd struct {
	MatchesAll []string `json:"matchesAll"`
}

// Helper entities for tag search
type TagOr struct {
	MatchesAny []*TagAnd `json:"matchesAny"`
}

// Input parameters for a call adding a tag to pool
type TagPoolInput struct {
	TagID  int `json:"tagId"`
	PoolID int `json:"poolId"`
}

// Output of adding a specific tag to a pool
type TagPoolPayload struct {
	Tag *ent.Tag `json:"tag"`
}

// Input parameters for a call removing a tag from pool
type UntagPoolInput struct {
	TagID  int `json:"tagId"`
	PoolID int `json:"poolId"`
}

// Output of removing a specific tag from a pool
type UntagPoolPayload struct {
	Tag *ent.Tag `json:"tag"`
}

// Output of updating the alternative id of a resource
type UpdateResourceAltID struct {
	AlternativeID map[string]interface{} `json:"AlternativeId"`
	Resource      map[string]interface{} `json:"Resource"`
}

// Input parameters updating the name of a resource-type
type UpdateResourceTypeNameInput struct {
	ResourceTypeID int    `json:"resourceTypeId"`
	ResourceName   string `json:"resourceName"`
}

// Output of updating the name of a resource-type
type UpdateResourceTypeNamePayload struct {
	ResourceTypeID int `json:"resourceTypeId"`
}

// Input parameters for updating an existing tag
type UpdateTagInput struct {
	TagID   int    `json:"tagId"`
	TagText string `json:"tagText"`
}

// Output of updating a tag
type UpdateTagPayload struct {
	Tag *ent.Tag `json:"tag"`
}
